package domainmodel

type GameState int

const (
	StateWaitingPlayers GameState = iota
	StatePlayerTurn
	StateDraw
	StateWon
	StateArchived
)
