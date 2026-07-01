package activity

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
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
	sqlDB.Exec(`INSERT INTO res_user (id, username, name, email) VALUES
		('u1', 'alice', 'Alice Wong', 'alice@test.com'),
		('u2', 'bob', '', 'bob@test.com'),
		('admin', 'admin', 'Administrator', 'admin@test.com')`)
	return db
}

func insertAuditLog(t *testing.T, db *gorm.DB, id, userID, action, model, recordID string, createdAt time.Time) {
	t.Helper()
	db.Exec(`INSERT INTO audit_logs (id, user_id, action, model_name, record_id, created_at, updated_at, active)
		VALUES (?, ?, ?, ?, ?, ?, ?, 1)`, id, userID, action, model, recordID, createdAt, createdAt)
}

func TestService_ListActivity_ActionMapping(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	insertAuditLog(t, db, "a1", "u1", "create", "contact", "rec-1", now)
	insertAuditLog(t, db, "a2", "u2", "delete", "contact", "rec-2", now.Add(-1*time.Hour))
	insertAuditLog(t, db, "a3", "admin", "login", "", "", now.Add(-2*time.Hour))

	entries, total, err := svc.ListActivity(ctx, ActivityQueryOptions{Limit: 50})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	byID := make(map[string]ActivityEntry, len(entries))
	for _, e := range entries {
		byID[e.ID] = e
	}

	create := byID["a1"]
	if create.Verb != "created" {
		t.Errorf("expected verb 'created', got %q", create.Verb)
	}
	if create.Variant != VariantSuccess {
		t.Errorf("expected variant success, got %q", create.Variant)
	}
	if create.Actor != "Alice Wong" {
		t.Errorf("expected actor 'Alice Wong' (resolved from name), got %q", create.Actor)
	}

	del := byID["a2"]
	if del.Verb != "deleted" {
		t.Errorf("expected verb 'deleted', got %q", del.Verb)
	}
	if del.Variant != VariantDestructive {
		t.Errorf("expected variant destructive, got %q", del.Variant)
	}
	if del.Actor != "bob" {
		t.Errorf("expected actor 'bob' (username fallback, name empty), got %q", del.Actor)
	}

	login := byID["a3"]
	if login.Verb != "logged in" {
		t.Errorf("expected verb 'logged in', got %q", login.Verb)
	}
	if login.Variant != VariantInfo {
		t.Errorf("expected variant info, got %q", login.Variant)
	}
	if login.Actor != "Administrator" {
		t.Errorf("expected actor 'Administrator', got %q", login.Actor)
	}
}

func TestService_ListActivity_BatchUserResolution(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		insertAuditLog(t, db, uuid.New().String(), "u1", "update", "contact", "rec-1", now.Add(-time.Duration(i)*time.Minute))
	}
	for i := 0; i < 5; i++ {
		insertAuditLog(t, db, uuid.New().String(), "u2", "update", "contact", "rec-2", now.Add(-time.Duration(i)*time.Minute))
	}

	entries, _, err := svc.ListActivity(ctx, ActivityQueryOptions{Limit: 50})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}

	aliceCount := 0
	bobCount := 0
	for _, e := range entries {
		if e.Actor == "Alice Wong" {
			aliceCount++
		}
		if e.Actor == "bob" {
			bobCount++
		}
	}
	if aliceCount != 5 || bobCount != 5 {
		t.Errorf("expected 5 Alice + 5 Bob, got %d Alice + %d Bob", aliceCount, bobCount)
	}
}

func TestService_ListActivity_FilterByModelAndRecord(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	insertAuditLog(t, db, "a1", "u1", "create", "contact", "rec-1", now)
	insertAuditLog(t, db, "a2", "u2", "update", "contact", "rec-2", now)
	insertAuditLog(t, db, "a3", "u1", "create", "lead", "rec-3", now)

	entries, total, err := svc.ListActivity(ctx, ActivityQueryOptions{Model: "contact", RecordID: "rec-1"})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != "a1" {
		t.Errorf("expected entry a1, got %s", entries[0].ID)
	}
}

func TestService_ListActivity_FilterByUser(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	insertAuditLog(t, db, "a1", "u1", "create", "contact", "rec-1", now)
	insertAuditLog(t, db, "a2", "u2", "create", "contact", "rec-2", now)
	insertAuditLog(t, db, "a3", "u1", "update", "contact", "rec-3", now)

	entries, total, err := svc.ListActivity(ctx, ActivityQueryOptions{UserID: "u1"})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	for _, e := range entries {
		if e.ActorID != "u1" {
			t.Errorf("expected all entries for u1, got actor_id %q", e.ActorID)
		}
	}
}

func TestService_ListActivity_FilterBySince(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	base := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	insertAuditLog(t, db, "old1", "u1", "create", "contact", "rec-1", base)
	insertAuditLog(t, db, "old2", "u1", "update", "contact", "rec-1", base.Add(1*time.Hour))
	insertAuditLog(t, db, "new1", "u1", "delete", "contact", "rec-2", base.Add(3*time.Hour))

	since := base.Add(2 * time.Hour)
	entries, total, err := svc.ListActivity(ctx, ActivityQueryOptions{Since: since})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1 (only new1 after since), got %d", total)
	}
	if len(entries) != 1 || entries[0].ID != "new1" {
		t.Errorf("expected only new1, got %v", entries)
	}
}

func TestService_ListActivity_Pagination(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	for i := 0; i < 10; i++ {
		insertAuditLog(t, db, uuid.New().String(), "u1", "update", "contact", "rec-1", now.Add(-time.Duration(i)*time.Minute))
	}

	page1, total, err := svc.ListActivity(ctx, ActivityQueryOptions{Limit: 5, Offset: 0})
	if err != nil {
		t.Fatalf("page1 failed: %v", err)
	}
	if total != 10 {
		t.Fatalf("expected total 10, got %d", total)
	}
	if len(page1) != 5 {
		t.Fatalf("expected 5 entries on page 1, got %d", len(page1))
	}

	page2, _, err := svc.ListActivity(ctx, ActivityQueryOptions{Limit: 5, Offset: 5})
	if err != nil {
		t.Fatalf("page2 failed: %v", err)
	}
	if len(page2) != 5 {
		t.Fatalf("expected 5 entries on page 2, got %d", len(page2))
	}

	seen := make(map[string]bool)
	for _, e := range page1 {
		seen[e.ID] = true
	}
	for _, e := range page2 {
		if seen[e.ID] {
			t.Errorf("entry %s appeared on both pages", e.ID)
		}
	}
}

func TestService_ListActivity_LimitClamping(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	entries, _, err := svc.ListActivity(ctx, ActivityQueryOptions{Limit: 0})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	_ = entries

	entries, _, err = svc.ListActivity(ctx, ActivityQueryOptions{Limit: 1000})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	_ = entries
}

func TestService_ListActivity_EmptyResult(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	entries, total, err := svc.ListActivity(ctx, ActivityQueryOptions{Model: "nonexistent"})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestService_ListActivity_GlobalTimeline(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	insertAuditLog(t, db, "a1", "u1", "create", "contact", "rec-1", now)
	insertAuditLog(t, db, "a2", "u2", "login", "", "", now.Add(-1*time.Hour))
	insertAuditLog(t, db, "a3", "admin", "register", "", "", now.Add(-2*time.Hour))

	entries, total, err := svc.ListActivity(ctx, ActivityQueryOptions{})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3 for global timeline, got %d", total)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Timestamp == "" {
		t.Error("expected non-empty timestamp on first entry")
	}
}

func TestService_ListActivity_TargetBuilding(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	insertAuditLog(t, db, "a1", "u1", "create", "contact", "rec-1", now)
	insertAuditLog(t, db, "a2", "u1", "request", "", "", now.Add(-1*time.Hour))
	db.Exec(`UPDATE audit_logs SET request_path = '/api/v1/contacts' WHERE id = ?`, "a2")

	entries, _, err := svc.ListActivity(ctx, ActivityQueryOptions{Limit: 50})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}

	byID := make(map[string]ActivityEntry)
	for _, e := range entries {
		byID[e.ID] = e
	}

	if byID["a1"].Target != "contact #rec-1" {
		t.Errorf("expected target 'contact #rec-1', got %q", byID["a1"].Target)
	}
	if byID["a2"].Target != "/api/v1/contacts" {
		t.Errorf("expected target '/api/v1/contacts' for request entry, got %q", byID["a2"].Target)
	}
}

func TestService_ListActivity_SystemActor(t *testing.T) {
	db := setupActivityTestDB(t)
	svc := NewService(db)
	ctx := context.Background()

	now := time.Now().UTC()
	db.Exec(`INSERT INTO audit_logs (id, user_id, action, model_name, record_id, created_at, updated_at, active)
		VALUES (?, NULL, 'create', 'contact', 'rec-1', ?, ?, 1)`, "sys1", now, now)

	entries, _, err := svc.ListActivity(ctx, ActivityQueryOptions{Limit: 50})
	if err != nil {
		t.Fatalf("ListActivity failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Actor != "system" {
		t.Errorf("expected actor 'system' for null user_id, got %q", entries[0].Actor)
	}
}

func TestDefaultVerbResolver_AllActions(t *testing.T) {
	cases := []struct {
		action      string
		wantVerb    string
		wantVariant Variant
	}{
		{"create", "created", VariantSuccess},
		{"update", "updated", VariantDefault},
		{"delete", "deleted", VariantDestructive},
		{"login", "logged in", VariantInfo},
		{"logout", "logged out", VariantInfo},
		{"register", "registered", VariantSuccess},
		{"request", "requested", VariantDefault},
		{"custom_action", "custom_action", VariantDefault},
	}
	for _, tc := range cases {
		verb, variant := defaultVerbResolver(tc.action)
		if verb != tc.wantVerb {
			t.Errorf("action %q: expected verb %q, got %q", tc.action, tc.wantVerb, verb)
		}
		if variant != tc.wantVariant {
			t.Errorf("action %q: expected variant %q, got %q", tc.action, tc.wantVariant, variant)
		}
	}
}
