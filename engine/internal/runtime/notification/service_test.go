package notification

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type mockBroadcaster struct {
	messages []struct {
		userID  string
		channel string
		data    any
	}
}

func (m *mockBroadcaster) Broadcast(channel string, data any) {}
func (m *mockBroadcaster) BroadcastToUser(userID string, channel string, data any) {
	m.messages = append(m.messages, struct {
		userID  string
		channel string
		data    any
	}{userID, channel, data})
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
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
	return db
}

func TestCreate_RequiresUserID(t *testing.T) {
	svc := NewService(setupTestDB(t))
	_, err := svc.Create(context.Background(), CreateInput{Type: "info"})
	if err == nil {
		t.Error("expected error for missing user_id")
	}
}

func TestCreate_InvalidType(t *testing.T) {
	svc := NewService(setupTestDB(t))
	_, err := svc.Create(context.Background(), CreateInput{UserID: "u-1", Type: "bogus"})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestCreate_DefaultsTypeToInfo(t *testing.T) {
	svc := NewService(setupTestDB(t))
	rec, err := svc.Create(context.Background(), CreateInput{UserID: "u-1", Title: "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Type != "info" {
		t.Errorf("expected type=info, got %s", rec.Type)
	}
	if rec.ID == "" {
		t.Error("expected non-empty ID")
	}
	if rec.ReadAt != nil {
		t.Error("expected nil ReadAt for new notification")
	}
}

func TestCreate_BroadcastsToUser(t *testing.T) {
	mb := &mockBroadcaster{}
	svc := NewService(setupTestDB(t))
	svc.SetBroadcaster(mb)
	_, err := svc.Create(context.Background(), CreateInput{UserID: "u-1", Title: "Hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mb.messages) != 1 {
		t.Fatalf("expected 1 broadcast, got %d", len(mb.messages))
	}
	if mb.messages[0].userID != "u-1" {
		t.Errorf("expected broadcast to u-1, got %s", mb.messages[0].userID)
	}
	if mb.messages[0].channel != "notification:user:u-1" {
		t.Errorf("expected channel notification:user:u-1, got %s", mb.messages[0].channel)
	}
}

func TestCreate_SerializesData(t *testing.T) {
	svc := NewService(setupTestDB(t))
	rec, err := svc.Create(context.Background(), CreateInput{
		UserID: "u-1",
		Data:   map[string]any{"key": "value", "num": 42},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Data == nil {
		t.Fatal("expected non-nil Data")
	}
}

func TestList_UnreadFirst(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	oldRead := now.Add(-2 * time.Hour)
	recentUnread := now.Add(-30 * time.Minute)
	olderUnread := now.Add(-1 * time.Hour)

	seedNotifications(db, []map[string]any{
		{"id": "n-read-old", "user_id": "u-1", "type": "info", "title": "Read Old", "read_at": oldRead, "active": 1, "created_at": oldRead, "updated_at": oldRead},
		{"id": "n-unread-recent", "user_id": "u-1", "type": "info", "title": "Unread Recent", "read_at": nil, "active": 1, "created_at": recentUnread, "updated_at": recentUnread},
		{"id": "n-unread-older", "user_id": "u-1", "type": "info", "title": "Unread Older", "read_at": nil, "active": 1, "created_at": olderUnread, "updated_at": olderUnread},
		{"id": "n-read-recent", "user_id": "u-1", "type": "info", "title": "Read Recent", "read_at": now, "active": 1, "created_at": now, "updated_at": now},
	})

	result, err := svc.List(ctx, ListOptions{UserID: "u-1", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Records) != 4 {
		t.Fatalf("expected 4 records, got %d", len(result.Records))
	}

	if result.Records[0].ReadAt != nil {
		t.Error("expected first record to be unread (read_at IS NULL sorts first)")
	}
	if result.Records[1].ReadAt != nil {
		t.Error("expected second record to be unread")
	}
	if result.Records[2].ReadAt == nil {
		t.Error("expected third record to be read")
	}

	if result.Records[0].Title != "Unread Recent" {
		t.Errorf("expected most recent unread first, got %s", result.Records[0].Title)
	}
	if result.Records[1].Title != "Unread Older" {
		t.Errorf("expected older unread second, got %s", result.Records[1].Title)
	}

	if result.UnreadCount != 2 {
		t.Errorf("expected unread count 2, got %d", result.UnreadCount)
	}
	if result.Total != 4 {
		t.Errorf("expected total 4, got %d", result.Total)
	}
}

func TestList_UserScopeIsolation(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateInput{UserID: "u-1", Title: "User 1 Notif"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = svc.Create(ctx, CreateInput{UserID: "u-2", Title: "User 2 Notif"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := svc.List(ctx, ListOptions{UserID: "u-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("expected 1 record for u-1, got %d", len(result.Records))
	}
	if result.Records[0].UserID != "u-1" {
		t.Errorf("expected user_id=u-1, got %s", result.Records[0].UserID)
	}
}

func TestList_UnreadOnlyFilter(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	rec, err := svc.Create(ctx, CreateInput{UserID: "u-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rec2, err := svc.Create(ctx, CreateInput{UserID: "u-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := svc.MarkRead(ctx, "u-1", rec2.ID); err != nil {
		t.Fatalf("unexpected error marking read: %v", err)
	}

	result, err := svc.List(ctx, ListOptions{UserID: "u-1", UnreadOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("expected 1 unread record, got %d", len(result.Records))
	}
	if result.Records[0].ID != rec.ID {
		t.Errorf("expected %s, got %s", rec.ID, result.Records[0].ID)
	}
}

func TestList_Pagination(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, err := svc.Create(ctx, CreateInput{UserID: "u-1", Title: "Notif"})
		if err != nil {
			t.Fatalf("unexpected error on create %d: %v", i, err)
		}
	}

	result, err := svc.List(ctx, ListOptions{UserID: "u-1", Page: 1, PageSize: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Records) != 2 {
		t.Fatalf("expected 2 records on page 1, got %d", len(result.Records))
	}
	if result.Total != 5 {
		t.Errorf("expected total 5, got %d", result.Total)
	}
	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}
}

func TestUnreadCount(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := svc.Create(ctx, CreateInput{UserID: "u-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	rec4, _ := svc.Create(ctx, CreateInput{UserID: "u-1"})
	svc.MarkRead(ctx, "u-1", rec4.ID)

	count, err := svc.UnreadCount(ctx, "u-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 unread, got %d", count)
	}
}

func TestMarkRead(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	rec, err := svc.Create(ctx, CreateInput{UserID: "u-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = svc.MarkRead(ctx, "u-1", rec.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count, _ := svc.UnreadCount(ctx, "u-1")
	if count != 0 {
		t.Errorf("expected 0 unread after mark read, got %d", count)
	}

	result, _ := svc.List(ctx, ListOptions{UserID: "u-1"})
	if len(result.Records) != 1 && result.Records[0].ReadAt == nil {
		t.Error("expected read_at to be set")
	}
}

func TestMarkRead_OtherUserDenied(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	rec, err := svc.Create(ctx, CreateInput{UserID: "u-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = svc.MarkRead(ctx, "u-2", rec.ID)
	if err == nil {
		t.Error("expected error when marking another user's notification")
	}
}

func TestMarkAllRead(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := svc.Create(ctx, CreateInput{UserID: "u-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	count, err := svc.MarkAllRead(ctx, "u-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 marked read, got %d", count)
	}

	unread, _ := svc.UnreadCount(ctx, "u-1")
	if unread != 0 {
		t.Errorf("expected 0 unread after mark all, got %d", unread)
	}
}

func TestMarkAllRead_DoesNotAffectOtherUsers(t *testing.T) {
	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateInput{UserID: "u-1"})
	_, _ = svc.Create(ctx, CreateInput{UserID: "u-2"})

	_, err := svc.MarkAllRead(ctx, "u-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unread2, _ := svc.UnreadCount(ctx, "u-2")
	if unread2 != 1 {
		t.Errorf("expected u-2 to still have 1 unread, got %d", unread2)
	}
}

func TestValidType(t *testing.T) {
	valid := []string{"info", "warning", "error", "success", "mention", "message"}
	for _, v := range valid {
		if !ValidType(v) {
			t.Errorf("expected %s to be valid", v)
		}
	}
	if ValidType("bogus") {
		t.Error("expected 'bogus' to be invalid")
	}
}

func TestCreate_WithNotifiable(t *testing.T) {
	svc := NewService(setupTestDB(t))
	rec, err := svc.Create(context.Background(), CreateInput{
		UserID:         "u-1",
		Type:           "mention",
		NotifiableType: "comment",
		NotifiableID:   "rec-123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.NotifiableType != "comment" {
		t.Errorf("expected notifiable_type=comment, got %s", rec.NotifiableType)
	}
	if rec.NotifiableID != "rec-123" {
		t.Errorf("expected notifiable_id=rec-123, got %s", rec.NotifiableID)
	}
}

func seedNotifications(db *gorm.DB, rows []map[string]any) {
	for _, row := range rows {
		normalized := normalizeSeedRow(row)
		db.Table("notifications").Create(&normalized)
	}
}

func normalizeSeedRow(row map[string]any) map[string]any {
	out := make(map[string]any, len(row))
	for k, v := range row {
		out[k] = v
	}
	if _, ok := out["active"]; !ok {
		out["active"] = 1
	}
	if _, ok := out["created_at"]; !ok {
		out["created_at"] = time.Now().UTC()
	}
	if _, ok := out["updated_at"]; !ok {
		out["updated_at"] = time.Now().UTC()
	}
	return out
}
