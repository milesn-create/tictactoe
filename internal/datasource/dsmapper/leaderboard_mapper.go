package dsmapper

import (
	"tictactoe/datasource/dsmodel"
	"tictactoe/domain/domainmodel"
)

func LeaderboardToDomain(row dsmodel.LeaderboardRow) domainmodel.LeaderboardEntry {
	return domainmodel.LeaderboardEntry{UserID: row.UserID, Login: row.Login, WinRate: row.WinRate}
}
func LeaderboardListToDomain(rows []dsmodel.LeaderboardRow) []domainmodel.LeaderboardEntry {
	result := make([]domainmodel.LeaderboardEntry, 0, len(rows))
	for i := range len(rows) {
		result = append(result, LeaderboardToDomain(rows[i]))
	}
	return result
}
