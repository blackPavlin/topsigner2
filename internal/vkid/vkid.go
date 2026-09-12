package vkid

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bboykiv/topsigner/gen/external/vkid/httpclient"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
)

type Client struct {
	config *config.Config
	client *httpclient.ClientWithResponses
}

func NewClient(config *config.Config) (*Client, error) {
	// todo: добавить логгирование
	// todo: добавить метрики

	client, err := httpclient.NewClientWithResponses(config.VKID.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("create new vkid client with responses: %w", err)
	}

	return &Client{config: config, client: client}, nil
}

func (c *Client) GenerateOAuthURL(challenge, state string) (string, error) {
	u, err := url.Parse(c.config.VKID.BaseURL)
	if err != nil {
		return "", fmt.Errorf("parse vkid base url: %w", err)
	}

	u = u.JoinPath("authorize")

	q := u.Query()

	q.Set("response_type", "code")
	q.Set("client_id", c.config.VKID.ClientID)
	q.Set("redirect_uri", c.config.VKID.RedirectURL)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	q.Set("state", state)
	q.Set("scope", c.config.VKID.Scope)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *Client) ExchangeOAuthToken(
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

	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return nil, errors.New(resp.JSON400.Error)
		case http.StatusInternalServerError:
			return nil, errors.New(resp.JSON500.Error)
		default:
			// todo: улучшить обработку ошибок
		}
	}

	return tokenResponseToOAuthToken(resp.JSON200), nil
}

func (c *Client) RefreshOAuthToken(
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

	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return nil, errors.New(resp.JSON400.Error)
		case http.StatusInternalServerError:
			return nil, errors.New(resp.JSON500.Error)
		default:
			// todo: улучшить обработку ошибок
		}
	}

	return tokenResponseToOAuthToken(resp.JSON200), nil
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

	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return errors.New(resp.JSON401.Error)
		case http.StatusInternalServerError:
			return errors.New(resp.JSON500.Error)
		default:
			// todo: улучшить обработку ошибок
		}
	}

	return nil
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
