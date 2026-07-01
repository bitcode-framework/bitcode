package builder

import (
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
)

type stubProcessResolver struct {
	processes map[string]*parser.ProcessDefinition
}

func (s *stubProcessResolver) LoadProcess(name string) (*parser.ProcessDefinition, error) {
	if p, ok := s.processes[name]; ok {
		return p, nil
	}
	return nil, &processNotFoundError{name: name}
}

type processNotFoundError struct{ name string }

func (e *processNotFoundError) Error() string { return "process not found: " + e.name }

func testProcessResolver() *stubProcessResolver {
	return &stubProcessResolver{
		processes: map[string]*parser.ProcessDefinition{
			"send_email": {
				Name: "send_email",
				Steps: []parser.StepDefinition{
					{Name: "emit_sent", Type: parser.StepEmit, Event: "email_sent"},
				},
			},
		},
	}
}

func validSequentialProcessJSON() []byte {
	return []byte(`{
		"name": "create_lead",
		"steps": [
			{"name": "validate_input", "type": "validate", "rules": {}},
			{"name": "create_record", "type": "create", "model": "lead"},
			{"name": "emit_created", "type": "emit", "event": "lead_created"}
		]
	}`)
}

func validDAGProcessJSON() []byte {
	return []byte(`{
		"name": "approval_flow",
		"nodes": {
			"start":    {"type": "query", "model": "invoice"},
			"validate": {"type": "validate", "rules": {}},
			"approve":  {"type": "update", "model": "invoice"},
			"notify":   {"type": "emit", "event": "approved"}
		},
		"edges": [
			{"from": "start", "to": "validate"},
			{"from": "validate", "to": "approve"},
			{"from": "approve", "to": "notify"}
		]
	}`)
}

// --- Structure validation ---

func TestProcessValidator_ValidSequential(t *testing.T) {
	v := NewProcessValidator(testProcessResolver())
	result := v.Validate(validSequentialProcessJSON())

	if !result.Valid {
		t.Fatalf("expected valid=true, got false; errors: %+v", result.Errors)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected 0 errors, got %d: %+v", len(result.Errors), result.Errors)
	}
}

func TestProcessValidator_ValidDAG(t *testing.T) {
	v := NewProcessValidator(testProcessResolver())
	result := v.Validate(validDAGProcessJSON())

	if !result.Valid {
		t.Fatalf("expected valid=true, got false; errors: %+v", result.Errors)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected 0 errors, got %d: %+v", len(result.Errors), result.Errors)
	}
}

func TestProcessValidator_MissingName(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{"steps": [{"type": "query", "model": "lead"}]}`)

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

func TestProcessValidator_NameWithSpaces(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{"name": "my process", "steps": [{"type": "query", "model": "lead"}]}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for name with spaces")
	}
}

func TestProcessValidator_NoSteps(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{"name": "empty_process"}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for no steps")
	}
}

// --- Step type validation ---

func TestProcessValidator_InvalidStepType(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "bad_type",
		"steps": [
			{"type": "nonexistent_type"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid step type")
	}
	foundTypeError := false
	for _, e := range result.Errors {
		if e.Field == "steps[0].type" {
			foundTypeError = true
		}
	}
	if !foundTypeError {
		t.Fatalf("expected error about invalid step type, got: %+v", result.Errors)
	}
}

func TestProcessValidator_MissingStepType(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "no_type",
		"steps": [
			{"name": "untitled"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for missing step type")
	}
}

// --- Step config validation ---

func TestProcessValidator_QueryStepMissingModel(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "missing_model",
		"steps": [
			{"type": "query"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for query step without model")
	}
}

func TestProcessValidator_IfStepMissingCondition(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "missing_condition",
		"steps": [
			{"type": "if"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for if step without condition")
	}
}

func TestProcessValidator_EmitStepMissingEvent(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "missing_event",
		"steps": [
			{"type": "emit"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for emit step without event")
	}
}

func TestProcessValidator_HTTPStepMissingURL(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "missing_url",
		"steps": [
			{"type": "http", "method": "GET"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for http step without url")
	}
}

// --- Transition validation (invalid step ref / orphan transition) ---

func TestProcessValidator_InvalidStepRef_Then(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "bad_then",
		"steps": [
			{"name": "check", "type": "if", "condition": "x > 0", "then": "nonexistent_step", "else": "done"},
			{"name": "done", "type": "log", "message": "done"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid then ref")
	}
	foundRefError := false
	for _, e := range result.Errors {
		if e.Field == `steps["check"].then` {
			foundRefError = true
		}
	}
	if !foundRefError {
		t.Fatalf("expected error about invalid then transition target, got: %+v", result.Errors)
	}
}

func TestProcessValidator_InvalidStepRef_Else(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "bad_else",
		"steps": [
			{"name": "check", "type": "if", "condition": "x > 0", "then": "done", "else": "ghost_step"},
			{"name": "done", "type": "log", "message": "done"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid else ref")
	}
}

func TestProcessValidator_InvalidStepRef_Case(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "bad_case",
		"steps": [
			{"name": "branch", "type": "switch", "field": "status", "cases": {"active": "ghost_step", "closed": "done"}},
			{"name": "done", "type": "log", "message": "done"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for invalid case ref")
	}
	foundCaseError := false
	for _, e := range result.Errors {
		if e.Field == `steps["branch"].cases["active"]` {
			foundCaseError = true
		}
	}
	if !foundCaseError {
		t.Fatalf("expected error about invalid case transition target, got: %+v", result.Errors)
	}
}

func TestProcessValidator_ValidStepRefs(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "valid_refs",
		"steps": [
			{"name": "check", "type": "if", "condition": "x > 0", "then": "yes", "else": "no"},
			{"name": "yes", "type": "log", "message": "yes"},
			{"name": "no", "type": "log", "message": "no"}
		]
	}`)

	result := v.Validate(data)

	if !result.Valid {
		t.Fatalf("expected valid=true for valid refs, got false; errors: %+v", result.Errors)
	}
}

func TestProcessValidator_OrphanTransition_DAG(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "dag_bad_edge",
		"nodes": {
			"start":  {"type": "query", "model": "lead"},
			"save":   {"type": "create", "model": "lead"}
		},
		"edges": [
			{"from": "start", "to": "ghost_node"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for DAG with orphan edge (ghost_node)")
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

func TestProcessValidator_DAG_OrphanNodeWarning(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "dag_isolated",
		"nodes": {
			"start":  {"type": "query", "model": "lead"},
			"save":   {"type": "create", "model": "lead"},
			"orphan": {"type": "log", "message": "nobody calls me"}
		},
		"edges": [
			{"from": "start", "to": "save"}
		]
	}`)

	result := v.Validate(data)

	if !result.Valid {
		t.Fatalf("expected valid=true (orphan is a warning), got false; errors: %+v", result.Errors)
	}
	foundWarning := false
	for _, w := range result.Warnings {
		if w.Field == `nodes["orphan"]` {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Fatalf("expected warning about isolated node, got: %+v", result.Warnings)
	}
}

func TestProcessValidator_DAG_EdgeMissingFrom(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "dag_no_from",
		"nodes": {
			"a": {"type": "query", "model": "lead"},
			"b": {"type": "create", "model": "lead"}
		},
		"edges": [
			{"to": "b"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for edge missing from")
	}
}

// --- Call step cross-process ref validation ---

func TestProcessValidator_CallStepValidRef(t *testing.T) {
	v := NewProcessValidator(testProcessResolver())
	data := []byte(`{
		"name": "lead_with_email",
		"steps": [
			{"type": "create", "model": "lead"},
			{"type": "call", "process": "send_email"}
		]
	}`)

	result := v.Validate(data)

	if !result.Valid {
		t.Fatalf("expected valid=true for call to existing process, got false; errors: %+v", result.Errors)
	}
}

func TestProcessValidator_CallStepInvalidRef(t *testing.T) {
	v := NewProcessValidator(testProcessResolver())
	data := []byte(`{
		"name": "lead_with_bad_call",
		"steps": [
			{"type": "create", "model": "lead"},
			{"name": "email", "type": "call", "process": "nonexistent_process"}
		]
	}`)

	result := v.Validate(data)

	if result.Valid {
		t.Fatal("expected valid=false for call to nonexistent process")
	}
	foundCallError := false
	for _, e := range result.Errors {
		if e.Field == `steps["email"].process` {
			foundCallError = true
		}
	}
	if !foundCallError {
		t.Fatalf("expected error about nonexistent process call, got: %+v", result.Errors)
	}
}

func TestProcessValidator_CallStepNilRegistry(t *testing.T) {
	v := NewProcessValidator(nil)
	data := []byte(`{
		"name": "call_no_registry",
		"steps": [
			{"type": "create", "model": "lead"},
			{"type": "call", "process": "any_process"}
		]
	}`)

	result := v.Validate(data)

	if !result.Valid {
		t.Fatalf("expected valid=true with nil registry (skips call ref checks), got false; errors: %+v", result.Errors)
	}
}

// --- Invalid JSON ---

func TestProcessValidator_InvalidJSON(t *testing.T) {
	v := NewProcessValidator(nil)
	result := v.Validate([]byte(`{invalid json`))

	if result.Valid {
		t.Fatal("expected valid=false for invalid JSON")
	}
	if len(result.Errors) == 0 {
		t.Fatal("expected at least one error")
	}
}

func TestProcessValidator_NilRegistry(t *testing.T) {
	v := NewProcessValidator(nil)
	result := v.Validate(validSequentialProcessJSON())

	if !result.Valid {
		t.Fatalf("expected valid=true with nil registry, got false; errors: %+v", result.Errors)
	}
}
