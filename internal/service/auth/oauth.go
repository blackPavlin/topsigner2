package auth

type OAuthExchangeTokenParams struct {
	Code         string
	CodeVerifier string
	DeviceID     string
	State        string
	IP           string
	UserAgent    string
}

type OAuthToken struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	TokenType    string
	ExpiresIn    int
	UserID       int64
	State        string
	Scope        string
}
