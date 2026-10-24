package webmapper

import (
	"errors"
	"strings"
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmodel"
)

var (
	ErrInvalidLogin    = errors.New("invalid login input")
	ErrInvalidPassword = errors.New("invalid password input")
)

func SignUpToDomain(request webmodel.SignUpRequest) (domainmodel.SignUpRequest, error) {
	login := strings.TrimSpace(request.Login)
	if len(login) == 0 || strings.Contains(login, ":") {
		return domainmodel.SignUpRequest{}, ErrInvalidLogin
	}
	if len(strings.TrimSpace(request.Password)) == 0 || len(request.Password) > 72 {
		return domainmodel.SignUpRequest{}, ErrInvalidPassword
	}
	return domainmodel.SignUpRequest{Login: login, Password: request.Password}, nil
}
