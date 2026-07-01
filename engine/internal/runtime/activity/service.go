package activity

import (
	"context"
	"fmt"
	"time"

	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

const (
	userTable = "res_user"

	maxLimit = 200
	defLimit = 50
)

type Variant string

const (
	VariantDefault     Variant = "default"
	VariantSuccess     Variant = "success"
	VariantWarning     Variant = "warning"
	VariantDestructive Variant = "destructive"
	VariantInfo        Variant = "info"
)

type ActivityEntry struct {
	ID        string  `json:"id"`
	Actor     string  `json:"actor"`
	ActorID   string  `json:"actor_id,omitempty"`
	Action    string  `json:"action"`
	Verb      string  `json:"verb"`
	Target    string  `json:"target,omitempty"`
	Model     string  `json:"model,omitempty"`
	RecordID  string  `json:"record_id,omitempty"`
	Timestamp string  `json:"timestamp"`
	Variant   Variant `json:"variant"`
	Icon      string  `json:"icon,omitempty"`
	IPAddress string  `json:"ip_address,omitempty"`
	Changes   string  `json:"changes,omitempty"`
}

type ActivityQueryOptions struct {
	Model    string
	RecordID string
	UserID   string
	Since    time.Time
	Limit    int
	Offset   int
}

type Service struct {
	db       *gorm.DB
	repo     *persistence.AuditLogRepository
	userTbl  string
	verber   VerbResolver
	renderer ActorRenderer
}

type VerbResolver func(action string) (verb string, variant Variant)
type ActorRenderer func(userID string) string

func NewService(db *gorm.DB) *Service {
	repo := persistence.NewAuditLogRepository(db)
	return &Service{
		db:      db,
		repo:    repo,
		userTbl: userTable,
		verber:  defaultVerbResolver,
		renderer: func(userID string) string {
			return userID
		},
	}
}

func (s *Service) SetUserTable(name string) {
	s.userTbl = name
}

func (s *Service) SetVerbResolver(fn VerbResolver) {
	if fn != nil {
		s.verber = fn
	}
}

func (s *Service) SetActorRenderer(fn ActorRenderer) {
	if fn != nil {
		s.renderer = fn
	}
}

func (s *Service) ListActivity(ctx context.Context, opts ActivityQueryOptions) ([]ActivityEntry, int64, error) {
	if opts.Limit <= 0 {
		opts.Limit = defLimit
	}
	if opts.Limit > maxLimit {
		opts.Limit = maxLimit
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	rows, total, err := s.repo.FindActivity(opts.Model, opts.RecordID, opts.UserID, opts.Since, opts.Limit, opts.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("activity: query audit_log: %w", err)
	}

	userIDs := collectUserIDs(rows)
	nameByID := s.resolveActorNames(ctx, userIDs)

	entries := make([]ActivityEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, s.transform(row, nameByID))
	}
	return entries, total, nil
}

func (s *Service) transform(row map[string]any, nameByID map[string]string) ActivityEntry {
	userID := asString(row["user_id"])
	action := asString(row["action"])
	verb, variant := s.verber(action)

	actor := s.renderActor(userID, nameByID)
	target := buildTarget(row)

	return ActivityEntry{
		ID:        asString(row["id"]),
		Actor:     actor,
		ActorID:   userID,
		Action:    action,
		Verb:      verb,
		Target:    target,
		Model:     asString(row["model_name"]),
		RecordID:  asString(row["record_id"]),
		Timestamp: formatTimestamp(asTime(row["created_at"])),
		Variant:   variant,
		Icon:      iconForVariant(variant),
		IPAddress: asString(row["ip_address"]),
		Changes:   asString(row["changes"]),
	}
}

func (s *Service) renderActor(userID string, nameByID map[string]string) string {
	if userID == "" {
		return "system"
	}
	if name, ok := nameByID[userID]; ok && name != "" {
		return name
	}
	return s.renderer(userID)
}

func (s *Service) resolveActorNames(ctx context.Context, userIDs []string) map[string]string {
	out := make(map[string]string, len(userIDs))
	if len(userIDs) == 0 {
		return out
	}
	var rows []struct {
		ID       string `gorm:"column:id"`
		Username string `gorm:"column:username"`
		Name     string `gorm:"column:name"`
	}
	err := s.db.Table(s.userTbl).
		Select("id, username, name").
		Where("id IN ?", userIDs).
		Find(&rows).Error
	if err != nil {
		return out
	}
	for _, r := range rows {
		display := r.Name
		if display == "" {
			display = r.Username
		}
		if display == "" {
			display = r.ID
		}
		out[r.ID] = display
	}
	return out
}

func collectUserIDs(rows []map[string]any) []string {
	seen := make(map[string]struct{})
	var ids []string
	for _, row := range rows {
		if uid := asString(row["user_id"]); uid != "" {
			if _, ok := seen[uid]; !ok {
				seen[uid] = struct{}{}
				ids = append(ids, uid)
			}
		}
	}
	return ids
}

func buildTarget(row map[string]any) string {
	model := asString(row["model_name"])
	recordID := asString(row["record_id"])
	if model == "" && recordID == "" {
		path := asString(row["request_path"])
		if path != "" {
			return path
		}
		return ""
	}
	if model != "" && recordID != "" {
		return fmt.Sprintf("%s #%s", model, recordID)
	}
	if model != "" {
		return model
	}
	return recordID
}

func defaultVerbResolver(action string) (string, Variant) {
	switch action {
	case "create":
		return "created", VariantSuccess
	case "update":
		return "updated", VariantDefault
	case "delete":
		return "deleted", VariantDestructive
	case "login":
		return "logged in", VariantInfo
	case "logout":
		return "logged out", VariantInfo
	case "register":
		return "registered", VariantSuccess
	case "request":
		return "requested", VariantDefault
	default:
		return action, VariantDefault
	}
}

func iconForVariant(v Variant) string {
	switch v {
	case VariantSuccess:
		return "+"
	case VariantDestructive:
		return "×"
	case VariantWarning:
		return "!"
	case VariantInfo:
		return "i"
	default:
		return "•"
	}
}

func formatTimestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
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
