package application

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/adapter/crypto"
	"github.com/bboykiv/topsigner/internal/adapter/postgres"
	"github.com/bboykiv/topsigner/internal/adapter/redis"
	"github.com/bboykiv/topsigner/internal/adapter/s3"
	"github.com/bboykiv/topsigner/internal/adapter/vk"
	"github.com/bboykiv/topsigner/internal/adapter/vkid"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
	"github.com/bboykiv/topsigner/internal/service/group"
	"github.com/bboykiv/topsigner/internal/service/health"
	"github.com/bboykiv/topsigner/internal/service/image"
	"github.com/bboykiv/topsigner/internal/service/user"
	"github.com/bboykiv/topsigner/internal/transport/httptransport"
)

func New() fx.Option {
	return fx.Options(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: logger.WithOptions(
					zap.AddCallerSkip(1),
					zap.IncreaseLevel(zap.ErrorLevel),
				),
			}
		}),
		postgres.Module,
		redis.Module,
		s3.Module,
		vk.Module,
		vkid.Module,
		crypto.Module,
		fx.Provide(
			NewLogger,
			config.New,
			auth.New,
			user.New,
			font.New,
			image.New,
			group.New,
			health.New,
			NewHTTPServer,
			httptransport.NewHandler,
		),
		fx.Invoke(func(userService *user.Service) error {
			return userService.CreateDefault(context.Background())
		}),
		fx.Invoke(func(*http.Server) {}),
	)
}
