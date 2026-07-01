package comment

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
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
		data TEXT,
		notifiable_type TEXT,
		notifiable_id TEXT,
		read_at DATETIME,
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
	return db
}

func TestResolveMentions_ResolvesRealUser(t *testing.T) {
	db := setupCommentTestDB(t)
	sqlDB, _ := db.DB()
	sqlDB.Exec(`INSERT INTO res_user (id, username) VALUES ('u-1', 'alice')`)

	svc := NewService(db)
	mentions := ParseMentions("Hi @alice and @nobody")
	resolved := svc.ResolveMentions(context.Background(), mentions)

	if len(resolved) != 2 {
		t.Fatalf("expected 2 mentions, got %d", len(resolved))
	}

	aliceFound := false
	nobodyFound := false
	for _, m := range resolved {
		if m.Username == "alice" {
			aliceFound = true
			if !m.Resolved {
				t.Error("alice should be resolved")
			}
			if m.UserID != "u-1" {
				t.Errorf("expected alice UserID u-1, got %q", m.UserID)
			}
		}
		if m.Username == "nobody" {
			nobodyFound = true
			if m.Resolved {
				t.Error("nobody should NOT be resolved")
			}
		}
	}
	if !aliceFound {
		t.Error("alice mention not found")
	}
	if !nobodyFound {
		t.Error("nobody mention not found")
	}
}

func TestCreate_ParsesMentionsAndCreatesNotification(t *testing.T) {
	db := setupCommentTestDB(t)
	sqlDB, _ := db.DB()
	sqlDB.Exec(`INSERT INTO res_user (id, username) VALUES ('u-1', 'alice')`)
	sqlDB.Exec(`INSERT INTO res_user (id, username) VALUES ('u-2', 'bob')`)

	svc := NewService(db)
	created, err := svc.Create(context.Background(), CreateCommentInput{
		Body:            "Hey @alice check this out",
		CommentableType: "crm.contact",
		CommentableID:   "rec-1",
		AuthorID:        "u-2",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID == "" {
		t.Error("expected non-empty comment ID")
	}
	if len(created.Mentions) != 1 {
		t.Fatalf("expected 1 mention, got %d", len(created.Mentions))
	}
	if created.Mentions[0].Username != "alice" {
		t.Errorf("expected mention alice, got %q", created.Mentions[0].Username)
	}
	if !created.Mentions[0].Resolved {
		t.Error("alice mention should be resolved")
	}

	var notifCount int64
	db.Table("notifications").Where("user_id = ? AND notifiable_id = ?", "u-1", created.CommentableID).Count(&notifCount)
	if notifCount != 1 {
		t.Errorf("expected 1 notification for alice, got %d", notifCount)
	}
}

func TestListByRecord_ReturnsOnlyMatchingComments(t *testing.T) {
	db := setupCommentTestDB(t)
	svc := NewService(db)

	ctx := context.Background()
	_, _ = svc.Create(ctx, CreateCommentInput{
		Body:            "Comment on rec-1",
		CommentableType: "crm.contact",
		CommentableID:   "rec-1",
		AuthorID:        "u-1",
	})
	_, _ = svc.Create(ctx, CreateCommentInput{
		Body:            "Comment on rec-2",
		CommentableType: "crm.contact",
		CommentableID:   "rec-2",
		AuthorID:        "u-1",
	})

	records, total, err := svc.ListByRecord(ctx, "crm.contact", "rec-1", 1, 50)
	if err != nil {
		t.Fatalf("ListByRecord failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].CommentableID != "rec-1" {
		t.Errorf("expected rec-1, got %q", records[0].CommentableID)
	}
}
