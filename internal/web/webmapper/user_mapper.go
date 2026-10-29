package webmapper

import (
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmodel"
)

func ToUserResponse(user domainmodel.User) webmodel.UserResponse {
	return webmodel.UserResponse{ID: user.ID, Login: user.Login}
}
