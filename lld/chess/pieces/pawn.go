package pieces

import (
	"fmt"

	"github.com/lld/chess/utils"
)

func NewPawn(c PieceColor, p Position) *Pawn {
	return &Pawn{
		BasePiece: BasePiece{color: c, position: p, ptype: PawnType},
		hasMoved:  false,
	}
}

func (p *Pawn) CanMove(board BoardView, from, to Position) bool {
	rowDiff := to.GetRow() - from.GetRow()
	colDiff := to.GetCol() - from.GetCol()

	direction := p.GetDirection() // +1 or -1

	// 1. Forward move (1 step)
	if colDiff == 0 && rowDiff == direction {
		if board.GetPiece(to.GetRow(), to.GetCol()) == nil {
			return true
		}
	}

	// 2. First move (2 steps)
	if colDiff == 0 && rowDiff == 2*direction && !p.hasMoved {
		midRow := from.GetRow() + direction

		// both cells must be empty
		if board.GetPiece(midRow, from.GetCol()) == nil &&
			board.GetPiece(to.GetRow(), to.GetCol()) == nil {
			return true
		}
	}

	// 3. Capture (diagonal)
	if utils.Abs(colDiff) == 1 && rowDiff == direction {
		dest := board.GetPiece(to.GetRow(), to.GetCol())
		if dest != nil && dest.GetColor() != p.color {
			return true
		}
	}
	fmt.Println("Invalid move...")
	return false
}

func (p *Pawn) CanAttack(board BoardView, from, to Position) bool {
	rowDiff := to.GetRow() - from.GetRow()
	colDiff := utils.Abs(to.GetCol() - from.GetCol())

	direction := p.GetDirection()

	return rowDiff == direction && colDiff == 1
}

func (p *Pawn) GetDirection() int {
	if p.color == White {
		return -1 // moving up
	}
	return 1 // moving down
}

func (p *Pawn) HasMoved() bool {
	return p.hasMoved
}

func (p *Pawn) MarkMoved() {
	p.hasMoved = true
}

func (p *Pawn) GetMoves(board BoardView) []Position {
	from := p.GetPosition()

	dir := p.GetDirection()
	var moves []Position

	// 1. forward move
	r := from.GetRow() + dir
	c := from.GetCol()

	if board.IsInsideBoard(NewPosition(r, c)) && board.IsEmpty(NewPosition(r, c)) {
		moves = append(moves, NewPosition(r, c))

		// 2. double move
		if !p.hasMoved {
			r2 := from.GetRow() + 2*dir
			if board.IsInsideBoard(NewPosition(r2, c)) && board.IsEmpty(NewPosition(r2, c)) {
				moves = append(moves, NewPosition(r2, c))
			}
		}
	}

	// 3. diagonal capture candidates
	for _, dc := range []int{-1, 1} {
		r := from.GetRow() + dir
		c := from.GetCol() + dc

		if board.IsInsideBoard(NewPosition(r, c)) {
			moves = append(moves, NewPosition(r, c))
		}
	}

	// NOTE:
	// - En passant NOT handled here
	// - Capture validation done via CanMove / Game

	return moves
}
