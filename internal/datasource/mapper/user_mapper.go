package mapper

import (
	"tictactoe/datasource/model"
	"tictactoe/domain"
)

func UserToDomain(user model.User) domain.User {
	return domain.User{ID: user.ID, Login: user.Login, PasswordHash: user.PasswordHash}
}

func UserToDatasource(user domain.User) model.User {
	return model.User{ID: user.ID, Login: user.Login, PasswordHash: user.PasswordHash}
}
