package board

import "github.com/lld/chess/pieces"

type Board struct {
	grid [8][8]pieces.Piece
}
