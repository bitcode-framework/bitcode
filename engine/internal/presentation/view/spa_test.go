package view

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	tmplEngine "github.com/bitcode-framework/bitcode/internal/presentation/template"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupSPATestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	sqlDB, _ := db.DB()
	if _, err := sqlDB.Exec(`CREATE TABLE contacts (
		id TEXT PRIMARY KEY,
		name TEXT,
		status TEXT,
		created_at DATETIME,
		updated_at DATETIME
	)`); err != nil {
		t.Fatalf("create contacts: %v", err)
	}
	if _, err := sqlDB.Exec(`CREATE TABLE leads (
		id TEXT PRIMARY KEY,
		name TEXT,
		status TEXT,
		company TEXT,
		created_at DATETIME,
		updated_at DATETIME
	)`); err != nil {
		t.Fatalf("create leads: %v", err)
	}
	if _, err := sqlDB.Exec(`INSERT INTO contacts (id, name, status) VALUES ('c1', 'Alice', 'active'), ('c2', 'Bob', 'inactive')`); err != nil {
		t.Fatalf("seed contacts: %v", err)
	}
	if _, err := sqlDB.Exec(`INSERT INTO leads (id, name, status, company) VALUES ('l1', 'Deal A', 'new', 'Acme'), ('l2', 'Deal B', 'won', 'Globex')`); err != nil {
		t.Fatalf("seed leads: %v", err)
	}
	return db
}

func newSPARenderer(db *gorm.DB) *Renderer {
	r := NewRenderer(db, tmplEngine.NewEngine())
	r.SetTableNameResolver(spaTestTableResolver{})
	return r
}

// spaTestTableResolver maps model names to the test table names (pluralized).
type spaTestTableResolver struct{}

func (spaTestTableResolver) TableName(model string) string {
	switch model {
	case "contact":
		return "contacts"
	case "lead":
		return "leads"
	}
	return model
}

// ---- Parser: render flag validation ----

func TestParseView_RenderModeValidation(t *testing.T) {
	t.Run("spa mode accepted", func(t *testing.T) {
		data := []byte(`{"name":"v","type":"list","model":"contact","render":"spa"}`)
		v, err := parser.ParseView(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Render != parser.RenderSPA {
			t.Errorf("render = %q, want %q", v.Render, parser.RenderSPA)
		}
	})

	t.Run("ssr mode accepted", func(t *testing.T) {
		data := []byte(`{"name":"v","type":"list","model":"contact","render":"ssr"}`)
		v, err := parser.ParseView(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Render != parser.RenderSSR {
			t.Errorf("render = %q, want %q", v.Render, parser.RenderSSR)
		}
	})

	t.Run("empty render defaults to ssr (backwards compat)", func(t *testing.T) {
		data := []byte(`{"name":"v","type":"list","model":"contact"}`)
		v, err := parser.ParseView(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Render != "" {
			t.Errorf("render = %q, want empty", v.Render)
		}
		if IsSPAMode(v) {
			t.Error("empty render must not be SPA")
		}
	})

	t.Run("invalid render mode rejected", func(t *testing.T) {
		data := []byte(`{"name":"v","type":"list","model":"contact","render":"weird"}`)
		_, err := parser.ParseView(data)
		if err == nil {
			t.Fatal("expected error for invalid render mode, got nil")
		}
		if !strings.Contains(err.Error(), "invalid render mode") {
			t.Errorf("error = %q, want it to contain 'invalid render mode'", err.Error())
		}
	})
}

// ---- IsSPAMode ----

func TestIsSPAMode(t *testing.T) {
	if IsSPAMode(nil) {
		t.Error("nil view must not be SPA")
	}
	if IsSPAMode(&parser.ViewDefinition{Render: ""}) {
		t.Error("empty render must not be SPA")
	}
	if IsSPAMode(&parser.ViewDefinition{Render: parser.RenderSSR}) {
		t.Error("ssr render must not be SPA")
	}
	if !IsSPAMode(&parser.ViewDefinition{Render: parser.RenderSPA}) {
		t.Error("spa render must be SPA")
	}
}

// ---- SPA envelope: list view ----

func TestRenderSPAEnvelope_ListView(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	viewDef := &parser.ViewDefinition{
		Name:   "contact_list",
		Type:   parser.ViewList,
		Model:  "contact",
		Title:  "Contacts",
		Fields: []string{"name", "status"},
		Render: parser.RenderSPA,
	}

	env, err := r.RenderSPAEnvelope(context.Background(), viewDef, RenderOptions{Module: "base", ViewPath: "base/contact_list"})
	if err != nil {
		t.Fatalf("RenderSPAEnvelope error: %v", err)
	}

	// Envelope structure
	if env.View.Name != "contact_list" {
		t.Errorf("view.name = %q, want contact_list", env.View.Name)
	}
	if env.View.Type != parser.ViewList {
		t.Errorf("view.type = %q, want list", env.View.Type)
	}
	if env.View.Render != parser.RenderSPA {
		t.Errorf("view.render = %q, want spa", env.View.Render)
	}
	if env.Meta.Module != "base" {
		t.Errorf("meta.module = %q, want base", env.Meta.Module)
	}
	if env.Meta.ViewType != "list" {
		t.Errorf("meta.viewType = %q, want list", env.Meta.ViewType)
	}
	if env.Meta.APISource != "/api/v1/base/contact" {
		t.Errorf("meta.apiSource = %q, want /api/v1/base/contact", env.Meta.APISource)
	}

	// List data exposes API source + fields (datatable fetches client-side,
	// mirroring SSR which also passes no records to the datatable).
	if env.Data["apiSource"] != "/api/v1/base/contact" {
		t.Errorf("data.apiSource = %v, want /api/v1/base/contact", env.Data["apiSource"])
	}

	// JSON-serializable (the envelope must round-trip to the client).
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	if !strings.Contains(string(raw), `"view"`) {
		t.Error("envelope JSON must contain 'view' key")
	}
	if !strings.Contains(string(raw), `"data"`) {
		t.Error("envelope JSON must contain 'data' key")
	}
	if !strings.Contains(string(raw), `"meta"`) {
		t.Error("envelope JSON must contain 'meta' key")
	}
}

// ---- SPA envelope: form view with record (data parity) ----

func TestRenderSPAEnvelope_FormView_DataParity(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	viewDef := &parser.ViewDefinition{
		Name:   "contact_form",
		Type:   parser.ViewForm,
		Model:  "contact",
		Title:  "Contact",
		Render: parser.RenderSPA,
	}

	opt := RenderOptions{Module: "base", RecordID: "c1", ViewPath: "base/contact_form"}

	env, err := r.RenderSPAEnvelope(context.Background(), viewDef, opt)
	if err != nil {
		t.Fatalf("RenderSPAEnvelope error: %v", err)
	}

	// The SPA envelope must contain the same record the SSR renderer loads.
	rec, ok := env.Data["record"].(map[string]any)
	if !ok {
		t.Fatalf("record must be a map, got %T", env.Data["record"])
	}
	if rec["id"] != "c1" {
		t.Errorf("record.id = %v, want c1", rec["id"])
	}
	if rec["name"] != "Alice" {
		t.Errorf("record.name = %v, want Alice", rec["name"])
	}
	if rec["status"] != "active" {
		t.Errorf("record.status = %v, want active", rec["status"])
	}

	// Meta carries record id + module.
	if env.Meta.RecordID != "c1" {
		t.Errorf("meta.recordId = %q, want c1", env.Meta.RecordID)
	}
	if env.Meta.Module != "base" {
		t.Errorf("meta.module = %q, want base", env.Meta.Module)
	}
	if env.Meta.ViewType != "form" {
		t.Errorf("meta.viewType = %q, want form", env.Meta.ViewType)
	}
}

// ---- SPA envelope: form view without record (new form) ----

func TestRenderSPAEnvelope_FormView_NewRecord(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	viewDef := &parser.ViewDefinition{
		Name:   "contact_form",
		Type:   parser.ViewForm,
		Model:  "contact",
		Render: parser.RenderSPA,
	}

	env, err := r.RenderSPAEnvelope(context.Background(), viewDef, RenderOptions{Module: "base"})
	if err != nil {
		t.Fatalf("RenderSPAEnvelope error: %v", err)
	}

	rec, ok := env.Data["record"].(map[string]any)
	if !ok {
		t.Fatalf("record must be a map, got %T", env.Data["record"])
	}
	if len(rec) != 0 {
		t.Errorf("record must be empty for new form, got %v", rec)
	}
}

// ---- SPA envelope: kanban view ----

func TestRenderSPAEnvelope_KanbanView(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	viewDef := &parser.ViewDefinition{
		Name:    "lead_kanban",
		Type:    parser.ViewKanban,
		Model:   "lead",
		Title:   "Pipeline",
		Fields:  []string{"name", "company"},
		GroupBy: "status",
		Render:  parser.RenderSPA,
	}

	env, err := r.RenderSPAEnvelope(context.Background(), viewDef, RenderOptions{Module: "base"})
	if err != nil {
		t.Fatalf("RenderSPAEnvelope error: %v", err)
	}

	records, ok := env.Data["records"].([]map[string]any)
	if !ok {
		t.Fatalf("kanban data must include records, got %T", env.Data["records"])
	}
	if len(records) != 2 {
		t.Errorf("records len = %d, want 2", len(records))
	}
	if env.Data["groupBy"] != "status" {
		t.Errorf("data.groupBy = %v, want status", env.Data["groupBy"])
	}
}

// ---- SPA envelope: nil view def errors ----

func TestRenderSPAEnvelope_NilView(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	_, err := r.RenderSPAEnvelope(context.Background(), nil, RenderOptions{})
	if err == nil {
		t.Fatal("expected error for nil view, got nil")
	}
}

// ---- SPA envelope: kanban missing group_by errors ----

func TestRenderSPAEnvelope_KanbanMissingGroupBy(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	viewDef := &parser.ViewDefinition{
		Name:   "lead_kanban",
		Type:   parser.ViewKanban,
		Model:  "lead",
		Render: parser.RenderSPA,
	}

	_, err := r.RenderSPAEnvelope(context.Background(), viewDef, RenderOptions{Module: "base"})
	if err == nil {
		t.Fatal("expected error for kanban missing group_by, got nil")
	}
	if !strings.Contains(err.Error(), "group_by") {
		t.Errorf("error = %q, want it to contain 'group_by'", err.Error())
	}
}

// ---- Data parity: SPA form record == SSR form record ----
// Guarantees ADR-014 D3: SPA data is identical to SSR data for the same view.
// Both SPA and SSR call loadFormData → FindByID, so the records match.

func TestDataParity_SPAEqualsSSR_FormView(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	viewDef := &parser.ViewDefinition{
		Name:   "contact_form",
		Type:   parser.ViewForm,
		Model:  "contact",
		Render: parser.RenderSPA,
	}
	opt := RenderOptions{Module: "base", RecordID: "c1"}

	// SPA envelope record
	env, err := r.RenderSPAEnvelope(context.Background(), viewDef, opt)
	if err != nil {
		t.Fatalf("RenderSPAEnvelope error: %v", err)
	}
	spaRecord := env.Data["record"].(map[string]any)

	// SSR form render loads the same record via the same repo path. We extract
	// the record using the same data-loading helper to prove parity.
	ssrData, err := r.loadFormData(context.Background(), viewDef, opt)
	if err != nil {
		t.Fatalf("loadFormData error: %v", err)
	}
	ssrRecord := ssrData["record"].(map[string]any)

	// The records must be equal (same id, name, status).
	if ssrRecord["id"] != spaRecord["id"] {
		t.Errorf("id: spa=%v ssr=%v", spaRecord["id"], ssrRecord["id"])
	}
	if ssrRecord["name"] != spaRecord["name"] {
		t.Errorf("name: spa=%v ssr=%v", spaRecord["name"], ssrRecord["name"])
	}
	if ssrRecord["status"] != spaRecord["status"] {
		t.Errorf("status: spa=%v ssr=%v", spaRecord["status"], ssrRecord["status"])
	}
}

// ---- Backwards compat: SSR render produces HTML, not JSON ----

func TestSSR_RenderProducesHTML(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	// View WITHOUT render flag → SSR (backwards compat). Renderer.RenderView
	// produces HTML strings, not JSON envelopes.
	viewDef := &parser.ViewDefinition{
		Name:   "contact_list",
		Type:   parser.ViewList,
		Model:  "contact",
		Title:  "Contacts",
		Fields: []string{"name", "status"},
	}

	html, err := r.RenderView(context.Background(), viewDef, RenderOptions{Module: "base", UseLayout: false})
	if err != nil {
		t.Fatalf("RenderView error: %v", err)
	}
	if html == "" {
		t.Fatal("expected non-empty HTML output")
	}
	// SSR output is HTML (contains tags), not a JSON envelope.
	if !strings.Contains(html, "<") {
		t.Error("SSR output must contain HTML tags")
	}
	// Must NOT be a JSON envelope (no {view:...} shape).
	if strings.Contains(html, `"view":`) {
		t.Error("SSR output must not be a JSON envelope")
	}
}

// ---- Backwards compat: render:"ssr" explicitly also produces HTML ----

func TestSSR_ExplicitRenderMode_ProducesHTML(t *testing.T) {
	db := setupSPATestDB(t)
	r := newSPARenderer(db)

	viewDef := &parser.ViewDefinition{
		Name:   "contact_list",
		Type:   parser.ViewList,
		Model:  "contact",
		Title:  "Contacts",
		Fields: []string{"name", "status"},
		Render: parser.RenderSSR,
	}

	if IsSPAMode(viewDef) {
		t.Error("render:ssr must not be SPA mode")
	}

	html, err := r.RenderView(context.Background(), viewDef, RenderOptions{Module: "base", UseLayout: false})
	if err != nil {
		t.Fatalf("RenderView error: %v", err)
	}
	if html == "" {
		t.Fatal("expected non-empty HTML output")
	}
	if !strings.Contains(html, "<") {
		t.Error("SSR output must contain HTML tags")
	}
}
