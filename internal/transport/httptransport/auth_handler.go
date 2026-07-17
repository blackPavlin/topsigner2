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
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return openapi.AuthLogout401JSONResponse{
			UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	var refreshToken *string

	if r.Body != nil && r.Body.RefreshToken != nil {
		refreshToken = r.Body.RefreshToken
	}

	if err := h.authService.Logout(ctx, user.ID, refreshToken); err != nil {
		return openapi.AuthLogout500JSONResponse{
			InternalErrorJSONResponse: openapi.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	return openapi.AuthLogout204Response{}, nil
}

// Refresh auth tokens
// (POST /api/v1/auth/refresh)
func (h *AuthHandler) AuthRefresh(
	ctx context.Context,
	r openapi.AuthRefreshRequestObject,
) (openapi.AuthRefreshResponseObject, error) {
	token, err := h.authService.Refresh(ctx, r.Body.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSessionNotFound):
			return openapi.AuthRefresh401JSONResponse{
				UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
			}, nil
		case errors.Is(err, auth.ErrTokenIsExpired):
			return openapi.AuthRefresh401JSONResponse{
				UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
			}, nil
		default:
			return openapi.AuthRefresh500JSONResponse{
				InternalErrorJSONResponse: openapi.InternalErrorJSONResponse{
					Message: "internal server error",
				},
			}, nil
		}
	}

	return openapi.AuthRefresh200JSONResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
