package pieces

import (
	"fmt"

	"github.com/lld/chess/utils"
)

func NewPosition(r, c int) Position {
	return Position{row: r, col: c}
}

func (p Position) GetRow() int {
	return p.row
}

func (p Position) GetCol() int {
	return p.col
}

func (p Position) String() string {
	return fmt.Sprintf("%c%d", 'a'+p.col, 8-p.row)
}

func (b BasePiece) GetColor() PieceColor {
	return b.color
}

func (b BasePiece) GetPieceType() PieceType {
	return b.ptype
}

func (b BasePiece) GetPosition() Position {
	return b.position
}

func (b *BasePiece) SetPosition(p Position) {
	b.position = p
}

func isPathBlocked(board BoardView, from, to Position) bool {
	rowStep := utils.Sign(to.GetRow() - from.GetRow())
	colStep := utils.Sign(to.GetCol() - from.GetCol())

	currRow := from.GetRow() + rowStep
	currCol := from.GetCol() + colStep

	for currRow != to.GetRow() || currCol != to.GetCol() {
		if board.GetPiece(currRow, currCol) != nil {
			return true // blocked
		}

		currRow += rowStep
		currCol += colStep
	}

	return false
}
