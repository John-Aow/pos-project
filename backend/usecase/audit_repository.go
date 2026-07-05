package usecase

import (
	"context"

	"pos-project/backend/entity"
)

// AuditRepository records manager actions for traceability.
type AuditRepository interface {
	RecordMenuItemAuditEntry(ctx context.Context, entry entity.MenuItemAuditEntry) error
}
