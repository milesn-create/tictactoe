package web

import "github.com/google/uuid"

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type SignInResponse struct {
	ID uuid.UUID `json:"id"`
}
