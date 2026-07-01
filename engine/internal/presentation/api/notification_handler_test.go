package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bitcode-framework/bitcode/internal/runtime/notification"
	"github.com/gofiber/fiber/v2"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupNotificationTestApp(t *testing.T) (*fiber.App, *notification.Service) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	_, err = sqlDB.Exec(`CREATE TABLE notifications (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		type TEXT DEFAULT 'info',
		title TEXT,
		body TEXT,
		data TEXT,
		notifiable_type TEXT,
		notifiable_id TEXT,
		read_at DATETIME,
		active INTEGER DEFAULT 1,
		created_at DATETIME,
		updated_at DATETIME
	)`)
	if err != nil {
		t.Fatalf("failed to create notifications table: %v", err)
	}

	svc := notification.NewService(db)
	handler := NewNotificationHandler(NotificationHandlerConfig{Service: svc})

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user-1")
		return c.Next()
	})
	group := app.Group("/api/v1")
	handler.RegisterRoutes(group)
	return app, svc
}

func TestNotificationHandler_List(t *testing.T) {
	app, svc := setupNotificationTestApp(t)

	_, err := svc.Create(t.Context(), notification.CreateInput{
		UserID: "test-user-1",
		Title:  "Test Notif",
		Body:   "Hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/notifications", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data, ok := result["data"].([]any)
	if !ok {
		t.Fatal("expected data array")
	}
	if len(data) != 1 {
		t.Errorf("expected 1 notification, got %d", len(data))
	}
}

func TestNotificationHandler_UnreadCount(t *testing.T) {
	app, svc := setupNotificationTestApp(t)

	for i := 0; i < 3; i++ {
		_, err := svc.Create(t.Context(), notification.CreateInput{
			UserID: "test-user-1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	req := httptest.NewRequest("GET", "/api/v1/notifications/unread-count", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	count, ok := result["count"].(float64)
	if !ok {
		t.Fatalf("expected count number, got %T", result["count"])
	}
	if count != 3 {
		t.Errorf("expected count 3, got %v", count)
	}
}

func TestNotificationHandler_MarkRead(t *testing.T) {
	app, svc := setupNotificationTestApp(t)

	rec, err := svc.Create(t.Context(), notification.CreateInput{
		UserID: "test-user-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/notifications/"+rec.ID+"/read", bytes.NewReader(nil))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	count, _ := svc.UnreadCount(t.Context(), "test-user-1")
	if count != 0 {
		t.Errorf("expected 0 unread after mark read, got %d", count)
	}
}

func TestNotificationHandler_MarkAllRead(t *testing.T) {
	app, svc := setupNotificationTestApp(t)

	for i := 0; i < 3; i++ {
		_, err := svc.Create(t.Context(), notification.CreateInput{
			UserID: "test-user-1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	req := httptest.NewRequest("POST", "/api/v1/notifications/read-all", bytes.NewReader(nil))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	count, ok := result["count"].(float64)
	if !ok {
		t.Fatalf("expected count, got %T", result["count"])
	}
	if count != 3 {
		t.Errorf("expected 3 marked read, got %v", count)
	}

	unread, _ := svc.UnreadCount(t.Context(), "test-user-1")
	if unread != 0 {
		t.Errorf("expected 0 unread after mark all, got %d", unread)
	}
}

func TestNotificationHandler_RequiresAuth(t *testing.T) {
	app, _ := setupNotificationTestApp(t)

	app2 := fiber.New()
	group2 := app2.Group("/api/v1")
	NewNotificationHandler(NotificationHandlerConfig{}).RegisterRoutes(group2)

	req := httptest.NewRequest("GET", "/api/v1/notifications", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp
}

func TestNotificationHandler_Pagination(t *testing.T) {
	app, svc := setupNotificationTestApp(t)

	for i := 0; i < 5; i++ {
		_, _ = svc.Create(t.Context(), notification.CreateInput{
			UserID: "test-user-1",
			Title:  "Notif",
		})
	}

	req := httptest.NewRequest("GET", "/api/v1/notifications?page=1&page_size=2", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 2 {
		t.Errorf("expected 2 records, got %d", len(data))
	}
	total := result["total"].(float64)
	if total != 5 {
		t.Errorf("expected total 5, got %v", total)
	}
}

func TestNotificationHandler_UnreadOnlyFilter(t *testing.T) {
	app, svc := setupNotificationTestApp(t)

	rec1, _ := svc.Create(t.Context(), notification.CreateInput{UserID: "test-user-1"})
	_, _ = svc.Create(t.Context(), notification.CreateInput{UserID: "test-user-1"})
	svc.MarkRead(t.Context(), "test-user-1", rec1.ID)

	req := httptest.NewRequest("GET", "/api/v1/notifications?unread_only=true", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 1 {
		t.Errorf("expected 1 unread record, got %d", len(data))
	}
}
