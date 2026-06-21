package application_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/application"
)

func TestValidateOptions(t *testing.T) {
	require.NoError(t, fx.ValidateApp(application.New()))
}
