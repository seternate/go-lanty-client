package game

import "context"

type ExtensionRepository interface {
	ListBySlug(ctx context.Context, slug string) ([]Extension, error)
	Store(ctx context.Context, ext Extension) error
}
