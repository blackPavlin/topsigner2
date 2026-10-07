package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
)

func NewHTTPServer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	config *config.Config,
	handler http.Handler,
) *http.Server {
	addr := fmt.Sprintf(":%d", config.HTTP.Port)

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       config.HTTP.ReadTimeout,
		ReadHeaderTimeout: config.HTTP.ReadHeaderTimeout,
		WriteTimeout:      config.HTTP.WriteTimeout,
		IdleTimeout:       config.HTTP.IdleTimeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("starting http server", zap.String("addr", addr))

			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("start http server error", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed shutdown http server: %w", err)
			}

			logger.Info("stopping http server")

			return nil
		},
	})

	return server
}
