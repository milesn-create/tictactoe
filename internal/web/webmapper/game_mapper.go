package webmapper

import (
	"errors"
	"tictactoe/domain/domainmodel"
	"tictactoe/web/webmodel"

	"github.com/google/uuid"
)

var (
	ErrInvalidBoardSize = errors.New("board must be 3x3")
	ErrInvalidCellValue = errors.New("cell must be -1, 0 or 1")
)

func ToDomain(game webmodel.Game) (domainmodel.Game, error) {
	if len(game.Board) != 3 {
		return domainmodel.Game{}, ErrInvalidBoardSize
	}
	var board domainmodel.Board
	for i := 0; i < 3; i++ {
		if len(game.Board[i]) != 3 {
			return domainmodel.Game{}, ErrInvalidBoardSize
		}
		for j := 0; j < 3; j++ {
			if game.Board[i][j] != int(domainmodel.Nought) && game.Board[i][j] != int(domainmodel.Cross) && game.Board[i][j] != int(domainmodel.Empty) {
				return domainmodel.Game{}, ErrInvalidCellValue
			}
			board[i][j] = domainmodel.Cell(game.Board[i][j])
		}
	}
	return domainmodel.Game{ID: game.ID, Board: board}, nil

}
func ToWeb(game domainmodel.Game) webmodel.Game {
	board := make([][]int, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = int(game.Board[i][j])
		}
	}
	return webmodel.Game{ID: game.ID, Board: board}
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
