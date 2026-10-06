package vkid

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
)

var Module = fx.Module("vkid", fx.Provide(
	fx.Annotate(New, fx.As(new(auth.VKIDClient))),
))

type Params struct {
	fx.In
	Config *config.Config
}

func New(params Params) (*Client, error) {
	client, err := NewClient(params.Config)
	if err != nil {
		return nil, fmt.Errorf("create vkid client: %w", err)
	}

	return client, nil
}
