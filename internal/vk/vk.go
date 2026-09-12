package vk

import (
	"context"
	"fmt"

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

func (c *Client) GetGroups(ctx context.Context) ([]*group.Group, error) {
	// todo: пагинация ???
	params := &httpclient.GetGroupsParams{
		V: httpclient.GetGroupsParamsVN5199,
	}

	_, err := c.client.GetGroupsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get vk groups: %w", err)
	}

	return nil, nil
}
