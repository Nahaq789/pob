package auth

type TokenResponse struct {
	JwtToken string `json:"access_token"`
}
