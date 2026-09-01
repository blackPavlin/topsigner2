package vkid_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/gen/external/vkid/httpclient"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/vkid"
)

func TestVKID_GenerateOAuthURL_Success(t *testing.T) {
	cfg := &config.Config{
		VKID: config.VKIDConfig{
			BaseURL:     "https://id.vk.com",
			ClientID:    "client123",
			RedirectURL: "https://example.com/callback",
			Scope:       "email",
		},
	}

	client, err := vkid.NewClient(cfg)
	require.NoError(t, err)

	got, err := client.GenerateOAuthURL("challenge123", "state456")
	require.NoError(t, err)

	u, err := url.Parse(got)
	require.NoError(t, err)

	require.Equal(t, "id.vk.com", u.Host)
	require.Equal(t, "/authorize", u.Path)
	require.Equal(t, "client123", u.Query().Get("client_id"))
	require.Equal(t, "S256", u.Query().Get("code_challenge_method"))
}

func TestVKID_ExchangeOAuthToken_Success(t *testing.T) {
	const (
		clientID    = "client123"
		redirectURL = "https://example.com/callback"
		scope       = "email"
	)

	tokenResponse := httpclient.TokenResponse{
		AccessToken:  "access_token",
		ExpiresIn:    600,
		IDToken:      "id_token",
		RefreshToken: "refresh_token",
		Scope:        scope,
		State:        "state456",
		TokenType:    "Bearer",
		UserID:       123,
	}

	params := &auth.OAuthExchangeTokenParams{
		Code:         "code123",
		CodeVerifier: "verifier123",
		DeviceID:     "device1",
		State:        "state1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		err := r.ParseForm()
		require.NoError(t, err)

		require.Equal(t, string(httpclient.AuthorizationCode), r.Form.Get("grant_type"))
		require.Equal(t, clientID, r.Form.Get("client_id"))
		require.Equal(t, params.Code, r.Form.Get("code"))
		require.Equal(t, params.CodeVerifier, r.Form.Get("code_verifier"))
		require.Equal(t, redirectURL, r.Form.Get("redirect_uri"))
		require.Equal(t, params.DeviceID, r.Form.Get("device_id"))
		require.Equal(t, params.State, r.Form.Get("state"))

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(tokenResponse)
		require.NoError(t, err)
	}))
	defer server.Close()

	cfg := &config.Config{
		VKID: config.VKIDConfig{
			BaseURL:     server.URL,
			ClientID:    clientID,
			RedirectURL: redirectURL,
			Scope:       scope,
		},
	}

	client, err := vkid.NewClient(cfg)
	require.NoError(t, err)

	got, err := client.ExchangeOAuthToken(t.Context(), params)
	require.NoError(t, err)
	require.Equal(t, tokenResponse.AccessToken, got.AccessToken)
	require.Equal(t, tokenResponse.ExpiresIn, got.ExpiresIn)
	require.Equal(t, tokenResponse.IDToken, got.IDToken)
	require.Equal(t, tokenResponse.RefreshToken, got.RefreshToken)
	require.Equal(t, tokenResponse.Scope, got.Scope)
	require.Equal(t, tokenResponse.State, got.State)
	require.Equal(t, tokenResponse.TokenType, got.TokenType)
	require.Equal(t, tokenResponse.UserID, got.UserID)
}
