package webmapper

import (
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmodel"

	"github.com/google/uuid"
)

func ToBoard(board [][]int) (domainmodel.Board, error) {
	if len(board) != 3 {
		return domainmodel.Board{}, domainerror.ErrInvalidBoardSize
	}
	var boardDomain domainmodel.Board
	for i := 0; i < 3; i++ {
		if len(board[i]) != 3 {
			return domainmodel.Board{}, domainerror.ErrInvalidBoardSize
		}
		for j := 0; j < 3; j++ {
			if board[i][j] != int(domainmodel.Cross) && board[i][j] != int(domainmodel.Nought) && board[i][j] != int(domainmodel.Empty) {
				return domainmodel.Board{}, domainerror.ErrInvalidCellValue
			}
			boardDomain[i][j] = domainmodel.Cell(board[i][j])
		}
	}
	return boardDomain, nil
}
func ToResponse(game domainmodel.Game) webmodel.GameResponse {
	board := make([][]int, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = int(game.Board[i][j])
		}
	}
	state := stateToString(game.State)
	return webmodel.GameResponse{ID: game.ID, Board: board, State: state, VsComputer: game.VsComputer, TurnPlayerID: uuidToPtr(game.TurnPlayerID), WinnerID: uuidToPtr(game.WinnerID), CrossPlayerID: uuidToPtr(game.CrossPlayerID), NoughtPlayerID: uuidToPtr(game.NoughtPlayerID)}

}
func stateToString(state domainmodel.GameState) string {
	switch state {
	case domainmodel.StateWaitingPlayers:
		return "waiting_players"
	case domainmodel.StatePlayerTurn:
		return "player_turn"
	case domainmodel.StateDraw:
		return "draw"
	case domainmodel.StateWon:
		return "won"
	case domainmodel.StateArchived:
		return "archived"
	default:
		return "unknown"
	}

}
func uuidToPtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func ToResponseList(games []domainmodel.Game) []webmodel.GameResponse {
	gamesResp := make([]webmodel.GameResponse, 0, len(games))
	for i := range len(games) {
		gamesResp = append(gamesResp, ToResponse(games[i]))
	}
	return gamesResp
}
