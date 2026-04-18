package player

import (
	"github.com/lld/chess/pieces"
)

type Player struct {
	id      string
	name    string
	color   pieces.PieceColor
	isHuman bool
}
