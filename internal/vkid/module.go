package vkid

import (
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
	return NewClient(params.Config), nil
}
