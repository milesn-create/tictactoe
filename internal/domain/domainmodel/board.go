package domainmodel

type Cell int
type Board [3][3]Cell

const (
	Empty  Cell = 0
	Cross  Cell = 1
	Nought Cell = -1
)
