package repositories

import (
	"context"
)

// Storer store contract
type Storer interface {
	Store(ctx context.Context, param any) (int64, error)
}

// Updater update contract
type Updater interface {
	Update(ctx context.Context, input any, where any) (int64, error)
}

// Deleter delete contract
type Deleter interface {
	Delete(ctx context.Context, param any) (int64, error)
}

// Counter count contract
type Counter interface {
	Count(ctx context.Context, p any) (total int64, err error)
}

