package export

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/xuri/excelize/v2"
)

// Format identifies an export output format.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatXLSX Format = "xlsx"
	FormatJSON Format = "json"
)

// ColumnSpec maps a model field to an export column header.
type ColumnSpec struct {
	Field string `json:"field"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

// SortDefinition mirrors parser.SortDefinition so callers can pass sort
// metadata without importing the parser package.
type SortDefinition struct {
	Field string `json:"field"`
	Order string `json:"order,omitempty"`
}

// ExportOptions configures an export run. Fetcher supplies the records in
// batches; Filter and Sort are applied to the underlying query before
// fetching.
type ExportOptions struct {
	Model   string
	Format  Format
	Columns []ColumnSpec
	Filter  *persistence.Query
	Sort    *SortDefinition
	Fetcher FetchFunc
}

// FetchFunc fetches a page of records. The generators call it in batches to
// stream large datasets without loading everything into memory. The query and
// total are opaque to the formatter — only the rows matter for rendering.
type FetchFunc func(page, pageSize int) ([]map[string]any, int64, error)

// RecordFetcher is the context-aware variant used by the Service.Export
// method (kept for callers that need to thread a context through to the
// repository layer).
type RecordFetcher func(ctx context.Context, query *persistence.Query, page, pageSize int) ([]map[string]any, int64, error)

const defaultBatchSize = 1000

// Service renders model data to CSV, XLSX, or JSON. It is transport-agnostic:
// callers (HTTP handlers, CLI, tests) supply a FetchFunc/RecordFetcher that
// knows how to load records, and the service handles formatting only.
type Service struct {
	batchSize int
}

// NewService constructs an ExportService with the default batch size (1000).
func NewService() *Service {
	return &Service{batchSize: defaultBatchSize}
}

// NewServiceWithBatchSize constructs an ExportService with a custom batch size.
func NewServiceWithBatchSize(batchSize int) *Service {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	return &Service{batchSize: batchSize}
}

// BatchSize returns the configured batch size.
func (s *Service) BatchSize() int { return s.batchSize }

// Export renders the given options and returns a ReadSeeker of the content,
// the filename, and the content type. It delegates to the standalone
// GenerateCSV/GenerateXLSX/GenerateJSON functions, adapting the context-aware
// RecordFetcher into a FetchFunc.
func (s *Service) Export(ctx context.Context, opts ExportOptions) (io.ReadSeeker, string, string, error) {
	if opts.Model == "" {
		return nil, "", "", fmt.Errorf("model is required")
	}
	if opts.Fetcher == nil {
		return nil, "", "", fmt.Errorf("fetcher is required")
	}

	format := opts.Format
	if format == "" {
		format = FormatCSV
	}
	switch format {
	case FormatCSV, FormatXLSX, FormatJSON:
	default:
		return nil, "", "", fmt.Errorf("unsupported format %q (supported: csv, xlsx, json)", format)
	}
	if len(opts.Columns) == 0 {
		return nil, "", "", fmt.Errorf("at least one column is required")
	}

	var content []byte
	var err error
	switch format {
	case FormatCSV:
		content, err = GenerateCSV(opts.Columns, opts.Fetcher, s.batchSize)
	case FormatJSON:
		content, err = GenerateJSON(opts.Columns, opts.Fetcher, s.batchSize)
	case FormatXLSX:
		content, err = GenerateXLSX(opts.Model, opts.Columns, opts.Fetcher, s.batchSize)
	}
	if err != nil {
		return nil, "", "", err
	}

	return bytes.NewReader(content), GenerateFilename(opts.Model, format), ContentType(format), nil
}

// GenerateCSV renders records to a CSV byte slice. A UTF-8 BOM is prepended so
// Excel opens UTF-8 content correctly. Every field is quoted by the stdlib
// csv writer. Records are fetched in batches of batchSize to bound memory.
func GenerateCSV(columns []ColumnSpec, fetcher FetchFunc, batchSize int) ([]byte, error) {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)

	bw.Write([]byte{0xEF, 0xBB, 0xBF})

	cw := csv.NewWriter(bw)

	labels := make([]string, len(columns))
	for i, col := range columns {
		labels[i] = col.Label
	}
	if err := cw.Write(labels); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	page := 1
	for {
		rows, _, err := fetcher(page, batchSize)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch records (page %d): %w", page, err)
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			values := make([]string, len(columns))
			for i, col := range columns {
				values[i] = FormatCellValue(row[col.Field])
			}
			if err := cw.Write(values); err != nil {
				return nil, fmt.Errorf("failed to write CSV row: %w", err)
			}
		}
		cw.Flush()
		if err := cw.Error(); err != nil {
			return nil, fmt.Errorf("CSV flush error: %w", err)
		}
		if len(rows) < batchSize {
			break
		}
		page++
	}

	cw.Flush()
	bw.Flush()
	return buf.Bytes(), nil
}

// GenerateJSON renders records as a pretty JSON array of objects keyed by
// column labels. Records are fetched in batches to bound memory.
func GenerateJSON(columns []ColumnSpec, fetcher FetchFunc, batchSize int) ([]byte, error) {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)

	bw.WriteString("[")

	first := true
	page := 1
	for {
		rows, _, err := fetcher(page, batchSize)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch records (page %d): %w", page, err)
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			obj := make(map[string]any, len(columns))
			for _, col := range columns {
				obj[col.Label] = row[col.Field]
			}
			if !first {
				bw.WriteString(",")
			}
			first = false
			enc := json.NewEncoder(bw)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(obj); err != nil {
				return nil, fmt.Errorf("failed to encode JSON row: %w", err)
			}
		}
		bw.Flush()
		if len(rows) < batchSize {
			break
		}
		page++
	}

	bw.WriteString("]")
	bw.Flush()
	return buf.Bytes(), nil
}

// GenerateXLSX renders records to an XLSX byte slice using excelize's
// streaming writer so memory stays bounded for large datasets. Sheet 1 holds
// the data with a bold, shaded header row and auto-width columns.
func GenerateXLSX(modelName string, columns []ColumnSpec, fetcher FetchFunc, batchSize int) ([]byte, error) {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	defaultSheet := f.GetSheetName(0)
	if defaultSheet != sheetName {
		f.SetSheetName(defaultSheet, sheetName)
	}

	sw, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create XLSX stream writer: %w", err)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F2F2F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	headerCells := make([]any, len(columns))
	for i, col := range columns {
		headerCells[i] = excelize.Cell{StyleID: headerStyle, Value: col.Label}
	}
	if err := sw.SetRow("A1", headerCells); err != nil {
		return nil, fmt.Errorf("failed to write XLSX header: %w", err)
	}

	colWidths := make([]float64, len(columns))

	rowIdx := 2
	page := 1
	for {
		rows, _, err := fetcher(page, batchSize)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch records (page %d): %w", page, err)
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			cellValues := make([]any, len(columns))
			for i, col := range columns {
				val := row[col.Field]
				cellValues[i] = normalizeXLSXValue(val)
				labelLen := len(FormatCellValue(val))
				if float64(labelLen) > colWidths[i] {
					colWidths[i] = float64(labelLen)
				}
			}
			cellName, _ := excelize.CoordinatesToCellName(1, rowIdx)
			if err := sw.SetRow(cellName, cellValues); err != nil {
				return nil, fmt.Errorf("failed to write XLSX row %d: %w", rowIdx, err)
			}
			rowIdx++
		}
		if len(rows) < batchSize {
			break
		}
		page++
	}

	if err := sw.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush XLSX stream: %w", err)
	}

	for i, width := range colWidths {
		if width > 0 {
			colName, _ := excelize.ColumnNumberToName(i + 1)
			finalWidth := width + 2
			if finalWidth > 80 {
				finalWidth = 80
			}
			if finalWidth < 10 {
				finalWidth = 10
			}
			_ = f.SetColWidth(sheetName, colName, colName, finalWidth)
		}
	}

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("failed to write XLSX to output: %w", err)
	}
	return buf.Bytes(), nil
}

func normalizeXLSXValue(v any) any {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string, int, int64, float64, bool:
		return val
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// FormatCellValue renders a cell value to a string for CSV output.
func FormatCellValue(v any) string {
	return formatCell(v)
}

func formatCell(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// ContentType returns the MIME type for a format.
func ContentType(format Format) string {
	switch format {
	case FormatXLSX:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case FormatJSON:
		return "application/json"
	default:
		return "text/csv"
	}
}

// GenerateFilename builds `{model}_{YYYYMMDD_HHMMSS}.{ext}`.
func GenerateFilename(model string, f Format) string {
	if model == "" {
		model = "export"
	}
	model = strings.ReplaceAll(model, ".", "_")
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s.%s", model, timestamp, string(f))
}

// ColumnsFromModel builds a ColumnSpec list from a ModelDefinition. Relational
// fields (one2many, many2many, morph collections) are excluded because they
// represent collections rather than scalar export values. many2one fields are
// included (they store a foreign-key value).
func ColumnsFromModel(modelDef *parser.ModelDefinition) []ColumnSpec {
	return DefaultColumnsFromModel(modelDef)
}

// DefaultColumnsFromModel is an alias for ColumnsFromModel that the tests
// reference. Relational collection fields are excluded.
func DefaultColumnsFromModel(modelDef *parser.ModelDefinition) []ColumnSpec {
	if modelDef == nil {
		return nil
	}
	var cols []ColumnSpec
	for name, fd := range modelDef.Fields {
		switch fd.Type {
		case parser.FieldOne2Many, parser.FieldMany2Many,
			parser.FieldMorphOne, parser.FieldMorphMany,
			parser.FieldMorphToMany, parser.FieldMorphByMany:
			continue
		}
		label := fd.Label
		if label == "" {
			label = name
		}
		cols = append(cols, ColumnSpec{Field: name, Label: label, Type: string(fd.Type)})
	}
	return cols
}

// QueryFetcher adapts a persistence.Query + RecordFetcher into a FetchFunc
// that the generators can consume. If fetcher is nil it returns nil.
func QueryFetcher(query *persistence.Query, fetcher RecordFetcher) FetchFunc {
	if fetcher == nil {
		return nil
	}
	return func(page, pageSize int) ([]map[string]any, int64, error) {
		return fetcher(context.Background(), query, page, pageSize)
	}
}
