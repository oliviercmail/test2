package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	dataPrivacy struct {
		actionlog actionlog.Recorder

		ac       dataPrivacyAccessController
		eventbus eventDispatcher

		store store.Storer
	}

	dataPrivacyAccessController interface {
		CanSearchDataPrivacyRequests(context.Context) bool
		CanCreateDataPrivacyRequest(context.Context) bool
		CanReadDataPrivacyRequest(context.Context, *types.DataPrivacyRequest) bool
		CanApproveDataPrivacyRequest(context.Context, *types.DataPrivacyRequest) bool
	}

	DataPrivacyService interface {
		FindRequestByID(ctx context.Context, requestID uint64) (*types.DataPrivacyRequest, error)
		FindRequests(context.Context, types.DataPrivacyRequestFilter) (types.DataPrivacyRequestSet, types.DataPrivacyRequestFilter, error)
		CreateRequest(ctx context.Context, request *types.DataPrivacyRequest) (*types.DataPrivacyRequest, error)
		UpdateRequestStatus(ctx context.Context, request *types.DataPrivacyRequest) (*types.DataPrivacyRequest, error)

		FindRequestComments(ctx context.Context, filter types.DataPrivacyRequestCommentFilter) (rr types.DataPrivacyRequestCommentSet, f types.DataPrivacyRequestCommentFilter, err error)
		CreateRequestComment(ctx context.Context, new *types.DataPrivacyRequestComment) (r *types.DataPrivacyRequestComment, err error)
	}
)

func DataPrivacy(s store.Storer, ac dataPrivacyAccessController, al actionlog.Recorder, eb eventDispatcher) *dataPrivacy {
	return &dataPrivacy{
		actionlog: al,
		ac:        ac,
		eventbus:  eb,
		store:     s,
	}
}

func (svc dataPrivacy) FindRequestByID(ctx context.Context, requestID uint64) (r *types.DataPrivacyRequest, err error) {
	var (
		raProps = &dataPrivacyActionProps{dataPrivacyRequest: &types.DataPrivacyRequest{ID: requestID}}
	)

	err = func() error {
		if requestID == 0 {
			return DataPrivacyErrInvalidID()
		}

		r, err = store.LookupDataPrivacyRequestByID(ctx, svc.store, requestID)
		if r, err = svc.procRequest(ctx, r, err); err != nil {
			return err
		}

		raProps.setDataPrivacyRequest(r)

		if !svc.ac.CanReadDataPrivacyRequest(ctx, r) {
			return DataPrivacyErrNotAllowedToRead()
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, DataPrivacyActionLookup, err)
}

func (svc dataPrivacy) procRequest(_ context.Context, r *types.DataPrivacyRequest, err error) (*types.DataPrivacyRequest, error) {
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, DataPrivacyErrNotFound()
		}

		return nil, err
	}

	return r, nil
}

func (svc dataPrivacy) FindRequests(ctx context.Context, filter types.DataPrivacyRequestFilter) (rr types.DataPrivacyRequestSet, f types.DataPrivacyRequestFilter, err error) {
	var (
		raProps = &dataPrivacyActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(req *types.DataPrivacyRequest) (bool, error) {
		if !svc.ac.CanReadDataPrivacyRequest(ctx, req) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchDataPrivacyRequests(ctx) {
			return DataPrivacyErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchDataPrivacyRequests(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, DataPrivacyActionSearch, err)
}

func (svc dataPrivacy) CreateRequest(ctx context.Context, new *types.DataPrivacyRequest) (r *types.DataPrivacyRequest, err error) {
	var (
		raProps = &dataPrivacyActionProps{new: new}
	)

	err = func() (err error) {
		if len(new.Kind.String()) == 0 {
			return DataPrivacyErrInvalidKind()
		}

		if !svc.ac.CanCreateDataPrivacyRequest(ctx) {
			return DataPrivacyErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.DataPrivacyRequestBeforeCreate(new, r)); err != nil {
			return
		}

		new.ID = nextID()
		new.Status = types.RequestStatusPending
		new.RequestedAt = *now()
		new.RequestedBy = a.GetIdentityFromContext(ctx).Identity()
		new.CreatedAt = *now()

		if err = store.CreateDataPrivacyRequest(ctx, svc.store, new); err != nil {
			return
		}

		r = new

		svc.eventbus.Dispatch(ctx, event.DataPrivacyRequestAfterCreate(new, r))
		return
	}()

	return r, svc.recordAction(ctx, raProps, DataPrivacyActionCreate, err)
}

func (svc dataPrivacy) UpdateRequestStatus(ctx context.Context, upd *types.DataPrivacyRequest) (r *types.DataPrivacyRequest, err error) {
	var (
		raProps = &dataPrivacyActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return DataPrivacyErrInvalidID()
		}

		if len(upd.Status.String()) == 0 {
			return DataPrivacyErrInvalidStatus()
		}

		if upd.Status == types.RequestStatusPending {
			return DataPrivacyErrInvalidStatus()
		}

		if upd.Status == types.RequestStatusApproved || upd.Status == types.RequestStatusRejected {
			if !svc.ac.CanApproveDataPrivacyRequest(ctx, upd) {
				return DataPrivacyErrNotAllowedToApprove()
			}
		}

		if r, err = store.LookupDataPrivacyRequestByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setDataPrivacyRequest(r)

		if err = svc.eventbus.WaitFor(ctx, event.DataPrivacyRequestBeforeUpdate(upd, r)); err != nil {
			return
		}

		r.Status = upd.Status
		r.CompletedAt = now()
		r.CompletedBy = a.GetIdentityFromContext(ctx).Identity()
		r.UpdatedAt = now()

		// Assign changed values
		if err = store.UpdateDataPrivacyRequest(ctx, svc.store, r); err != nil {
			return err
		}

		svc.eventbus.Dispatch(ctx, event.DataPrivacyRequestAfterUpdate(upd, r))

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, DataPrivacyActionApprove, err)
}

func (svc dataPrivacy) FindRequestComments(ctx context.Context, filter types.DataPrivacyRequestCommentFilter) (rr types.DataPrivacyRequestCommentSet, f types.DataPrivacyRequestCommentFilter, err error) {
	err = func() error {
		if rr, f, err = store.SearchDataPrivacyRequestComments(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, err
}

func (svc dataPrivacy) CreateRequestComment(ctx context.Context, new *types.DataPrivacyRequestComment) (r *types.DataPrivacyRequestComment, err error) {
	err = func() (err error) {

		_, err = svc.FindRequestByID(ctx, new.RequestID)
		if err != nil {
			return
		}

		new.ID = nextID()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
		new.CreatedAt = *now()

		if err = store.CreateDataPrivacyRequestComment(ctx, svc.store, new); err != nil {
			return
		}

		r = new

		return
	}()

	return r, err
}
