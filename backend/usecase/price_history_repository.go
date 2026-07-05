package usecase

import (
	"context"

	"pos-project/backend/entity"
)

// PriceHistoryRepository records immutable menu price changes.
type PriceHistoryRepository interface {
	RecordPriceHistory(ctx context.Context, history entity.PriceHistory) error
}
