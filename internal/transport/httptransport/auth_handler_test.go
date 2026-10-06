package httptransport_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func TestAuthHandler_LoginUser_EmptyBody(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'required' tag",
	}

	resp, err := client.LoginUserWithResponse(t.Context(), httpserver.LoginUserJSONRequestBody{})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}

func TestAuthHandler_LoginUser_EmptyEmail(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'required' tag",
	}

	resp, err := client.LoginUserWithResponse(t.Context(), httpserver.LoginUserJSONRequestBody{
		Password: "password",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}

func TestAuthHandler_LoginUser_InvalidEmail(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'email' tag",
	}

	resp, err := client.LoginUserWithResponse(t.Context(), httpserver.LoginUserJSONRequestBody{
		Email:    "invalid-email",
		Password: "password",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}

func TestAuthHandler_LoginUser_EmptyPassword(t *testing.T) {
	t.Parallel()

	message := &httpserver.BadRequest{
		Message: "field validation for 'password' failed on the 'required' tag",
	}

	resp, err := client.LoginUserWithResponse(t.Context(), httpserver.LoginUserJSONRequestBody{
		Email: "test@email.com",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Equal(t, message, resp.JSON400)
}
