package comment

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var mentionRegex = regexp.MustCompile(`(?:^|[\s])@([a-zA-Z0-9_.\-]+)`)

type Mention struct {
	Username string `json:"username"`
	UserID   string `json:"user_id,omitempty"`
	Resolved bool   `json:"resolved"`
}

type CommentRecord struct {
	ID              string    `json:"id"`
	Body            string    `json:"body"`
	CommentableType string    `json:"commentable_type"`
	CommentableID   string    `json:"commentable_id"`
	ParentID        *string   `json:"parent_id,omitempty"`
	AuthorID        *string   `json:"author_id,omitempty"`
	Mentions        []Mention `json:"mentions,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WebSocketBroadcaster interface {
	Broadcast(channel string, data any)
}

type Service struct {
	db          *gorm.DB
	tableName   string
	notifTable  string
	userTable   string
	broadcaster WebSocketBroadcaster
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:         db,
		tableName:  "comments",
		notifTable: "notifications",
		userTable:  "res_user",
	}
}

func (s *Service) SetBroadcaster(b WebSocketBroadcaster) {
	s.broadcaster = b
}

func (s *Service) SetTableName(name string) {
	s.tableName = name
}

func (s *Service) SetNotificationTable(name string) {
	s.notifTable = name
}

func (s *Service) SetUserTable(name string) {
	s.userTable = name
}

func ParseMentions(body string) []Mention {
	matches := mentionRegex.FindAllStringSubmatch(body, -1)
	seen := make(map[string]bool)
	var mentions []Mention
	for _, m := range matches {
		username := m[1]
		if seen[username] {
			continue
		}
		seen[username] = true
		mentions = append(mentions, Mention{Username: username})
	}
	return mentions
}

func (s *Service) ResolveMentions(ctx context.Context, mentions []Mention) []Mention {
	if len(mentions) == 0 {
		return mentions
	}
	for i := range mentions {
		var userID string
		err := s.db.Table(s.userTable).
			Select("id").
			Where("username = ? OR email = ?", mentions[i].Username, mentions[i].Username).
			Limit(1).
			Scan(&userID).Error
		if err == nil && userID != "" {
			mentions[i].UserID = userID
			mentions[i].Resolved = true
		}
	}
	return mentions
}

type CreateCommentInput struct {
	Body            string  `json:"body"`
	CommentableType string  `json:"commentable_type"`
	CommentableID   string  `json:"commentable_id"`
	ParentID        *string `json:"parent_id,omitempty"`
	AuthorID        string  `json:"author_id"`
}

func (s *Service) Create(ctx context.Context, input CreateCommentInput) (*CommentRecord, error) {
	if strings.TrimSpace(input.Body) == "" {
		return nil, fmt.Errorf("body is required")
	}
	if input.CommentableType == "" || input.CommentableID == "" {
		return nil, fmt.Errorf("commentable_type and commentable_id are required")
	}

	mentions := s.ResolveMentions(ctx, ParseMentions(input.Body))
	var mentionsJSON any
	if len(mentions) > 0 {
		data, _ := json.Marshal(mentions)
		mentionsJSON = string(data)
	} else {
		mentionsJSON = "[]"
	}

	now := time.Now().UTC()
	record := map[string]any{
		"id":               uuid.New().String(),
		"body":             input.Body,
		"commentable_type": input.CommentableType,
		"commentable_id":   input.CommentableID,
		"parent_id":        nilIfEmptyPtr(input.ParentID),
		"author_id":        nilIfEmpty(input.AuthorID),
		"mentions":         mentionsJSON,
		"created_at":       now,
		"updated_at":       now,
		"active":           true,
	}

	if err := s.db.Table(s.tableName).Create(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	created := s.recordFromMap(record, mentions)

	for _, m := range mentions {
		if m.Resolved {
			s.createNotification(ctx, m.UserID, input.CommentableType, input.CommentableID, input.Body, created.ID)
		}
	}

	if s.broadcaster != nil {
		s.broadcaster.Broadcast("comment:"+input.CommentableType+":"+input.CommentableID, map[string]any{
			"type": "comment.created",
			"data": created,
		})
	}

	return &created, nil
}

func (s *Service) createNotification(ctx context.Context, userID, commentableType, commentableID, body, commentID string) {
	now := time.Now().UTC()
	excerpt := body
	if len(excerpt) > 200 {
		excerpt = excerpt[:200]
	}
	notif := map[string]any{
		"id":          uuid.New().String(),
		"user_id":     userID,
		"type":        "mention",
		"title":       "You were mentioned",
		"body":        excerpt,
		"model_name":  commentableType,
		"record_id":   commentableID,
		"comment_id":  commentID,
		"read":        false,
		"created_at":  now,
		"updated_at":  now,
		"active":      true,
	}
	if err := s.db.Table(s.notifTable).Create(&notif).Error; err != nil {
		return
	}
	if s.broadcaster != nil {
		s.broadcaster.Broadcast("notification:user:"+userID, map[string]any{
			"type": "notification.created",
			"data": notif,
		})
	}
}

func (s *Service) ListByRecord(ctx context.Context, commentableType, commentableID string, page, pageSize int) ([]CommentRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	var total int64
	countErr := s.db.Table(s.tableName).
		Where("commentable_type = ? AND commentable_id = ? AND active = ?", commentableType, commentableID, true).
		Count(&total).Error
	if countErr != nil {
		return nil, 0, countErr
	}

	var rows []map[string]any
	err := s.db.Table(s.tableName).
		Where("commentable_type = ? AND commentable_id = ? AND active = ?", commentableType, commentableID, true).
		Order("created_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	records := make([]CommentRecord, 0, len(rows))
	for _, row := range rows {
		records = append(records, s.recordFromMap(row, nil))
	}
	return records, total, nil
}

func (s *Service) recordFromMap(row map[string]any, preResolvedMentions []Mention) CommentRecord {
	rec := CommentRecord{
		ID:              asString(row["id"]),
		Body:            asString(row["body"]),
		CommentableType: asString(row["commentable_type"]),
		CommentableID:   asString(row["commentable_id"]),
		CreatedAt:       asTime(row["created_at"]),
		UpdatedAt:       asTime(row["updated_at"]),
	}
	if v := asString(row["parent_id"]); v != "" {
		rec.ParentID = &v
	}
	if v := asString(row["author_id"]); v != "" {
		rec.AuthorID = &v
	}
	if preResolvedMentions != nil {
		rec.Mentions = preResolvedMentions
	} else if raw := asString(row["mentions"]); raw != "" && raw != "[]" {
		var ms []Mention
		if err := json.Unmarshal([]byte(raw), &ms); err == nil {
			rec.Mentions = ms
		}
	}
	return rec
}

func (s *Service) EnsureTables() error {
	if !s.db.Migrator().HasTable(s.notifTable) {
		if err := s.createNotificationTable(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) createNotificationTable() error {
	dialect := s.db.Dialector.Name()
	var typeCol string
	switch dialect {
	case "postgres":
		typeCol = "VARCHAR(50)"
	case "mysql":
		typeCol = "VARCHAR(50)"
	default:
		typeCol = "TEXT"
	}
	sqlStr := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		type %s DEFAULT 'mention',
		title TEXT,
		body TEXT,
		model_name TEXT,
		record_id TEXT,
		comment_id TEXT,
		read INTEGER DEFAULT 0,
		active INTEGER DEFAULT 1,
		created_at DATETIME,
		updated_at DATETIME
	)`, s.notifTable, typeCol)
	return s.db.Exec(sqlStr).Error
}

func (s *Service) GetByID(ctx context.Context, id string) (*CommentRecord, error) {
	var row map[string]any
	err := s.db.Table(s.tableName).Where("id = ? AND active = ?", id, true).First(&row).Error
	if err != nil {
		return nil, err
	}
	rec := s.recordFromMap(row, nil)
	return &rec, nil
}

func (s *Service) Delete(ctx context.Context, id, userID string) error {
	result := s.db.Table(s.tableName).
		Where("id = ? AND (author_id = ? OR ? = '')", id, userID, userID).
		Update("active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nilIfEmptyPtr(s *string) any {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func asTime(v any) time.Time {
	if v == nil {
		return time.Time{}
	}
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		parsed, err := time.Parse(time.RFC3339, t)
		if err == nil {
			return parsed
		}
		return time.Time{}
	default:
		return time.Time{}
	}
}
