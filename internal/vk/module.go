package vk

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/group"
)

var Module = fx.Module("vk", fx.Provide(
	fx.Annotate(New, fx.As(new(group.VKClient))),
))

type Params struct {
	fx.In
	Config *config.Config
}

func New(params Params) (*Client, error) {
	client, err := NewClient(params.Config)
	if err != nil {
		return nil, fmt.Errorf("create vk client: %w", err)
	}

	return client, nil
}
