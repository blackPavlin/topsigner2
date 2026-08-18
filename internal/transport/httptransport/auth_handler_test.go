package httptransport_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func TestAuthLogin(t *testing.T) {
	t.Parallel()

	type wants struct {
		statusCode int
		json200    *httpserver.AuthTokens
		json400    *httpserver.BadRequest
	}

	type testCase struct {
		name  string
		args  httpserver.AuthLoginJSONRequestBody
		wants wants
	}

	testCases := []testCase{
		{
			name: "400 empty body",
			args: httpserver.AuthLoginJSONRequestBody{},
			wants: wants{
				statusCode: 400,
				json400: &httpserver.BadRequest{
					Message: "field validation for 'email' failed on the 'required' tag",
				},
			},
		},
		{
			name: "400 empty email",
			args: httpserver.AuthLoginJSONRequestBody{
				Email:    "",
				Password: "password",
			},
			wants: wants{
				statusCode: 400,
				json400: &httpserver.BadRequest{
					Message: "field validation for 'email' failed on the 'required' tag",
				},
			},
		},
		{
			name: "400 invalid email",
			args: httpserver.AuthLoginJSONRequestBody{
				Email:    "not-valid-email",
				Password: "password",
			},
			wants: wants{
				statusCode: 400,
				json400: &httpserver.BadRequest{
					Message: "field validation for 'email' failed on the 'email' tag",
				},
			},
		},
		{
			name: "400 empty password",
			args: httpserver.AuthLoginJSONRequestBody{
				Email:    "test@email.com",
				Password: "",
			},
			wants: wants{
				statusCode: 400,
				json400: &httpserver.BadRequest{
					Message: "field validation for 'password' failed on the 'required' tag",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := client.AuthLoginWithResponse(t.Context(), tc.args)
			require.NoError(t, err)
			require.Equal(t, tc.wants.statusCode, resp.StatusCode())
			require.Equal(t, tc.wants.json200, resp.JSON200)
			require.Equal(t, tc.wants.json400, resp.JSON400)
		})
	}
}
