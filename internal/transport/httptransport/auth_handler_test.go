package httptransport_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func TestAuthHandler_AuthLogin_EmptyBody(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'required' tag",
	}

	resp, err := client.AuthLoginWithResponse(t.Context(), httpserver.AuthLoginJSONRequestBody{})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}

func TestAuthHandler_AuthLogin_EmptyEmail(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'required' tag",
	}

	resp, err := client.AuthLoginWithResponse(t.Context(), httpserver.AuthLoginJSONRequestBody{
		Password: "password",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}

func TestAuthHandler_AuthLogin_InvalidEmail(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'email' tag",
	}

	resp, err := client.AuthLoginWithResponse(t.Context(), httpserver.AuthLoginJSONRequestBody{
		Email:    "invalid-email",
		Password: "password",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}

func TestAuthHandler_AuthLogin_EmptyPassword(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'password' failed on the 'required' tag",
	}

	resp, err := client.AuthLoginWithResponse(t.Context(), httpserver.AuthLoginJSONRequestBody{
		Email: "test@email.com",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}
