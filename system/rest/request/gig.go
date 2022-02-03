package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/conv"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	GigCreate struct {
		// Worker POST parameter
		//
		// Gig worker
		Worker string

		// Preprocessors POST parameter
		//
		// Worker preprocessing to do
		Preprocessors conv.ParamWrapSet

		// Postprocessors POST parameter
		//
		// Output postprocessing to do
		Postprocessors conv.ParamWrapSet

		// Completion POST parameter
		//
		// Specifies when the gig is marked as completed
		Completion string
	}

	GigGo struct {
		// Worker POST parameter
		//
		// Gig worker
		Worker string

		// Preprocessors POST parameter
		//
		// Worker preprocessing to do
		Preprocessors conv.ParamWrapSet

		// Postprocessors POST parameter
		//
		// Output postprocessing to do
		Postprocessors conv.ParamWrapSet
	}

	GigRead struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigUpdate struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`

		// Decoders POST parameter
		//
		// Decoders to apply to the sources
		Decoders conv.ParamWrapSet

		// Preprocessors POST parameter
		//
		// Worker preprocessing to do
		Preprocessors conv.ParamWrapSet

		// Postprocessors POST parameter
		//
		// Output postprocessing to do
		Postprocessors conv.ParamWrapSet

		// Completion POST parameter
		//
		// Specifies when the gig is marked as completed
		Completion string
	}

	GigDelete struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigUndelete struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigAddSource struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`

		// Upload POST parameter
		//
		// File source to add
		Upload *multipart.FileHeader

		// Uri POST parameter
		//
		// Source location to add
		Uri string

		// Decoders POST parameter
		//
		// Decoders to apply to the sources
		Decoders conv.ParamWrapSet
	}

	GigRemoveSource struct {
		// GigID PATH parameter
		//
		// Gig ID
		GigID uint64 `json:",string"`

		// SourceID PATH parameter
		//
		// Source ID
		SourceID uint64 `json:",string"`
	}

	GigPrepare struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigExec struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigOutput struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigOutputAll struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigOutputSpecific struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`

		// SourceID PATH parameter
		//
		// Source ID
		SourceID uint64 `json:",string"`
	}

	GigState struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigStatus struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigComplete struct {
		// GigID PATH parameter
		//
		// ID
		GigID uint64 `json:",string"`
	}

	GigWorkers struct {
	}

	GigTasks struct {
	}
)

// NewGigCreate request
func NewGigCreate() *GigCreate {
	return &GigCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"worker":         r.Worker,
		"preprocessors":  r.Preprocessors,
		"postprocessors": r.Postprocessors,
		"completion":     r.Completion,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) GetWorker() string {
	return r.Worker
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) GetPreprocessors() conv.ParamWrapSet {
	return r.Preprocessors
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) GetPostprocessors() conv.ParamWrapSet {
	return r.Postprocessors
}

// Auditable returns all auditable/loggable parameters
func (r GigCreate) GetCompletion() string {
	return r.Completion
}

// Fill processes request and fills internal variables
func (r *GigCreate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["worker"]; ok && len(val) > 0 {
				r.Worker, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["preprocessors[]"]; ok {
				r.Preprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["preprocessors"]; ok {
				r.Preprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["postprocessors[]"]; ok {
				r.Postprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["postprocessors"]; ok {
				r.Postprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["completion"]; ok && len(val) > 0 {
				r.Completion, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["worker"]; ok && len(val) > 0 {
			r.Worker, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["preprocessors[]"]; ok {
			r.Preprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["preprocessors"]; ok {
			r.Preprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["postprocessors[]"]; ok {
			r.Postprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["postprocessors"]; ok {
			r.Postprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["completion"]; ok && len(val) > 0 {
			r.Completion, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewGigGo request
func NewGigGo() *GigGo {
	return &GigGo{}
}

// Auditable returns all auditable/loggable parameters
func (r GigGo) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"worker":         r.Worker,
		"preprocessors":  r.Preprocessors,
		"postprocessors": r.Postprocessors,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigGo) GetWorker() string {
	return r.Worker
}

// Auditable returns all auditable/loggable parameters
func (r GigGo) GetPreprocessors() conv.ParamWrapSet {
	return r.Preprocessors
}

// Auditable returns all auditable/loggable parameters
func (r GigGo) GetPostprocessors() conv.ParamWrapSet {
	return r.Postprocessors
}

// Fill processes request and fills internal variables
func (r *GigGo) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["worker"]; ok && len(val) > 0 {
				r.Worker, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["preprocessors[]"]; ok {
				r.Preprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["preprocessors"]; ok {
				r.Preprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["postprocessors[]"]; ok {
				r.Postprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["postprocessors"]; ok {
				r.Postprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["worker"]; ok && len(val) > 0 {
			r.Worker, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["preprocessors[]"]; ok {
			r.Preprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["preprocessors"]; ok {
			r.Preprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["postprocessors[]"]; ok {
			r.Postprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["postprocessors"]; ok {
			r.Postprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewGigRead request
func NewGigRead() *GigRead {
	return &GigRead{}
}

// Auditable returns all auditable/loggable parameters
func (r GigRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigRead) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigUpdate request
func NewGigUpdate() *GigUpdate {
	return &GigUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID":          r.GigID,
		"decoders":       r.Decoders,
		"preprocessors":  r.Preprocessors,
		"postprocessors": r.Postprocessors,
		"completion":     r.Completion,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetGigID() uint64 {
	return r.GigID
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetDecoders() conv.ParamWrapSet {
	return r.Decoders
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetPreprocessors() conv.ParamWrapSet {
	return r.Preprocessors
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetPostprocessors() conv.ParamWrapSet {
	return r.Postprocessors
}

// Auditable returns all auditable/loggable parameters
func (r GigUpdate) GetCompletion() string {
	return r.Completion
}

// Fill processes request and fills internal variables
func (r *GigUpdate) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["decoders[]"]; ok {
				r.Decoders, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["decoders"]; ok {
				r.Decoders, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["preprocessors[]"]; ok {
				r.Preprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["preprocessors"]; ok {
				r.Preprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["postprocessors[]"]; ok {
				r.Postprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["postprocessors"]; ok {
				r.Postprocessors, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["completion"]; ok && len(val) > 0 {
				r.Completion, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["decoders[]"]; ok {
			r.Decoders, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["decoders"]; ok {
			r.Decoders, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["preprocessors[]"]; ok {
			r.Preprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["preprocessors"]; ok {
			r.Preprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["postprocessors[]"]; ok {
			r.Postprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["postprocessors"]; ok {
			r.Postprocessors, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["completion"]; ok && len(val) > 0 {
			r.Completion, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigDelete request
func NewGigDelete() *GigDelete {
	return &GigDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r GigDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigDelete) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigUndelete request
func NewGigUndelete() *GigUndelete {
	return &GigUndelete{}
}

// Auditable returns all auditable/loggable parameters
func (r GigUndelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigUndelete) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigUndelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigAddSource request
func NewGigAddSource() *GigAddSource {
	return &GigAddSource{}
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID":    r.GigID,
		"upload":   r.Upload,
		"uri":      r.Uri,
		"decoders": r.Decoders,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetGigID() uint64 {
	return r.GigID
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetUpload() *multipart.FileHeader {
	return r.Upload
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetUri() string {
	return r.Uri
}

// Auditable returns all auditable/loggable parameters
func (r GigAddSource) GetDecoders() conv.ParamWrapSet {
	return r.Decoders
}

// Fill processes request and fills internal variables
func (r *GigAddSource) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			// Ignoring upload as its handled in the POST params section

			if val, ok := req.MultipartForm.Value["uri"]; ok && len(val) > 0 {
				r.Uri, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["decoders[]"]; ok {
				r.Decoders, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			} else if val, ok := req.MultipartForm.Value["decoders"]; ok {
				r.Decoders, err = conv.ParseParamWrap(val)
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if _, r.Upload, err = req.FormFile("upload"); err != nil {
			return fmt.Errorf("error processing uploaded file: %w", err)
		}

		if val, ok := req.Form["uri"]; ok && len(val) > 0 {
			r.Uri, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["decoders[]"]; ok {
			r.Decoders, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["decoders"]; ok {
			r.Decoders, err = conv.ParseParamWrap(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigRemoveSource request
func NewGigRemoveSource() *GigRemoveSource {
	return &GigRemoveSource{}
}

// Auditable returns all auditable/loggable parameters
func (r GigRemoveSource) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID":    r.GigID,
		"sourceID": r.SourceID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigRemoveSource) GetGigID() uint64 {
	return r.GigID
}

// Auditable returns all auditable/loggable parameters
func (r GigRemoveSource) GetSourceID() uint64 {
	return r.SourceID
}

// Fill processes request and fills internal variables
func (r *GigRemoveSource) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "sourceID")
		r.SourceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigPrepare request
func NewGigPrepare() *GigPrepare {
	return &GigPrepare{}
}

// Auditable returns all auditable/loggable parameters
func (r GigPrepare) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigPrepare) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigPrepare) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigExec request
func NewGigExec() *GigExec {
	return &GigExec{}
}

// Auditable returns all auditable/loggable parameters
func (r GigExec) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigExec) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigExec) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigOutput request
func NewGigOutput() *GigOutput {
	return &GigOutput{}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutput) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutput) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigOutput) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigOutputAll request
func NewGigOutputAll() *GigOutputAll {
	return &GigOutputAll{}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutputAll) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutputAll) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigOutputAll) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigOutputSpecific request
func NewGigOutputSpecific() *GigOutputSpecific {
	return &GigOutputSpecific{}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutputSpecific) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID":    r.GigID,
		"sourceID": r.SourceID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigOutputSpecific) GetGigID() uint64 {
	return r.GigID
}

// Auditable returns all auditable/loggable parameters
func (r GigOutputSpecific) GetSourceID() uint64 {
	return r.SourceID
}

// Fill processes request and fills internal variables
func (r *GigOutputSpecific) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

		val = chi.URLParam(req, "sourceID")
		r.SourceID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigState request
func NewGigState() *GigState {
	return &GigState{}
}

// Auditable returns all auditable/loggable parameters
func (r GigState) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigState) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigState) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigStatus request
func NewGigStatus() *GigStatus {
	return &GigStatus{}
}

// Auditable returns all auditable/loggable parameters
func (r GigStatus) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigStatus) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigStatus) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigComplete request
func NewGigComplete() *GigComplete {
	return &GigComplete{}
}

// Auditable returns all auditable/loggable parameters
func (r GigComplete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"gigID": r.GigID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r GigComplete) GetGigID() uint64 {
	return r.GigID
}

// Fill processes request and fills internal variables
func (r *GigComplete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "gigID")
		r.GigID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewGigWorkers request
func NewGigWorkers() *GigWorkers {
	return &GigWorkers{}
}

// Auditable returns all auditable/loggable parameters
func (r GigWorkers) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *GigWorkers) Fill(req *http.Request) (err error) {

	return err
}

// NewGigTasks request
func NewGigTasks() *GigTasks {
	return &GigTasks{}
}

// Auditable returns all auditable/loggable parameters
func (r GigTasks) Auditable() map[string]interface{} {
	return map[string]interface{}{}
}

// Fill processes request and fills internal variables
func (r *GigTasks) Fill(req *http.Request) (err error) {

	return err
}
