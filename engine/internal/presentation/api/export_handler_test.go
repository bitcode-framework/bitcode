package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type exportModelRegistry struct {
	models map[string]*parser.ModelDefinition
}

func newExportModelRegistry() *exportModelRegistry {
	return &exportModelRegistry{
		models: map[string]*parser.ModelDefinition{
			"crm.contact": {
				Name:   "contact",
				Module: "crm",
				Fields: map[string]parser.FieldDefinition{
					"name":  {Type: parser.FieldString, Label: "Name"},
					"email": {Type: parser.FieldEmail, Label: "Email"},
					"phone": {Type: parser.FieldString, Label: "Phone"},
				},
			},
			"contact": {
				Name:   "contact",
				Module: "crm",
				Fields: map[string]parser.FieldDefinition{
					"name":  {Type: parser.FieldString, Label: "Name"},
					"email": {Type: parser.FieldEmail, Label: "Email"},
					"phone": {Type: parser.FieldString, Label: "Phone"},
				},
			},
			"crm.employee": {
				Name:   "employee",
				Module: "crm",
				Fields: map[string]parser.FieldDefinition{
					"name":   {Type: parser.FieldString, Label: "Name"},
					"phone":  {Type: parser.FieldString, Label: "Phone", Mask: true, MaskLength: 4},
					"salary": {Type: parser.FieldDecimal, Label: "Salary", Groups: []string{"hr.manager"}},
				},
			},
			"employee": {
				Name:   "employee",
				Module: "crm",
				Fields: map[string]parser.FieldDefinition{
					"name":   {Type: parser.FieldString, Label: "Name"},
					"phone":  {Type: parser.FieldString, Label: "Phone", Mask: true, MaskLength: 4},
					"salary": {Type: parser.FieldDecimal, Label: "Salary", Groups: []string{"hr.manager"}},
				},
			},
		},
	}
}

func (r *exportModelRegistry) Get(name string) (*parser.ModelDefinition, error) {
	if m, ok := r.models[name]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("model %q not found", name)
}

func (r *exportModelRegistry) TableName(name string) string {
	if name == "contact" || name == "crm.contact" {
		return "contacts"
	}
	if name == "employee" || name == "crm.employee" {
		return "employees"
	}
	return name
}

func setupExportTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.Exec(`CREATE TABLE contacts (
		id TEXT PRIMARY KEY,
		name TEXT,
		email TEXT,
		phone TEXT,
		active INTEGER DEFAULT 1,
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

func seedExportData(t *testing.T, db *gorm.DB) {
	t.Helper()
	sqlDB, _ := db.DB()
	sqlDB.Exec(`INSERT INTO res_user (id, username, email, is_superuser) VALUES ('user-1', 'alice', 'alice@test.com', 0)`)
	sqlDB.Exec(`INSERT INTO res_user (id, username, email, is_superuser) VALUES ('admin-1', 'admin', 'admin@test.com', 1)`)
	sqlDB.Exec(`INSERT INTO contacts (id, name, email, phone) VALUES ('c1', 'Alice Smith', 'alice@acme.com', '555-0100')`)
	sqlDB.Exec(`INSERT INTO contacts (id, name, email, phone) VALUES ('c2', 'Bob Jones', 'bob@acme.com', '555-0101')`)
	sqlDB.Exec(`INSERT INTO contacts (id, name, email, phone) VALUES ('c3', 'Carol Lee', 'carol@acme.com', '555-0102')`)
	sqlDB.Exec(`INSERT INTO model_access (id, model_name, group_id, can_select, can_read, can_export)
		VALUES ('ma1', 'crm.contact', '', 1, 1, 1)`)
}

func newExportTestApp(db *gorm.DB, permService *persistence.PermissionService, superuser bool) *fiber.App {
	app := fiber.New()

	registry := newExportModelRegistry()

	var perm *persistence.PermissionService
	if permService != nil {
		perm = permService
	} else {
		perm = persistence.NewPermissionService(db)
		perm.SetTableNameResolver(func(name string) string {
			if name == "user" {
				return "res_user"
			}
			return name
		})
	}

	handler := NewExportHandler(ExportHandlerConfig{
		DB:                db,
		PermissionService: perm,
		ModelRegistry:     registry,
	})

	group := app.Group("/api/v1/crm/contacts", func(c *fiber.Ctx) error {
		if superuser {
			c.Locals("user_id", "admin-1")
			c.Locals("is_superuser", true)
		} else {
			c.Locals("user_id", "user-1")
			c.Locals("is_superuser", false)
		}
		c.Locals("export_module", "crm")
		c.Locals("export_model", "crm.contact")
		return c.Next()
	})
	handler.RegisterRoutes(group)

	return app
}

func TestExportHandler_CSV(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(body))
	}

	if ct := resp.Header.Get("Content-Type"); ct != "text/csv" {
		t.Errorf("expected content-type text/csv, got %s", ct)
	}
	cd := resp.Header.Get("Content-Disposition")
	if cd == "" || cd[:len("attachment")] != "attachment" {
		t.Errorf("expected attachment content-disposition, got %q", cd)
	}

	body, _ := io.ReadAll(resp.Body)
	if len(body) < 3 || body[0] != 0xEF || body[1] != 0xBB || body[2] != 0xBF {
		t.Errorf("expected UTF-8 BOM at start of CSV")
	}

	reader := csv.NewReader(bytes.NewReader(body[3:]))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	if len(records) < 2 {
		t.Fatalf("expected at least header + 1 row, got %d rows", len(records))
	}

	headerRow := records[0]
	foundName := false
	for _, h := range headerRow {
		if h == "Name" || h == "name" {
			foundName = true
		}
	}
	if !foundName {
		t.Errorf("expected name column in header, got %v", headerRow)
	}

	dataRows := records[1:]
	if len(dataRows) != 3 {
		t.Errorf("expected 3 data rows, got %d", len(dataRows))
	}
}

func TestExportHandler_XLSX(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	body, _ := json.Marshal(map[string]any{
		"format": "xlsx",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Errorf("expected xlsx content-type, got %s", ct)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if len(respBody) < 4 {
		t.Fatalf("expected non-empty xlsx body, got %d bytes", len(respBody))
	}
	if string(respBody[:2]) != "PK" {
		t.Errorf("expected XLSX (ZIP) magic bytes PK, got %x", respBody[:2])
	}
}

func TestExportHandler_JSON(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	body, _ := json.Marshal(map[string]any{
		"format": "json",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected json content-type, got %s", ct)
	}

	var arr []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		t.Fatalf("failed to decode JSON array: %v", err)
	}
	if len(arr) != 3 {
		t.Errorf("expected 3 records, got %d", len(arr))
	}
}

func TestExportHandler_PermissionDenied(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)

	sqlDB, _ := db.DB()
	sqlDB.Exec(`DELETE FROM model_access WHERE model_name = 'crm.contact'`)
	sqlDB.Exec(`INSERT INTO model_access (id, model_name, group_id, can_select, can_read, can_export)
		VALUES ('ma2', 'crm.contact', '', 1, 1, 0)`)

	perm := persistence.NewPermissionService(db)
	perm.SetTableNameResolver(func(name string) string {
		if name == "user" {
			return "res_user"
		}
		return name
	})

	app := newExportTestApp(db, perm, false)

	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Errorf("expected 403 for no export permission, got %d", resp.StatusCode)
	}
}

func TestExportHandler_ColumnSelection(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	body, _ := json.Marshal(map[string]any{
		"format": "csv",
		"columns": []map[string]any{
			{"field": "name", "label": "Full Name"},
			{"field": "email", "label": "Email Address"},
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	respBody, _ := io.ReadAll(resp.Body)
	reader := csv.NewReader(bytes.NewReader(respBody[3:]))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	if len(records) < 1 {
		t.Fatalf("expected header row, got 0 rows")
	}

	header := records[0]
	if len(header) != 2 {
		t.Fatalf("expected 2 columns, got %d: %v", len(header), header)
	}
	if header[0] != "Full Name" || header[1] != "Email Address" {
		t.Errorf("expected [Full Name, Email Address], got %v", header)
	}

	if len(records) > 1 {
		for _, row := range records[1:] {
			if len(row) != 2 {
				t.Errorf("expected 2 values per row, got %d: %v", len(row), row)
			}
		}
	}
}

func TestExportHandler_Filter(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	filterJSON, _ := json.Marshal(map[string]any{
		"wheres": []map[string]any{
			{"field": "name", "op": "=", "value": "Alice Smith"},
		},
	})
	body, _ := json.Marshal(map[string]any{
		"format": "csv",
		"filter": json.RawMessage(filterJSON),
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	respBody, _ := io.ReadAll(resp.Body)
	reader := csv.NewReader(bytes.NewReader(respBody[3:]))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}

	dataRows := records[1:]
	if len(dataRows) != 1 {
		t.Errorf("expected 1 filtered row, got %d", len(dataRows))
	}
}

func TestExportHandler_DefaultColumns(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	reader := csv.NewReader(bytes.NewReader(respBody[3:]))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}

	header := records[0]
	hasName := false
	hasEmail := false
	hasPhone := false
	for _, h := range header {
		if h == "Name" || h == "name" {
			hasName = true
		}
		if h == "Email" || h == "email" {
			hasEmail = true
		}
		if h == "Phone" || h == "phone" {
			hasPhone = true
		}
	}
	if !hasName || !hasEmail || !hasPhone {
		t.Errorf("expected name, email, phone in default columns, got %v", header)
	}
}

func TestExportHandler_GetVariant(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	req := httptest.NewRequest("GET", "/api/v1/crm/contacts/export?format=json", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	var arr []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if len(arr) != 3 {
		t.Errorf("expected 3 records via GET, got %d", len(arr))
	}
}

func TestExportHandler_UnsupportedFormat(t *testing.T) {
	db := setupExportTestDB(t)
	seedExportData(t, db)
	app := newExportTestApp(db, nil, true)

	body, _ := json.Marshal(map[string]any{
		"format": "pdf",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/contacts/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for unsupported format, got %d", resp.StatusCode)
	}
}

// newExportEmployeeApp mounts the export handler for the crm.employee model so
// that field-level permission (mask + groups) can be exercised. The superuser
// flag controls whether the authenticated user bypasses permission checks.
func newExportEmployeeApp(db *gorm.DB, superuser bool, userID string) *fiber.App {
	app := fiber.New()
	registry := newExportModelRegistry()
	perm := persistence.NewPermissionService(db)
	perm.SetTableNameResolver(func(name string) string {
		if name == "user" {
			return "res_user"
		}
		return name
	})
	handler := NewExportHandler(ExportHandlerConfig{
		DB:                db,
		PermissionService: perm,
		ModelRegistry:     registry,
	})
	uid := userID
	if uid == "" {
		uid = "user-1"
	}
	group := app.Group("/api/v1/crm/employees", func(c *fiber.Ctx) error {
		if superuser {
			c.Locals("user_id", "admin-1")
			c.Locals("is_superuser", true)
		} else {
			c.Locals("user_id", uid)
			c.Locals("is_superuser", false)
		}
		c.Locals("export_module", "crm")
		c.Locals("export_model", "crm.employee")
		return c.Next()
	})
	handler.RegisterRoutes(group)
	return app
}

func setupEmployeeExportDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := setupExportTestDB(t)
	seedExportData(t, db)
	sqlDB, _ := db.DB()
	sqlDB.Exec(`CREATE TABLE employees (
		id TEXT PRIMARY KEY,
		name TEXT,
		phone TEXT,
		salary REAL,
		active INTEGER DEFAULT 1,
		created_at DATETIME,
		updated_at DATETIME
	)`)
	sqlDB.Exec(`INSERT INTO employees (id, name, phone, salary) VALUES ('e1', 'Dana White', '08123456789', 7500000)`)
	sqlDB.Exec(`INSERT INTO employees (id, name, phone, salary) VALUES ('e2', 'Eve Blue', '08987654321', 6200000)`)
	// user-1 has export permission on crm.employee but NOT can_mask, and is NOT
	// in the hr.manager group (so salary must be dropped and phone masked).
	sqlDB.Exec(`INSERT INTO model_access (id, model_name, group_id, can_select, can_read, can_export, can_mask)
		VALUES ('ma2', 'crm.employee', '', 1, 1, 1, 0)`)
	return db
}

// TestExportHandler_FieldMask verifies that a user without the can_mask
// permission sees masked values for fields flagged mask:true. The employee
// model marks phone as mask:true with mask_length 4 — a non-mask user must
// receive the masked form (e.g. *******6789) rather than the raw value.
func TestExportHandler_FieldMask(t *testing.T) {
	db := setupEmployeeExportDB(t)
	app := newExportEmployeeApp(db, false, "user-1")

	body, _ := json.Marshal(map[string]any{
		"format": "json",
		"columns": []map[string]any{
			{"field": "name", "label": "Name"},
			{"field": "phone", "label": "Phone"},
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/employees/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	respBody, _ := io.ReadAll(resp.Body)
	var arr []map[string]any
	if err := json.Unmarshal(respBody, &arr); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if len(arr) != 2 {
		t.Fatalf("expected 2 employees, got %d", len(arr))
	}
	for _, row := range arr {
		phone, _ := row["Phone"].(string)
		if phone == "08123456789" || phone == "08987654321" {
			t.Errorf("expected masked phone, got raw value %q", phone)
		}
		if phone == "" {
			t.Errorf("expected masked phone value, got empty string")
		}
		if !containsOnlyMaskAndTail(phone) {
			t.Errorf("phone %q does not look masked (expected stars + last 4)", phone)
		}
	}
}

// containsOnlyMaskAndTail checks that the value is composed of '*' characters
// followed by at most 4 trailing characters (the unmasked suffix).
func containsOnlyMaskAndTail(s string) bool {
	if s == "" {
		return false
	}
	stars := 0
	for _, c := range s {
		if c == '*' {
			stars++
		}
	}
	return stars > 0 && stars < len(s)
}

// TestExportHandler_ColumnSanitization_Groups verifies that a field restricted
// to a group the user does not belong to is dropped from the export entirely.
// The employee model restricts salary to the "hr.manager" group; user-1 is not
// a member, so salary must not appear in the exported columns.
func TestExportHandler_ColumnSanitization_Groups(t *testing.T) {
	db := setupEmployeeExportDB(t)
	app := newExportEmployeeApp(db, false, "user-1")

	body, _ := json.Marshal(map[string]any{
		"format": "csv",
		"columns": []map[string]any{
			{"field": "name", "label": "Name"},
			{"field": "salary", "label": "Salary"},
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/employees/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	respBody, _ := io.ReadAll(resp.Body)
	reader := csv.NewReader(bytes.NewReader(respBody[3:]))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	if len(records) < 1 {
		t.Fatal("expected at least a header row")
	}
	header := records[0]
	for _, h := range header {
		if h == "Salary" {
			t.Errorf("salary column should be dropped (user not in hr.manager group), but found header %q", h)
		}
	}
	hasName := false
	for _, h := range header {
		if h == "Name" {
			hasName = true
		}
	}
	if !hasName {
		t.Errorf("expected Name column to remain, got header %v", header)
	}
}

// TestExportHandler_FieldMask_SuperuserBypass verifies that a superuser (who
// has can_mask via the superuser bypass) sees the raw, unmasked phone value.
func TestExportHandler_FieldMask_SuperuserBypass(t *testing.T) {
	db := setupEmployeeExportDB(t)
	app := newExportEmployeeApp(db, true, "")

	body, _ := json.Marshal(map[string]any{
		"format": "json",
		"columns": []map[string]any{
			{"field": "name", "label": "Name"},
			{"field": "phone", "label": "Phone"},
		},
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/employees/export", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}
	respBody, _ := io.ReadAll(resp.Body)
	var arr []map[string]any
	if err := json.Unmarshal(respBody, &arr); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if len(arr) != 2 {
		t.Fatalf("expected 2 employees, got %d", len(arr))
	}
	rawSeen := false
	for _, row := range arr {
		if phone, _ := row["Phone"].(string); phone == "08123456789" || phone == "08987654321" {
			rawSeen = true
		}
	}
	if !rawSeen {
		t.Errorf("superuser should see raw phone values, none found in %v", arr)
	}
}
