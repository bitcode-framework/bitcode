package pivot

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
)

// AggFunc identifies an aggregation function applied to a value field.
type AggFunc string

const (
	AggSum   AggFunc = "sum"
	AggCount AggFunc = "count"
	AggAvg   AggFunc = "avg"
	AggMin   AggFunc = "min"
	AggMax   AggFunc = "max"
)

// ValidAggFuncs is the set of supported aggregation identifiers.
var ValidAggFuncs = map[AggFunc]bool{
	AggSum:   true,
	AggCount: true,
	AggAvg:   true,
	AggMin:   true,
	AggMax:   true,
}

// ValueField specifies a single measure in the pivot request.
type ValueField struct {
	Field string  `json:"field"`
	Agg   AggFunc `json:"agg"`
	Alias string  `json:"alias,omitempty"`
}

// Request defines a pivot computation. Rows and Cols are dimension field
// names. Values is the list of measures (field + aggregation). Filter is an
// optional query filter applied to the source rows before aggregation.
type Request struct {
	Rows   []string         `json:"rows"`
	Cols   []string         `json:"cols"`
	Values []ValueField     `json:"values"`
	Filter map[string]any   `json:"filter,omitempty"`
}

// Result is the 2D pivot matrix returned by the service.
//
// RowHeaders and ColHeaders are ordered arrays of composite dimension keys
// (joined by " / " for multi-dimension). Cells is indexed
// [rowIndex][colIndex] and each entry has one float64 per ValueField.
// RowTotals/ColTotals/GrandTotal follow the same value-index ordering. For
// avg aggregations the totals are computed from the raw data (not from
// averaged cell values) so they remain mathematically correct.
type Result struct {
	RowHeaders  []string    `json:"rowHeaders"`
	ColHeaders  []string    `json:"colHeaders"`
	Cells       [][][]float64 `json:"cells"`
	RowTotals   [][]float64 `json:"rowTotals"`
	ColTotals   [][]float64 `json:"colTotals"`
	GrandTotal  []float64   `json:"grandTotal"`
	ValueFields []ValueField `json:"valueFields"`
}

// RecordFetcher loads raw rows matching the query. The service is
// transport-agnostic: callers supply the fetcher (HTTP handler, CLI, tests).
type RecordFetcher func(ctx context.Context, query *persistence.Query) ([]map[string]any, error)

// Service computes OLAP pivot aggregation over raw records. It performs
// GROUP BY rows+cols and aggregate measures in Go, so it works with any
// database backend without raw SQL.
type Service struct{}

// NewService constructs a PivotService.
func NewService() *Service {
	return &Service{}
}

// aggState accumulates raw values for a single group so that any aggregation
// function can be derived from it in a single pass.
type aggState struct {
	sum    float64
	count  int64
	min    float64
	max    float64
	hasVal bool
}

func newAggState() *aggState {
	return &aggState{}
}

func (a *aggState) accumulate(agg AggFunc, raw any) {
	if agg == AggCount {
		a.count++
		return
	}
	v, ok := toFloat(raw)
	if !ok || math.IsNaN(v) || math.IsInf(v, 0) {
		return
	}
	a.sum += v
	a.count++
	if !a.hasVal {
		a.min = v
		a.max = v
		a.hasVal = true
	} else {
		if v < a.min {
			a.min = v
		}
		if v > a.max {
			a.max = v
		}
	}
}

func (a *aggState) mergeFrom(other *aggState, agg AggFunc) {
	if other == nil || other.count == 0 {
		return
	}
	a.count += other.count
	if agg != AggCount {
		a.sum += other.sum
		if other.hasVal {
			if !a.hasVal {
				a.min = other.min
				a.max = other.max
				a.hasVal = true
			} else {
				if other.min < a.min {
					a.min = other.min
				}
				if other.max > a.max {
					a.max = other.max
				}
			}
		}
	}
}

func (a *aggState) finalize(agg AggFunc) float64 {
	if a == nil || a.count == 0 {
		return 0
	}
	switch agg {
	case AggSum:
		return a.sum
	case AggCount:
		return float64(a.count)
	case AggAvg:
		return a.sum / float64(a.count)
	case AggMin:
		if !a.hasVal {
			return 0
		}
		return a.min
	case AggMax:
		if !a.hasVal {
			return 0
		}
		return a.max
	}
	return 0
}

// cellKey identifies a unique row+col dimension combination.
type cellKey struct {
	row string
	col string
}

// Validate checks the request for structural correctness against the model
// definition. Every row, col, and value field must exist on the model.
// Returns an error describing the first violation.
func (s *Service) Validate(req *Request, modelDef *parser.ModelDefinition) error {
	if modelDef == nil {
		return fmt.Errorf("model definition is required")
	}
	if len(req.Rows) == 0 {
		return fmt.Errorf("at least one row dimension is required")
	}
	if len(req.Values) == 0 {
		return fmt.Errorf("at least one value field is required")
	}

	for _, f := range req.Rows {
		if !fieldExists(modelDef, f) {
			return fmt.Errorf("row field %q not found on model %q", f, modelDef.Name)
		}
	}
	for _, f := range req.Cols {
		if !fieldExists(modelDef, f) {
			return fmt.Errorf("col field %q not found on model %q", f, modelDef.Name)
		}
	}
	for _, vf := range req.Values {
		if !fieldExists(modelDef, vf.Field) {
			return fmt.Errorf("value field %q not found on model %q", vf.Field, modelDef.Name)
		}
		if !ValidAggFuncs[vf.Agg] {
			return fmt.Errorf("unsupported aggregation %q for field %q (supported: sum, count, avg, min, max)", vf.Agg, vf.Field)
		}
	}

	return nil
}

// Compute performs the pivot aggregation. It fetches raw rows via the
// RecordFetcher (with filter applied), then groups by row+col dimension
// combinations and aggregates each measure. The result includes row totals,
// column totals, and a grand total computed from raw data.
func (s *Service) Compute(ctx context.Context, req *Request, modelDef *parser.ModelDefinition, fetcher RecordFetcher) (*Result, error) {
	if err := s.Validate(req, modelDef); err != nil {
		return nil, err
	}

	query := buildQuery(req.Filter)
	rows, err := fetcher(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("pivot data fetch failed: %w", err)
	}

	return s.Aggregate(req, rows), nil
}

// Aggregate computes the pivot matrix from already-loaded rows. This is the
// pure computation core with no I/O or side effects.
//
// For each row, four levels of aggregation state are maintained in a single
// pass so that avg/min/max totals remain mathematically correct:
//
//	cellAggs[rowKey][colKey][valueIndex] — the intersection cells
//	rowAggs[rowKey][valueIndex]          — row totals (collapse cols)
//	colAggs[colKey][valueIndex]          — column totals (collapse rows)
//	grandAggs[valueIndex]                — grand total (collapse both)
func (s *Service) Aggregate(req *Request, rows []map[string]any) *Result {
	nValues := len(req.Values)

	cellAggs := make(map[cellKey]map[int]*aggState)
	rowAggs := make(map[string]map[int]*aggState)
	colAggs := make(map[string]map[int]*aggState)
	grandAggs := make(map[int]*aggState)

	rowSet := make(map[string]bool)
	colSet := make(map[string]bool)

	ensureMap := func(m map[cellKey]map[int]*aggState, key cellKey) map[int]*aggState {
		if v, ok := m[key]; ok {
			return v
		}
		v := make(map[int]*aggState)
		m[key] = v
		return v
	}
	ensureDimMap := func(m map[string]map[int]*aggState, key string) map[int]*aggState {
		if key == "" {
			return nil
		}
		if v, ok := m[key]; ok {
			return v
		}
		v := make(map[int]*aggState)
		m[key] = v
		return v
	}
	ensureValState := func(m map[int]*aggState, vi int) *aggState {
		if st, ok := m[vi]; ok {
			return st
		}
		st := newAggState()
		m[vi] = st
		return st
	}

	for _, row := range rows {
		rk := dimKey(row, req.Rows)
		ck := dimKey(row, req.Cols)

		rowSet[rk] = true
		if len(req.Cols) > 0 {
			colSet[ck] = true
		}

		cellMap := ensureMap(cellAggs, cellKey{row: rk, col: ck})
		rowMap := ensureDimMap(rowAggs, rk)
		colMap := ensureDimMap(colAggs, ck)

		for vi, vf := range req.Values {
			raw := row[vf.Field]

			if cellMap != nil {
				st := ensureValState(cellMap, vi)
				st.accumulate(vf.Agg, raw)
			}
			if rowMap != nil {
				st := ensureValState(rowMap, vi)
				st.accumulate(vf.Agg, raw)
			}
			if colMap != nil {
				st := ensureValState(colMap, vi)
				st.accumulate(vf.Agg, raw)
			}
			if grandSt, ok := grandAggs[vi]; ok {
				grandSt.accumulate(vf.Agg, raw)
			} else {
				grandSt := newAggState()
				grandSt.accumulate(vf.Agg, raw)
				grandAggs[vi] = grandSt
			}
		}
	}

	rowOrder := sortedKeys(rowSet)
	colOrder := sortedKeys(colSet)
	if len(colOrder) == 0 {
		colOrder = []string{""}
	}

	nRows := len(rowOrder)
	nCols := len(colOrder)

	cells := make([][][]float64, nRows)
	for ri := range cells {
		cells[ri] = make([][]float64, nCols)
		for ci := range cells[ri] {
			cells[ri][ci] = make([]float64, nValues)
		}
	}

	for ri, rk := range rowOrder {
		for ci, ck := range colOrder {
			key := cellKey{row: rk, col: ck}
			if cellMap, ok := cellAggs[key]; ok {
				for vi := 0; vi < nValues; vi++ {
					if st, ok := cellMap[vi]; ok {
						cells[ri][ci][vi] = st.finalize(req.Values[vi].Agg)
					}
				}
			}
		}
	}

	rowTotals := make([][]float64, nRows)
	for ri, rk := range rowOrder {
		rowTotals[ri] = make([]float64, nValues)
		if rowMap, ok := rowAggs[rk]; ok {
			for vi := 0; vi < nValues; vi++ {
				if st, ok := rowMap[vi]; ok {
					rowTotals[ri][vi] = st.finalize(req.Values[vi].Agg)
				}
			}
		}
	}

	colTotals := make([][]float64, nCols)
	for ci, ck := range colOrder {
		colTotals[ci] = make([]float64, nValues)
		if colMap, ok := colAggs[ck]; ok {
			for vi := 0; vi < nValues; vi++ {
				if st, ok := colMap[vi]; ok {
					colTotals[ci][vi] = st.finalize(req.Values[vi].Agg)
				}
			}
		}
	}

	grandTotal := make([]float64, nValues)
	for vi := 0; vi < nValues; vi++ {
		if st, ok := grandAggs[vi]; ok {
			grandTotal[vi] = st.finalize(req.Values[vi].Agg)
		}
	}

	valueFields := make([]ValueField, len(req.Values))
	for i, vf := range req.Values {
		alias := vf.Alias
		if alias == "" {
			alias = fmt.Sprintf("%s_%s", vf.Agg, vf.Field)
		}
		valueFields[i] = ValueField{Field: vf.Field, Agg: vf.Agg, Alias: alias}
	}

	return &Result{
		RowHeaders:  rowOrder,
		ColHeaders:  colOrder,
		Cells:       cells,
		RowTotals:   rowTotals,
		ColTotals:   colTotals,
		GrandTotal:  grandTotal,
		ValueFields: valueFields,
	}
}

func dimKey(row map[string]any, dims []string) string {
	if len(dims) == 0 {
		return ""
	}
	parts := make([]string, len(dims))
	for i, d := range dims {
		val := row[d]
		parts[i] = fmt.Sprintf("%v", val)
	}
	return strings.Join(parts, " / ")
}

func toFloat(v any) (float64, bool) {
	if v == nil {
		return 0, false
	}
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case bool:
		if val {
			return 1, true
		}
		return 0, true
	case string:
		var f float64
		_, err := fmt.Sscanf(val, "%f", &f)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

func fieldExists(modelDef *parser.ModelDefinition, field string) bool {
	if modelDef == nil {
		return false
	}
	_, ok := modelDef.Fields[field]
	return ok
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func buildQuery(filter map[string]any) *persistence.Query {
	if len(filter) > 0 {
		return persistence.ParseQueryFromMap(filter)
	}
	return persistence.NewQuery()
}
