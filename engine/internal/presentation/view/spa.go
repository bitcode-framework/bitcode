package view

import (
	"context"
	"fmt"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
)

// SPAEnvelope is the JSON contract returned when a view has render:"spa".
// Instead of server-rendered HTML, the engine returns the view definition
// plus the data the client runtime (BcRenderer + ReactiveStore) needs to
// render the view on web (Mitosis) or Flutter.
//
//   - View:  the full ViewDefinition (serializable as JSON).
//   - Data:  view-type-specific payload (records for list/kanban/calendar,
//            single record for form, data-source results for chart/custom).
//   - Meta:  rendering metadata (module, view path, record id, page, locale,
//            API source URLs the client may need).
type SPAEnvelope struct {
	View parser.ViewDefinition `json:"view"`
	Data map[string]any        `json:"data"`
	Meta SPAMeta               `json:"meta"`
}

// SPAMeta carries rendering metadata for the client runtime.
type SPAMeta struct {
	Module    string `json:"module"`
	ViewPath  string `json:"viewPath"`
	ViewType  string `json:"viewType"`
	RecordID  string `json:"recordId,omitempty"`
	Page      int    `json:"page,omitempty"`
	Locale    string `json:"locale,omitempty"`
	Username  string `json:"username,omitempty"`
	ActiveView string `json:"activeView,omitempty"`
	// APISource is the convention-driven CRUD endpoint the client uses to
	// fetch list records (empty for non-list views).
	APISource string `json:"apiSource,omitempty"`
}

// RenderSPAEnvelope produces the {view, data, meta} envelope for a view with
// render:"spa". It reuses the SAME data-loading code paths as SSR rendering
// (renderForm, renderKanban, renderCalendar, renderChart, renderCustom) so
// that SPA data is identical to SSR data (data-parity guarantee, ADR-014 D3).
//
// For list views, SSR emits a datatable that fetches records client-side via
// the convention-driven CRUD API; the SPA envelope mirrors this by exposing
// the API source URL in meta.apiSource (no records in data, matching SSR
// which also passes no records to the datatable).
func (r *Renderer) RenderSPAEnvelope(ctx context.Context, viewDef *parser.ViewDefinition, opt RenderOptions) (*SPAEnvelope, error) {
	if viewDef == nil {
		return nil, fmt.Errorf("view definition is required for SPA envelope")
	}

	data, err := r.loadSPAData(ctx, viewDef, opt)
	if err != nil {
		return nil, err
	}

	meta := SPAMeta{
		Module:     opt.Module,
		ViewPath:   opt.ViewPath,
		ViewType:   string(viewDef.Type),
		RecordID:   opt.RecordID,
		Page:       opt.Page,
		Locale:     opt.Locale,
		Username:   opt.Username,
		ActiveView: opt.ActiveView,
	}
	if viewDef.Model != "" && opt.Module != "" {
		meta.APISource = fmt.Sprintf("/api/v1/%s/%s", opt.Module, viewDef.Model)
	}

	return &SPAEnvelope{
		View: *viewDef,
		Data: data,
		Meta: meta,
	}, nil
}

// loadSPAData loads the view-type-specific data payload, reusing the same
// data-loading logic as the SSR renderers to guarantee data parity.
func (r *Renderer) loadSPAData(ctx context.Context, viewDef *parser.ViewDefinition, opt RenderOptions) (map[string]any, error) {
	switch viewDef.Type {
	case parser.ViewForm:
		return r.loadFormData(ctx, viewDef, opt)
	case parser.ViewList:
		return r.loadListData(ctx, viewDef, opt)
	case parser.ViewKanban:
		return r.loadKanbanData(ctx, viewDef)
	case parser.ViewCalendar:
		return r.loadCalendarData(ctx, viewDef)
	case parser.ViewChart:
		return r.loadChartData(ctx, viewDef)
	case parser.ViewCustom:
		return r.loadCustomData(ctx, viewDef, opt.Module)
	case parser.ViewActivity:
		return map[string]any{
			"recordId": opt.RecordID,
			"model":    viewDef.Model,
		}, nil
	default:
		return map[string]any{}, nil
	}
}

// loadFormData mirrors renderForm's record loading (data parity).
func (r *Renderer) loadFormData(ctx context.Context, viewDef *parser.ViewDefinition, opt RenderOptions) (map[string]any, error) {
	var record map[string]any
	recordID := opt.RecordID

	if recordID != "" && viewDef.Model != "" {
		repo := persistence.NewGenericRepository(r.db, r.resolveTable(viewDef.Model))
		rec, err := repo.FindByID(ctx, recordID)
		if err == nil && rec != nil {
			record = rec
			if r.hydrator != nil && r.modelRegistry != nil {
				if modelDef, mErr := r.modelRegistry.Get(viewDef.Model); mErr == nil {
					r.hydrator.HydrateRecord(ctx, modelDef, record)
				}
			}
		}
	}
	if record == nil {
		record = make(map[string]any)
	}

	listURL := ""
	formAction := ""
	if opt.Module != "" {
		listURL = fmt.Sprintf("/app/%s/%s/list", opt.Module, viewDef.Model)
		formAction = fmt.Sprintf("/app/%s/%s/form", opt.Module, viewDef.Model)
		if recordID != "" {
			formAction += "?id=" + recordID
		}
	}

	return map[string]any{
		"record":     record,
		"recordId":   recordID,
		"listUrl":    listURL,
		"formAction": formAction,
	}, nil
}

// loadListData mirrors renderList (datatable fetches client-side; SSR passes
// no records). The SPA envelope exposes the API source URL in meta so the
// client runtime can fetch records the same way the SSR datatable does.
func (r *Renderer) loadListData(ctx context.Context, viewDef *parser.ViewDefinition, opt RenderOptions) (map[string]any, error) {
	data := map[string]any{
		"fields":   viewDef.Fields,
		"exportable": viewDef.Exportable,
	}
	if opt.Module != "" && viewDef.Model != "" {
		data["apiSource"] = fmt.Sprintf("/api/v1/%s/%s", opt.Module, viewDef.Model)
		data["createUrl"] = fmt.Sprintf("/app/%s/%s/new", opt.Module, viewDef.Model)
		data["detailUrl"] = fmt.Sprintf("/app/%s/%s/:id", opt.Module, viewDef.Model)
	}
	return data, nil
}

// loadKanbanData mirrors renderKanban's record loading (data parity).
func (r *Renderer) loadKanbanData(ctx context.Context, viewDef *parser.ViewDefinition) (map[string]any, error) {
	if viewDef.Model == "" || viewDef.GroupBy == "" {
		return nil, fmt.Errorf("kanban view requires model and group_by")
	}
	repo := persistence.NewGenericRepository(r.db, r.resolveTable(viewDef.Model))
	records, _, err := repo.FindAll(ctx, nil, 1, 200)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"records": records,
		"fields":  viewDef.Fields,
		"groupBy": viewDef.GroupBy,
	}, nil
}

// loadCalendarData mirrors renderCalendar's record loading (data parity).
func (r *Renderer) loadCalendarData(ctx context.Context, viewDef *parser.ViewDefinition) (map[string]any, error) {
	if viewDef.Model == "" {
		return nil, fmt.Errorf("calendar view requires a model")
	}
	repo := persistence.NewGenericRepository(r.db, r.resolveTable(viewDef.Model))
	records, _, err := repo.FindAll(ctx, nil, 1, 100)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"records":  records,
		"fields":   viewDef.Fields,
		"dateField": viewDef.DateField,
	}, nil
}

// loadChartData mirrors renderChart's data-source loading (data parity).
func (r *Renderer) loadChartData(ctx context.Context, viewDef *parser.ViewDefinition) (map[string]any, error) {
	data := make(map[string]any)
	for name, ds := range viewDef.DataSources {
		if ds.Model != "" {
			repo := persistence.NewGenericRepository(r.db, r.resolveTable(ds.Model))
			records, _, err := repo.FindAll(ctx, persistence.QueryFromDomain(ds.Domain), 1, 1000)
			if err != nil || records == nil {
				data[name] = []map[string]any{}
				continue
			}
			data[name] = records
		} else if ds.Process != "" && r.processExecutor != nil {
			result, err := r.processExecutor.ExecuteByName(ctx, ds.Process, nil)
			if err != nil {
				data[name] = []map[string]any{}
				continue
			}
			data[name] = result
		}
	}
	return data, nil
}

// loadCustomData mirrors renderCustom's data-source loading (data parity).
func (r *Renderer) loadCustomData(ctx context.Context, viewDef *parser.ViewDefinition, moduleName string) (map[string]any, error) {
	data := make(map[string]any)
	for name, ds := range viewDef.DataSources {
		if ds.Model != "" {
			repo := persistence.NewGenericRepository(r.db, r.resolveTable(ds.Model))
			records, _, err := repo.FindAll(ctx, persistence.QueryFromDomain(ds.Domain), 1, 100)
			if err != nil {
				data[name] = []map[string]any{}
				continue
			}
			if records == nil {
				records = []map[string]any{}
			}
			data[name] = records
		} else if ds.Process != "" && r.processExecutor != nil {
			result, err := r.processExecutor.ExecuteByName(ctx, ds.Process, nil)
			if err != nil {
				data[name] = []map[string]any{}
				continue
			}
			data[name] = result
		}
	}
	return data, nil
}

// IsSPAMode returns true when the view definition opts into SPA rendering.
// A view is SPA when render == "spa". Empty or "ssr" means SSR (default,
// backwards compatible per ADR-014 D3).
func IsSPAMode(viewDef *parser.ViewDefinition) bool {
	if viewDef == nil {
		return false
	}
	return viewDef.Render == parser.RenderSPA
}
