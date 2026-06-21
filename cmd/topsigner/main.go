package main

import (
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/application"
)

func main() {
	fx.New(application.New()).Run()
}
