package comment

import (
	"context"
	"testing"
)

func TestParseMentions(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    []string
	}{
		{"no mentions", "Hello world", nil},
		{"single mention", "Hi @alice", []string{"alice"}},
		{"multiple mentions", "Hi @alice and @bob", []string{"alice", "bob"}},
		{"duplicate mentions deduped", "Hi @alice @alice", []string{"alice"}},
		{"mention at string boundary", "@alice", []string{"alice"}},
		{"mention followed by punctuation", "Hello @alice!", []string{"alice"}},
		{"mention with dots and dashes", "Hi @alice.doe and @bob-smith", []string{"alice.doe", "bob-smith"}},
		{"email not parsed as mention", "Send to alice@example.com", nil},
		{"empty body", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mentions := ParseMentions(tt.body)
			if len(mentions) != len(tt.want) {
				t.Fatalf("expected %d mentions, got %d: %v", len(tt.want), len(mentions), mentions)
			}
			seen := make(map[string]bool)
			for _, m := range mentions {
				seen[m.Username] = true
			}
			for _, expected := range tt.want {
				if !seen[expected] {
					t.Errorf("expected mention %q not found", expected)
				}
			}
		})
	}
}

func TestParseMentions_EdgeCases(t *testing.T) {
	mentions := ParseMentions("Thanks @user1 @user2 @user1 @user3")
	if len(mentions) != 3 {
		t.Fatalf("expected 3 unique mentions, got %d", len(mentions))
	}
	usernames := make(map[string]bool)
	for _, m := range mentions {
		usernames[m.Username] = true
		if m.Resolved {
			t.Errorf("mention %q should not be resolved before ResolveMentions", m.Username)
		}
	}
	for _, expected := range []string{"user1", "user2", "user3"} {
		if !usernames[expected] {
			t.Errorf("expected mention %q", expected)
		}
	}
}

func TestResolveMentions_NoUsers(t *testing.T) {
	db := setupCommentTestDB(t)
	svc := NewService(db)

	mentions := ParseMentions("Hi @nonexistent")
	resolved := svc.ResolveMentions(context.Background(), mentions)

	if len(resolved) != 1 {
		t.Fatalf("expected 1 mention, got %d", len(resolved))
	}
	if resolved[0].Resolved {
		t.Error("mention should not be resolved for non-existent user")
	}
	if resolved[0].UserID != "" {
		t.Errorf("expected empty UserID, got %q", resolved[0].UserID)
	}
}
