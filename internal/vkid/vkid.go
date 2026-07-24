package vkid

import (
	"net/url"

	"github.com/bboykiv/topsigner/internal/config"
)

type Client struct {
	config *config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) GenerateAuthorizationURL(challenge, state string) string {
	u := url.URL{Scheme: "https", Host: "id.vk.ru", Path: "/authorize"}
	q := u.Query()

	q.Set("response_type", "code")
	q.Set("client_id", c.config.VKID.ClientID)
	q.Set("redirect_uri", c.config.VKID.RedirectURL)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	q.Set("state", state)
	q.Set("scope", c.config.VKID.Scope)

	u.RawQuery = q.Encode()

	return u.String()
}
