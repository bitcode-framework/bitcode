package export

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/xuri/excelize/v2"
)

func sampleColumns() []ColumnSpec {
	return []ColumnSpec{
		{Field: "name", Label: "Name", Type: "string"},
		{Field: "email", Label: "Email", Type: "email"},
		{Field: "age", Label: "Age", Type: "integer"},
	}
}

func sampleFetcher() FetchFunc {
	return func(page, pageSize int) ([]map[string]any, int64, error) {
		if page > 1 {
			return nil, 0, nil
		}
		return []map[string]any{
			{"name": "Alice", "email": "alice@test.com", "age": int64(30)},
			{"name": "Bob", "email": "bob@test.com", "age": int64(25)},
			{"name": "Charlie", "email": "charlie@test.com", "age": int64(35)},
		}, 3, nil
	}
}

func sampleModelDef() *parser.ModelDefinition {
	return &parser.ModelDefinition{
		Name: "contact",
		Fields: map[string]parser.FieldDefinition{
			"name":  {Type: parser.FieldString, Label: "Name"},
			"email": {Type: parser.FieldEmail, Label: "Email"},
			"age":   {Type: parser.FieldInteger, Label: "Age"},
		},
	}
}

func TestGenerateCSV(t *testing.T) {
	cols := sampleColumns()
	content, err := GenerateCSV(cols, sampleFetcher(), 1000)
	if err != nil {
		t.Fatalf("GenerateCSV failed: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected non-empty CSV content")
	}

	// Strip BOM
	body := content
	if bytes.HasPrefix(body, []byte{0xEF, 0xBB, 0xBF}) {
		body = body[3:]
	}

	reader := csv.NewReader(bytes.NewReader(body))
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("csv parse failed: %v", err)
	}
	if len(rows) != 4 {
		t.Fatalf("expected 4 rows (1 header + 3 data), got %d", len(rows))
	}
	if rows[0][0] != "Name" || rows[0][1] != "Email" || rows[0][2] != "Age" {
		t.Errorf("unexpected header: %v", rows[0])
	}
	if rows[1][0] != "Alice" {
		t.Errorf("expected Alice in row 1, got %s", rows[1][0])
	}
	if rows[2][0] != "Bob" {
		t.Errorf("expected Bob in row 2, got %s", rows[2][0])
	}
}

func TestGenerateCSV_BOM(t *testing.T) {
	content, err := GenerateCSV(sampleColumns(), sampleFetcher(), 1000)
	if err != nil {
		t.Fatalf("GenerateCSV failed: %v", err)
	}
	if !bytes.HasPrefix(content, []byte{0xEF, 0xBB, 0xBF}) {
		t.Error("expected UTF-8 BOM at start of CSV for Excel compatibility")
	}
}

func TestGenerateCSV_Quoting(t *testing.T) {
	cols := []ColumnSpec{{Field: "note", Label: "Note", Type: "string"}}
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		if page > 1 {
			return nil, 0, nil
		}
		return []map[string]any{
			{"note": "hello, world"},
			{"note": `has "quotes"`},
		}, 2, nil
	}
	content, err := GenerateCSV(cols, fetcher, 1000)
	if err != nil {
		t.Fatalf("GenerateCSV failed: %v", err)
	}
	body := content
	if bytes.HasPrefix(body, []byte{0xEF, 0xBB, 0xBF}) {
		body = body[3:]
	}
	reader := csv.NewReader(bytes.NewReader(body))
	rows, _ := reader.ReadAll()
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[1][0] != "hello, world" {
		t.Errorf("expected 'hello, world', got %q", rows[1][0])
	}
	if rows[2][0] != `has "quotes"` {
		t.Errorf("expected quoted value, got %q", rows[2][0])
	}
}

func TestGenerateCSV_Empty(t *testing.T) {
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		return nil, 0, nil
	}
	content, err := GenerateCSV(sampleColumns(), fetcher, 1000)
	if err != nil {
		t.Fatalf("GenerateCSV failed: %v", err)
	}
	body := content
	if bytes.HasPrefix(body, []byte{0xEF, 0xBB, 0xBF}) {
		body = body[3:]
	}
	rows, _ := csv.NewReader(bytes.NewReader(body)).ReadAll()
	if len(rows) != 1 {
		t.Errorf("expected 1 header row for empty data, got %d", len(rows))
	}
}

func TestGenerateXLSX(t *testing.T) {
	cols := sampleColumns()
	content, err := GenerateXLSX("contact", cols, sampleFetcher(), 1000)
	if err != nil {
		t.Fatalf("GenerateXLSX failed: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected non-empty XLSX content")
	}

	reader := bytes.NewReader(content)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		t.Fatalf("failed to open generated xlsx: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		t.Fatal("expected at least one sheet")
	}
	if sheets[0] != "Sheet1" {
		t.Errorf("expected Sheet1, got %s", sheets[0])
	}

	header, err := f.GetCellValue("Sheet1", "A1")
	if err != nil {
		t.Fatalf("failed to read A1: %v", err)
	}
	if header != "Name" {
		t.Errorf("expected header 'Name' in A1, got %q", header)
	}

	emailHeader, _ := f.GetCellValue("Sheet1", "B1")
	if emailHeader != "Email" {
		t.Errorf("expected header 'Email' in B1, got %q", emailHeader)
	}

	dataCell, _ := f.GetCellValue("Sheet1", "A2")
	if dataCell != "Alice" {
		t.Errorf("expected 'Alice' in A2, got %q", dataCell)
	}
}

func TestGenerateXLSX_Empty(t *testing.T) {
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		return nil, 0, nil
	}
	content, err := GenerateXLSX("contact", sampleColumns(), fetcher, 1000)
	if err != nil {
		t.Fatalf("GenerateXLSX failed: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected non-empty XLSX even for empty data (header row)")
	}
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	defer f.Close()
	header, _ := f.GetCellValue("Sheet1", "A1")
	if header != "Name" {
		t.Errorf("expected header row even for empty data, got A1=%q", header)
	}
}

func TestGenerateJSON(t *testing.T) {
	cols := sampleColumns()
	content, err := GenerateJSON(cols, sampleFetcher(), 1000)
	if err != nil {
		t.Fatalf("GenerateJSON failed: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected non-empty JSON content")
	}

	trimmed := bytes.TrimSpace(content)
	if !bytes.HasPrefix(trimmed, []byte("[")) {
		t.Errorf("expected JSON to start with '[', got: %s", string(trimmed[:min(20, len(trimmed))]))
	}
	if !bytes.HasSuffix(trimmed, []byte("]")) {
		t.Errorf("expected JSON to end with ']'")
	}

	var result []map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("invalid JSON: %v\ncontent: %s", err, string(content))
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 records, got %d", len(result))
	}
	if result[0]["Name"] != "Alice" {
		t.Errorf("expected Name=Alice, got %v", result[0]["Name"])
	}
	if result[1]["Email"] != "bob@test.com" {
		t.Errorf("expected Email=bob@test.com, got %v", result[1]["Email"])
	}
}

func TestGenerateJSON_Keys(t *testing.T) {
	cols := []ColumnSpec{
		{Field: "name", Label: "Full Name", Type: "string"},
		{Field: "phone", Label: "Phone Number", Type: "string"},
	}
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		if page > 1 {
			return nil, 0, nil
		}
		return []map[string]any{
			{"name": "Alice", "phone": "123-456"},
		}, 1, nil
	}
	content, err := GenerateJSON(cols, fetcher, 1000)
	if err != nil {
		t.Fatalf("GenerateJSON failed: %v", err)
	}
	var result []map[string]any
	json.Unmarshal(content, &result)
	if len(result) != 1 {
		t.Fatalf("expected 1 record, got %d", len(result))
	}
	if _, ok := result[0]["Full Name"]; !ok {
		t.Errorf("expected key 'Full Name', got keys: %v", mapKeys(result[0]))
	}
	if _, ok := result[0]["Phone Number"]; !ok {
		t.Errorf("expected key 'Phone Number', got keys: %v", mapKeys(result[0]))
	}
}

func TestGenerateJSON_Empty(t *testing.T) {
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		return nil, 0, nil
	}
	content, err := GenerateJSON(sampleColumns(), fetcher, 1000)
	if err != nil {
		t.Fatalf("GenerateJSON failed: %v", err)
	}
	var result []map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("invalid JSON for empty: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty array, got %d", len(result))
	}
}

func TestDefaultColumnsFromModel(t *testing.T) {
	md := &parser.ModelDefinition{
		Name: "contact",
		Fields: map[string]parser.FieldDefinition{
			"name":    {Type: parser.FieldString, Label: "Name"},
			"lines":   {Type: parser.FieldOne2Many, Label: "Lines"},
			"tags":    {Type: parser.FieldMany2Many, Label: "Tags"},
			"image":   {Type: parser.FieldMorphMany, Label: "Images"},
			"related": {Type: parser.FieldMorphToMany, Label: "Related"},
			"email":   {Type: parser.FieldEmail, Label: ""},
		},
	}
	cols := DefaultColumnsFromModel(md)
	if len(cols) != 2 {
		t.Fatalf("expected 2 exportable columns (string+email), got %d", len(cols))
	}
	found := map[string]bool{}
	for _, c := range cols {
		found[c.Field] = true
	}
	if !found["name"] {
		t.Error("expected 'name' column")
	}
	if !found["email"] {
		t.Error("expected 'email' column")
	}
	for _, c := range cols {
		if c.Field == "email" && c.Label != "email" {
			t.Errorf("expected label to default to field name 'email', got %q", c.Label)
		}
		if c.Field == "name" && c.Label != "Name" {
			t.Errorf("expected label 'Name', got %q", c.Label)
		}
	}
}

func TestGenerateFilename(t *testing.T) {
	filename := GenerateFilename("contact", FormatCSV)
	if !strings.HasSuffix(filename, ".csv") {
		t.Errorf("expected .csv extension, got %s", filename)
	}
	if !strings.HasPrefix(filename, "contact_") {
		t.Errorf("expected 'contact_' prefix, got %s", filename)
	}

	filename2 := GenerateFilename("crm.contact", FormatXLSX)
	if !strings.HasSuffix(filename2, ".xlsx") {
		t.Errorf("expected .xlsx extension, got %s", filename2)
	}
	if strings.Contains(filename2, ".") && !strings.HasSuffix(filename2, ".xlsx") {
		t.Errorf("model dot should be replaced, got %s", filename2)
	}

	filename3 := GenerateFilename("order", FormatJSON)
	if !strings.HasSuffix(filename3, ".json") {
		t.Errorf("expected .json extension, got %s", filename3)
	}
}

func TestGenerateFilename_TimestampFormat(t *testing.T) {
	filename := GenerateFilename("contact", FormatCSV)
	// Should be contact_YYYYMMDD_HHMMSS.csv
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 {
		t.Fatalf("expected split by first underscore, got %s", filename)
	}
	if parts[0] != "contact" {
		t.Errorf("expected 'contact' prefix, got %q", parts[0])
	}
	ts := strings.TrimSuffix(parts[1], ".csv")
	if len(ts) != len("20060102_150405") {
		t.Errorf("expected timestamp length %d, got %d (%s)", len("20060102_150405"), len(ts), ts)
	}
	if _, err := time.Parse("20060102_150405", ts); err != nil {
		t.Errorf("invalid timestamp format %q: %v", ts, err)
	}
}

func TestFormatCell(t *testing.T) {
	if formatCell(nil) != "" {
		t.Error("nil should format as empty string")
	}
	if formatCell("hello") != "hello" {
		t.Error("string should pass through")
	}
	if formatCell(true) != "true" {
		t.Error("true should be 'true'")
	}
	if formatCell(false) != "false" {
		t.Error("false should be 'false'")
	}
	if formatCell(int64(42)) != "42" {
		t.Error("int64 42 should be '42'")
	}
	if formatCell(3.14) != "3.14" {
		t.Error("float 3.14 should be '3.14'")
	}
	ts := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	result := formatCell(ts)
	if !strings.Contains(result, "2026-01-15T10:30:00") {
		t.Errorf("time should format as RFC3339, got %s", result)
	}
}

func TestExport(t *testing.T) {
	svc := NewService()
	opts := ExportOptions{
		Model:   "contact",
		Format:  FormatCSV,
		Columns: sampleColumns(),
		Fetcher: sampleFetcher(),
	}
	reader, filename, contentType, err := svc.Export(nil, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if reader == nil {
		t.Fatal("expected non-nil reader")
	}
	if !strings.HasSuffix(filename, ".csv") {
		t.Errorf("expected .csv filename, got %s", filename)
	}
	if contentType != "text/csv" {
		t.Errorf("expected text/csv, got %s", contentType)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	svc := NewService()
	opts := ExportOptions{
		Model:   "contact",
		Format:  "xml",
		Columns: sampleColumns(),
		Fetcher: sampleFetcher(),
	}
	_, _, _, err := svc.Export(nil, opts)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected unsupported format error, got %v", err)
	}
}

func TestExport_JSON(t *testing.T) {
	svc := NewService()
	opts := ExportOptions{
		Model:   "contact",
		Format:  FormatJSON,
		Columns: sampleColumns(),
		Fetcher: sampleFetcher(),
	}
	reader, _, contentType, err := svc.Export(nil, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if contentType != "application/json" {
		t.Errorf("expected application/json, got %s", contentType)
	}
	buf := new(bytes.Buffer)
	io_Copy(buf, reader)
	var result []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 records, got %d", len(result))
	}
}

func TestExport_XLSX(t *testing.T) {
	svc := NewService()
	opts := ExportOptions{
		Model:   "contact",
		Format:  FormatXLSX,
		Columns: sampleColumns(),
		Fetcher: sampleFetcher(),
	}
	reader, _, contentType, err := svc.Export(nil, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Errorf("expected xlsx content type, got %s", contentType)
	}
	buf := new(bytes.Buffer)
	io_Copy(buf, reader)
	if buf.Len() == 0 {
		t.Error("expected non-empty xlsx body")
	}
	f, err := excelize.OpenReader(buf)
	if err != nil {
		t.Fatalf("failed to open xlsx: %v", err)
	}
	defer f.Close()
	cell, _ := f.GetCellValue("Sheet1", "A1")
	if cell != "Name" {
		t.Errorf("expected A1=Name, got %q", cell)
	}
}

func TestGenerateCSV_MultiPage(t *testing.T) {
	calls := 0
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		calls++
		if page == 1 {
			rows := make([]map[string]any, 2)
			rows[0] = map[string]any{"name": "A"}
			rows[1] = map[string]any{"name": "B"}
			return rows, 3, nil
		}
		if page == 2 {
			return []map[string]any{{"name": "C"}}, 3, nil
		}
		return nil, 0, nil
	}
	cols := []ColumnSpec{{Field: "name", Label: "Name", Type: "string"}}
	content, err := GenerateCSV(cols, fetcher, 2)
	if err != nil {
		t.Fatalf("GenerateCSV failed: %v", err)
	}
	body := content
	if bytes.HasPrefix(body, []byte{0xEF, 0xBB, 0xBF}) {
		body = body[3:]
	}
	rows, _ := csv.NewReader(bytes.NewReader(body)).ReadAll()
	if len(rows) != 4 {
		t.Errorf("expected 4 rows (1 header + 3 data), got %d", len(rows))
	}
	if calls < 2 {
		t.Errorf("expected at least 2 fetch calls for pagination, got %d", calls)
	}
}

func TestQueryFetcher_Integration(t *testing.T) {
	q := persistence.NewQuery().WhereEq("status", "active").Order("name", "asc")
	called := false
	rf := func(ctx context.Context, query *persistence.Query, page, pageSize int) ([]map[string]any, int64, error) {
		called = true
		if query != q {
			t.Error("expected query to be passed through")
		}
		return nil, 0, nil
	}
	ff := QueryFetcher(q, rf)
	if ff == nil {
		t.Fatal("expected non-nil fetch function")
	}
	ff(1, 1000)
	if !called {
		t.Error("expected underlying RecordFetcher to be called")
	}
}

func TestQueryFetcher_NilFetcher(t *testing.T) {
	q := persistence.NewQuery()
	if ff := QueryFetcher(q, nil); ff != nil {
		t.Error("expected nil FetchFunc when RecordFetcher is nil")
	}
}

func TestGenerateCSV_NilValue(t *testing.T) {
	cols := []ColumnSpec{{Field: "name", Label: "Name", Type: "string"}}
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		if page > 1 {
			return nil, 0, nil
		}
		return []map[string]any{
			{"name": nil},
			{"name": "Value"},
		}, 2, nil
	}
	content, err := GenerateCSV(cols, fetcher, 1000)
	if err != nil {
		t.Fatalf("GenerateCSV failed: %v", err)
	}
	body := content
	if bytes.HasPrefix(body, []byte{0xEF, 0xBB, 0xBF}) {
		body = body[3:]
	}
	// The nil value becomes an empty string field, which the CSV writer
	// emits as a blank line. Go's csv.Reader skips blank lines by default,
	// so we configure FieldsPerRecord=-1 and read line-by-line to verify
	// the structure rather than relying on ReadAll's skip behavior.
	reader := csv.NewReader(bytes.NewReader(body))
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("csv parse failed: %v", err)
	}
	// csv.ReadAll still drops truly blank lines (single empty field), so
	// verify via raw line count instead.
	lineCount := strings.Count(string(body), "\n")
	if lineCount != 3 {
		t.Errorf("expected 3 lines (header + nil row + value row), got %d", lineCount)
	}
	if len(rows) < 2 {
		t.Fatalf("expected at least 2 parsed rows, got %d", len(rows))
	}
	if rows[0][0] != "Name" {
		t.Errorf("expected header 'Name', got %q", rows[0][0])
	}
	lastRow := rows[len(rows)-1]
	if len(lastRow) > 0 && lastRow[0] != "Value" {
		t.Errorf("expected last data row 'Value', got %q", lastRow[0])
	}
}

// --- helpers ---

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func io_Copy(dst *bytes.Buffer, src interface{ Read([]byte) (int, error) }) {
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			dst.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

// suppress unused import warnings when test shape changes
var _ = bufio.NewWriter
var _ = bytes.NewBuffer
var _ = time.Now
