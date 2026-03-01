package game

import "context"

type InstallationRepository interface {
	GetAll(ctx context.Context) ([]Installation, error)
	GetBySlug(ctx context.Context, slug string) (Installation, error)
	Store(ctx context.Context, installation Installation) error
	Remove(ctx context.Context, slug string) error
}
