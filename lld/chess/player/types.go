package player

import (
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

type Player struct {
	id      string
	name    string
	color   pieces.PieceColor
	isHuman bool
}
