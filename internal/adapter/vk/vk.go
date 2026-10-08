package vk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bboykiv/topsigner/gen/external/vk/httpclient"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/group"
)

var (
	ErrInvalidAuthCode = errors.New("invalid authorization code")
	ErrUpstream        = errors.New("upstream provider error")
)

type Client struct {
	config      *config.Config
	client      *httpclient.ClientWithResponses
	oauthClient *httpclient.ClientWithResponses
}

// todo: добавить логгирование
// todo: добавить метрики
func NewClient(config *config.Config) (*Client, error) {
	client, err := httpclient.NewClientWithResponses(
		config.VK.BaseURL,
		httpclient.WithHTTPClient(&http.Client{Timeout: config.VKID.Timeout}),
	)
	if err != nil {
		return nil, fmt.Errorf("create new vk client with responses: %w", err)
	}

	oauthClient, err := httpclient.NewClientWithResponses(
		config.VK.OAuthBaseURL,
		httpclient.WithHTTPClient(&http.Client{Timeout: config.VKID.Timeout}),
	)
	if err != nil {
		return nil, fmt.Errorf("create new vk oauth client with responses: %w", err)
	}

	return &Client{
		config:      config,
		client:      client,
		oauthClient: oauthClient,
	}, nil
}

func (c *Client) GenerateConnectGroupURL(groupID int64, state string) (string, error) {
	params := &httpclient.AuthorizeParams{
		ClientID:     c.config.VKID.ClientID,
		RedirectURI:  c.config.VK.OAuthRedirectURL,
		GroupIDs:     strconv.FormatInt(groupID, 10),
		Scope:        c.config.VK.OAuthScope,
		ResponseType: httpclient.Code,
		Display:      new(httpclient.Page),
		State:        new(state),
		V:            httpclient.AuthorizeParamsVN5199,
	}

	req, err := httpclient.NewAuthorizeRequest(c.config.VKID.BaseURL, params)
	if err != nil {
		return "", fmt.Errorf("create authorization request: %w", err)
	}

	return req.URL.String(), nil
}

func (c *Client) ExchangeGroupCode(ctx context.Context, code string) (int64, string, error) {
	params := &httpclient.ExchangeGroupCodeParams{
		ClientID:     c.config.VKID.ClientID,
		ClientSecret: c.config.VKID.SecretKey,
		RedirectURI:  c.config.VK.OAuthRedirectURL,
		Code:         code,
	}

	resp, err := c.oauthClient.ExchangeGroupCodeWithResponse(ctx, params)
	if err != nil {
		return 0, "", fmt.Errorf("exchange group code: %w", err)
	}

	switch {
	case resp.JSON200 != nil && len(resp.JSON200.Groups) != 0:
		return resp.JSON200.Groups[0].GroupID, resp.JSON200.Groups[0].AccessToken, nil
	case resp.JSON200 != nil && resp.JSON200.Error != nil:
		return 0, "", fmt.Errorf("%w: %s", ErrInvalidAuthCode, resp.JSON200.Error.ErrorMsg)
	default:
		return 0, "", ErrUpstream
	}
}

func (c *Client) GetGroups(ctx context.Context, token string) ([]*group.Group, error) {
	// todo: пагинация ???
	params := &httpclient.GetGroupsParams{
		V:        httpclient.GetGroupsParamsVN5199,
		Extended: new(httpclient.GetGroupsParamsExtendedN1),
		Filter:   new([]httpclient.GroupsFilter{httpclient.Admin}),
	}

	resp, err := c.client.GetGroupsWithResponse(ctx, params, setAccessToken(token))
	if err != nil {
		return nil, fmt.Errorf("get vk groups: %w", err)
	}

	switch {
	case resp.JSON200 != nil && resp.JSON200.Response != nil:
		groups := make([]*group.Group, 0, resp.JSON200.Response.Count)

		for _, item := range resp.JSON200.Response.Items {
			groups = append(groups, &group.Group{
				ID:         item.ID,
				Name:       item.Name,
				ScreenName: item.ScreenName,
			})
		}

		return groups, nil
	case resp.JSON200 != nil && resp.JSON200.Error != nil:
		return nil, fmt.Errorf("get groups: %s", resp.JSON200.Error.ErrorMsg)
	default:
		return nil, ErrUpstream
	}
}

func setAccessToken(token string) httpclient.RequestEditorFn {
	return func(ctx context.Context, r *http.Request) error {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		return nil
	}
}
