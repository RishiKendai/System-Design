package common

import (
	"github.com/lld/chess/pieces"
)

type Move struct {
	Piece    pieces.Piece
	From     pieces.Position
	To       pieces.Position
	Captured pieces.Piece
	IsCheck  bool
	IsMate   bool
}

type GameState struct {
	LastMove    *Move
	CurrentTurn pieces.PieceColor
	Status      string
	Board       func()
}

type Board struct {
	Grid [8][8]pieces.Piece
}
