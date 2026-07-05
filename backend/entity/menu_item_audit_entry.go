package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// Menu item audit actions used by manager maintenance features.
	AuditActionCreate     = "create"
	AuditActionUpdate     = "update"
	AuditActionDeactivate = "deactivate"
	AuditActionPrice      = "price"
	AuditActionStock      = "stock"
	AuditActionThreshold  = "threshold"
)

// MenuItemAuditEntry records a traced manager action.
type MenuItemAuditEntry struct {
	ID            int64
	MenuItemID    int64
	ActionType    string
	PreviousValue string
	NewValue      string
	Actor         string
	Reason        string
	CreatedAt     time.Time
}

// Validate ensures traceability data is present.
func (e MenuItemAuditEntry) Validate() error {
	if e.MenuItemID <= 0 {
		return errors.New("audit entry menu item id is required")
	}
	if strings.TrimSpace(e.ActionType) == "" {
		return errors.New("audit entry action type is required")
	}
	if strings.TrimSpace(e.Actor) == "" {
		return errors.New("audit entry actor is required")
	}
	if e.CreatedAt.IsZero() {
		return fmt.Errorf("audit entry created_at is required")
	}
	return nil
}
