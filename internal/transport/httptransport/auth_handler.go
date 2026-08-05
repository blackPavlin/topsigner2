package httptransport

import (
	"context"
	"errors"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/bboykiv/topsigner/gen/httpserver"
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
	r httpserver.AuthLoginRequestObject,
) (httpserver.AuthLoginResponseObject, error) {
	token, err := h.authService.Login(ctx, &auth.LoginInput{
		IP:        middleware.GetClientIP(ctx),
		UserAgent: "",
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return httpserver.AuthLogin401JSONResponse{
				UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
			}, nil
		}

		return httpserver.AuthLogin500JSONResponse{
			InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	return httpserver.AuthLogin200JSONResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// Logout user
// (POST /api/v1/auth/logout)
func (h *AuthHandler) AuthLogout(
	ctx context.Context,
	r httpserver.AuthLogoutRequestObject,
) (httpserver.AuthLogoutResponseObject, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return httpserver.AuthLogout401JSONResponse{
			UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	var refreshToken *string

	if r.Body != nil && r.Body.RefreshToken != nil {
		refreshToken = r.Body.RefreshToken
	}

	if err := h.authService.Logout(ctx, user.ID, refreshToken); err != nil {
		return httpserver.AuthLogout500JSONResponse{
			InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	return httpserver.AuthLogout204Response{}, nil
}

// Refresh auth tokens
// (POST /api/v1/auth/refresh)
func (h *AuthHandler) AuthRefresh(
	ctx context.Context,
	r httpserver.AuthRefreshRequestObject,
) (httpserver.AuthRefreshResponseObject, error) {
	token, err := h.authService.Refresh(ctx, r.Body.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSessionNotFound):
			return httpserver.AuthRefresh401JSONResponse{
				UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
			}, nil
		case errors.Is(err, auth.ErrTokenIsExpired):
			return httpserver.AuthRefresh401JSONResponse{
				UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
			}, nil
		default:
			return httpserver.AuthRefresh500JSONResponse{
				InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
					Message: "internal server error",
				},
			}, nil
		}
	}

	return httpserver.AuthRefresh200JSONResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// Get OAuth authorization URL
// (GET /api/v1/auth/{provider})
func (h *AuthHandler) AuthOAuthRedirect(
	ctx context.Context,
	r httpserver.AuthOAuthRedirectRequestObject,
) (httpserver.AuthOAuthRedirectResponseObject, error) {
	var (
		authorizationURL string
		err              error
	)

	switch r.Provider {
	case httpserver.Vkontakte:
		authorizationURL, err = h.authService.GenerateVKIDAuthorizationURL()
		if err != nil {
			return httpserver.AuthOAuthRedirect500JSONResponse{
				InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
					Message: "internal server error",
				},
			}, nil
		}
	default:
		return httpserver.AuthOAuthRedirect400JSONResponse{
			BadRequestJSONResponse: httpserver.BadRequestJSONResponse{Message: "invalid provider"},
		}, nil
	}

	return httpserver.AuthOAuthRedirect200Response{}, nil
}
