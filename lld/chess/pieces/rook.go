package pieces

func NewRook(c PieceColor, p Position) *Rook {
	return &Rook{
		BasePiece:  BasePiece{color: c, position: p, ptype: RookType},
		hasCastled: false,
	}
}

func (r *Rook) CanMove(board BoardView, from, to Position) bool {
	// 1. Movement pattern
	if from.GetRow() != to.GetRow() && from.GetCol() != to.GetCol() {
		return false
	}

	// 2. Path blocking
	if isPathBlocked(board, from, to) {
		return false
	}

	// 3. Destination check
	dest := board.GetPiece(to.GetRow(), to.GetCol())
	if dest != nil && dest.GetColor() == r.color {
		return false
	}
	return true
}

func (r *Rook) CanAttack(board BoardView, from, to Position) bool {
	return r.CanMove(board, from, to)
}

func (r *Rook) HasMoved() bool {
	return r.hasMoved
}

func (r *Rook) MarkMoved() {
	r.hasMoved = true
}

func (r *Rook) GetMoves(board BoardView) []Position {
	from := r.GetPosition()

	deltas := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	var moves []Position

	for _, d := range deltas {
		row, col := from.GetRow(), from.GetCol()

		for {
			row += d[0]
			col += d[1]

			if !board.IsInsideBoard(NewPosition(row, col)) {
				break
			}

			pos := Position{row, col}
			moves = append(moves, pos)

			if !board.IsEmpty(pos) {
				break
			}
		}
	}
	return moves
}
