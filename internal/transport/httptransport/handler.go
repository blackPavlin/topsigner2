package httptransport

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
	"github.com/bboykiv/topsigner/internal/service/group"
	"github.com/bboykiv/topsigner/internal/service/health"
	"github.com/bboykiv/topsigner/internal/service/image"
	"github.com/bboykiv/topsigner/internal/transport/httptransport/middleware"
)

var _ httpserver.StrictServerInterface = (*strictServer)(nil)

type strictServer struct {
	*AuthHandler
	*ImageHandler
	*FontHandler
	*GroupHandler
	*HealthHandler
}

func NewHandler(
	logger *zap.Logger,
	config *config.Config,
	authService *auth.Service,
	imageService *image.Service,
	fontService *font.Service,
	groupService *group.Service,
	healthService *health.Service,
) http.Handler {
	middlewares := []httpserver.MiddlewareFunc{
		middleware.RequestID(),
		middleware.RequestLogger(logger),
		middleware.Recoverer(),
		middleware.Cors(config),
		middleware.NoCache(),
		middleware.IP(),
		middleware.UserAgent(),
		middleware.BearerAuth(authService),
	}

	options := httpserver.ChiServerOptions{
		Middlewares:      middlewares,
		ErrorHandlerFunc: requestErrorHandlerFunc,
	}

	server := &strictServer{
		AuthHandler:   NewAuthHandler(authService),
		ImageHandler:  NewImageHandler(imageService),
		FontHandler:   NewFontHandler(fontService),
		GroupHandler:  NewGroupHandler(groupService),
		HealthHandler: NewHealthHandler(healthService),
	}

	return httpserver.HandlerWithOptions(
		httpserver.NewStrictHandlerWithOptions(server, nil, httpserver.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  requestErrorHandlerFunc,
			ResponseErrorHandlerFunc: responseErrorHandlerFunc,
		}),
		options,
	)
}
