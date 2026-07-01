package pivot

import (
	"context"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
)

func sampleModel() *parser.ModelDefinition {
	return &parser.ModelDefinition{
		Name: "sale",
		Fields: map[string]parser.FieldDefinition{
			"region":  {Type: parser.FieldString, Label: "Region"},
			"product": {Type: parser.FieldString, Label: "Product"},
			"channel": {Type: parser.FieldString, Label: "Channel"},
			"amount":  {Type: parser.FieldDecimal, Label: "Amount"},
			"qty":     {Type: parser.FieldInteger, Label: "Quantity"},
		},
	}
}

func salesData() []map[string]any {
	return []map[string]any{
		{"region": "North", "product": "Widget", "channel": "Online", "amount": float64(100), "qty": int64(2)},
		{"region": "North", "product": "Widget", "channel": "Retail", "amount": float64(50), "qty": int64(1)},
		{"region": "North", "product": "Gadget", "channel": "Online", "amount": float64(200), "qty": int64(4)},
		{"region": "South", "product": "Widget", "channel": "Online", "amount": float64(300), "qty": int64(6)},
		{"region": "South", "product": "Gadget", "channel": "Retail", "amount": float64(150), "qty": int64(3)},
	}
}

func TestPivotSingleAgg(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}
	result := svc.Aggregate(req, salesData())

	if len(result.RowHeaders) != 2 {
		t.Fatalf("expected 2 row headers, got %d: %v", len(result.RowHeaders), result.RowHeaders)
	}
	if result.RowHeaders[0] != "North" || result.RowHeaders[1] != "South" {
		t.Errorf("unexpected row headers: %v", result.RowHeaders)
	}
	if len(result.ColHeaders) != 2 {
		t.Fatalf("expected 2 col headers, got %d: %v", len(result.ColHeaders), result.ColHeaders)
	}
	if result.ColHeaders[0] != "Gadget" || result.ColHeaders[1] != "Widget" {
		t.Errorf("unexpected col headers: %v", result.ColHeaders)
	}

	northIdx := indexOf(result.RowHeaders, "North")
	southIdx := indexOf(result.RowHeaders, "South")
	widgetIdx := indexOf(result.ColHeaders, "Widget")
	gadgetIdx := indexOf(result.ColHeaders, "Gadget")

	if result.Cells[northIdx][widgetIdx][0] != 150 {
		t.Errorf("North/Widget sum: expected 150, got %v", result.Cells[northIdx][widgetIdx][0])
	}
	if result.Cells[northIdx][gadgetIdx][0] != 200 {
		t.Errorf("North/Gadget sum: expected 200, got %v", result.Cells[northIdx][gadgetIdx][0])
	}
	if result.Cells[southIdx][widgetIdx][0] != 300 {
		t.Errorf("South/Widget sum: expected 300, got %v", result.Cells[southIdx][widgetIdx][0])
	}
	if result.Cells[southIdx][gadgetIdx][0] != 150 {
		t.Errorf("South/Gadget sum: expected 150, got %v", result.Cells[southIdx][gadgetIdx][0])
	}
}

func TestPivotMultiAgg(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows: []string{"region"},
		Cols: []string{"product"},
		Values: []ValueField{
			{Field: "amount", Agg: AggSum},
			{Field: "amount", Agg: AggAvg},
			{Field: "qty", Agg: AggCount},
		},
	}
	result := svc.Aggregate(req, salesData())

	northIdx := indexOf(result.RowHeaders, "North")
	widgetIdx := indexOf(result.ColHeaders, "Widget")

	cell := result.Cells[northIdx][widgetIdx]
	if cell[0] != 150 {
		t.Errorf("sum: expected 150, got %v", cell[0])
	}
	if cell[1] != 75 {
		t.Errorf("avg: expected 75, got %v", cell[1])
	}
	if cell[2] != 2 {
		t.Errorf("count: expected 2, got %v", cell[2])
	}
}

func TestPivotMultiRows(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region", "product"},
		Cols:   []string{"channel"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}
	result := svc.Aggregate(req, salesData())

	if len(result.RowHeaders) != 4 {
		t.Fatalf("expected 4 row headers for 2x2 dims, got %d: %v", len(result.RowHeaders), result.RowHeaders)
	}

	northGadgetOnline := indexOf(result.RowHeaders, "North / Gadget")
	onlineIdx := indexOf(result.ColHeaders, "Online")
	if result.Cells[northGadgetOnline][onlineIdx][0] != 200 {
		t.Errorf("North/Gadget/Online: expected 200, got %v", result.Cells[northGadgetOnline][onlineIdx][0])
	}
}

func TestPivotMultiCols(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product", "channel"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}
	result := svc.Aggregate(req, salesData())

	if len(result.ColHeaders) != 4 {
		t.Fatalf("expected 4 col headers for 2x2 dims, got %d: %v", len(result.ColHeaders), result.ColHeaders)
	}

	northIdx := indexOf(result.RowHeaders, "North")
	northWidgetOnline := indexOf(result.ColHeaders, "Widget / Online")
	if result.Cells[northIdx][northWidgetOnline][0] != 100 {
		t.Errorf("North/Widget/Online: expected 100, got %v", result.Cells[northIdx][northWidgetOnline][0])
	}
}

func TestPivotNullHandling(t *testing.T) {
	data := []map[string]any{
		{"region": "North", "product": "Widget", "amount": float64(100)},
		{"region": "North", "product": "Widget", "amount": nil},
		{"region": "North", "product": "Widget", "amount": float64(200)},
		{"region": "South", "product": "Widget", "amount": nil},
	}

	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}
	result := svc.Aggregate(req, data)

	northIdx := indexOf(result.RowHeaders, "North")
	southIdx := indexOf(result.RowHeaders, "South")
	widgetIdx := indexOf(result.ColHeaders, "Widget")

	if result.Cells[northIdx][widgetIdx][0] != 300 {
		t.Errorf("North/Widget sum (with null): expected 300, got %v", result.Cells[northIdx][widgetIdx][0])
	}
	if result.Cells[southIdx][widgetIdx][0] != 0 {
		t.Errorf("South/Widget sum (all null): expected 0, got %v", result.Cells[southIdx][widgetIdx][0])
	}
}

func TestPivotNullAvgHandling(t *testing.T) {
	data := []map[string]any{
		{"region": "North", "amount": float64(100)},
		{"region": "North", "amount": nil},
		{"region": "North", "amount": float64(300)},
	}

	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Values: []ValueField{{Field: "amount", Agg: AggAvg}},
	}
	result := svc.Aggregate(req, data)

	northIdx := indexOf(result.RowHeaders, "North")
	grand := result.RowTotals[northIdx][0]
	if grand != 200 {
		t.Errorf("North avg (null excluded): expected 200, got %v", grand)
	}
}

func TestPivotCountWithNulls(t *testing.T) {
	data := []map[string]any{
		{"region": "North", "amount": float64(100)},
		{"region": "North", "amount": nil},
		{"region": "North", "amount": float64(300)},
	}

	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Values: []ValueField{
			{Field: "amount", Agg: AggCount},
		},
	}
	result := svc.Aggregate(req, data)

	northIdx := indexOf(result.RowHeaders, "North")
	countVal := result.RowTotals[northIdx][0]
	if countVal != 3 {
		t.Errorf("count(*): expected 3 (all rows), got %v", countVal)
	}
}

func TestPivotTotals(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}
	result := svc.Aggregate(req, salesData())

	northIdx := indexOf(result.RowHeaders, "North")
	southIdx := indexOf(result.RowHeaders, "South")
	widgetIdx := indexOf(result.ColHeaders, "Widget")

	northTotal := result.RowTotals[northIdx][0]
	if northTotal != 350 {
		t.Errorf("North row total: expected 350, got %v", northTotal)
	}
	southTotal := result.RowTotals[southIdx][0]
	if southTotal != 450 {
		t.Errorf("South row total: expected 450, got %v", southTotal)
	}

	widgetColTotal := result.ColTotals[widgetIdx][0]
	if widgetColTotal != 450 {
		t.Errorf("Widget col total: expected 450, got %v", widgetColTotal)
	}

	if result.GrandTotal[0] != 800 {
		t.Errorf("Grand total: expected 800, got %v", result.GrandTotal[0])
	}
}

func TestPivotAvgTotalCorrectness(t *testing.T) {
	data := []map[string]any{
		{"region": "North", "product": "Widget", "amount": float64(100)},
		{"region": "North", "product": "Widget", "amount": float64(200)},
		{"region": "North", "product": "Gadget", "amount": float64(300)},
		{"region": "North", "product": "Gadget", "amount": float64(500)},
	}

	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product"},
		Values: []ValueField{{Field: "amount", Agg: AggAvg}},
	}
	result := svc.Aggregate(req, data)

	northIdx := indexOf(result.RowHeaders, "North")
	rowTotal := result.RowTotals[northIdx][0]
	if rowTotal != 275 {
		t.Errorf("North row avg total: expected 275 (sum=1100/count=4), got %v", rowTotal)
	}

	if result.GrandTotal[0] != 275 {
		t.Errorf("Grand avg: expected 275, got %v", result.GrandTotal[0])
	}
}

func TestPivotEmptyData(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}
	result := svc.Aggregate(req, []map[string]any{})

	if len(result.RowHeaders) != 0 {
		t.Errorf("expected 0 row headers for empty data, got %d", len(result.RowHeaders))
	}
	if result.GrandTotal[0] != 0 {
		t.Errorf("grand total for empty data: expected 0, got %v", result.GrandTotal[0])
	}
}

func TestPivotNoColsDimension(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}
	result := svc.Aggregate(req, salesData())

	if len(result.ColHeaders) != 1 || result.ColHeaders[0] != "" {
		t.Fatalf("expected single empty col header, got %v", result.ColHeaders)
	}

	northIdx := indexOf(result.RowHeaders, "North")
	if result.Cells[northIdx][0][0] != 350 {
		t.Errorf("North total (no cols): expected 350, got %v", result.Cells[northIdx][0][0])
	}
}

func TestPivotMinMaxAgg(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product"},
		Values: []ValueField{
			{Field: "amount", Agg: AggMin},
			{Field: "amount", Agg: AggMax},
		},
	}
	result := svc.Aggregate(req, salesData())

	northIdx := indexOf(result.RowHeaders, "North")
	widgetIdx := indexOf(result.ColHeaders, "Widget")

	if result.Cells[northIdx][widgetIdx][0] != 50 {
		t.Errorf("North/Widget min: expected 50, got %v", result.Cells[northIdx][widgetIdx][0])
	}
	if result.Cells[northIdx][widgetIdx][1] != 100 {
		t.Errorf("North/Widget max: expected 100, got %v", result.Cells[northIdx][widgetIdx][1])
	}
}

func TestPivotValidate(t *testing.T) {
	svc := NewService()
	model := sampleModel()

	tests := []struct {
		name    string
		req     *Request
		wantErr bool
	}{
		{
			name: "valid request",
			req: &Request{
				Rows:   []string{"region"},
				Cols:   []string{"product"},
				Values: []ValueField{{Field: "amount", Agg: AggSum}},
			},
			wantErr: false,
		},
		{
			name: "missing rows",
			req: &Request{
				Values: []ValueField{{Field: "amount", Agg: AggSum}},
			},
			wantErr: true,
		},
		{
			name: "missing values",
			req: &Request{
				Rows: []string{"region"},
			},
			wantErr: true,
		},
		{
			name: "invalid row field",
			req: &Request{
				Rows:   []string{"nonexistent"},
				Values: []ValueField{{Field: "amount", Agg: AggSum}},
			},
			wantErr: true,
		},
		{
			name: "invalid agg func",
			req: &Request{
				Rows:   []string{"region"},
				Values: []ValueField{{Field: "amount", Agg: AggFunc("median")}},
			},
			wantErr: true,
		},
		{
			name: "invalid value field",
			req: &Request{
				Rows:   []string{"region"},
				Values: []ValueField{{Field: "nonexistent", Agg: AggSum}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Validate(tt.req, model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestPivotComputeWithFetcher(t *testing.T) {
	svc := NewService()
	req := &Request{
		Rows:   []string{"region"},
		Cols:   []string{"product"},
		Values: []ValueField{{Field: "amount", Agg: AggSum}},
	}

	fetcher := func(ctx context.Context, query *persistence.Query) ([]map[string]any, error) {
		return salesData(), nil
	}

	result, err := svc.Compute(context.Background(), req, sampleModel(), fetcher)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}
	if result.GrandTotal[0] != 800 {
		t.Errorf("Grand total: expected 800, got %v", result.GrandTotal[0])
	}
}

func indexOf(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}
