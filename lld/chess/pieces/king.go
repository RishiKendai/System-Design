package pieces

import "github.com/lld/chess/utils"

func NewKing(c PieceColor, p Position) *King {
	return &King{BasePiece: BasePiece{color: c, position: p, ptype: KingType}, hasCastled: false}
}

func (k *King) CanMove(board BoardView, from, to Position) bool {
	rowDiff := utils.Abs(to.GetRow() - from.GetRow())
	colDiff := utils.Abs(to.GetCol() - from.GetCol())

	// 1. Must move only 1 step
	if rowDiff > 1 || colDiff > 1 {
		return false
	}

	// 2. Cannot stay in same position
	if rowDiff == 0 && colDiff == 0 {
		return false
	}

	// 3. Cannot capture own piece
	dest := board.GetPiece(to.GetRow(), to.GetCol())
	if dest != nil && dest.GetColor() == k.color {
		return false
	}

	return true
}

func (k *King) CanAttack(board BoardView, from, to Position) bool {
	rowDiff := utils.Abs(to.GetRow() - from.GetRow())
	colDiff := utils.Abs(to.GetCol() - from.GetCol())

	return rowDiff <= 1 && colDiff <= 1 && !(rowDiff == 0 && colDiff == 0)
}

func isCastlingMove(from, to Position) bool {
	return utils.Abs(to.GetCol()-from.GetCol()) == 2
}

func (k *King) HasMoved() bool {
	return k.hasMoved
}

func (k *King) MarkMoved() {
	k.hasMoved = true
}

func (k *King) GetMoves(board BoardView) []Position {
	from := k.GetPosition()

	deltas := [][]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}

	var moves []Position

	for _, d := range deltas {
		r := from.GetRow() + d[0]
		c := from.GetCol() + d[1]

		if board.IsInsideBoard(NewPosition(r, c)) {
			moves = append(moves, NewPosition(r, c))
		}
	}

	// Castling NOT included here (Game handles it)
	return moves
}
