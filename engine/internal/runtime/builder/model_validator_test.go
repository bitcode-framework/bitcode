package builder

import (
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
)

func validModelJSON() []byte {
	return []byte(`{
		"name": "product",
		"fields": {
			"name": {"type": "string", "label": "Name"},
			"price": {"type": "float", "label": "Price"},
			"active": {"type": "boolean", "label": "Active"}
		}
	}`)
}

func TestModelValidator_ValidModel(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate(validModelJSON())
	if !result.Valid {
		t.Fatalf("expected valid=true for valid model, got false; errors: %+v", result.Errors)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected 0 errors, got %d: %+v", len(result.Errors), result.Errors)
	}
}

func TestModelValidator_InvalidJSON(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{bad json`))
	if result.Valid {
		t.Fatal("expected valid=false for invalid JSON")
	}
	if len(result.Errors) == 0 {
		t.Fatal("expected at least one error for invalid JSON")
	}
}

func TestModelValidator_MissingName(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"fields": {
			"name": {"type": "string"}
		}
	}`))
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

func TestModelValidator_NoFields(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "empty_model",
		"fields": {}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for no fields")
	}
	foundFieldsError := false
	for _, e := range result.Errors {
		if e.Field == "fields" {
			foundFieldsError = true
		}
	}
	if !foundFieldsError {
		t.Fatalf("expected error about missing fields, got: %+v", result.Errors)
	}
}

func TestModelValidator_InvalidFieldType(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "bad_type_model",
		"fields": {
			"quantity": {"type": "nonexistent_type", "label": "Quantity"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for invalid field type")
	}
	foundTypeError := false
	for _, e := range result.Errors {
		if e.Field == `fields["quantity"].type` {
			foundTypeError = true
		}
	}
	if !foundTypeError {
		t.Fatalf("expected error about invalid field type, got: %+v", result.Errors)
	}
}

func TestModelValidator_SelectionWithoutOptions(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "selection_model",
		"fields": {
			"status": {"type": "selection", "label": "Status"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for selection without options")
	}
	foundOptionsError := false
	for _, e := range result.Errors {
		if e.Field == `fields["status"].options` {
			foundOptionsError = true
		}
	}
	if !foundOptionsError {
		t.Fatalf("expected error about missing options, got: %+v", result.Errors)
	}
}

func TestModelValidator_SelectionWithOptions(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "selection_model",
		"fields": {
			"status": {"type": "selection", "label": "Status", "options": ["draft", "published"]}
		}
	}`))
	if !result.Valid {
		t.Fatalf("expected valid=true for selection with options, got false; errors: %+v", result.Errors)
	}
}

func TestModelValidator_ReservedNamespace(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "reserved_model",
		"fields": {
			"ctx": {"type": "string", "label": "Context"},
			"input": {"type": "string", "label": "Input"},
			"valid_field": {"type": "string", "label": "Valid"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for reserved namespace fields")
	}
	errorFields := make(map[string]bool)
	for _, e := range result.Errors {
		errorFields[e.Field] = true
	}
	if !errorFields[`fields["ctx"]`] {
		t.Errorf("expected error for 'ctx' field, got: %+v", result.Errors)
	}
	if !errorFields[`fields["input"]`] {
		t.Errorf("expected error for 'input' field, got: %+v", result.Errors)
	}
}

func TestModelValidator_AllReservedNamespaces(t *testing.T) {
	reservedNames := []string{"ctx", "input", "old", "session", "self"}
	for _, name := range reservedNames {
		v := NewModelValidator(nil)
		result := v.Validate([]byte(`{
			"name": "test_model",
			"fields": {
				"` + name + `": {"type": "string", "label": "Reserved"}
			}
		}`))
		if result.Valid {
			t.Errorf("expected valid=false for reserved namespace %q", name)
		}
		found := false
		for _, e := range result.Errors {
			if e.Field == `fields["`+name+`"]` {
				found = true
			}
		}
		if !found {
			t.Errorf("expected error for reserved namespace %q, got: %+v", name, result.Errors)
		}
	}
}

func TestModelValidator_RelationConfig_Many2OneWithoutModel(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "bad_relation_model",
		"fields": {
			"category": {"type": "many2one", "label": "Category"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for many2one without model")
	}
	foundModelError := false
	for _, e := range result.Errors {
		if e.Field == `fields["category"].model` {
			foundModelError = true
		}
	}
	if !foundModelError {
		t.Fatalf("expected error about missing model on many2one, got: %+v", result.Errors)
	}
}

func TestModelValidator_RelationConfig_One2ManyWithoutInverse(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "bad_o2m_model",
		"fields": {
			"items": {"type": "one2many", "model": "item", "label": "Items"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for one2many without inverse")
	}
	foundInverseError := false
	for _, e := range result.Errors {
		if e.Field == `fields["items"].inverse` {
			foundInverseError = true
		}
	}
	if !foundInverseError {
		t.Fatalf("expected error about missing inverse, got: %+v", result.Errors)
	}
}

func TestModelValidator_RelationTarget_NotRegistered(t *testing.T) {
	registry := &stubModelResolver{
		models: map[string]*parser.ModelDefinition{
			"category": {Name: "category"},
		},
	}
	v := NewModelValidator(registry)
	result := v.Validate([]byte(`{
		"name": "product",
		"fields": {
			"category": {"type": "many2one", "model": "category", "label": "Category"},
			"tags": {"type": "many2many", "model": "nonexistent", "label": "Tags"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for unregistered relation target")
	}
	foundTargetError := false
	for _, e := range result.Errors {
		if e.Field == `fields["tags"].model` {
			foundTargetError = true
		}
	}
	if !foundTargetError {
		t.Fatalf("expected error about unregistered model, got: %+v", result.Errors)
	}
}

func TestModelValidator_RelationTarget_Registered(t *testing.T) {
	registry := &stubModelResolver{
		models: map[string]*parser.ModelDefinition{
			"category": {Name: "category"},
		},
	}
	v := NewModelValidator(registry)
	result := v.Validate([]byte(`{
		"name": "product",
		"fields": {
			"category": {"type": "many2one", "model": "category", "label": "Category"}
		}
	}`))
	if !result.Valid {
		t.Fatalf("expected valid=true for registered relation target, got false; errors: %+v", result.Errors)
	}
}

func TestModelValidator_SelfReference_NoCycle(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "tree_node",
		"fields": {
			"name": {"type": "string", "label": "Name"},
			"parent": {"type": "many2one", "model": "tree_node", "label": "Parent"},
			"children": {"type": "one2many", "model": "tree_node", "inverse": "parent", "label": "Children"}
		}
	}`))
	if !result.Valid {
		t.Fatalf("expected valid=true for self-reference (tree) without cycle, got false; errors: %+v", result.Errors)
	}
}

func TestModelValidator_CircularRelation(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "circular_model",
		"fields": {
			"name": {"type": "string", "label": "Name"},
			"items": {"type": "one2many", "model": "circular_model", "inverse": "owner", "label": "Items"},
			"owner": {"type": "one2many", "model": "circular_model", "inverse": "items", "label": "Owner"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for circular relation")
	}
	foundCircularError := false
	for _, e := range result.Errors {
		if contains(e.Message, "circular") {
			foundCircularError = true
		}
	}
	if !foundCircularError {
		t.Fatalf("expected error about circular relation, got: %+v", result.Errors)
	}
}

func TestModelValidator_AllValidFieldTypes(t *testing.T) {
	fieldTypes := map[string]string{
		"f_string":      `"string"`,
		"f_text":        `"text"`,
		"f_integer":     `"integer"`,
		"f_decimal":     `"decimal"`,
		"f_boolean":     `"boolean"`,
		"f_date":        `"date"`,
		"f_datetime":    `"datetime"`,
		"f_email":       `"email"`,
		"f_json":        `"json"`,
		"f_file":        `"file"`,
		"f_smalltext":   `"smalltext"`,
		"f_richtext":    `"richtext"`,
		"f_markdown":    `"markdown"`,
		"f_html":        `"html"`,
		"f_code":        `"code"`,
		"f_password":    `"password"`,
		"f_float":       `"float"`,
		"f_currency":    `"currency"`,
		"f_percent":     `"percent"`,
		"f_toggle":      `"toggle"`,
		"f_time":        `"time"`,
		"f_duration":    `"duration"`,
		"f_image":       `"image"`,
		"f_signature":   `"signature"`,
		"f_barcode":     `"barcode"`,
		"f_color":       `"color"`,
		"f_geolocation": `"geolocation"`,
		"f_rating":      `"rating"`,
		"f_uuid":        `"uuid"`,
		"f_ip":          `"ip"`,
		"f_ipv6":        `"ipv6"`,
		"f_year":        `"year"`,
		"f_binary":      `"binary"`,
		"f_jsonobj":     `"json:object"`,
		"f_jsonarr":     `"json:array"`,
	}

	fieldsJSON := ""
	first := true
	for name, typ := range fieldTypes {
		if !first {
			fieldsJSON += ","
		}
		first = false
		fieldsJSON += `"` + name + `": {"type": ` + typ + `, "label": "` + name + `"}`
	}

	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "all_types_model",
		"fields": {
			` + fieldsJSON + `
		}
	}`))
	if !result.Valid {
		t.Fatalf("expected valid=true for all valid field types, got false; errors: %+v", result.Errors)
	}
}

func TestModelValidator_MorphFieldWithoutModel(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "morph_model",
		"fields": {
			"imageable": {"type": "morph_one", "morph": "imageable", "label": "Imageable"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for morph_one without model")
	}
	foundError := false
	for _, e := range result.Errors {
		if e.Field == `fields["imageable"].model` {
			foundError = true
		}
	}
	if !foundError {
		t.Fatalf("expected error about missing model on morph_one, got: %+v", result.Errors)
	}
}

func TestModelValidator_MorphFieldWithoutMorph(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "morph_model",
		"fields": {
			"imageable": {"type": "morph_many", "model": "image", "label": "Imageable"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for morph_many without morph")
	}
	foundError := false
	for _, e := range result.Errors {
		if e.Field == `fields["imageable"].morph` {
			foundError = true
		}
	}
	if !foundError {
		t.Fatalf("expected error about missing morph name, got: %+v", result.Errors)
	}
}

func TestModelValidator_NameWithSpaces(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "model with spaces",
		"fields": {
			"name": {"type": "string", "label": "Name"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for name with spaces")
	}
	foundNameError := false
	for _, e := range result.Errors {
		if e.Field == "name" && contains(e.Message, "spaces") {
			foundNameError = true
		}
	}
	if !foundNameError {
		t.Fatalf("expected error about name containing spaces, got: %+v", result.Errors)
	}
}

func TestModelValidator_VectorWithoutDimensions(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "vector_model",
		"fields": {
			"embedding": {"type": "vector", "label": "Embedding"}
		}
	}`))
	if result.Valid {
		t.Fatal("expected valid=false for vector without dimensions")
	}
	foundError := false
	for _, e := range result.Errors {
		if e.Field == `fields["embedding"].dimensions` {
			foundError = true
		}
	}
	if !foundError {
		t.Fatalf("expected error about missing dimensions, got: %+v", result.Errors)
	}
}

func TestModelValidator_VectorWithDimensions(t *testing.T) {
	v := NewModelValidator(nil)
	result := v.Validate([]byte(`{
		"name": "vector_model",
		"fields": {
			"embedding": {"type": "vector", "dimensions": 1536, "label": "Embedding"}
		}
	}`))
	if !result.Valid {
		t.Fatalf("expected valid=true for vector with dimensions, got false; errors: %+v", result.Errors)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
