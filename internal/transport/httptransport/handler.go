package httptransport

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
	"github.com/bboykiv/topsigner/internal/service/image"
	mw "github.com/bboykiv/topsigner/internal/transport/httptransport/middleware"
)

var _ httpserver.StrictServerInterface = (*strictServer)(nil)

type strictServer struct {
	*AuthHandler
	*ImageHandler
	*FontHandler
}

func NewHandler(
	logger *zap.Logger,
	config *config.Config,
	authService *auth.Service,
	imageService *image.Service,
	fontService *font.Service,
) http.Handler {
	corsOptions := cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedMethods:   config.Cors.AllowedMethods,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		ExposedHeaders:   config.Cors.ExposedHeaders,
		AllowCredentials: config.Cors.AllowCredentials,
		MaxAge:           config.Cors.MaxAge,
	}

	options := httpserver.ChiServerOptions{
		Middlewares: []httpserver.MiddlewareFunc{
			middleware.NoCache,
			middleware.RequestID,
			middleware.Recoverer,
			// todo: ClientIPFromRemoteAddr не будет корректно работать если сервер стоит за reverse proxy
			middleware.ClientIPFromRemoteAddr,
			cors.Handler(corsOptions),
			mw.RequestLogger(logger),
			mw.BearerAuth(authService),
			mw.UserAgent(),
		},
	}

	server := &strictServer{
		AuthHandler:  NewAuthHandler(authService),
		ImageHandler: NewImageHandler(imageService),
		FontHandler:  NewFontHandler(fontService),
	}

	return httpserver.HandlerWithOptions(httpserver.NewStrictHandler(server, nil), options)
}
