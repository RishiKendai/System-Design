package pieces

import "github.com/RishiKendai/System-Design/lld/chess/utils"

func NewBishop(c PieceColor, p Position) *Bishop {
	return &Bishop{BasePiece: BasePiece{color: c, position: p, ptype: BishopType}}
}

func (b *Bishop) CanMove(board BoardView, from, to Position) bool {
	rowDiff := utils.Abs(to.GetRow() - from.GetRow())
	colDiff := utils.Abs(to.GetCol() - from.GetCol())

	// 1. Must move diagonally
	if rowDiff != colDiff {
		return false
	}

	// 2. Path must be clear
	if isPathBlocked(board, from, to) {
		return false
	}

	// 3. Cannot capture own piece
	dest := board.GetPiece(to.GetRow(), to.GetCol())
	if dest != nil && dest.GetColor() == b.color {
		return false
	}

	return true
}

func (b *Bishop) CanAttack(board BoardView, from, to Position) bool {
	return b.CanMove(board, from, to)
}

func (b *Bishop) GetMoves(board BoardView) []Position {
	from := b.GetPosition()

	deltas := [][]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	var moves []Position

	for _, d := range deltas {
		r := from.GetRow() + d[0]
		c := from.GetCol() + d[1]

		for {
			r += d[0]
			c += d[1]

			if !board.IsInsideBoard(NewPosition(r, c)) {
				break
			}

			pos := NewPosition(r, c)
			moves = append(moves, pos)

			if !board.IsEmpty(pos) {
				break // stop at first piece
			}
		}
	}
	return moves
}
