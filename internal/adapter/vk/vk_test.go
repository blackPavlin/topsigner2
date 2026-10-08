package vk_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/internal/adapter/vk"
	"github.com/bboykiv/topsigner/internal/config"
)

func TestVK_GenerateConnectGroupURL_Success(t *testing.T) {
	cfg := &config.Config{
		VK: config.VKConfig{
			OAuthBaseURL:     "ttps://oauth.vk.com",
			OAuthRedirectURL: "https://example.com/callback",
			OAuthScope:       "manage,messages",
		},
	}

	var (
		groupID int64 = 123
		state         = "state456"
	)

	client, err := vk.NewClient(cfg)
	require.NoError(t, err)

	got, err := client.GenerateConnectGroupURL(groupID, state)
	require.NoError(t, err)

	u, err := url.Parse(got)
	require.NoError(t, err)

	require.Equal(t, "oauth.vk.com", u.Host)
	require.Equal(t, "/authorize", u.Path)
	require.Equal(t, state, u.Query().Get("state"))
	require.Equal(t, cfg.VK.OAuthScope, u.Query().Get("scope"))
	require.Equal(t, strconv.FormatInt(groupID, 10), u.Query().Get("group_ids"))
	require.Equal(t, cfg.VK.OAuthRedirectURL, u.Query().Get("redirect_uri"))
}
