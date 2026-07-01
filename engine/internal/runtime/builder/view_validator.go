package builder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
)

type ValidationIssue struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationIssue `json:"errors"`
	Warnings []ValidationIssue `json:"warnings"`
}

type ModelResolver interface {
	Get(name string) (*parser.ModelDefinition, error)
}

type ViewValidator struct {
	modelRegistry ModelResolver
}

func NewViewValidator(modelRegistry ModelResolver) *ViewValidator {
	return &ViewValidator{modelRegistry: modelRegistry}
}

var validViewTypes = map[parser.ViewType]bool{
	parser.ViewList:     true,
	parser.ViewForm:     true,
	parser.ViewKanban:   true,
	parser.ViewCalendar: true,
	parser.ViewChart:    true,
	parser.ViewCustom:   true,
	parser.ViewGantt:    true,
	parser.ViewMap:      true,
	parser.ViewTree:     true,
	parser.ViewActivity: true,
	parser.ViewReport:   true,
	parser.ViewPivot:    true,
}

func (v *ViewValidator) Validate(data []byte) *ValidationResult {
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

	var view parser.ViewDefinition
	if err := json.Unmarshal(data, &view); err != nil {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "root",
			Message: fmt.Sprintf("cannot parse view definition: %s", err.Error()),
		})
		result.Valid = false
		return result
	}

	v.validateStructure(&view, raw, result)
	v.validateDataSources(&view, result)
	v.validateModelReference(&view, result)
	v.validateFieldReferences(&view, result)

	result.Valid = len(result.Errors) == 0
	return result
}

func (v *ViewValidator) validateStructure(view *parser.ViewDefinition, raw map[string]any, result *ValidationResult) {
	if view.Name == "" {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "name",
			Message: "view name is required",
		})
	}

	if strings.Contains(view.Name, " ") {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "name",
			Message: "view name must not contain spaces",
		})
	}

	if view.Type == "" {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "type",
			Message: "view type is required",
		})
	} else if !validViewTypes[view.Type] {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "type",
			Message: fmt.Sprintf("invalid view type %q: must be one of list, form, kanban, calendar, chart, custom, gantt, map, tree, activity, report, pivot", view.Type),
		})
	}

	if view.Type == parser.ViewCustom && view.Template == "" {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "template",
			Message: "custom view requires a template",
		})
	}

	if view.Render != "" && view.Render != parser.RenderSSR && view.Render != parser.RenderSPA {
		result.Errors = append(result.Errors, ValidationIssue{
			Field:   "render",
			Message: fmt.Sprintf("invalid render mode %q: must be %q or %q", view.Render, parser.RenderSSR, parser.RenderSPA),
		})
	}

	if len(view.Layout) == 0 && view.Type == parser.ViewForm {
		result.Warnings = append(result.Warnings, ValidationIssue{
			Field:   "layout",
			Message: "form view has no layout defined; all fields will be auto-generated",
		})
	}
}

func (v *ViewValidator) validateDataSources(view *parser.ViewDefinition, result *ValidationResult) {
	for dsName, ds := range view.DataSources {
		if ds.Model != "" && ds.Process != "" {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("data_sources.%s", dsName),
				Message: fmt.Sprintf("data source %q cannot have both model and process", dsName),
			})
		}
	}
}

func (v *ViewValidator) validateModelReference(view *parser.ViewDefinition, result *ValidationResult) {
	if view.Model == "" {
		return
	}

	if v.modelRegistry == nil {
		return
	}

	modelDef, err := v.modelRegistry.Get(view.Model)
	if err != nil {
		_, errQualified := v.modelRegistry.Get(view.Model)
		if errQualified != nil {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   "model",
				Message: fmt.Sprintf("model %q referenced by view not found in registry", view.Model),
			})
			return
		}
		modelDef, _ = v.modelRegistry.Get(view.Model)
	}

	if modelDef == nil {
		return
	}
}

func (v *ViewValidator) validateFieldReferences(view *parser.ViewDefinition, result *ValidationResult) {
	if v.modelRegistry == nil || view.Model == "" {
		return
	}

	modelDef, err := v.modelRegistry.Get(view.Model)
	if err != nil {
		return
	}

	for _, field := range view.Fields {
		if _, ok := modelDef.Fields[field]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("fields[%q]", field),
				Message: fmt.Sprintf("field %q does not exist on model %q", field, view.Model),
			})
		}
	}

	for _, f := range view.Filters {
		if _, ok := modelDef.Fields[f]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("filters[%q]", f),
				Message: fmt.Sprintf("filter field %q does not exist on model %q", f, view.Model),
			})
		}
	}

	if view.Sort != nil && view.Sort.Field != "" {
		if _, ok := modelDef.Fields[view.Sort.Field]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("sort.field"),
				Message: fmt.Sprintf("sort field %q does not exist on model %q", view.Sort.Field, view.Model),
			})
		}
	}

	if view.DateField != "" {
		if _, ok := modelDef.Fields[view.DateField]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   "date_field",
				Message: fmt.Sprintf("date_field %q does not exist on model %q", view.DateField, view.Model),
			})
		}
	}
	if view.StartField != "" {
		if _, ok := modelDef.Fields[view.StartField]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   "start_field",
				Message: fmt.Sprintf("start_field %q does not exist on model %q", view.StartField, view.Model),
			})
		}
	}
	if view.EndField != "" {
		if _, ok := modelDef.Fields[view.EndField]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   "end_field",
				Message: fmt.Sprintf("end_field %q does not exist on model %q", view.EndField, view.Model),
			})
		}
	}
	if view.ParentField != "" {
		if _, ok := modelDef.Fields[view.ParentField]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   "parent_field",
				Message: fmt.Sprintf("parent_field %q does not exist on model %q", view.ParentField, view.Model),
			})
		}
	}
	if view.GroupBy != "" {
		if _, ok := modelDef.Fields[view.GroupBy]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   "group_by",
				Message: fmt.Sprintf("group_by field %q does not exist on model %q", view.GroupBy, view.Model),
			})
		}
	}

	v.validateLayoutFields(view.Layout, modelDef, result)
}

func (v *ViewValidator) validateLayoutFields(layout []parser.LayoutItem, modelDef *parser.ModelDefinition, result *ValidationResult) {
	for i := range layout {
		v.validateLayoutItem(&layout[i], modelDef, result)
	}
}

func (v *ViewValidator) validateLayoutItem(item *parser.LayoutItem, modelDef *parser.ModelDefinition, result *ValidationResult) {
	for _, row := range item.Row {
		if row.Field == "" {
			continue
		}
		if _, ok := modelDef.Fields[row.Field]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("layout.row[%q]", row.Field),
				Message: fmt.Sprintf("layout field %q does not exist on model %q", row.Field, modelDef.Name),
			})
		}
	}

	for _, tab := range item.Tabs {
		for _, tabField := range tab.Fields {
			if tabField == "" {
				continue
			}
			if _, ok := modelDef.Fields[tabField]; !ok {
				result.Errors = append(result.Errors, ValidationIssue{
					Field:   fmt.Sprintf("layout.tab[%q].fields[%q]", tab.Label, tabField),
					Message: fmt.Sprintf("tab field %q does not exist on model %q", tabField, modelDef.Name),
				})
			}
		}
	}

	if item.ChildTable != nil && item.ChildTable.Field != "" {
		if _, ok := modelDef.Fields[item.ChildTable.Field]; !ok {
			result.Errors = append(result.Errors, ValidationIssue{
				Field:   fmt.Sprintf("layout.child_table[%q]", item.ChildTable.Field),
				Message: fmt.Sprintf("child table field %q does not exist on model %q", item.ChildTable.Field, modelDef.Name),
			})
		}
	}

	for i := range item.Rows {
		v.validateLayoutItem(&item.Rows[i], modelDef, result)
	}
}
