package webmodel

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type SignInResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
