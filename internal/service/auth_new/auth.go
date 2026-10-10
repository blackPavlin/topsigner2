package auth

import "net/netip"

type Login struct {
	Email     string
	Password  string
	DeviceID  string
	IP        netip.Addr
	UserAgent string
}

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}
