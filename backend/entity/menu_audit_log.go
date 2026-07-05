package entity

import (
	"errors"
	"strings"
	"time"
)

// MenuAuditLog records manager menu, price, stock, and threshold actions.
type MenuAuditLog struct {
	ID            int64
	MenuItemID    int64
	ActionType    string
	PreviousValue string
	NewValue      string
	Actor         string
	Reason        string
	CreatedAt     time.Time
}

// Validate ensures sensitive manager actions remain traceable.
func (l MenuAuditLog) Validate() error {
	if l.MenuItemID <= 0 {
		return errors.New("menu audit log menu item id is required")
	}
	if strings.TrimSpace(l.ActionType) == "" {
		return errors.New("menu audit log action type is required")
	}
	if strings.TrimSpace(l.Actor) == "" {
		return errors.New("menu audit log actor is required")
	}
	if l.CreatedAt.IsZero() {
		return errors.New("menu audit log created_at is required")
	}
	return nil
}
