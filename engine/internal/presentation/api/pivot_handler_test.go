package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pivotModelRegistry struct {
	models map[string]*parser.ModelDefinition
}

func newPivotModelRegistry() *pivotModelRegistry {
	return &pivotModelRegistry{
		models: map[string]*parser.ModelDefinition{
			"sales.sale": {
				Name:   "sale",
				Module: "sales",
				Fields: map[string]parser.FieldDefinition{
					"region":  {Type: parser.FieldString, Label: "Region"},
					"product": {Type: parser.FieldString, Label: "Product"},
					"amount":  {Type: parser.FieldDecimal, Label: "Amount"},
					"qty":     {Type: parser.FieldInteger, Label: "Quantity"},
					"salary":  {Type: parser.FieldDecimal, Label: "Salary", Groups: []string{"hr.manager"}},
				},
			},
			"sale": {
				Name:   "sale",
				Module: "sales",
				Fields: map[string]parser.FieldDefinition{
					"region":  {Type: parser.FieldString, Label: "Region"},
					"product": {Type: parser.FieldString, Label: "Product"},
					"amount":  {Type: parser.FieldDecimal, Label: "Amount"},
					"qty":     {Type: parser.FieldInteger, Label: "Quantity"},
					"salary":  {Type: parser.FieldDecimal, Label: "Salary", Groups: []string{"hr.manager"}},
				},
			},
		},
	}
}

func (r *pivotModelRegistry) Get(name string) (*parser.ModelDefinition, error) {
	if m, ok := r.models[name]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("model %q not found", name)
}

func (r *pivotModelRegistry) TableName(name string) string {
	if name == "sale" || name == "sales.sale" {
		return "sales"
	}
	return name
}

func setupPivotTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.Exec(`CREATE TABLE sales (
		id TEXT PRIMARY KEY,
		region TEXT,
		product TEXT,
		amount REAL,
		qty INTEGER,
		salary REAL,
		deleted_at DATETIME,
		created_at DATETIME,
		updated_at DATETIME
	)`)
	sqlDB.Exec(`CREATE TABLE res_user (
		id TEXT PRIMARY KEY,
		username TEXT,
		email TEXT,
		is_superuser INTEGER DEFAULT 0
	)`)
	sqlDB.Exec(`CREATE TABLE model_access (
		id TEXT PRIMARY KEY,
		model_name TEXT,
		group_id TEXT DEFAULT '',
		can_select INTEGER DEFAULT 0,
		can_read INTEGER DEFAULT 0,
		can_write INTEGER DEFAULT 0,
		can_create INTEGER DEFAULT 0,
		can_delete INTEGER DEFAULT 0,
		can_print INTEGER DEFAULT 0,
		can_email INTEGER DEFAULT 0,
		can_report INTEGER DEFAULT 0,
		can_export INTEGER DEFAULT 0,
		can_import INTEGER DEFAULT 0,
		can_mask INTEGER DEFAULT 0,
		can_clone INTEGER DEFAULT 0
	)`)
	sqlDB.Exec(`CREATE TABLE "group" (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		display_name TEXT,
		category TEXT
	)`)
	sqlDB.Exec(`CREATE TABLE res_user_group (
		user_id TEXT,
		group_id TEXT,
		PRIMARY KEY (user_id, group_id)
	)`)
	sqlDB.Exec(`CREATE TABLE group_implies (
		group_id TEXT,
		implied_group_id TEXT,
		PRIMARY KEY (group_id, implied_group_id)
	)`)
	return db
}

func seedPivotData(t *testing.T, db *gorm.DB) {
	t.Helper()
	sqlDB, _ := db.DB()
	sqlDB.Exec(`INSERT INTO res_user (id, username, email, is_superuser) VALUES ('user-1', 'alice', 'alice@test.com', 0)`)
	sqlDB.Exec(`INSERT INTO res_user (id, username, email, is_superuser) VALUES ('admin-1', 'admin', 'admin@test.com', 1)`)

	sqlDB.Exec(`INSERT INTO sales (id, region, product, amount, qty) VALUES ('s1', 'North', 'Widget', 100.0, 2)`)
	sqlDB.Exec(`INSERT INTO sales (id, region, product, amount, qty) VALUES ('s2', 'North', 'Widget', 50.0, 1)`)
	sqlDB.Exec(`INSERT INTO sales (id, region, product, amount, qty) VALUES ('s3', 'North', 'Gadget', 200.0, 4)`)
	sqlDB.Exec(`INSERT INTO sales (id, region, product, amount, qty) VALUES ('s4', 'South', 'Widget', 300.0, 6)`)
	sqlDB.Exec(`INSERT INTO sales (id, region, product, amount, qty) VALUES ('s5', 'South', 'Gadget', 150.0, 3)`)

	sqlDB.Exec(`INSERT INTO model_access (id, model_name, group_id, can_select, can_read)
		VALUES ('ma1', 'sales.sale', '', 1, 1)`)
}

func newPivotTestApp(db *gorm.DB, superuser bool) *fiber.App {
	app := fiber.New()

	registry := newPivotModelRegistry()
	perm := persistence.NewPermissionService(db)
	perm.SetTableNameResolver(func(name string) string {
		if name == "user" {
			return "res_user"
		}
		return name
	})

	handler := NewPivotHandler(PivotHandlerConfig{
		DB:                db,
		PermissionService: perm,
		ModelRegistry:     registry,
	})

	group := app.Group("/api/v1/sales/sale", func(c *fiber.Ctx) error {
		if superuser {
			c.Locals("user_id", "admin-1")
			c.Locals("is_superuser", true)
		} else {
			c.Locals("user_id", "user-1")
			c.Locals("is_superuser", false)
		}
		c.Locals("pivot_module", "sales")
		c.Locals("pivot_model", "sales.sale")
		return c.Next()
	})
	handler.RegisterRoutes(group)

	return app
}

func TestPivotHandler_SingleAgg(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := newPivotTestApp(db, true)

	body := `{"rows":["region"],"cols":["product"],"values":[{"field":"amount","agg":"sum"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			RowHeaders  []string     `json:"rowHeaders"`
			ColHeaders  []string     `json:"colHeaders"`
			Cells       [][][]float64 `json:"cells"`
			GrandTotal  []float64    `json:"grandTotal"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(envelope.Data.RowHeaders) != 2 {
		t.Errorf("expected 2 row headers, got %d", len(envelope.Data.RowHeaders))
	}
	if len(envelope.Data.ColHeaders) != 2 {
		t.Errorf("expected 2 col headers, got %d", len(envelope.Data.ColHeaders))
	}
	if envelope.Data.GrandTotal[0] != 800 {
		t.Errorf("grand total: expected 800, got %v", envelope.Data.GrandTotal[0])
	}
}

func TestPivotHandler_MultiAgg(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := newPivotTestApp(db, true)

	body := `{"rows":["region"],"cols":["product"],"values":[{"field":"amount","agg":"sum"},{"field":"qty","agg":"count"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			ValueFields []struct {
				Field string `json:"field"`
				Agg   string `json:"agg"`
			} `json:"valueFields"`
			GrandTotal []float64 `json:"grandTotal"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(envelope.Data.ValueFields) != 2 {
		t.Errorf("expected 2 value fields, got %d", len(envelope.Data.ValueFields))
	}
}

func TestPivotHandler_MissingRows(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := newPivotTestApp(db, true)

	body := `{"cols":["product"],"values":[{"field":"amount","agg":"sum"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for missing rows, got %d", resp.StatusCode)
	}
}

func TestPivotHandler_MissingValues(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := newPivotTestApp(db, true)

	body := `{"rows":["region"],"cols":["product"]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for missing values, got %d", resp.StatusCode)
	}
}

func TestPivotHandler_InvalidAgg(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := newPivotTestApp(db, true)

	body := `{"rows":["region"],"values":[{"field":"amount","agg":"median"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for invalid agg, got %d", resp.StatusCode)
	}
}

func TestPivotHandler_PermissionDenied(t *testing.T) {
	db := setupPivotTestDB(t)
	sqlDB, _ := db.DB()
	sqlDB.Exec(`UPDATE model_access SET can_select=0, can_read=0 WHERE model_name='sales.sale'`)

	app := newPivotTestApp(db, false)

	body := `{"rows":["region"],"values":[{"field":"amount","agg":"sum"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 403 {
		t.Errorf("expected 403 for no read permission, got %d", resp.StatusCode)
	}
}

func TestPivotHandler_FieldLevelPermission(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := newPivotTestApp(db, false)

	body := `{"rows":["region"],"values":[{"field":"salary","agg":"sum"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for group-restricted field without permission, got %d (field should be dropped leaving no values)", resp.StatusCode)
	}
}

func TestPivotHandler_AuthRequired(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := fiber.New()

	registry := newPivotModelRegistry()
	perm := persistence.NewPermissionService(db)
	perm.SetTableNameResolver(func(name string) string { return name })

	handler := NewPivotHandler(PivotHandlerConfig{
		DB:                db,
		PermissionService: perm,
		ModelRegistry:     registry,
	})

	group := app.Group("/api/v1/sales/sale", func(c *fiber.Ctx) error {
		c.Locals("user_id", "")
		c.Locals("is_superuser", false)
		c.Locals("pivot_module", "sales")
		c.Locals("pivot_model", "sales.sale")
		return c.Next()
	})
	handler.RegisterRoutes(group)

	body := `{"rows":["region"],"values":[{"field":"amount","agg":"sum"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 for no auth, got %d", resp.StatusCode)
	}
}

func TestPivotHandler_NoCols(t *testing.T) {
	db := setupPivotTestDB(t)
	seedPivotData(t, db)
	app := newPivotTestApp(db, true)

	body := `{"rows":["region"],"values":[{"field":"amount","agg":"sum"}]}`
	req := httptest.NewRequest("POST", "/api/v1/sales/sale/pivot", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			ColHeaders []string `json:"colHeaders"`
			RowTotals  [][]float64 `json:"rowTotals"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&envelope)
	if len(envelope.Data.ColHeaders) != 1 || envelope.Data.ColHeaders[0] != "" {
		t.Errorf("expected single empty col header, got %v", envelope.Data.ColHeaders)
	}
}
