package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/bitcode-framework/bitcode/internal/runtime/comment"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupCommentTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.Exec(`CREATE TABLE comments (
		id TEXT PRIMARY KEY,
		body TEXT,
		commentable_type TEXT,
		commentable_id TEXT,
		parent_id TEXT,
		author_id TEXT,
		mentions TEXT,
		active INTEGER DEFAULT 1,
		created_at DATETIME,
		updated_at DATETIME
	)`)
	sqlDB.Exec(`CREATE TABLE notifications (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		type TEXT DEFAULT 'mention',
		title TEXT,
		body TEXT,
		model_name TEXT,
		record_id TEXT,
		comment_id TEXT,
		read INTEGER DEFAULT 0,
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
	sqlDB.Exec(`CREATE TABLE contacts (
		id TEXT PRIMARY KEY,
		name TEXT,
		created_by TEXT,
		updated_by TEXT,
		active INTEGER DEFAULT 1
	)`)
	return db
}

type mockCommentModelRegistry struct {
	models map[string]*parser.ModelDefinition
}

func newCommentHandlerApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	svc := comment.NewService(db)
	svc.SetTableName("comments")
	svc.SetNotificationTable("notifications")
	svc.SetUserTable("res_user")
	handler := NewCommentHandler(CommentHandlerConfig{
		Service: svc,
	})
	api := app.Group("/api/v1")
	g := api.Group("/:module", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-1")
		c.Locals("is_superuser", false)
		return c.Next()
	})
	handler.RegisterRoutes(g)
	return app
}

func TestCommentHandler_CreateAndList(t *testing.T) {
	db := setupCommentTestDB(t)
	sqlDB, _ := db.DB()
	sqlDB.Exec(`INSERT INTO res_user (id, username, email) VALUES ('user-1', 'alice', 'alice@test.com')`)
	sqlDB.Exec(`INSERT INTO res_user (id, username, email) VALUES ('user-2', 'bob', 'bob@test.com')`)
	sqlDB.Exec(`INSERT INTO contacts (id, name) VALUES ('rec-1', 'Test Contact')`)

	app := newCommentHandlerApp(db)

	body, _ := json.Marshal(map[string]any{
		"body":      "Hello @bob, check this out!",
		"record_id": "rec-1",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	if resp.StatusCode != 201 {
		rb, _ := json.Marshal(map[string]any{})
		_ = rb
		var errBody map[string]any
		json.NewDecoder(resp.Body).Decode(&errBody)
		t.Fatalf("expected 201, got %d: %v", resp.StatusCode, errBody)
	}

	var createResult map[string]any
	json.NewDecoder(resp.Body).Decode(&createResult)
	data, ok := createResult["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %v", createResult)
	}
	if data["body"] != "Hello @bob, check this out!" {
		t.Errorf("unexpected body: %v", data["body"])
	}
	mentions, ok := data["mentions"].([]any)
	if !ok || len(mentions) != 1 {
		t.Errorf("expected 1 mention, got %v", data["mentions"])
	} else {
		m := mentions[0].(map[string]any)
		if m["username"] != "bob" {
			t.Errorf("expected mention username bob, got %v", m["username"])
		}
		if m["resolved"] != true {
			t.Errorf("expected mention resolved true, got %v", m["resolved"])
		}
	}

	var notifCount int64
	db.Table("notifications").Where("user_id = ?", "user-2").Count(&notifCount)
	if notifCount != 1 {
		t.Errorf("expected 1 notification for user-2, got %d", notifCount)
	}

	req2 := httptest.NewRequest("GET", "/api/v1/crm/comments?record_id=rec-1", nil)
	resp2, err := app.Test(req2, -1)
	if err != nil {
		t.Fatalf("list request failed: %v", err)
	}
	if resp2.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
	var listResult map[string]any
	json.NewDecoder(resp2.Body).Decode(&listResult)
	thread, ok := listResult["data"].([]any)
	if !ok || len(thread) != 1 {
		t.Errorf("expected 1 thread item, got %v", listResult["data"])
	}
}

func TestCommentHandler_Threading(t *testing.T) {
	db := setupCommentTestDB(t)
	sqlDB, _ := db.DB()
	sqlDB.Exec(`INSERT INTO res_user (id, username, email) VALUES ('user-1', 'alice', 'alice@test.com')`)
	sqlDB.Exec(`INSERT INTO contacts (id, name) VALUES ('rec-1', 'Test Contact')`)

	app := newCommentHandlerApp(db)

	body, _ := json.Marshal(map[string]any{
		"body":      "Parent comment",
		"record_id": "rec-1",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201 for parent, got %d", resp.StatusCode)
	}
	var parentResult map[string]any
	json.NewDecoder(resp.Body).Decode(&parentResult)
	parentData := parentResult["data"].(map[string]any)
	parentID := parentData["id"].(string)

	replyBody, _ := json.Marshal(map[string]any{
		"body":      "Reply to parent",
		"record_id": "rec-1",
		"parent_id": parentID,
	})
	req2 := httptest.NewRequest("POST", "/api/v1/crm/comments", bytes.NewReader(replyBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2, -1)
	if resp2.StatusCode != 201 {
		t.Fatalf("expected 201 for reply, got %d", resp2.StatusCode)
	}

	req3 := httptest.NewRequest("GET", "/api/v1/crm/comments?record_id=rec-1", nil)
	resp3, _ := app.Test(req3, -1)
	var listResult map[string]any
	json.NewDecoder(resp3.Body).Decode(&listResult)
	thread := listResult["data"].([]any)
	if len(thread) != 1 {
		t.Fatalf("expected 1 top-level thread, got %d", len(thread))
	}
	topMsg := thread[0].(map[string]any)
	replies, ok := topMsg["replies"].([]any)
	if !ok || len(replies) != 1 {
		t.Errorf("expected 1 reply in thread, got %v", topMsg["replies"])
	}
}

func TestCommentHandler_MissingRecordID(t *testing.T) {
	db := setupCommentTestDB(t)
	app := newCommentHandlerApp(db)

	body, _ := json.Marshal(map[string]any{
		"body": "No record",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for missing record_id, got %d", resp.StatusCode)
	}
}

func TestCommentHandler_EmptyBody(t *testing.T) {
	db := setupCommentTestDB(t)
	app := newCommentHandlerApp(db)

	body, _ := json.Marshal(map[string]any{
		"body":      "",
		"record_id": "rec-1",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for empty body, got %d", resp.StatusCode)
	}
}

func TestCommentHandler_PermissionDenied(t *testing.T) {
	db := setupCommentTestDB(t)

	sqlDB, _ := db.DB()
	sqlDB.Exec(`CREATE TABLE record_rule (
		id TEXT PRIMARY KEY,
		name TEXT,
		model_name TEXT,
		group_names TEXT DEFAULT '',
		domain_filter TEXT,
		domain_filter_expr TEXT,
		can_read INTEGER DEFAULT 1,
		can_create INTEGER DEFAULT 1,
		can_write INTEGER DEFAULT 1,
		can_delete INTEGER DEFAULT 0,
		is_global INTEGER DEFAULT 0,
		active INTEGER DEFAULT 1,
		module TEXT,
		modified_source TEXT DEFAULT 'json'
	)`)
	sqlDB.Exec(`CREATE TABLE record_rule_groups (
		record_rule_id TEXT,
		group_id TEXT,
		PRIMARY KEY (record_rule_id, group_id)
	)`)
	sqlDB.Exec(`CREATE TABLE "group" (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		display_name TEXT,
		category TEXT
	)`)
	sqlDB.Exec(`CREATE TABLE user_group (
		user_id TEXT,
		group_id TEXT,
		PRIMARY KEY (user_id, group_id)
	)`)
	sqlDB.Exec(`CREATE TABLE group_implies (
		group_id TEXT,
		implied_group_id TEXT,
		PRIMARY KEY (group_id, implied_group_id)
	)`)
	sqlDB.Exec(`CREATE TABLE "user" (
		id TEXT PRIMARY KEY,
		username TEXT,
		email TEXT,
		is_superuser INTEGER DEFAULT 0
	)`)
	sqlDB.Exec(`CREATE TABLE crm (
		id TEXT PRIMARY KEY,
		name TEXT
	)`)
	sqlDB.Exec(`INSERT INTO "user" (id, username, email) VALUES ('user-1', 'alice', 'alice@test.com')`)
	sqlDB.Exec(`INSERT INTO crm (id, name) VALUES ('rec-1', 'Test Contact')`)
	sqlDB.Exec(`INSERT INTO record_rule (id, name, model_name, domain_filter, can_read, can_write, is_global, active, modified_source)
		VALUES ('r1', 'deny-crm', 'crm', '[["id","=","nonexistent"]]', 1, 1, 1, 1, 'json')`)

	svc := comment.NewService(db)
	svc.SetTableName("comments")
	svc.SetNotificationTable("notifications")
	svc.SetUserTable("user")
	if err := svc.EnsureTables(); err != nil {
		t.Fatalf("EnsureTables failed: %v", err)
	}

	rrSvc := persistence.NewRecordRuleService(db)
	rrSvc.SetTableNameResolver(func(name string) string {
		if name == "comment" {
			return "comments"
		}
		return name
	})

	handler := NewCommentHandler(CommentHandlerConfig{
		Service:          svc,
		RecordRuleService: rrSvc,
	})

	app := fiber.New()
	api := app.Group("/api/v1/:module", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-1")
		return c.Next()
	})
	handler.RegisterRoutes(api)

	body, _ := json.Marshal(map[string]any{
		"body":      "Should be denied",
		"record_id": "rec-1",
	})
	req := httptest.NewRequest("POST", "/api/v1/crm/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 403 {
		t.Errorf("expected 403 for record-rule denied access, got %d", resp.StatusCode)
	}
}
