package vk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bboykiv/topsigner/gen/external/vk/httpclient"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/group"
)

type Client struct {
	config      *config.Config
	client      *httpclient.ClientWithResponses
	oauthClient *httpclient.ClientWithResponses
}

func NewClient(config *config.Config) (*Client, error) {
	// todo: добавить логгирование
	// todo: добавить метрики

	client, err := httpclient.NewClientWithResponses(config.VK.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("create new vk client with responses: %w", err)
	}

	oauthClient, err := httpclient.NewClientWithResponses(config.VK.OAuthBaseURL)
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
	u, err := url.Parse(c.config.VK.OAuthBaseURL)
	if err != nil {
		return "", fmt.Errorf("parse vk oauth base url: %w", err)
	}

	u = u.JoinPath("authorize")

	q := u.Query()

	q.Set("client_id", c.config.VKID.ClientID)
	q.Set("redirect_uri", c.config.VK.OAuthRedirectURL)
	q.Set("group_ids", strconv.FormatInt(groupID, 10))
	q.Set("display", "page")
	q.Set("scope", c.config.VK.OAuthScope)
	q.Set("response_type", "code")
	q.Set("v", string(httpclient.GetGroupsParamsVN5199))
	q.Set("state", state)

	u.RawQuery = q.Encode()

	return u.String(), nil
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

	if resp.JSON200.Error != nil {
		return nil, fmt.Errorf("get groups: %s", resp.JSON200.Error.ErrorMsg)
	}

	groups := make([]*group.Group, 0, resp.JSON200.Response.Count)

	for _, item := range resp.JSON200.Response.Items {
		groups = append(groups, &group.Group{
			ID:         item.ID,
			Name:       item.Name,
			ScreenName: item.ScreenName,
		})
	}

	return groups, nil
}

func (c *Client) ExchangGroupCode(ctx context.Context, code string) (int64, string, error) {
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

	if resp.JSON200.Error != nil {
		return 0, "", fmt.Errorf("exchange group token: %s", resp.JSON200.Error.ErrorMsg)
	}

	return resp.JSON200.Groups[0].GroupID, resp.JSON200.Groups[0].AccessToken, nil
}

func setAccessToken(token string) httpclient.RequestEditorFn {
	return func(ctx context.Context, r *http.Request) error {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		return nil
	}
}
