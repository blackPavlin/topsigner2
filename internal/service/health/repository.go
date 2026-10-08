package health

//go:generate go tool mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock -typed

import "context"

type DatabaseChecker interface {
	Ping(ctx context.Context) error
}

type StorageChecker interface {
	Ping(ctx context.Context) error
}

type CacheChecker interface {
	Ping(ctx context.Context) error
}
