package builder

import (
	"encoding/json"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
)

type stubModelResolver struct {
	models map[string]*parser.ModelDefinition
}

func (s *stubModelResolver) Get(name string) (*parser.ModelDefinition, error) {
	if m, ok := s.models[name]; ok {
		return m, nil
	}
	return nil, &modelNotFoundError{name: name}
}

type modelNotFoundError struct{ name string }

func (e *modelNotFoundError) Error() string { return "model not found: " + e.name }

func testModelRegistry() *stubModelResolver {
	return &stubModelResolver{
		models: map[string]*parser.ModelDefinition{
			"contact": {
				Name: "contact",
				Fields: map[string]parser.FieldDefinition{
					"name":  {Type: "string", Label: "Name"},
					"email": {Type: "string", Label: "Email"},
					"phone": {Type: "string", Label: "Phone"},
				},
			},
		},
	}
}

func validViewJSON() []byte {
	return []byte(`{
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
	}`)
}

func TestViewValidator_ValidView(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	result := v.Validate(validViewJSON())

	if !result.Valid {
		t.Fatalf("expected valid=true, got false; errors: %+v", result.Errors)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected 0 errors, got %d: %+v", len(result.Errors), result.Errors)
	}
}

func TestViewValidator_InvalidFieldRef(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "contact_form",
		"type": "form",
		"model": "contact",
		"fields": ["name", "nonexistent_field"]
	}`)

	result := v.Validate(data)

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

func TestViewValidator_MissingRequired(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"type": "form",
		"model": "contact"
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for missing name")
	}
	foundNameError := false
	for _, e := range result.Errors {
		if e.Field == "name" {
			foundNameError = true
		}
	}
	if !foundNameError {
		t.Fatalf("expected error about missing name, got: %+v", result.Errors)
	}
}

func TestViewValidator_MissingType(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "some_view",
		"model": "contact"
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for missing type")
	}
	foundTypeError := false
	for _, e := range result.Errors {
		if e.Field == "type" {
			foundTypeError = true
		}
	}
	if !foundTypeError {
		t.Fatalf("expected error about missing type, got: %+v", result.Errors)
	}
}

func TestViewValidator_InvalidViewType(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "test_view",
		"type": "nonexistent_type",
		"model": "contact"
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid view type")
	}
}

func TestViewValidator_ModelNotFound(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "test_view",
		"type": "form",
		"model": "nonexistent_model"
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for model not found")
	}
	foundModelError := false
	for _, e := range result.Errors {
		if e.Field == "model" {
			foundModelError = true
		}
	}
	if !foundModelError {
		t.Fatalf("expected error about model not found, got: %+v", result.Errors)
	}
}

func TestViewValidator_InvalidJSON(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	result := v.Validate([]byte(`{invalid json`))

	if result.Valid {
		t.Fatal("expected valid=false for invalid JSON")
	}
	if len(result.Errors) == 0 {
		t.Fatal("expected at least one error")
	}
}

func TestViewValidator_NilModelRegistry(t *testing.T) {
	v := NewViewValidator(nil)
	result := v.Validate(validViewJSON())

	if !result.Valid {
		t.Fatalf("expected valid=true with nil registry (skips field ref checks), got false; errors: %+v", result.Errors)
	}
}

func TestViewValidator_FormNoLayoutWarning(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "contact_form",
		"type": "form",
		"model": "contact",
		"fields": ["name"]
	}`)

	result := v.Validate(data)

	if !result.Valid {
		t.Fatalf("expected valid=true (warning only), got false; errors: %+v", result.Errors)
	}
	if len(result.Warnings) == 0 {
		t.Fatal("expected at least one warning for form without layout")
	}
}

func TestViewValidator_LayoutFieldRef(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "contact_form",
		"type": "form",
		"model": "contact",
		"layout": [
			{
				"row": [
					{"field": "bad_field"}
				]
			}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for layout field ref error")
	}
	foundLayoutError := false
	for _, e := range result.Errors {
		if e.Field == `layout.row["bad_field"]` {
			foundLayoutError = true
		}
	}
	if !foundLayoutError {
		t.Fatalf("expected layout field ref error, got: %+v", result.Errors)
	}
}

func TestViewValidator_SortFieldRef(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "contact_list",
		"type": "list",
		"model": "contact",
		"sort": {"field": "bad_sort_field", "direction": "asc"}
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid sort field")
	}
}

func TestViewValidator_FilterFieldRef(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "contact_list",
		"type": "list",
		"model": "contact",
		"filters": ["bad_filter_field"]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid filter field")
	}
}

func TestViewValidator_DataSourceMutualExclusion(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "test_view",
		"type": "list",
		"model": "contact",
		"data_sources": {
			"ds1": {
				"model": "contact",
				"process": "some_process"
			}
		}
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for data source with both model and process")
	}
}

func TestViewValidator_RenderModeValidation(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "test_view",
		"type": "list",
		"model": "contact",
		"render": "invalid_mode"
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid render mode")
	}
}

func TestViewValidator_CustomRequiresTemplate(t *testing.T) {
	v := NewViewValidator(testModelRegistry())
	data := []byte(`{
		"name": "custom_view",
		"type": "custom"
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for custom view without template")
	}
}

func TestValidationResult_JSON(t *testing.T) {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationIssue{},
		Warnings: []ValidationIssue{
			{Field: "layout", Message: "no layout"},
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed["valid"] != true {
		t.Fatal("expected valid=true in JSON")
	}
	if _, ok := parsed["errors"]; !ok {
		t.Fatal("expected errors key in JSON")
	}
	if _, ok := parsed["warnings"]; !ok {
		t.Fatal("expected warnings key in JSON")
	}
}
