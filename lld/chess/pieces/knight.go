package pieces

import "github.com/lld/chess/utils"

func NewKnight(c PieceColor, p Position) *Knight {
	return &Knight{BasePiece: BasePiece{color: c, position: p, ptype: KnightType}}
}

func (k *Knight) CanMove(board BoardView, from, to Position) bool {
	rowDiff := utils.Abs(to.GetRow() - from.GetRow())
	colDiff := utils.Abs(to.GetCol() - from.GetCol())

	// 1. L-shape check
	if !((rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2)) {
		return false
	}

	// 2. Destination check (cannot capture own piece)
	dest := board.GetPiece(to.GetRow(), to.GetCol())
	if dest != nil && dest.GetColor() == k.color {
		return false
	}

	return true
}

func (k *Knight) CanAttack(board BoardView, from, to Position) bool {
	return k.CanMove(board, from, to)
}

func (k *Knight) GetMoves(board BoardView) []Position {
	from := k.GetPosition()

	deltas := [][]int{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}

	var moves []Position

	for _, d := range deltas {
		r := from.GetRow() + d[0]
		c := from.GetCol() + d[1]

		if board.IsInsideBoard(NewPosition(r, c)) {
			moves = append(moves, NewPosition(r, c))
		}
	}

	return moves
}
