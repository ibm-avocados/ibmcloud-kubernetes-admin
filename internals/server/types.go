package server

type AccountEmailBody struct {
	AccountID string   `json:"accountID"`
	Email     []string `json:"email"`
}

type StatusOK struct {
	Message string `json:"message"`
}

type AccountLogin struct {
	OTP string `json:"otp"`
}

type Bill struct {
	Bill string `json:"bill"`
}

type OauthSettings struct {
	authURL, tokenURL, clientID, clientSecret, redirectURI string
}

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   int64  `json:"expiration"`
}