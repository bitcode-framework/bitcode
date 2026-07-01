package builder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/domain/model"
)

// ModelValidator validates model.json definitions: structure, field type
// validity, relation integrity (including circular relation detection), and
// reserved namespace checks.
type ModelValidator struct {
	// modelRegistry is used to resolve relational targets (many2one, one2many,
	// many2many, morph variants). May be nil to skip cross-model relation checks.
	modelRegistry ModelResolver
}

// NewModelValidator creates a validator. modelRegistry may be nil to skip
// cross-model relational target checks (e.g. when validating a standalone
// model.json without the full registry loaded).
func NewModelValidator(modelRegistry ModelResolver) *ModelValidator {
	return &ModelValidator{modelRegistry: modelRegistry}
}

// validFieldTypes is the complete set of field types accepted by the engine
// parser (see parser/model.go FieldType constants).
var validFieldTypes = map[parser.FieldType]bool{
	parser.FieldString:        true,
	parser.FieldText:          true,
	parser.FieldInteger:       true,
	parser.FieldDecimal:       true,
	parser.FieldBoolean:       true,
	parser.FieldDate:          true,
	parser.FieldDatetime:      true,
	parser.FieldSelection:     true,
	parser.FieldEmail:         true,
	parser.FieldMany2One:      true,
	parser.FieldOne2Many:      true,
	parser.FieldMany2Many:     true,
	parser.FieldJSON:          true,
	parser.FieldFile:          true,
	parser.FieldComputed:      true,
	parser.FieldSmallText:     true,
	parser.FieldRichText:      true,
	parser.FieldMarkdown:      true,
	parser.FieldHTML:          true,
	parser.FieldCode:          true,
	parser.FieldPassword:      true,
	parser.FieldFloat:         true,
	parser.FieldCurrency:      true,
	parser.FieldPercent:       true,
	parser.FieldToggle:        true,
	parser.FieldRadio:         true,
	parser.FieldDynamicLink:   true,
	parser.FieldTime:          true,
	parser.FieldDuration:      true,
	parser.FieldImage:         true,
	parser.FieldSignature:     true,
	parser.FieldBarcode:       true,
	parser.FieldColor:         true,
	parser.FieldGeolocation:   true,
	parser.FieldRating:        true,
	parser.FieldUUID:          true,
	parser.FieldIP:            true,
	parser.FieldIPv6:          true,
	parser.FieldYear:          true,
	parser.FieldVector:        true,
	parser.FieldBinary:        true,
	parser.FieldJSONObject:    true,
	parser.FieldJSONArray:     true,
	parser.FieldMorphTo:       true,
	parser.FieldMorphOne:      true,
	parser.FieldMorphMany:     true,
	parser.FieldMorphToMany:   true,
	parser.FieldMorphByMany:   true,
}

// relationalFieldTypes are field types that reference another model via the
// `model` property. Used for circular relation detection and target validation.
var relationalFieldTypes = map[parser.FieldType]bool{
	parser.FieldMany2One:    true,
	parser.FieldOne2Many:    true,
	parser.FieldMany2Many:   true,
	parser.FieldMorphOne:    true,
	parser.FieldMorphMany:   true,
	parser.FieldMorphToMany: true,
}

// Validate validates a model definition JSON blob and returns a structured
// result with errors and warnings.
func (v *ModelValidator) Validate(data []byte) *ValidationResult {
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

	var modelDef parser.ModelDefinition
	if err := json.Unmarshal(data, &modelDef); err != nil {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "root",
			Message: fmt.Sprintf("cannot parse model definition: %s", err.Error()),
		})
		result.Valid = false
		return result
	}

	v.validateStructure(&modelDef, result)
	v.validateFields(&modelDef, result)
	v.validateReservedNamespaces(&modelDef, result)
	v.validateRelations(&modelDef, result)

	result.Valid = len(result.Errors) == 0
	return result
}

func (v *ModelValidator) validateStructure(m *parser.ModelDefinition, result *ValidationResult) {
	if m.Name == "" {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "name",
			Message: "model name is required",
		})
	}

	if strings.Contains(m.Name, " ") {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "name",
			Message: "model name must not contain spaces",
		})
	}

	if len(m.Fields) == 0 {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "fields",
			Message: "model must have at least one field",
		})
	}
}

func (v *ModelValidator) validateFields(m *parser.ModelDefinition, result *ValidationResult) {
	for fieldName, field := range m.Fields {
		if field.Type == "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("fields[%q].type", fieldName),
				Message: fmt.Sprintf("field %q must have a type", fieldName),
			})
			continue
		}

		normalizedType := field.Type
		switch string(field.Type) {
		case "json:object":
			normalizedType = parser.FieldJSONObject
		case "json:array":
			normalizedType = parser.FieldJSONArray
		case "ip:v4":
			normalizedType = parser.FieldIP
		case "ip:v6":
			normalizedType = parser.FieldIPv6
		}

		if !validFieldTypes[normalizedType] {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("fields[%q].type", fieldName),
				Message: fmt.Sprintf("invalid field type %q for field %q: must be a supported type", field.Type, fieldName),
			})
		}

		if normalizedType == parser.FieldSelection && len(field.Options) == 0 {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("fields[%q].options", fieldName),
				Message: fmt.Sprintf("selection field %q must have options", fieldName),
			})
		}

		if normalizedType == parser.FieldRadio && len(field.Options) == 0 {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("fields[%q].options", fieldName),
				Message: fmt.Sprintf("radio field %q must have options", fieldName),
			})
		}

		if normalizedType == parser.FieldVector && field.Dimensions == 0 {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("fields[%q].dimensions", fieldName),
				Message: fmt.Sprintf("vector field %q must specify dimensions", fieldName),
			})
		}
	}
}

func (v *ModelValidator) validateReservedNamespaces(m *parser.ModelDefinition, result *ValidationResult) {
	for fieldName := range m.Fields {
		for _, reserved := range model.ReservedNamespaces {
			if fieldName == reserved {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("fields[%q]", fieldName),
					Message: fmt.Sprintf("field name %q is a reserved expression namespace (%s) and cannot be used", fieldName, strings.Join(model.ReservedNamespaces, ", ")),
				})
			}
		}
	}
}

func (v *ModelValidator) validateRelations(m *parser.ModelDefinition, result *ValidationResult) {
	v.validateRelationConfig(m, result)
	v.validateRelationTargets(m, result)
	v.detectCircularRelations(m, result)
}

func (v *ModelValidator) validateRelationConfig(m *parser.ModelDefinition, result *ValidationResult) {
	for fieldName, field := range m.Fields {
		switch field.Type {
		case parser.FieldMany2One, parser.FieldMany2Many:
			if field.Model == "" {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("fields[%q].model", fieldName),
					Message: fmt.Sprintf("%s field %q must specify a target model", field.Type, fieldName),
				})
			}
		case parser.FieldOne2Many:
			if field.Model == "" {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("fields[%q].model", fieldName),
					Message: fmt.Sprintf("one2many field %q must specify a target model", fieldName),
				})
			}
			if field.Inverse == "" {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("fields[%q].inverse", fieldName),
					Message: fmt.Sprintf("one2many field %q must specify an inverse field", fieldName),
				})
			}
		case parser.FieldMorphOne, parser.FieldMorphMany, parser.FieldMorphToMany, parser.FieldMorphByMany:
			if field.Model == "" {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("fields[%q].model", fieldName),
					Message: fmt.Sprintf("%s field %q must specify a target model", field.Type, fieldName),
				})
			}
			if field.Morph == "" {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("fields[%q].morph", fieldName),
					Message: fmt.Sprintf("%s field %q must specify a morph name", field.Type, fieldName),
				})
			}
		}
	}
}

func (v *ModelValidator) validateRelationTargets(m *parser.ModelDefinition, result *ValidationResult) {
	if v.modelRegistry == nil {
		return
	}

	for fieldName, field := range m.Fields {
		if !relationalFieldTypes[field.Type] {
			continue
		}
		if field.Model == "" {
			continue
		}
		if field.Model == m.Name {
			continue
		}

		if _, err := v.modelRegistry.Get(field.Model); err != nil {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("fields[%q].model", fieldName),
				Message: fmt.Sprintf("%s field %q references model %q which is not registered", field.Type, fieldName, field.Model),
			})
		}
	}
}

// detectCircularRelations checks for circular self-referential relations.
// A one2many or many2many field that targets its own model (self-ref) is valid
// (e.g. parent/child tree), but if the inverse also points back creating a
// cycle, it is an error. A cycle exists when following the inverse chain
// returns to the starting field.
func (v *ModelValidator) detectCircularRelations(m *parser.ModelDefinition, result *ValidationResult) {
	reportedCycles := make(map[string]bool)

	for fieldName, field := range m.Fields {
		if field.Type != parser.FieldOne2Many && field.Type != parser.FieldMany2Many {
			continue
		}
		if field.Model != m.Name {
			continue
		}
		if field.Inverse == "" {
			continue
		}

		visited := make(map[string]bool)
		recursionStack := make(map[string]bool)
		path := []string{fieldName}
		if cyclePath := v.findInverseCycle(m, fieldName, visited, recursionStack, path); cyclePath != nil {
			cycleKey := strings.Join(cyclePath, " -> ")
			normalizedKey := normalizeCycleKey(cyclePath)
			if !reportedCycles[normalizedKey] {
				reportedCycles[normalizedKey] = true
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("fields[%q]", fieldName),
					Message: fmt.Sprintf("circular relation detected: %s", cycleKey),
				})
			}
		}
	}
}

// findInverseCycle follows the inverse chain to detect a cycle. Returns the
// cycle path if found, nil otherwise.
func (v *ModelValidator) findInverseCycle(m *parser.ModelDefinition, current string, visited, recursionStack map[string]bool, path []string) []string {
	if recursionStack[current] {
		return append(path, current)
	}
	if visited[current] {
		return nil
	}

	recursionStack[current] = true
	visited[current] = true

	field, ok := m.Fields[current]
	if !ok {
		recursionStack[current] = false
		return nil
	}

	target := field.Inverse
	if target != "" {
		if _, exists := m.Fields[target]; exists {
			if cycle := v.findInverseCycle(m, target, visited, recursionStack, append(path, target)); cycle != nil {
				recursionStack[current] = false
				return cycle
			}
		}
	}

	recursionStack[current] = false
	return nil
}

// normalizeCycleKey produces a canonical key for a cycle path so the same
// cycle reported from different starting points is deduplicated. It rotates
// the path so the lexicographically smallest element is first.
func normalizeCycleKey(path []string) string {
	if len(path) == 0 {
		return ""
	}
	minIdx := 0
	for i, s := range path {
		if s < path[minIdx] {
			minIdx = i
		}
	}
	rotated := make([]string, 0, len(path))
	rotated = append(rotated, path[minIdx:]...)
	rotated = append(rotated, path[:minIdx]...)
	return strings.Join(rotated, " -> ")
}
