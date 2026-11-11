package domainservice

import (
	"context"
	"tictactoe/domain/domainmodel"
)

type AuthService interface {
	SignUp(ctx context.Context, request domainmodel.SignUpRequest) error
	Authorize(ctx context.Context, request domainmodel.JwtRequest) (domainmodel.JwtResponse, error)
	UpdateAccessToken(ctx context.Context, request domainmodel.RefreshJwtRequest) (domainmodel.JwtResponse, error)
	UpdateRefreshToken(ctx context.Context, request domainmodel.RefreshJwtRequest) (domainmodel.JwtResponse, error)
}
