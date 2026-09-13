package vk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bboykiv/topsigner/gen/external/vk/httpclient"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/group"
)

type Client struct {
	config *config.Config
	client *httpclient.ClientWithResponses
}

func NewClient(config *config.Config) (*Client, error) {
	// todo: добавить логгирование
	// todo: добавить метрики

	client, err := httpclient.NewClientWithResponses(config.VK.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("create new vk client with responses: %w", err)
	}

	return &Client{config: config, client: client}, nil
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

	groupList, err := resp.JSON200.Response.Items.AsGroupList()
	if err != nil {
		return nil, fmt.Errorf("unmarshal group list: %w", err)
	}

	groups := make([]*group.Group, 0, resp.JSON200.Response.Count)

	for _, g := range groupList {
		groups = append(groups, &group.Group{
			ID:         g.ID,
			Name:       g.Name,
			ScreenName: g.ScreenName,
		})
	}

	return groups, nil
}

func setAccessToken(token string) httpclient.RequestEditorFn {
	return func(ctx context.Context, r *http.Request) error {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		return nil
	}
}
