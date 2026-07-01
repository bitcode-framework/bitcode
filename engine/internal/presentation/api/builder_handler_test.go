package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/runtime/builder"
	"github.com/gofiber/fiber/v2"
)

type testModelResolver struct {
	models map[string]*parser.ModelDefinition
}

func (s *testModelResolver) Get(name string) (*parser.ModelDefinition, error) {
	if m, ok := s.models[name]; ok {
		return m, nil
	}
	return nil, &notFoundErr{name: name}
}

type notFoundErr struct{ name string }

func (e *notFoundErr) Error() string { return "model not found: " + e.name }

type testProcessResolver struct {
	processes map[string]*parser.ProcessDefinition
}

func (s *testProcessResolver) LoadProcess(name string) (*parser.ProcessDefinition, error) {
	if p, ok := s.processes[name]; ok {
		return p, nil
	}
	return nil, &notFoundErr{name: name}
}

func setupBuilderTestApp(t *testing.T) *fiber.App {
	t.Helper()
	registry := &testModelResolver{
		models: map[string]*parser.ModelDefinition{
			"contact": {
				Name: "contact",
				Fields: map[string]parser.FieldDefinition{
					"name":  {Type: "string", Label: "Name"},
					"email": {Type: "string", Label: "Email"},
				},
			},
		},
	}
	procRegistry := &testProcessResolver{
		processes: map[string]*parser.ProcessDefinition{
			"send_email": {Name: "send_email"},
		},
	}
	viewValidator := builder.NewViewValidator(registry)
	processValidator := builder.NewProcessValidator(procRegistry)
	handler := NewBuilderHandler(BuilderHandlerConfig{
		ViewValidator:    viewValidator,
		ProcessValidator: processValidator,
	})

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "admin-user")
		return c.Next()
	})
	group := app.Group("/api/v1/admin")
	group.Post("/views/validate", handler.ValidateView)
	group.Post("/processes/validate", handler.ValidateProcess)
	return app
}

func TestBuilderHandler_ValidateView_Valid(t *testing.T) {
	app := setupBuilderTestApp(t)

	body := bytes.NewReader([]byte(`{
		"name": "contact_form",
		"type": "form",
		"model": "contact",
		"title": "Contact Form",
		"fields": ["name", "email"],
		"layout": [
			{
				"row": [
					{"field": "name", "width": 6},
					{"field": "email", "width": 6}
				]
			}
		]
	}`))

	req := httptest.NewRequest("POST", "/api/v1/admin/views/validate", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result builder.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !result.Valid {
		t.Fatalf("expected valid=true, got false; errors: %+v", result.Errors)
	}
}

func TestBuilderHandler_ValidateView_InvalidFieldRef(t *testing.T) {
	app := setupBuilderTestApp(t)

	body := bytes.NewReader([]byte(`{
		"name": "contact_form",
		"type": "form",
		"model": "contact",
		"fields": ["name", "nonexistent_field"]
	}`))

	req := httptest.NewRequest("POST", "/api/v1/admin/views/validate", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result builder.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Valid {
		t.Fatal("expected valid=false for invalid field ref")
	}
	foundFieldError := false
	for _, e := range result.Errors {
		if e.Field == `fields["nonexistent_field"]` {
			foundFieldError = true
		}
	}
	if !foundFieldError {
		t.Fatalf("expected error about nonexistent_field, got: %+v", result.Errors)
	}
}

func TestBuilderHandler_ValidateView_MissingRequired(t *testing.T) {
	app := setupBuilderTestApp(t)

	body := bytes.NewReader([]byte(`{
		"type": "form",
		"model": "contact"
	}`))

	req := httptest.NewRequest("POST", "/api/v1/admin/views/validate", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result builder.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Valid {
		t.Fatal("expected valid=false for missing name")
	}
}

func TestBuilderHandler_ValidateView_EmptyBody(t *testing.T) {
	app := setupBuilderTestApp(t)

	req := httptest.NewRequest("POST", "/api/v1/admin/views/validate", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400 for empty body, got %d", resp.StatusCode)
	}
}

func TestBuilderHandler_ValidateView_InvalidJSON(t *testing.T) {
	app := setupBuilderTestApp(t)

	req := httptest.NewRequest("POST", "/api/v1/admin/views/validate",
		bytes.NewReader([]byte(`{bad json`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result builder.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Valid {
		t.Fatal("expected valid=false for invalid JSON")
	}
}

func TestBuilderHandler_ValidateProcess_Valid(t *testing.T) {
	app := setupBuilderTestApp(t)

	body := `{
		"name": "create_lead",
		"steps": [
			{"name": "validate_input", "type": "validate", "rules": {}},
			{"name": "create_record", "type": "create", "model": "contact"},
			{"name": "emit_created", "type": "emit", "event": "lead_created"}
		]
	}`

	req := httptest.NewRequest("POST", "/api/v1/admin/processes/validate",
		bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result builder.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !result.Valid {
		t.Fatalf("expected valid=true for valid process, got false; errors: %+v", result.Errors)
	}
}

func TestBuilderHandler_ValidateProcess_InvalidStepRef(t *testing.T) {
	app := setupBuilderTestApp(t)

	body := `{
		"name": "bad_step_ref",
		"steps": [
			{"name": "check", "type": "if", "condition": "x > 0", "then": "ghost_step", "else": "done"},
			{"name": "done", "type": "log", "message": "finished"}
		]
	}`

	req := httptest.NewRequest("POST", "/api/v1/admin/processes/validate",
		bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result builder.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Valid {
		t.Fatal("expected valid=false for invalid step ref")
	}
	if len(result.Errors) == 0 {
		t.Fatal("expected at least one error")
	}
}

func TestBuilderHandler_ValidateProcess_OrphanTransition(t *testing.T) {
	app := setupBuilderTestApp(t)

	body := `{
		"name": "orphan_edge",
		"nodes": {
			"start": {"type": "query", "model": "contact"},
			"save":  {"type": "create", "model": "contact"}
		},
		"edges": [
			{"from": "start", "to": "ghost_node"}
		]
	}`

	req := httptest.NewRequest("POST", "/api/v1/admin/processes/validate",
		bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result builder.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Valid {
		t.Fatal("expected valid=false for orphan transition")
	}
	foundEdgeError := false
	for _, e := range result.Errors {
		if e.Field == "edges[0].to" {
			foundEdgeError = true
		}
	}
	if !foundEdgeError {
		t.Fatalf("expected edge ref error for ghost_node, got: %+v", result.Errors)
	}
}

func TestBuilderHandler_ValidateProcess_EmptyBody(t *testing.T) {
	app := setupBuilderTestApp(t)

	req := httptest.NewRequest("POST", "/api/v1/admin/processes/validate", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400 for empty body, got %d", resp.StatusCode)
	}
}
