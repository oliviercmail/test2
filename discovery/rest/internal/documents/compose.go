package documents

import (
	"context"
	"fmt"
	cmpService "github.com/cortezaproject/corteza-server/compose/service"
	cmpTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/discovery/service"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	composeResources struct {
		opt      options.DiscoveryOpt
		settings *sysTypes.AppSettings

		rbac interface {
			SignificantRoles(res rbac.Resource, op string) (aRR, dRR []uint64)
		}

		ac interface {
			CanReadModule(ctx context.Context, r *cmpTypes.Module) bool
			CanReadNamespace(ctx context.Context, r *cmpTypes.Namespace) bool
			CanReadRecord(ctx context.Context, r *cmpTypes.Record) bool
			CanReadRecordValue(ctx context.Context, r *cmpTypes.ModuleField) bool
			CanReadChart(ctx context.Context, r *cmpTypes.Chart) bool
			CanReadPage(ctx context.Context, r *cmpTypes.Page) bool
		}

		ns interface {
			FindByID(context.Context, uint64) (*cmpTypes.Namespace, error)
			Find(context.Context, cmpTypes.NamespaceFilter) (cmpTypes.NamespaceSet, cmpTypes.NamespaceFilter, error)
		}

		mod interface {
			FindByID(context.Context, uint64, uint64) (*cmpTypes.Module, error)
			Find(ctx context.Context, filter cmpTypes.ModuleFilter) (set cmpTypes.ModuleSet, f cmpTypes.ModuleFilter, err error)
		}

		rec interface {
			Find(ctx context.Context, filter cmpTypes.RecordFilter) (set cmpTypes.RecordSet, f cmpTypes.RecordFilter, err error)
		}

		page interface {
			Find(ctx context.Context, filter cmpTypes.PageFilter) (set cmpTypes.PageSet, f cmpTypes.PageFilter, err error)
		}
	}

	pageDetail struct {
		namespaceSlug string
		moduleID      uint64
		pageID        uint64
		recordID      uint64
	}
)

var (
	modulePages map[uint64]pageDetail
)

func ComposeResources() *composeResources {
	return &composeResources{
		opt:      service.DefaultOption,
		settings: sysService.CurrentSettings,
		rbac:     rbac.Global(),
		ac:       cmpService.DefaultAccessControl,
		ns:       cmpService.DefaultNamespace,
		mod:      cmpService.DefaultModule,
		rec:      cmpService.DefaultRecord,
		page:     cmpService.DefaultPage,
	}
}

func (d composeResources) Namespaces(ctx context.Context, limit uint, cur string) (rsp *Response, err error) {
	return rsp, func() (err error) {
		if !d.settings.Discovery.ComposeNamespaces.Enabled {
			return errors.Internal("compose namespace indexing disabled")
		}

		var (
			resourceType = "compose:namespace" // @todo use RBAC resourceType
			ps           cmpTypes.PageSet
			pf           cmpTypes.PageFilter
			nss          cmpTypes.NamespaceSet
			f            = cmpTypes.NamespaceFilter{
				Deleted: filter.StateInclusive,
			}
		)

		if f.Paging, err = filter.NewPaging(limit, cur); err != nil {
			return err
		}

		if nss, f, err = d.ns.Find(ctx, f); err != nil {
			return err
		}

		rsp = &Response{
			Documents: make([]Document, len(nss)),
			Filter: Filter{
				Limit:    limit,
				NextPage: f.NextPage,
			},
		}

		for i, ns := range nss {
			nsID := ns.ID
			nsSlug := ns.Slug
			rsp.Documents[i].ID = nsID
			pf.NamespaceID = nsID

			// namespace pages
			if ps, pf, err = d.page.Find(ctx, pf); err != nil {
				return err
			}

			if modulePages == nil {
				modulePages = make(map[uint64]pageDetail)
			}
			for _, p := range ps {
				moduleID := p.ModuleID
				if len(nsSlug) > 0 && moduleID > 0 {
					modulePages[moduleID] = pageDetail{
						namespaceSlug: nsSlug,
						moduleID:      moduleID,
						pageID:        p.ID,
					}
				}
			}

			// where should this link to?
			// namespace root page on the compose?
			// rsp.Documents[i].URL = ""  // added to source
			doc := &docComposeNamespace{
				ResourceType: resourceType,
				NamespaceID:  nsID,
				Name:         ns.Name,
				Handle:       nsSlug,
				Url:          d.getUrlToResource(pageDetail{namespaceSlug: nsSlug}),
				Meta: docPartialComposeNamespaceMeta{
					Subtitle:    ns.Meta.Subtitle,
					Description: ns.Meta.Description,
				},
				Created: makePartialChange(&ns.CreatedAt),
				Updated: makePartialChange(ns.UpdatedAt),
				Deleted: makePartialChange(ns.DeletedAt),
			}

			doc.Security.AllowedRoles, doc.Security.DeniedRoles = d.rbac.SignificantRoles(ns, "read")

			rsp.Documents[i].Source = doc
		}

		return nil
	}()
}

func (d composeResources) Modules(ctx context.Context, namespaceID uint64, limit uint, cur string) (rsp *Response, err error) {
	return rsp, func() (err error) {
		if !d.settings.Discovery.ComposeModules.Enabled {
			return errors.Internal("compose module indexing disabled")
		}

		var (
			ns *cmpTypes.Namespace
			mm cmpTypes.ModuleSet
			f  = cmpTypes.ModuleFilter{
				NamespaceID: namespaceID,
				Deleted:     filter.StateInclusive,
			}
		)

		if f.Paging, err = filter.NewPaging(limit, cur); err != nil {
			return
		}

		if mm, f, err = d.mod.Find(ctx, f); err != nil {
			return
		}

		rsp = &Response{
			Documents: make([]Document, len(mm)),
			Filter: Filter{
				Limit:    limit,
				NextPage: f.NextPage,
			},
		}

		if ns, err = d.ns.FindByID(ctx, namespaceID); err != nil {
			return
		}

		nsPartial := docPartialComposeNamespace{
			NamespaceID: namespaceID,
			Name:        ns.Name,
			Handle:      ns.Slug,
		}

		for i, mod := range mm {
			rsp.Documents[i].ID = mod.ID
			// Where should this link to?
			// module edit screen in the administration? does this make sense?
			//rsp.Documents[i].URL = "@todo"
			doc := &docComposeModule{
				ResourceType: "compose:module", // @todo use RBAC resourceType
				ModuleID:     mod.ID,
				Name:         mod.Name,
				Handle:       mod.Handle,
				Namespace:    nsPartial,
				Fields: func() []*docPartialComposeModuleField {
					out := make([]*docPartialComposeModuleField, len(mod.Fields))
					for i, f := range mod.Fields {
						out[i] = &docPartialComposeModuleField{
							Name:  f.Name,
							Label: f.Label,
						}
					}
					return out
				}(),
				Created: makePartialChange(&mod.CreatedAt),
				Updated: makePartialChange(mod.UpdatedAt),
				Deleted: makePartialChange(mod.DeletedAt),
			}

			doc.Security.AllowedRoles, doc.Security.DeniedRoles = d.rbac.SignificantRoles(mod, "read")

			rsp.Documents[i].Source = doc
		}

		return nil
	}()
}

func (d composeResources) Records(ctx context.Context, namespaceID, moduleID uint64, limit uint, cur string) (rsp *Response, err error) {
	return rsp, func() (err error) {
		if !d.settings.Discovery.ComposeRecords.Enabled {
			return errors.Internal("compose record indexing disabled")
		}

		var (
			ns  *cmpTypes.Namespace
			mod *cmpTypes.Module
			rr  cmpTypes.RecordSet
			f   = cmpTypes.RecordFilter{
				NamespaceID: namespaceID,
				ModuleID:    moduleID,
				Deleted:     filter.StateInclusive,
			}
		)

		if f.Paging, err = filter.NewPaging(limit, cur); err != nil {
			return err
		}

		if rr, f, err = d.rec.Find(ctx, f); err != nil {
			return err
		}

		rsp = &Response{
			Documents: make([]Document, len(rr)),
			Filter: Filter{
				Limit:    limit,
				NextPage: f.NextPage,
			},
		}

		// @todo handle unreadable (access-control) namespaces
		if ns, err = d.ns.FindByID(ctx, namespaceID); err != nil {
			return
		}

		nsPartial := docPartialComposeNamespace{
			NamespaceID: namespaceID,
			Name:        ns.Name,
			Handle:      ns.Slug,
		}

		// @todo handle unreadable (access-control) modules
		if mod, err = d.mod.FindByID(ctx, namespaceID, moduleID); err != nil {
			return
		}

		modPartial := docPartialComposeModule{
			ModuleID: f.ModuleID,
			Name:     mod.Name,
			Handle:   mod.Handle,
		}

		for i, rec := range rr {
			recordID := rec.ID
			rsp.Documents[i].ID = recordID
			// where should this link to? record page in the compose?
			// rsp.Documents[i].URL = "" // added to source
			doc := &docComposeRecord{
				ResourceType: "compose:record", // @todo use RBAC resourceType
				RecordID:     recordID,
				Namespace:    nsPartial,
				Module:       modPartial,
				Url:          d.getUrlToResource(pageDetail{moduleID: moduleID, recordID: recordID}),
				Values:       d.recordValues(ctx, rec),
				Created:      makePartialChange(&rec.CreatedAt),
				Updated:      makePartialChange(rec.UpdatedAt),
				Deleted:      makePartialChange(rec.DeletedAt),
			}

			doc.Security.AllowedRoles, doc.Security.DeniedRoles = d.rbac.SignificantRoles(rec, "read")

			rsp.Documents[i].Source = doc
		}

		return nil
	}()
}

func (d composeResources) recordValues(ctx context.Context, rec *cmpTypes.Record) map[string][]interface{} {
	var (
		rval = make(map[string][]interface{})
	)

	if rec.GetModule() == nil {
		return nil
	}

	_ = rec.GetModule().Fields.Walk(func(f *cmpTypes.ModuleField) error {
		if !d.ac.CanReadRecordValue(ctx, f) {
			return nil
		}

		var (
			rv = rec.Values.FilterByName(f.Name)
			vv = make([]interface{}, 0, len(rv))
		)

		if len(rv) == 0 {
			return nil
		}

		for _, val := range rv {
			// refs needs to be casted to string (json & unsigned 64-bit integers)!
			if f.IsRef() {
				vv = append(vv, fmt.Sprintf("%d", val.Ref))
				continue
			}

			if tmp, err := val.Cast(f); err == nil {
				vv = append(vv, tmp)
			}

		}

		if len(vv) == 0 {
			return nil
		}

		rval[f.Name] = vv

		return nil
	})

	return rval
}

// getUrlToResource construct page url for compose resources
func (d composeResources) getUrlToResource(page pageDetail) (url string) {
	var (
		host          = d.opt.CortezaDomain
		validNsSlung  = len(page.namespaceSlug) > 0
		validModuleID = page.moduleID > 0
		validRecord   = page.recordID > 0
		validPageID   = page.pageID > 0
	)

	if len(host) == 0 {
		return
	}
	if validModuleID && validRecord {
		// For record, take moduleID & recordID from params(page) then construct url from modulePages map
		if p, is := modulePages[page.moduleID]; is {
			url = fmt.Sprintf("%s/compose/ns/%s/pages/%d/records/%d", host, p.namespaceSlug, p.pageID, page.recordID)
		}
	} else if validNsSlung && validPageID {
		// for module? @todo, we need moduleID ref for module home page
		// we can check page.self_id is 0 that is modules home pages, but ref_module is also 0
		url = fmt.Sprintf("%s/compose/ns/%s/pages/%d", host, page.namespaceSlug, page.pageID)
	} else if validNsSlung {
		// For namespace, just send to /pages FE handles fallback to home(default) page
		url = fmt.Sprintf("%s/compose/ns/%s/pages", host, page.namespaceSlug)
	}
	return
}
