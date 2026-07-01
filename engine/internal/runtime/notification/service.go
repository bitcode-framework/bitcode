package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string

const (
	TypeInfo    NotificationType = "info"
	TypeWarning NotificationType = "warning"
	TypeError   NotificationType = "error"
	TypeSuccess NotificationType = "success"
	TypeMention NotificationType = "mention"
	TypeMessage NotificationType = "message"
)

func ValidType(t string) bool {
	switch NotificationType(t) {
	case TypeInfo, TypeWarning, TypeError, TypeSuccess, TypeMention, TypeMessage:
		return true
	}
	return false
}

type NotificationRecord struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	Type            string     `json:"type"`
	Title           string     `json:"title"`
	Body            string     `json:"body"`
	Data            any        `json:"data,omitempty"`
	NotifiableType  string     `json:"notifiable_type,omitempty"`
	NotifiableID    string     `json:"notifiable_id,omitempty"`
	ReadAt          *time.Time `json:"read_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type WebSocketBroadcaster interface {
	Broadcast(channel string, data any)
	BroadcastToUser(userID string, channel string, data any)
}

type noopBroadcaster struct{}

func (noopBroadcaster) Broadcast(string, any)            {}
func (noopBroadcaster) BroadcastToUser(string, string, any) {}

type Service struct {
	db          *gorm.DB
	tableName   string
	userTable   string
	broadcaster WebSocketBroadcaster
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:          db,
		tableName:   "notifications",
		userTable:   "res_user",
		broadcaster: noopBroadcaster{},
	}
}

func (s *Service) SetBroadcaster(b WebSocketBroadcaster) {
	if b == nil {
		s.broadcaster = noopBroadcaster{}
		return
	}
	s.broadcaster = b
}

func (s *Service) SetTableName(name string) {
	s.tableName = name
}

func (s *Service) SetUserTable(name string) {
	s.userTable = name
}

func (s *Service) TableName() string {
	return s.tableName
}

type CreateInput struct {
	UserID         string
	Type           string
	Title          string
	Body           string
	Data           any
	NotifiableType string
	NotifiableID   string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*NotificationRecord, error) {
	if input.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	notifType := input.Type
	if notifType == "" {
		notifType = string(TypeInfo)
	}
	if !ValidType(notifType) {
		return nil, fmt.Errorf("invalid notification type: %s", notifType)
	}

	var dataJSON any
	if input.Data != nil {
		raw, err := json.Marshal(input.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data: %w", err)
		}
		dataJSON = string(raw)
	} else {
		dataJSON = "{}"
	}

	now := time.Now().UTC()
	record := map[string]any{
		"id":              uuid.New().String(),
		"user_id":         input.UserID,
		"type":            notifType,
		"title":           input.Title,
		"body":            input.Body,
		"data":            dataJSON,
		"notifiable_type": nilIfEmpty(input.NotifiableType),
		"notifiable_id":   nilIfEmpty(input.NotifiableID),
		"read_at":         nil,
		"created_at":      now,
		"updated_at":      now,
		"active":          true,
	}

	if err := s.db.Table(s.tableName).Create(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	created := s.recordFromMap(record)

	s.broadcaster.BroadcastToUser(input.UserID, "notification:user:"+input.UserID, map[string]any{
		"type": "notification.created",
		"data": created,
	})
	s.broadcaster.Broadcast("notifications", map[string]any{
		"type":    "notification.created",
		"user_id": input.UserID,
		"data":    created,
	})

	return &created, nil
}

type ListOptions struct {
	UserID   string
	Page     int
	PageSize int
	UnreadOnly bool
}

type ListResult struct {
	Records   []NotificationRecord `json:"data"`
	Total     int64                `json:"total"`
	Page      int                  `json:"page"`
	PageSize  int                  `json:"page_size"`
	UnreadCount int64              `json:"unread_count"`
}

func (s *Service) List(ctx context.Context, opts ListOptions) (*ListResult, error) {
	if opts.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 200 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	baseQuery := s.db.Table(s.tableName).Where("user_id = ? AND active = ?", opts.UserID, true)
	if opts.UnreadOnly {
		baseQuery = baseQuery.Where("read_at IS NULL")
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count notifications: %w", err)
	}

	var unreadCount int64
	countErr := s.db.Table(s.tableName).
		Where("user_id = ? AND active = ? AND read_at IS NULL", opts.UserID, true).
		Count(&unreadCount).Error
	if countErr != nil {
		unreadCount = 0
	}

	var rows []map[string]any
	findQuery := s.db.Table(s.tableName).Where("user_id = ? AND active = ?", opts.UserID, true)
	if opts.UnreadOnly {
		findQuery = findQuery.Where("read_at IS NULL")
	}
	err := findQuery.
		Order("CASE WHEN read_at IS NULL THEN 0 ELSE 1 END, created_at DESC").
		Offset(offset).
		Limit(opts.PageSize).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	records := make([]NotificationRecord, 0, len(rows))
	for _, row := range rows {
		records = append(records, s.recordFromMap(row))
	}

	return &ListResult{
		Records:     records,
		Total:       total,
		Page:        opts.Page,
		PageSize:    opts.PageSize,
		UnreadCount: unreadCount,
	}, nil
}

func (s *Service) UnreadCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("user_id is required")
	}
	var count int64
	err := s.db.Table(s.tableName).
		Where("user_id = ? AND active = ? AND read_at IS NULL", userID, true).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}

func (s *Service) MarkRead(ctx context.Context, userID, notificationID string) error {
	if userID == "" || notificationID == "" {
		return fmt.Errorf("user_id and notification id are required")
	}
	now := time.Now().UTC()
	result := s.db.Table(s.tableName).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("read_at", now)
	if result.Error != nil {
		return fmt.Errorf("failed to mark notification read: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found or not owned by user")
	}

	s.broadcaster.BroadcastToUser(userID, "notification:user:"+userID, map[string]any{
		"type":             "notification.read",
		"notification_id":  notificationID,
	})
	return nil
}

func (s *Service) MarkAllRead(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("user_id is required")
	}
	now := time.Now().UTC()
	result := s.db.Table(s.tableName).
		Where("user_id = ? AND active = ? AND read_at IS NULL", userID, true).
		Updates(map[string]any{
			"read_at":    now,
			"updated_at": now,
		})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to mark all notifications read: %w", result.Error)
	}

	s.broadcaster.BroadcastToUser(userID, "notification:user:"+userID, map[string]any{
		"type": "notification.read_all",
		"count": result.RowsAffected,
	})
	return result.RowsAffected, nil
}

func (s *Service) recordFromMap(row map[string]any) NotificationRecord {
	rec := NotificationRecord{
		ID:             asString(row["id"]),
		UserID:         asString(row["user_id"]),
		Type:           asString(row["type"]),
		Title:          asString(row["title"]),
		Body:           asString(row["body"]),
		NotifiableType: asString(row["notifiable_type"]),
		NotifiableID:   asString(row["notifiable_id"]),
		CreatedAt:      asTime(row["created_at"]),
		UpdatedAt:      asTime(row["updated_at"]),
	}
	if raw := asString(row["data"]); raw != "" && raw != "{}" {
		var parsed any
		if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
			rec.Data = parsed
		} else {
			rec.Data = raw
		}
	}
	if v := row["read_at"]; v != nil {
		switch t := v.(type) {
		case time.Time:
			if !t.IsZero() {
				rec.ReadAt = &t
			}
		case string:
			if t != "" {
				parsed, err := time.Parse(time.RFC3339, t)
				if err == nil && !parsed.IsZero() {
					rec.ReadAt = &parsed
				}
			}
		}
	}
	return rec
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
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
