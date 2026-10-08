package health

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type Service struct {
	logger   *zap.Logger
	database DatabaseChecker
	storage  StorageChecker
	cache    CacheChecker
}

func New(
	logger *zap.Logger,
	database DatabaseChecker,
	storage StorageChecker,
	cache CacheChecker,
) *Service {
	return &Service{
		logger:   logger.Named("health-service"),
		database: database,
		storage:  storage,
		cache:    cache,
	}
}

func (s *Service) Check(ctx context.Context) *Report {
	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	var (
		wg                                sync.WaitGroup
		databaseErr, storageErr, cacheErr error
	)

	wg.Go(func() { databaseErr = s.database.Ping(ctx) })

	wg.Go(func() { storageErr = s.storage.Ping(ctx) })

	wg.Go(func() { cacheErr = s.cache.Ping(ctx) })

	wg.Wait()

	if databaseErr != nil {
		return &Report{Status: StatusDown}
	}

	if storageErr != nil || cacheErr != nil {
		return &Report{Status: StatusDegraded}
	}

	return &Report{Status: StausOK}
}
