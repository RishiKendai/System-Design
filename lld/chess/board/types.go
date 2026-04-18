package board

import "github.com/RishiKendai/System-Design/lld/chess/pieces"

type Board struct {
	grid [8][8]pieces.Piece
}
