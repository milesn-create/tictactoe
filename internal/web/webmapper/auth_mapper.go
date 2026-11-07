package webmapper

import (
	"strings"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmodel"
)

func SignUpToDomain(request webmodel.SignUpRequest) (domainmodel.SignUpRequest, error) {
	login := strings.TrimSpace(request.Login)
	if len(login) == 0 || strings.Contains(login, ":") {
		return domainmodel.SignUpRequest{}, domainerror.ErrInvalidLogin
	}
	if len(strings.TrimSpace(request.Password)) == 0 || len(request.Password) > 72 {
		return domainmodel.SignUpRequest{}, domainerror.ErrInvalidPassword
	}
	return domainmodel.SignUpRequest{Login: login, Password: request.Password}, nil
}
func SignInToDomain(request webmodel.JwtRequest) (domainmodel.JwtRequest, error) {
	login := strings.TrimSpace(request.Login)
	if len(login) == 0 || strings.Contains(login, ":") {
		return domainmodel.JwtRequest{}, domainerror.ErrInvalidLogin
	}
	if len(strings.TrimSpace(request.Password)) == 0 || len(request.Password) > 72 {
		return domainmodel.JwtRequest{}, domainerror.ErrInvalidPassword
	}
	return domainmodel.JwtRequest{Login: login, Password: request.Password}, nil

}
func RefreshToDomain(request webmodel.RefreshJwtRequest) (domainmodel.RefreshJwtRequest, error) {
	if request.RefreshToken == "" {
		return domainmodel.RefreshJwtRequest{}, domainerror.ErrEmptyRefreshToken
	}
	return domainmodel.RefreshJwtRequest{RefreshToken: request.RefreshToken}, nil

}

func JwtToWeb(response domainmodel.JwtResponse) webmodel.JwtResponse {
	return webmodel.JwtResponse{Type: response.Type, AccessToken: response.AccessToken, RefreshToken: response.RefreshToken}
}
