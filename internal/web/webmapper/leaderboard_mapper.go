package webmapper

import (
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmodel"
)

func LeaderboardToWeb(leaderboard domainmodel.LeaderboardEntry) webmodel.LeaderboardResponse {
	return webmodel.LeaderboardResponse{UserID: leaderboard.UserID, Login: leaderboard.Login, WinRate: leaderboard.WinRate}
}
func LeaderboardListToWeb(leaderboard []domainmodel.LeaderboardEntry) []webmodel.LeaderboardResponse {
	result := make([]webmodel.LeaderboardResponse, 0, len(leaderboard))
	for i := range len(leaderboard) {
		result = append(result, LeaderboardToWeb(leaderboard[i]))
	}
	return result
}
