package game

import (
	"slices"
	"strings"

	"github.com/RishiKendai/System-Design/lld/chess/common"
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

// PieceAt returns the piece at the given square, or nil if empty.
func (g *Game) PieceAt(pos pieces.Position) pieces.Piece {
	return g.getPieceAtPosition(pos)
}

// MovesSnapshot returns a copy of the move history for display.
func (g *Game) MovesSnapshot() []common.Move {
	out := make([]common.Move, len(g.movesHistory))
	copy(out, g.movesHistory)
	return out
}

// ActivePiecesSummary returns a space-separated emoji line of active pieces for a side.
func (g *Game) ActivePiecesSummary(color pieces.PieceColor) string {
	pieceSet := g.whitePieces
	if color == pieces.Black {
		pieceSet = g.blackPieces
	}
	list := make([]pieces.Piece, 0, len(pieceSet))
	for p := range pieceSet {
		list = append(list, p)
	}
	slices.SortFunc(list, func(a, b pieces.Piece) int {
		oa, ob := piecePrintOrder(a.GetPieceType()), piecePrintOrder(b.GetPieceType())
		if oa != ob {
			return oa - ob
		}
		pa, pb := a.GetPosition(), b.GetPosition()
		if d := pa.GetRow() - pb.GetRow(); d != 0 {
			return d
		}
		return pa.GetCol() - pb.GetCol()
	})
	var b strings.Builder
	for i, p := range list {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(common.GetPieceEmoji(p))
	}
	return b.String()
}
