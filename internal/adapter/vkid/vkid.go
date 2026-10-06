package vkid

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bboykiv/topsigner/gen/external/vkid/httpclient"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
)

var (
	ErrInvalidAuthCode = errors.New("invalid authorization code")
	ErrUpstream        = errors.New("upstream provider error")
)

type Client struct {
	config *config.Config
	client *httpclient.ClientWithResponses
}

// todo: добавить логгирование
// todo: добавить метрики
func NewClient(config *config.Config) (*Client, error) {
	client, err := httpclient.NewClientWithResponses(
		config.VKID.BaseURL,
		httpclient.WithHTTPClient(&http.Client{Timeout: config.VKID.Timeout}),
	)
	if err != nil {
		return nil, fmt.Errorf("create new vkid client with responses: %w", err)
	}

	return &Client{config: config, client: client}, nil
}

func (c *Client) GetAuthorizationURL(codeChallenge, state string) (string, error) {
	params := &httpclient.AuthorizeParams{
		ResponseType:        httpclient.Code,
		ClientID:            c.config.VKID.ClientID,
		RedirectURI:         c.config.VKID.RedirectURL,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: httpclient.S256,
		State:               state,
		Scope:               &c.config.VKID.Scope,
	}

	req, err := httpclient.NewAuthorizeRequest(c.config.VKID.BaseURL, params)
	if err != nil {
		return "", fmt.Errorf("create authorization request: %w", err)
	}

	return req.URL.String(), nil
}

func (c *Client) Exchange(
	ctx context.Context,
	params *auth.OAuthExchangeTokenParams,
) (*auth.OAuthToken, error) {
	body := httpclient.ExchangeTokenFormdataRequestBody{
		GrantType:    httpclient.AuthorizationCode,
		CodeVerifier: new(params.CodeVerifier),
		RedirectURI:  new(c.config.VKID.RedirectURL),
		Code:         new(params.Code),
		ClientID:     c.config.VKID.ClientID,
		DeviceID:     params.DeviceID,
		State:        params.State,
	}

	resp, err := c.client.ExchangeTokenWithFormdataBodyWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("exchane vkid oauth token: %w", err)
	}

	switch {
	case resp.StatusCode() == http.StatusOK && resp.JSON200 != nil:
		return tokenResponseToOAuthToken(resp.JSON200), nil
	case resp.StatusCode() == http.StatusBadRequest && resp.JSON400 != nil:
		return nil, fmt.Errorf("%w: %s", ErrInvalidAuthCode, resp.JSON400.Error)
	default:
		return nil, ErrUpstream
	}
}

func (c *Client) Refresh(
	ctx context.Context,
	params *auth.OAuthRefreshTokenParams,
) (*auth.OAuthToken, error) {
	body := httpclient.ExchangeTokenFormdataRequestBody{
		GrantType:    httpclient.RefreshToken,
		RefreshToken: new(params.RefreshToken),
		ClientID:     c.config.VKID.ClientID,
		DeviceID:     params.DeviceID,
		State:        params.State,
	}

	resp, err := c.client.ExchangeTokenWithFormdataBodyWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("refresh vkid oauth token: %w", err)
	}

	switch {
	case resp.StatusCode() == http.StatusOK && resp.JSON200 != nil:
		return tokenResponseToOAuthToken(resp.JSON200), nil
	case resp.StatusCode() == http.StatusBadRequest && resp.JSON400 != nil:
		return nil, fmt.Errorf("%w: %s", ErrInvalidAuthCode, resp.JSON400.Error)
	default:
		return nil, ErrUpstream
	}
}

func (c *Client) Logout(ctx context.Context, token string) error {
	body := httpclient.LogoutFormdataRequestBody{
		ClientID:    c.config.VKID.ClientID,
		AccessToken: token,
	}

	resp, err := c.client.LogoutWithFormdataBodyWithResponse(ctx, body)
	if err != nil {
		return fmt.Errorf("vkid logout: %w", err)
	}

	switch {
	case resp.StatusCode() == http.StatusOK && resp.JSON200 != nil:
		return nil
	default:
		return ErrUpstream
	}
}

func tokenResponseToOAuthToken(resp *httpclient.TokenResponse) *auth.OAuthToken {
	return &auth.OAuthToken{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		IDToken:      resp.IDToken,
		TokenType:    resp.TokenType,
		ExpiresIn:    resp.ExpiresIn,
		UserID:       resp.UserID,
		State:        resp.State,
		Scope:        resp.Scope,
	}
}
