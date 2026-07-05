package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Category groups menu items for browsing and maintenance.
type Category struct {
	ID        int64
	Name      string
	SortOrder int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate ensures the category can be persisted safely.
func (c Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("category name is required")
	}
	if c.SortOrder < 0 {
		return fmt.Errorf("category sort order must be zero or greater")
	}
	return nil
}
