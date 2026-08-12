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
		Email:     string(r.Body.Email),
		Password:  r.Body.Password,
		IP:        middleware.GetClientIP(ctx),
		UserAgent: "",
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return httpserver.AuthLogin401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		}

		return httpserver.AuthLogin500JSONResponse{
			InternalErrorJSONResponse: NewInternalError(),
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
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	var refreshToken *string

	if r.Body != nil && r.Body.RefreshToken != nil {
		refreshToken = r.Body.RefreshToken
	}

	if err := h.authService.Logout(ctx, user.ID, refreshToken); err != nil {
		return httpserver.AuthLogout500JSONResponse{
			InternalErrorJSONResponse: NewInternalError(),
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
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		case errors.Is(err, auth.ErrTokenIsExpired):
			return httpserver.AuthRefresh401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		default:
			return httpserver.AuthRefresh500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
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
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	default:
		return httpserver.AuthOAuthRedirect400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError("invalid provider"),
		}, nil
	}

	// https://topsigner.ru/api/v1/auth/vkontakte/callback?
	// 	code=vk2.a.EthCB9hsxZEeCMepeZZBNuoCwS2q2k0xVFhiF4zvxNXjRhicQrpDj834h1bVzbX7LZgsAKKpYLIVTzo_9nR_Iquo3TGHQnWDIQrTn0fqoa6WoZZBVTBc-Gsn8YZXrLR1jppWuuj3WYmvE7sKsAADRuoBiqzugfSzVQ3Wf8IUYxZD5ha_1DlHYLwZFuwSTexGgz--BGIyiy-kc7mYKI0hOw
	// 	&expires_in=600
	// 	&device_id=LwSCGhkkaDNsBX68hNcZ0oUFqWhn2yVYCDefK5LQpr5sfrG4fAa7mk5HHbO5OdzoFZN_RGxJHJQydQtHd2QZAA
	// 	&state=QCX4cW5z0RFo0HC0XUF75guGHcLr029E
	// 	&ext_id=Z9hH0mnsngWDJp50EkKa9IrsgcsqJ61kiRsKWyqEBC_1_YrdfcQ5ue4hm3p4t_JdWdLa8X86eUEc-ec4E-eUD-iFGnPZSLtdjCBUmXqQ1XlGvUWn0KwreEu2O9jTDHCzZZg6ArbW3w56mSzBS2xIUEyz4VSdIbZn4N5QYhg5YOzT
	// 	&type=code_v2

	return httpserver.AuthOAuthRedirect200JSONResponse{
		URL: authorizationURL,
	}, nil
}
