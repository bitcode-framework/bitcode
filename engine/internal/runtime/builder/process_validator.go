package builder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
)

// ProcessResolver resolves process names for validating call-step references.
type ProcessResolver interface {
	LoadProcess(name string) (*parser.ProcessDefinition, error)
}

// ProcessValidator validates process.json definitions: structure, step refs,
// and transition completeness (sequential + DAG modes).
type ProcessValidator struct {
	processRegistry ProcessResolver
}

// NewProcessValidator creates a validator. processRegistry may be nil to skip
// cross-process call-step reference checks.
func NewProcessValidator(processRegistry ProcessResolver) *ProcessValidator {
	return &ProcessValidator{processRegistry: processRegistry}
}

var validStepTypes = map[parser.StepType]bool{
	parser.StepValidate: true,
	parser.StepQuery:    true,
	parser.StepCreate:   true,
	parser.StepUpdate:   true,
	parser.StepDelete:   true,
	parser.StepIf:       true,
	parser.StepSwitch:   true,
	parser.StepLoop:     true,
	parser.StepEmit:     true,
	parser.StepCall:     true,
	parser.StepScript:   true,
	parser.StepHTTP:     true,
	parser.StepAssign:   true,
	parser.StepLog:      true,
	parser.StepUpsert:   true,
	parser.StepCount:    true,
	parser.StepSum:      true,
}

// Validate validates a process definition JSON blob and returns a structured
// result with errors and warnings.
func (v *ProcessValidator) Validate(data []byte) *ValidationResult {
	result := &ValidationResult{
		Errors:   []ValidationIssue{},
		Warnings: []ValidationIssue{},
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "root",
			Message: fmt.Sprintf("invalid JSON: %s", err.Error()),
		})
		result.Valid = false
		return result
	}

	var proc parser.ProcessDefinition
	if err := json.Unmarshal(data, &proc); err != nil {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "root",
			Message: fmt.Sprintf("cannot parse process definition: %s", err.Error()),
		})
		result.Valid = false
		return result
	}

	v.validateStructure(&proc, result)

	if proc.IsDAG() {
		v.validateDAG(&proc, result)
	} else {
		v.validateSequential(&proc, result)
	}

	result.Valid = len(result.Errors) == 0
	return result
}

func (v *ProcessValidator) validateStructure(proc *parser.ProcessDefinition, result *ValidationResult) {
	if proc.Name == "" {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "name",
			Message: "process name is required",
		})
	}

	if strings.Contains(proc.Name, " ") {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "name",
			Message: "process name must not contain spaces",
		})
	}

	if !proc.IsDAG() && len(proc.Steps) == 0 {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "steps",
			Message: "process must have at least one step",
		})
	}

	if proc.IsDAG() && len(proc.Nodes) == 0 {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "nodes",
			Message: "DAG process must have at least one node",
		})
	}
}

// --- Sequential process validation ---

func (v *ProcessValidator) validateSequential(proc *parser.ProcessDefinition, result *ValidationResult) {
	stepNames := make(map[string]bool)
	for i := range proc.Steps {
		step := &proc.Steps[i]
		v.validateStepType(step, i, result)
		if step.Name != "" {
			stepNames[step.Name] = true
		}
	}

	for i := range proc.Steps {
		step := &proc.Steps[i]
		v.validateStepConfig(step, i, proc.Name, result)
		v.validateStepTransitions(step, i, stepNames, result)
		v.validateCallReference(step, i, result)
	}
}

func (v *ProcessValidator) validateStepType(step *parser.StepDefinition, index int, result *ValidationResult) {
	if step.Type == "" {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   fmt.Sprintf("steps[%d].type", index),
			Message: fmt.Sprintf("step %d must have a type", index),
		})
		return
	}
	if !validStepTypes[step.Type] {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   fmt.Sprintf("steps[%d].type", index),
			Message: fmt.Sprintf("invalid step type %q: must be one of %s", step.Type, formatValidStepTypes()),
		})
	}
}

func (v *ProcessValidator) validateStepConfig(step *parser.StepDefinition, index int, procName string, result *ValidationResult) {
	fieldPrefix := fmt.Sprintf("steps[%d]", index)
	if step.Name != "" {
		fieldPrefix = fmt.Sprintf("steps[%q]", step.Name)
	}

	switch step.Type {
	case parser.StepQuery, parser.StepCreate, parser.StepUpdate, parser.StepDelete,
		parser.StepUpsert, parser.StepCount, parser.StepSum:
		if step.Model == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fieldPrefix + ".model",
				Message: fmt.Sprintf("step of type %q requires a model", step.Type),
			})
		}
	case parser.StepIf:
		if step.Condition == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fieldPrefix + ".condition",
				Message: "if step requires a condition",
			})
		}
	case parser.StepSwitch:
		if step.Field == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fieldPrefix + ".field",
				Message: "switch step requires a field",
			})
		}
	case parser.StepLoop:
		if step.Over == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fieldPrefix + ".over",
				Message: "loop step requires an over expression",
			})
		}
	case parser.StepEmit:
		if step.Event == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fieldPrefix + ".event",
				Message: "emit step requires an event name",
			})
		}
	case parser.StepHTTP:
		if step.URL == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fieldPrefix + ".url",
				Message: "http step requires a url",
			})
		}
	case parser.StepAssign:
		if step.Variable == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fieldPrefix + ".variable",
				Message: "assign step requires a variable name",
			})
		}
	}
}

// validateStepTransitions checks that then/else/case references point to
// existing step names (sequential goto-style transitions).
func (v *ProcessValidator) validateStepTransitions(step *parser.StepDefinition, index int, stepNames map[string]bool, result *ValidationResult) {
	fieldPrefix := fmt.Sprintf("steps[%d]", index)
	if step.Name != "" {
		fieldPrefix = fmt.Sprintf("steps[%q]", step.Name)
	}

	if step.Then != "" && !stepNames[step.Then] {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   fieldPrefix + ".then",
			Message: fmt.Sprintf("transition target %q does not match any step name", step.Then),
		})
	}
	if step.Else != "" && !stepNames[step.Else] {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   fieldPrefix + ".else",
			Message: fmt.Sprintf("transition target %q does not match any step name", step.Else),
		})
	}
	for caseValue, target := range step.Cases {
		if target != "" && !stepNames[target] {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("%s.cases[%q]", fieldPrefix, caseValue),
				Message: fmt.Sprintf("case %q references step %q which does not exist", caseValue, target),
			})
		}
	}
}

func (v *ProcessValidator) validateCallReference(step *parser.StepDefinition, index int, result *ValidationResult) {
	if step.Type != parser.StepCall || step.Process == "" || v.processRegistry == nil {
		return
	}

	fieldPrefix := fmt.Sprintf("steps[%d]", index)
	if step.Name != "" {
		fieldPrefix = fmt.Sprintf("steps[%q]", step.Name)
	}

	if _, err := v.processRegistry.LoadProcess(step.Process); err != nil {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   fieldPrefix + ".process",
			Message: fmt.Sprintf("call step references process %q which is not registered", step.Process),
		})
	}
}

// --- DAG process validation ---

func (v *ProcessValidator) validateDAG(proc *parser.ProcessDefinition, result *ValidationResult) {
	nodeIDs := make(map[string]bool)
	for id, node := range proc.Nodes {
		if node.Type == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("nodes[%q].type", id),
				Message: fmt.Sprintf("node %q must have a type", id),
			})
		} else if !validStepTypes[node.Type] {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("nodes[%q].type", id),
				Message: fmt.Sprintf("invalid node type %q: must be one of %s", node.Type, formatValidStepTypes()),
			})
		}
		nodeIDs[id] = true

		if node.Type == parser.StepCall && node.Process != "" && v.processRegistry != nil {
			if _, err := v.processRegistry.LoadProcess(node.Process); err != nil {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("nodes[%q].process", id),
					Message: fmt.Sprintf("call node references process %q which is not registered", node.Process),
				})
			}
		}
	}

	for i, edge := range proc.Edges {
		if edge.From == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("edges[%d].from", i),
				Message: fmt.Sprintf("edge %d must have a 'from' node", i),
			})
		} else if !nodeIDs[edge.From] {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("edges[%d].from", i),
				Message: fmt.Sprintf("edge %d references unknown node %q", i, edge.From),
			})
		}

		if edge.To == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("edges[%d].to", i),
				Message: fmt.Sprintf("edge %d must have a 'to' node", i),
			})
		} else if !nodeIDs[edge.To] {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("edges[%d].to", i),
				Message: fmt.Sprintf("edge %d references unknown node %q", i, edge.To),
			})
		}
	}

	v.validateDAGReachability(proc, nodeIDs, result)
}

// validateDAGReachability detects orphan nodes — nodes that are never targeted
// by any edge (unreachable except for the implicit start nodes). This is a
// warning-level issue, not an error, since a process may have multiple entry
// points.
func (v *ProcessValidator) validateDAGReachability(proc *parser.ProcessDefinition, nodeIDs map[string]bool, result *ValidationResult) {
	targeted := make(map[string]bool)
	for _, edge := range proc.Edges {
		targeted[edge.To] = true
	}

	for id := range nodeIDs {
		if !targeted[id] {
			isSource := false
			for _, edge := range proc.Edges {
				if edge.From == id {
					isSource = true
					break
				}
			}
			if !isSource {
				result.Warnings = append(result.Warnings, ValidationIssue{
					Field:   fmt.Sprintf("nodes[%q]", id),
					Message: fmt.Sprintf("node %q is isolated: no edges target or originate from it", id),
				})
			}
		}
	}
}

func formatValidStepTypes() string {
	types := make([]string, 0, len(validStepTypes))
	for t := range validStepTypes {
		types = append(types, string(t))
	}
	return strings.Join(types, ", ")
}
