package httptransport

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
	"github.com/bboykiv/topsigner/internal/service/image"
	"github.com/bboykiv/topsigner/internal/transport/httptransport/middleware"
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
	options := httpserver.ChiServerOptions{
		Middlewares: []httpserver.MiddlewareFunc{
			middleware.NoCache(),
			middleware.Recoverer(),
			middleware.IP(),
			middleware.UserAgent(),
			middleware.Cors(config),
			middleware.RequestLogger(logger),
			middleware.BearerAuth(authService),
		},
	}

	server := &strictServer{
		AuthHandler:  NewAuthHandler(authService),
		ImageHandler: NewImageHandler(imageService),
		FontHandler:  NewFontHandler(fontService),
	}

	return httpserver.HandlerWithOptions(httpserver.NewStrictHandler(server, nil), options)
}
