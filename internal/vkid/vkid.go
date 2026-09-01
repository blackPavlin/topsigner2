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

	client, err := httpclient.NewClientWithResponses(
		config.VKID.BaseURL,
		httpclient.WithBaseURL(config.VKID.BaseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("create new client with responses: %w", err)
	}

	return &Client{
		config: config,
		client: client,
	}, nil
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

	return &auth.OAuthToken{
		AccessToken:  resp.JSON200.AccessToken,
		RefreshToken: resp.JSON200.RefreshToken,
		IDToken:      resp.JSON200.IDToken,
		TokenType:    resp.JSON200.TokenType,
		ExpiresIn:    resp.JSON200.ExpiresIn,
		UserID:       resp.JSON200.UserID,
		State:        resp.JSON200.State,
		Scope:        resp.JSON200.Scope,
	}, nil
}
