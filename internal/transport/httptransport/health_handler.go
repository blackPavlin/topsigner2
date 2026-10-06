package httptransport

import (
	"context"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/service/health"
)

type HealthHandler struct {
	service *health.Service
}

func NewHealthHandler(service *health.Service) *HealthHandler {
	return &HealthHandler{service: service}
}

// HealthCheck Health check
// (GET /api/v1/health)
func (h *HealthHandler) HealthCheck(
	ctx context.Context,
	request httpserver.HealthCheckRequestObject,
) (httpserver.HealthCheckResponseObject, error) {
	return httpserver.HealthCheck200JSONResponse{}, nil
}
