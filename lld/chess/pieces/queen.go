package pieces

import "github.com/lld/chess/utils"

func NewQueen(c PieceColor, p Position) *Queen {
	return &Queen{BasePiece: BasePiece{color: c, position: p, ptype: QueenType}}
}

func (q *Queen) CanMove(board BoardView, from, to Position) bool {
	rowDiff := utils.Abs(to.GetRow() - from.GetRow())
	colDiff := utils.Abs(to.GetCol() - from.GetCol())

	// 1. Check valid movement pattern
	isStraight := from.GetRow() == to.GetRow() || from.GetCol() == to.GetCol()
	isDiagonal := rowDiff == colDiff

	if !isStraight && !isDiagonal {
		return false
	}

	// 2. Path must be clear
	if isPathBlocked(board, from, to) {
		return false
	}

	// 3. Cannot capture own piece
	dest := board.GetPiece(to.GetRow(), to.GetCol())
	if dest != nil && dest.GetColor() == q.color {
		return false
	}

	return true
}

func (q *Queen) CanAttack(board BoardView, from, to Position) bool {
	return q.CanMove(board, from, to)
}

func (q *Queen) GetMoves(board BoardView) []Position {
	from := q.GetPosition()

	deltas := [][]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, // rook
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // bishop
	}

	var moves []Position

	for _, d := range deltas {
		r, c := from.GetRow(), from.GetCol()

		for {
			r += d[0]
			c += d[1]

			if !board.IsInsideBoard(NewPosition(r, c)) {
				break
			}

			pos := NewPosition(r, c)
			moves = append(moves, pos)

			if !board.IsEmpty(pos) {
				break
			}
		}
	}
	return moves
}
