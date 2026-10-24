package mapper

import (
	"tictactoe/datasource/model"
	"tictactoe/domain/domainmodel"
)

func UserToDomain(user model.User) domainmodel.User {
	return domainmodel.User{ID: user.ID, Login: user.Login, PasswordHash: user.PasswordHash}
}

func UserToDatasource(user domainmodel.User) model.User {
	return model.User{ID: user.ID, Login: user.Login, PasswordHash: user.PasswordHash}
}
