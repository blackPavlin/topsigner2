package httptransport

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/gen/openapi"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/image"
	mw "github.com/bboykiv/topsigner/internal/transport/httptransport/middleware"
)

var _ openapi.StrictServerInterface = (*strictServer)(nil)

type strictServer struct {
	*ImageHandler
	*FontHandler
}

func NewHandler(
	logger *zap.Logger,
	config *config.Config,
	imageService *image.Service,
) http.Handler {
	corsOptions := cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedMethods:   config.Cors.AllowedMethods,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		ExposedHeaders:   config.Cors.ExposedHeaders,
		AllowCredentials: config.Cors.AllowCredentials,
		MaxAge:           config.Cors.MaxAge,
	}

	options := openapi.ChiServerOptions{
		Middlewares: []openapi.MiddlewareFunc{
			cors.Handler(corsOptions),
			mw.RequestLogger(logger),
			middleware.NoCache,
			middleware.RequestID,
			middleware.Recoverer,
		},
	}

	server := &strictServer{
		ImageHandler: NewImageHandler(imageService),
		FontHandler:  NewFontHandler(),
	}

	return openapi.HandlerWithOptions(openapi.NewStrictHandler(server, nil), options)
}
