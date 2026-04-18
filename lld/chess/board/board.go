package board

import (
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

func NewBoard() *Board {
	b := &Board{}

	// Pawns
	for col := range 8 {
		b.grid[1][col] = pieces.NewPawn(pieces.Black, pieces.NewPosition(1, col))
		b.grid[6][col] = pieces.NewPawn(pieces.White, pieces.NewPosition(6, col))
	}

	// Black pieces (row 0)
	b.grid[0][0] = pieces.NewRook(pieces.Black, pieces.NewPosition(0, 0))
	b.grid[0][1] = pieces.NewKnight(pieces.Black, pieces.NewPosition(0, 1))
	b.grid[0][2] = pieces.NewBishop(pieces.Black, pieces.NewPosition(0, 2))
	b.grid[0][3] = pieces.NewQueen(pieces.Black, pieces.NewPosition(0, 3))
	b.grid[0][4] = pieces.NewKing(pieces.Black, pieces.NewPosition(0, 4))
	b.grid[0][5] = pieces.NewBishop(pieces.Black, pieces.NewPosition(0, 5))
	b.grid[0][6] = pieces.NewKnight(pieces.Black, pieces.NewPosition(0, 6))
	b.grid[0][7] = pieces.NewRook(pieces.Black, pieces.NewPosition(0, 7))

	// White pieces (row 7)
	b.grid[7][0] = pieces.NewRook(pieces.White, pieces.NewPosition(7, 0))
	b.grid[7][1] = pieces.NewKnight(pieces.White, pieces.NewPosition(7, 1))
	b.grid[7][2] = pieces.NewBishop(pieces.White, pieces.NewPosition(7, 2))
	b.grid[7][3] = pieces.NewQueen(pieces.White, pieces.NewPosition(7, 3))
	b.grid[7][4] = pieces.NewKing(pieces.White, pieces.NewPosition(7, 4))
	b.grid[7][5] = pieces.NewBishop(pieces.White, pieces.NewPosition(7, 5))
	b.grid[7][6] = pieces.NewKnight(pieces.White, pieces.NewPosition(7, 6))
	b.grid[7][7] = pieces.NewRook(pieces.White, pieces.NewPosition(7, 7))

	return b
}

func (b *Board) IsPathBlocked(from, to pieces.Position) bool {
	panic("Not implemented.")
}

func (b *Board) GetPiece(r, c int) pieces.Piece {
	return b.grid[r][c]
}

func (b *Board) SetPiece(r, c int, piece pieces.Piece) {
	b.grid[r][c] = piece
}

func (b *Board) IsInsideBoard(pos pieces.Position) bool {
	return pos.GetRow() >= 0 && pos.GetRow() < 8 && pos.GetCol() >= 0 && pos.GetCol() < 8
}

func (b *Board) IsEmpty(pos pieces.Position) bool {
	return b.grid[pos.GetRow()][pos.GetCol()] == nil
}

func (b *Board) Clear(r, c int) {
	b.grid[r][c] = nil
}

func (b *Board) Apply(from, to pieces.Position) pieces.Piece {
	piece := b.grid[from.GetRow()][from.GetCol()]
	capturedPiece := b.grid[to.GetRow()][to.GetCol()]
	// move piece
	b.grid[to.GetRow()][to.GetCol()] = piece
	b.grid[from.GetRow()][from.GetCol()] = nil

	// update piece position
	piece.SetPosition(to)

	return capturedPiece
}
