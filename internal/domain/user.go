package domain

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
}
type SignUpRequest struct {
	Login    string
	Password string
}
