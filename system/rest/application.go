package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Application struct {
		application applicationService
		ac          applicationAccessController
	}

	applicationService interface {
		LookupByID(ctx context.Context, ID uint64) (app *types.Application, err error)
		Search(ctx context.Context, filter types.ApplicationFilter) (aa types.ApplicationSet, f types.ApplicationFilter, err error)
		Create(ctx context.Context, new *types.Application) (app *types.Application, err error)
		Update(ctx context.Context, upd *types.Application) (app *types.Application, err error)
		Delete(ctx context.Context, ID uint64) (err error)
		Undelete(ctx context.Context, ID uint64) (err error)
	}

	applicationAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool
	}

	applicationPayload struct {
		*types.Application

		CanGrant             bool `json:"canGrant"`
		CanUpdateApplication bool `json:"canUpdateApplication"`
		CanDeleteApplication bool `json:"canDeleteApplication"`
	}

	applicationSetPayload struct {
		Filter types.ApplicationFilter `json:"filter"`
		Set    []*applicationPayload   `json:"set"`
	}
)

func (Application) New() *Application {
	return &Application{
		application: service.DefaultApplication,
		ac:          service.DefaultAccessControl,
	}
}

func (ctrl *Application) List(ctx context.Context, r *request.ApplicationList) (interface{}, error) {
	var (
		err error
		f   = types.ApplicationFilter{
			Name:  r.Name,
			Query: r.Query,

			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.application.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Application) Create(ctx context.Context, r *request.ApplicationCreate) (interface{}, error) {
	var (
		err error
		app = &types.Application{
			Name:    r.Name,
			Enabled: r.Enabled,
		}
	)

	if r.Unify != nil {
		app.Unify = &types.ApplicationUnify{}
		if err := r.Unify.Unmarshal(app.Unify); err != nil {
			return nil, err
		}
	}

	app, err = ctrl.application.Create(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Application) Update(ctx context.Context, r *request.ApplicationUpdate) (interface{}, error) {
	var (
		err error
		app = &types.Application{
			ID:      r.ApplicationID,
			Name:    r.Name,
			Enabled: r.Enabled,
		}
	)

	if r.Unify != nil {
		app.Unify = &types.ApplicationUnify{}
		if err := r.Unify.Unmarshal(app.Unify); err != nil {
			return nil, err
		}
	}

	app, err = ctrl.application.Update(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Application) Read(ctx context.Context, r *request.ApplicationRead) (interface{}, error) {
	app, err := ctrl.application.LookupByID(ctx, r.ApplicationID)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Application) Delete(ctx context.Context, r *request.ApplicationDelete) (interface{}, error) {
	return api.OK(), ctrl.application.Delete(ctx, r.ApplicationID)
}

func (ctrl *Application) Undelete(ctx context.Context, r *request.ApplicationUndelete) (interface{}, error) {
	return api.OK(), ctrl.application.Undelete(ctx, r.ApplicationID)
}

func (ctrl *Application) TriggerScript(ctx context.Context, r *request.ApplicationTriggerScript) (rsp interface{}, err error) {
	var (
		application *types.Application
	)

	if application, err = ctrl.application.LookupByID(ctx, r.ApplicationID); err != nil {
		return
	}

	// @todo implement same behaviour as we have on record - Application+oldApplication
	err = corredor.Service().Exec(ctx, r.Script, event.ApplicationOnManual(application, application))
	return application, err
}

func (ctrl Application) makePayload(ctx context.Context, m *types.Application, err error) (*applicationPayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &applicationPayload{
		Application: m,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateApplication: ctrl.ac.CanUpdateApplication(ctx, m),
		CanDeleteApplication: ctrl.ac.CanDeleteApplication(ctx, m),
	}, nil
}

func (ctrl Application) makeFilterPayload(ctx context.Context, nn types.ApplicationSet, f types.ApplicationFilter, err error) (*applicationSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &applicationSetPayload{Filter: f, Set: make([]*applicationPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
