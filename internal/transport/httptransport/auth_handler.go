package httptransport

import (
	"context"
	"errors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/auth"
	mw "github.com/bboykiv/topsigner/internal/transport/httptransport/middleware"
	"github.com/bboykiv/topsigner/internal/transport/httptransport/validation"
)

type AuthHandler struct {
	authService *auth.Service
	validate    *validator.Validate
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validation.New(),
	}
}

// Login user
// (POST /api/v1/auth/login)
func (h *AuthHandler) AuthLogin(
	ctx context.Context,
	r httpserver.AuthLoginRequestObject,
) (httpserver.AuthLoginResponseObject, error) {
	if err := h.validate.Struct(r.Body); err != nil {
		return httpserver.AuthLogin400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(err),
		}, nil
	}

	token, err := h.authService.Login(ctx, &auth.LoginInput{
		Email:     r.Body.Email,
		Password:  r.Body.Password,
		IP:        middleware.GetClientIP(ctx),
		UserAgent: mw.GetUserAgentFromContext(ctx),
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
// (GET /api/v1/auth/oauth/{provider})
func (h *AuthHandler) AuthOAuthGetURL(
	ctx context.Context,
	r httpserver.AuthOAuthGetURLRequestObject,
) (httpserver.AuthOAuthGetURLResponseObject, error) {
	var (
		authorizationURL string
		err              error
	)

	switch r.Provider {
	case httpserver.Vkontakte:
		authorizationURL, err = h.authService.GenerateVKIDOAuthURL(ctx)
		if err != nil {
			return httpserver.AuthOAuthGetURL500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	default:
		return httpserver.AuthOAuthGetURL400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(ErrInvalidOAuthProvider),
		}, nil
	}

	return httpserver.AuthOAuthGetURL200JSONResponse{
		URL: authorizationURL,
	}, nil
}

// AuthOAuthCallback OAuth callback
// (GET /api/v1/auth/oauth/{provider}/callback)
func (h *AuthHandler) AuthOAuthCallback(
	ctx context.Context,
	r httpserver.AuthOAuthCallbackRequestObject,
) (httpserver.AuthOAuthCallbackResponseObject, error) {
	var (
		tokenPair *auth.TokenPair
		err       error
	)

	switch r.Provider {
	case httpserver.Vkontakte:
		tokenPair, err = h.authService.ExchangeVKIDOAuthToken(ctx, &auth.OAuthExchangeTokenParams{
			Code:      r.Params.Code,
			DeviceID:  r.Params.DeviceID,
			State:     r.Params.State,
			IP:        middleware.GetClientIP(ctx),
			UserAgent: mw.GetUserAgentFromContext(ctx),
		})
		if err != nil {
			return httpserver.AuthOAuthCallback401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		}
	default:
		return httpserver.AuthOAuthCallback400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(ErrInvalidOAuthProvider),
		}, nil
	}

	return httpserver.AuthOAuthCallback200JSONResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
