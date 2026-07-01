package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/bitcode-framework/bitcode/internal/runtime/activity"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupActivityTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.Exec(`CREATE TABLE audit_logs (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		action TEXT,
		model_name TEXT,
		record_id TEXT,
		changes TEXT,
		ip_address TEXT,
		user_agent TEXT,
		request_method TEXT,
		request_path TEXT,
		status_code INTEGER,
		duration_ms INTEGER,
		impersonated_by TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		active INTEGER DEFAULT 1
	)`)
	sqlDB.Exec(`CREATE TABLE res_user (
		id TEXT PRIMARY KEY,
		username TEXT,
		name TEXT,
		email TEXT,
		is_superuser INTEGER DEFAULT 0
	)`)
	sqlDB.Exec(`CREATE TABLE contacts (
		id TEXT PRIMARY KEY,
		name TEXT,
		created_by TEXT,
		updated_by TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		active INTEGER DEFAULT 1
	)`)
	sqlDB.Exec(`CREATE TABLE model_access (
		id TEXT PRIMARY KEY,
		group_id TEXT DEFAULT '',
		model_name TEXT,
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
	sqlDB.Exec(`CREATE TABLE record_rule (
		id TEXT PRIMARY KEY,
		name TEXT,
		model_name TEXT,
		domain_filter TEXT DEFAULT '',
		domain_filter_expr TEXT DEFAULT '',
		can_read INTEGER DEFAULT 0,
		can_create INTEGER DEFAULT 0,
		can_write INTEGER DEFAULT 0,
		can_delete INTEGER DEFAULT 0,
		is_global INTEGER DEFAULT 0,
		active INTEGER DEFAULT 1,
		group_names TEXT DEFAULT '',
		modified_source TEXT DEFAULT 'json'
	)`)
	sqlDB.Exec(`CREATE TABLE record_rule_groups (
		record_rule_id TEXT,
		group_id TEXT,
		PRIMARY KEY (record_rule_id, group_id)
	)`)
	sqlDB.Exec(`INSERT INTO res_user (id, username, name, email, is_superuser) VALUES
		('u1', 'alice', 'Alice Wong', 'alice@test.com', 0),
		('u2', 'bob', '', 'bob@test.com', 0),
		('admin', 'admin', 'Administrator', 'admin@test.com', 1)`)
	sqlDB.Exec(`INSERT INTO contacts (id, name) VALUES ('rec-1', 'Test Contact')`)
	return db
}

func seedActivityLogs(t *testing.T, db *gorm.DB) {
	t.Helper()
	now := time.Now().UTC()
	rows := []struct {
		id, userID, action, model, recordID string
		offset                              time.Duration
	}{
		{"a1", "u1", "create", "contact", "rec-1", 0},
		{"a2", "u2", "update", "contact", "rec-1", -1 * time.Hour},
		{"a3", "u1", "delete", "contact", "rec-2", -2 * time.Hour},
		{"a4", "admin", "login", "", "", -3 * time.Hour},
		{"a5", "u2", "create", "lead", "lead-1", -4 * time.Hour},
	}
	for _, r := range rows {
		ts := now.Add(r.offset)
		db.Exec(`INSERT INTO audit_logs (id, user_id, action, model_name, record_id, created_at, updated_at, active)
			VALUES (?, ?, ?, ?, ?, ?, ?, 1)`, r.id, r.userID, r.action, r.model, r.recordID, ts, ts)
	}
}

func setupActivityHandler(t *testing.T, db *gorm.DB, userID string, isSuperuser bool) *fiber.App {
	t.Helper()
	svc := activity.NewService(db)
	svc.SetUserTable("res_user")

	permSvc := persistence.NewPermissionService(db)
	permSvc.SetTableNameResolver(func(name string) string {
		if name == "user" {
			return "res_user"
		}
		return name
	})
	rrSvc := persistence.NewRecordRuleService(db)
	rrSvc.SetTableNameResolver(func(name string) string {
		if name == "user" {
			return "res_user"
		}
		return name
	})

	handler := NewActivityHandler(ActivityHandlerConfig{
		Service:           svc,
		PermissionService: permSvc,
		RecordRuleService: rrSvc,
	})

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		c.Locals("is_superuser", isSuperuser)
		return c.Next()
	})
	handler.RegisterRoutes(app.Group("/api/v1"))
	return app
}

func TestActivityHandler_List_HappyPath(t *testing.T) {
	db := setupActivityTestDB(t)
	seedActivityLogs(t, db)
	app := setupActivityHandler(t, db, "admin", true)

	req := httptest.NewRequest("GET", "/api/v1/activity?model=contact&record_id=rec-1", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		b := readActivityBody(t, resp)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, b)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := result["data"].([]any)
	if !ok {
		t.Fatalf("expected data array, got %T", result["data"])
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 entries for contact/rec-1, got %d", len(data))
	}

	total, _ := result["total"].(float64)
	if total != 2 {
		t.Errorf("expected total 2, got %v", result["total"])
	}

	pageSize, _ := result["page_size"].(float64)
	if pageSize != 50 {
		t.Errorf("expected page_size 50, got %v", result["page_size"])
	}

	first := data[0].(map[string]any)
	if first["actor"] != "Alice Wong" {
		t.Errorf("expected first actor 'Alice Wong', got %v", first["actor"])
	}
	if first["verb"] != "created" {
		t.Errorf("expected verb 'created', got %v", first["verb"])
	}
}

func TestActivityHandler_List_GlobalTimeline(t *testing.T) {
	db := setupActivityTestDB(t)
	seedActivityLogs(t, db)
	app := setupActivityHandler(t, db, "admin", true)

	req := httptest.NewRequest("GET", "/api/v1/activity", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 5 {
		t.Errorf("expected 5 entries in global timeline, got %d", len(data))
	}
}

func TestActivityHandler_List_FilterByUser(t *testing.T) {
	db := setupActivityTestDB(t)
	seedActivityLogs(t, db)
	app := setupActivityHandler(t, db, "admin", true)

	req := httptest.NewRequest("GET", "/api/v1/activity?user=u2", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 2 {
		t.Errorf("expected 2 entries for u2, got %d", len(data))
	}
	for _, d := range data {
		entry := d.(map[string]any)
		if entry["actor_id"] != "u2" {
			t.Errorf("expected actor_id u2, got %v", entry["actor_id"])
		}
	}
}

func TestActivityHandler_List_FilterBySince(t *testing.T) {
	db := setupActivityTestDB(t)
	seedActivityLogs(t, db)
	app := setupActivityHandler(t, db, "admin", true)

	now := time.Now().UTC()
	since := now.Add(-90 * time.Minute).Format(time.RFC3339)

	req := httptest.NewRequest("GET", "/api/v1/activity?since="+since, nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	total, _ := result["total"].(float64)
	if total != 2 {
		t.Errorf("expected 2 entries within last 90min, got %v", result["total"])
	}
}

func TestActivityHandler_List_InvalidSince(t *testing.T) {
	db := setupActivityTestDB(t)
	app := setupActivityHandler(t, db, "admin", true)

	req := httptest.NewRequest("GET", "/api/v1/activity?since=not-a-date", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for invalid since, got %d", resp.StatusCode)
	}
}

func TestActivityHandler_List_Pagination(t *testing.T) {
	db := setupActivityTestDB(t)
	now := time.Now().UTC()
	for i := 0; i < 15; i++ {
		ts := now.Add(-time.Duration(i) * time.Minute)
		id := "pag-" + string(rune('a'+i))
		db.Exec(`INSERT INTO audit_logs (id, user_id, action, model_name, record_id, created_at, updated_at, active)
			VALUES (?, ?, 'update', 'contact', 'rec-1', ?, ?, 1)`,
			id, "u1", ts, ts)
	}
	app := setupActivityHandler(t, db, "admin", true)

	req := httptest.NewRequest("GET", "/api/v1/activity?limit=10&page=1", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 10 {
		t.Errorf("expected 10 entries on page 1, got %d", len(data))
	}
	total, _ := result["total"].(float64)
	if total != 15 {
		t.Errorf("expected total 15, got %v", result["total"])
	}

	req2 := httptest.NewRequest("GET", "/api/v1/activity?limit=10&page=2", nil)
	resp2, err := app.Test(req2, -1)
	if err != nil {
		t.Fatalf("page2 request failed: %v", err)
	}
	var result2 map[string]any
	json.NewDecoder(resp2.Body).Decode(&result2)
	data2 := result2["data"].([]any)
	if len(data2) != 5 {
		t.Errorf("expected 5 entries on page 2, got %d", len(data2))
	}
}

func TestActivityHandler_List_EmptyResult(t *testing.T) {
	db := setupActivityTestDB(t)
	app := setupActivityHandler(t, db, "admin", true)

	req := httptest.NewRequest("GET", "/api/v1/activity?model=nonexistent", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data, ok := result["data"].([]any)
	if !ok {
		t.Fatalf("expected data array, got %T", result["data"])
	}
	if len(data) != 0 {
		t.Errorf("expected 0 entries, got %d", len(data))
	}
	total, _ := result["total"].(float64)
	if total != 0 {
		t.Errorf("expected total 0, got %v", result["total"])
	}
}

func TestActivityHandler_List_PermissionDenied_NoReadAccess(t *testing.T) {
	db := setupActivityTestDB(t)
	seedActivityLogs(t, db)
	app := setupActivityHandler(t, db, "u1", false)

	req := httptest.NewRequest("GET", "/api/v1/activity?model=contact&record_id=rec-1", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 403 {
		b := readActivityBody(t, resp)
		t.Fatalf("expected 403 (no read permission), got %d: %s", resp.StatusCode, b)
	}
}

func TestActivityHandler_List_Unauthenticated(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := activity.NewService(db)
	handler := NewActivityHandler(ActivityHandlerConfig{Service: svc})

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "")
		c.Locals("is_superuser", false)
		return c.Next()
	})
	handler.RegisterRoutes(app.Group("/api/v1"))

	req := httptest.NewRequest("GET", "/api/v1/activity", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 unauthenticated, got %d", resp.StatusCode)
	}
}

func TestActivityHandler_List_LimitClampedToMax(t *testing.T) {
	db := setupActivityTestDB(t)
	seedActivityLogs(t, db)
	app := setupActivityHandler(t, db, "admin", true)

	req := httptest.NewRequest("GET", "/api/v1/activity?limit=99999", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	pageSize, _ := result["page_size"].(float64)
	if pageSize != 200 {
		t.Errorf("expected page_size clamped to 200, got %v", result["page_size"])
	}
}

func readActivityBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	return string(b)
}
