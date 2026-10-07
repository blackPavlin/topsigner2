package crypto

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
)

var Module = fx.Module("crypto", fx.Provide(New))

type Params struct {
	fx.In
	Config *config.Config
}

type Result struct {
	fx.Out
	Encryptor *Encryptor
}

func New(params Params) (Result, error) {
	encryptor, err := NewEncryptor(params.Config.Auth.EncryptionKey)
	if err != nil {
		return Result{}, fmt.Errorf("new encryptor: %w", err)
	}

	return Result{Encryptor: encryptor}, nil
}
