package dsmapper

import (
	"tictactoe/datasource/dsmodel"
	"tictactoe/domain/domainmodel"
)

func UserToDomain(user dsmodel.User) domainmodel.User {
	return domainmodel.User{ID: user.ID, Login: user.Login, PasswordHash: user.PasswordHash}
}

func UserToDatasource(user domainmodel.User) dsmodel.User {
	return dsmodel.User{ID: user.ID, Login: user.Login, PasswordHash: user.PasswordHash}
}
