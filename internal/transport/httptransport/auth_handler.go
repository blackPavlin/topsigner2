package httptransport

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/transport/httptransport/middleware"
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
func (h *AuthHandler) LoginUser(
	ctx context.Context,
	r httpserver.LoginUserRequestObject,
) (httpserver.LoginUserResponseObject, error) {
	if err := h.validate.Struct(r.Body); err != nil {
		return httpserver.LoginUser400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(err),
		}, nil
	}

	token, err := h.authService.Login(ctx, &auth.LoginInput{
		Email:     r.Body.Email,
		Password:  r.Body.Password,
		IP:        middleware.GetClientIPFromContext(ctx),
		UserAgent: middleware.GetUserAgentFromContext(ctx),
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			return httpserver.LoginUser401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		case errors.Is(err, auth.ErrPasswordLoginNotAvailable):
			return httpserver.LoginUser401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		case errors.Is(err, auth.ErrInvalidPassword):
			return httpserver.LoginUser401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		default:
			return httpserver.LoginUser500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	}

	return httpserver.LoginUser200JSONResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    httpserver.Bearer,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

// Logout user
// (POST /api/v1/auth/logout)
func (h *AuthHandler) LogoutUser(
	ctx context.Context,
	r httpserver.LogoutUserRequestObject,
) (httpserver.LogoutUserResponseObject, error) {
	sessionID, ok := middleware.GetSessionIDFromContext(ctx)
	if !ok {
		return httpserver.LogoutUser401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return httpserver.LogoutUser401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	var allSessions = false

	if r.Body != nil && r.Body.AllSessions != nil {
		allSessions = *r.Body.AllSessions
	}

	if err := h.authService.Logout(ctx, user.ID, sessionID, allSessions); err != nil {
		return httpserver.LogoutUser500JSONResponse{
			InternalErrorJSONResponse: NewInternalError(),
		}, nil
	}

	return httpserver.LogoutUser204Response{}, nil
}

// Refresh auth tokens
// (POST /api/v1/auth/refresh)
func (h *AuthHandler) RefreshTokens(
	ctx context.Context,
	r httpserver.RefreshTokensRequestObject,
) (httpserver.RefreshTokensResponseObject, error) {
	if err := h.validate.Struct(r.Body); err != nil {
		return httpserver.RefreshTokens400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(err),
		}, nil
	}

	token, err := h.authService.Refresh(ctx, r.Body.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSessionNotFound):
			return httpserver.RefreshTokens401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		case errors.Is(err, auth.ErrTokenIsExpired):
			return httpserver.RefreshTokens401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		case errors.Is(err, model.ErrInvalidAuthCode):
			return httpserver.RefreshTokens401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		case errors.Is(err, model.ErrUpstreamUnavailable):
			return httpserver.RefreshTokens502JSONResponse{
				BadGatewayJSONResponse: NewBadGatewayError(),
			}, nil
		default:
			return httpserver.RefreshTokens500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	}

	return httpserver.RefreshTokens200JSONResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    httpserver.Bearer,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

// GetVKIDAuthorizationURL Get VK ID authorization URL
// (GET /api/v1/auth/vkid/authorize)
func (h *AuthHandler) GetVKIDAuthorizationURL(
	ctx context.Context,
	r httpserver.GetVKIDAuthorizationURLRequestObject,
) (httpserver.GetVKIDAuthorizationURLResponseObject, error) {
	authorizationURL, err := h.authService.GenerateVKIDOAuthURL(ctx)
	if err != nil {
		return httpserver.GetVKIDAuthorizationURL500JSONResponse{
			InternalErrorJSONResponse: NewInternalError(),
		}, nil
	}

	return httpserver.GetVKIDAuthorizationURL200JSONResponse{
		URL: authorizationURL,
	}, nil
}

// HandleVKIDCallback Exchange VK ID authorization code
// (GET /api/v1/auth/vkid/callback)
func (h *AuthHandler) HandleVKIDCallback(
	ctx context.Context,
	r httpserver.HandleVKIDCallbackRequestObject,
) (httpserver.HandleVKIDCallbackResponseObject, error) {
	// todo: валидация входных параметров

	token, err := h.authService.ExchangeVKIDOAuthToken(ctx, &auth.OAuthExchangeTokenParams{
		Code:      r.Params.Code,
		DeviceID:  r.Params.DeviceID,
		State:     r.Params.State,
		IP:        middleware.GetClientIPFromContext(ctx),
		UserAgent: middleware.GetUserAgentFromContext(ctx),
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrCodeVerifierNotFound):
			return httpserver.HandleVKIDCallback400JSONResponse{
				BadRequestJSONResponse: NewBadRequestError(err),
			}, nil
		case errors.Is(err, model.ErrInvalidAuthCode):
			return httpserver.HandleVKIDCallback401JSONResponse{
				UnauthorizedJSONResponse: NewUnauthorizedError(),
			}, nil
		case errors.Is(err, model.ErrUpstreamUnavailable):
			return httpserver.HandleVKIDCallback502JSONResponse{
				BadGatewayJSONResponse: NewBadGatewayError(),
			}, nil
		default:
			return httpserver.HandleVKIDCallback500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	}

	return httpserver.HandleVKIDCallback200JSONResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    httpserver.Bearer,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}
