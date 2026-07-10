package httptransport

import (
	"context"

	"github.com/bboykiv/topsigner/gen/openapi"
	"github.com/bboykiv/topsigner/internal/service/auth"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Authenticate user
// (POST /api/v1/auth/login)
func (h *AuthHandler) AuthLogin(
	ctx context.Context,
	r openapi.AuthLoginRequestObject,
) (openapi.AuthLoginResponseObject, error) {
	token, err := h.authService.Login(ctx)
	if err != nil {
		return openapi.AuthLogin500JSONResponse{
			InternalErrorJSONResponse: openapi.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	return openapi.AuthLogin200JSONResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
