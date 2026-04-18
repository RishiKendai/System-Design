package game

import (
	"github.com/lld/chess/board"
	"github.com/lld/chess/pieces"
)

func (g *Game) isValidMove(board *board.Board, from, to pieces.Position, mType MoveType) bool {
	if g.currentPlayer < 0 || g.currentPlayer >= len(g.players) {
		return false
	}
	side := g.players[g.currentPlayer].GetColor()
	return g.isValidMoveForSide(side, board, from, to, mType)
}

// isValidMoveForSide validates a move as if it were played by `side` (used for
// checkmate / check escape without depending on whose turn it is on the clock).
func (g *Game) isValidMoveForSide(side pieces.PieceColor, board *board.Board, from, to pieces.Position, mType MoveType) bool {
	currentPiece := board.GetPiece(from.GetRow(), from.GetCol())
	if currentPiece == nil || currentPiece.GetColor() != side {
		return false
	}

	if !board.IsInsideBoard(to) {
		return false
	}

	dest := board.GetPiece(to.GetRow(), to.GetCol())
	if dest != nil && dest.GetColor() == currentPiece.GetColor() {
		return false
	}

	switch mType {
	case Castling:
		return g.validateCastling(from, to, side)
	case EnPassant:
		if !g.validateEnPassant(from, to, currentPiece) {
			return false
		}
		return g.isSafeMoveEnPassant(from, to, currentPiece, side)
	case Normal:
		if !currentPiece.CanMove(g.board, from, to) {
			return false
		}
		return g.isSafeMove(from, to, currentPiece, side)
	}
	return false
}

func (g *Game) isSafeMove(from, to pieces.Position, piece pieces.Piece, player pieces.PieceColor) bool {
	undo := g.simulateMove(g.board, from, to)

	inCheck := g.isKingUnderAttack(player)

	g.undoMove(undo)

	return !inCheck
}

func (g *Game) simulateMove(board *board.Board, from, to pieces.Position) (undo MoveUndo) {
	moved := board.GetPiece(from.GetRow(), from.GetCol())
	captured := board.GetPiece(to.GetRow(), to.GetCol())

	undo = MoveUndo{
		move: struct {
			From pieces.Position
			To   pieces.Position
		}{From: from, To: to},
		movedPiece:    moved,
		capturedPiece: captured,
		prevPosition:  from,
	}
	if captured != nil {
		undo.capturedPosition = pieces.NewPosition(to.GetRow(), to.GetCol())
	}

	board.SetPiece(to.GetRow(), to.GetCol(), moved)
	board.SetPiece(from.GetRow(), from.GetCol(), nil)
	if moved != nil {
		moved.SetPosition(to)
	}

	return undo
}
