package httptransport

import (
	"context"
	"errors"

	"github.com/bboykiv/topsigner/gen/openapi"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/auth"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login user
// (POST /api/v1/auth/login)
func (h *AuthHandler) AuthLogin(
	ctx context.Context,
	r openapi.AuthLoginRequestObject,
) (openapi.AuthLoginResponseObject, error) {
	token, err := h.authService.Login(ctx)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return openapi.AuthLogin401JSONResponse{
				UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
			}, nil
		}

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

// Logout user
// (POST /api/v1/auth/logout)
func (h *AuthHandler) AuthLogout(
	ctx context.Context,
	r openapi.AuthLogoutRequestObject,
) (openapi.AuthLogoutResponseObject, error) {
	user, ok := model.GetUserFromContext(ctx)
	if !ok {
		return openapi.AuthLogout401JSONResponse{
			UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	if err := h.authService.Logout(ctx, user.ID); err != nil {
		return openapi.AuthLogout500JSONResponse{
			InternalErrorJSONResponse: openapi.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	return openapi.AuthLogout204Response{}, nil
}
