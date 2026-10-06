package httptransport_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func TestFontHandler_ListFonts_Unauthorized(t *testing.T) {
	t.Parallel()

	message := &httpserver.Unauthorized{
		Message: "unauthorized",
	}

	resp, err := client.ListFontsWithResponse(t.Context(), &httpserver.ListFontsParams{})
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	require.Equal(t, message, resp.JSON401)
}
