package application

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/database"
	"github.com/bboykiv/topsigner/internal/database/repository"
	"github.com/bboykiv/topsigner/internal/keyvalue"
	"github.com/bboykiv/topsigner/internal/s3"
	"github.com/bboykiv/topsigner/internal/s3/storage"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
	"github.com/bboykiv/topsigner/internal/service/image"
	"github.com/bboykiv/topsigner/internal/transport/httptransport"
	"github.com/bboykiv/topsigner/internal/vkid"
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
		database.Module,
		keyvalue.Module,
		s3.Module,
		vkid.Module,
		fx.Provide(
			NewLogger,
			config.New,
			auth.New,
			font.New,
			image.New,
			fx.Annotate(
				repository.NewFontRepository,
				fx.As(new(font.Repository)),
			),
			fx.Annotate(
				repository.NewImageRepository,
				fx.As(new(image.Repository)),
			),
			fx.Annotate(
				storage.NewImageStorage,
				fx.As(new(image.Storage)),
			),
			fx.Annotate(
				repository.NewUserRepository,
				fx.As(new(auth.UserRepository)),
			),
			fx.Annotate(
				repository.NewSessionRepository,
				fx.As(new(auth.SessionRepository)),
			),
			fx.Annotate(
				keyvalue.NewCodeVerifierRepository,
				fx.As(new(auth.CodeVerifierRepository)),
			),
			NewHttpServer,
			httptransport.NewHandler,
		),
		fx.Invoke(func(*http.Server) {}),
	)
}
