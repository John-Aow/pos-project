package usecase

import (
	"context"

	"pos-project/backend/entity"
)

// AuditLogRepository records manager menu maintenance audit events.
type AuditLogRepository interface {
	RecordMenuAuditLog(ctx context.Context, log entity.MenuAuditLog) error
}
