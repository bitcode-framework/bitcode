package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityQueryFilter struct {
	Model    string
	RecordID string
	UserID   string
	Since    time.Time
	Limit    int
	Offset   int
}

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func AutoMigrateAuditLog(db *gorm.DB) error {
	if !db.Migrator().HasTable("audit_logs") {
		return nil
	}
	dialect := DetectDialect(db)
	switch dialect {
	case DialectPostgres:
		db.Exec("ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS impersonated_by TEXT")
	case DialectMySQL:
		var count int64
		db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name='audit_logs' AND column_name='impersonated_by'").Scan(&count)
		if count == 0 {
			db.Exec("ALTER TABLE audit_logs ADD COLUMN impersonated_by TEXT")
		}
	default:
		db.Exec("ALTER TABLE audit_logs ADD COLUMN impersonated_by TEXT")
	}
	return nil
}

type AuditLogEntry struct {
	UserID         string
	Action         string
	ModelName      string
	RecordID       string
	Changes        string
	IPAddress      string
	UserAgent      string
	RequestMethod  string
	RequestPath    string
	StatusCode     int
	DurationMs     int
	ImpersonatedBy string
}

func (r *AuditLogRepository) Write(entry AuditLogEntry) error {
	record := map[string]any{
		"id":              uuid.New().String(),
		"user_id":         nilIfEmpty(entry.UserID),
		"action":          entry.Action,
		"model_name":      nilIfEmpty(entry.ModelName),
		"record_id":       nilIfEmpty(entry.RecordID),
		"changes":         nilIfEmpty(entry.Changes),
		"ip_address":      nilIfEmpty(entry.IPAddress),
		"user_agent":      nilIfEmpty(entry.UserAgent),
		"request_method":  nilIfEmpty(entry.RequestMethod),
		"request_path":    nilIfEmpty(entry.RequestPath),
		"status_code":     nilIfZero(entry.StatusCode),
		"duration_ms":     nilIfZero(entry.DurationMs),
		"impersonated_by": nilIfEmpty(entry.ImpersonatedBy),
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
		"active":          true,
	}

	return r.db.Table("audit_logs").Create(&record).Error
}

func (r *AuditLogRepository) WriteAsync(entry AuditLogEntry) {
	go func() {
		r.Write(entry)
	}()
}

func (r *AuditLogRepository) FindByRecord(modelName, recordID string, limit int) ([]map[string]any, error) {
	var results []map[string]any
	q := r.db.Table("audit_logs").
		Where("model_name = ? AND record_id = ?", modelName, recordID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		results = []map[string]any{}
	}
	return results, nil
}

func (r *AuditLogRepository) FindByAction(action string, limit int) ([]map[string]any, error) {
	var results []map[string]any
	q := r.db.Table("audit_logs").
		Where("action = ?", action).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		results = []map[string]any{}
	}
	return results, nil
}

func (r *AuditLogRepository) FindByUser(userID string, limit int) ([]map[string]any, error) {
	var results []map[string]any
	q := r.db.Table("audit_logs").
		Where("user_id = ?", userID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		results = []map[string]any{}
	}
	return results, nil
}

func (r *AuditLogRepository) FindLoginHistory(limit int) ([]map[string]any, error) {
	var results []map[string]any
	q := r.db.Table("audit_logs").
		Where("action IN ?", []string{"login", "logout", "register"}).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		results = []map[string]any{}
	}
	return results, nil
}

func (r *AuditLogRepository) FindRequests(limit int, methodFilter string) ([]map[string]any, error) {
	var results []map[string]any
	q := r.db.Table("audit_logs").
		Where("action = ?", "request").
		Order("created_at DESC")
	if methodFilter != "" {
		q = q.Where("request_method = ?", methodFilter)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		results = []map[string]any{}
	}
	return results, nil
}

// FindActivity queries audit_logs with flexible filtering for the activity
// feed. All filters are optional — pass empty values to skip them. When all
// filters are empty a global timeline is returned. Results are ordered
// newest-first. Returns the total count (ignoring limit/offset) alongside the
// rows so callers can paginate without a second round-trip.
func (r *AuditLogRepository) FindActivity(model, recordID, userID string, since time.Time, limit, offset int) ([]map[string]any, int64, error) {
	q := r.db.Table("audit_logs")
	if model != "" {
		q = q.Where("model_name = ?", model)
	}
	if recordID != "" {
		q = q.Where("record_id = ?", recordID)
	}
	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	if !since.IsZero() {
		q = q.Where("created_at >= ?", since)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if offset > 0 {
		q = q.Offset(offset)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}

	var results []map[string]any
	if err := q.Order("created_at DESC").Find(&results).Error; err != nil {
		return nil, 0, err
	}
	if results == nil {
		results = []map[string]any{}
	}
	return results, total, nil
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nilIfZero(n int) any {
	if n == 0 {
		return nil
	}
	return n
}
