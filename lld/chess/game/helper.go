package game

import (
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

func (g *Game) isSquareUnderAttack(pos pieces.Position, player pieces.PieceColor) bool {
	var opponentPieces map[pieces.Piece]bool

	opponentPieces = g.whitePieces
	if player == pieces.White {
		opponentPieces = g.blackPieces
	}

	for p := range opponentPieces {
		from := p.GetPosition()
		if g.board.GetPiece(from.GetRow(), from.GetCol()) != p {
			continue
		}
		if p.CanAttack(g.board, from, pos) {
			return true
		}
	}

	return false
}

func (g *Game) isKingUnderAttack(color pieces.PieceColor) bool {
	kingPos := g.findKing(color)

	var attackers map[pieces.Piece]bool
	if color == pieces.White {
		attackers = g.blackPieces
	} else {
		attackers = g.whitePieces
	}

	for p := range attackers {
		pos := p.GetPosition()
		if g.board.GetPiece(pos.GetRow(), pos.GetCol()) != p {
			continue
		}
		if p.CanAttack(g.board, pos, kingPos) {
			return true
		}
	}

	return false
}

func (g *Game) findKing(color pieces.PieceColor) pieces.Position {
	var piecesMap map[pieces.Piece]bool

	if color == pieces.White {
		piecesMap = g.whitePieces
	} else {
		piecesMap = g.blackPieces
	}

	for p := range piecesMap {
		if _, ok := p.(*pieces.King); ok {
			return p.GetPosition()
		}
	}
	return pieces.Position{}
}

func (g *Game) getPieceAtPosition(pos pieces.Position) pieces.Piece {
	for p := range g.whitePieces {
		if isSamePosition(p.GetPosition(), pos) {
			return p
		}
	}
	for p := range g.blackPieces {
		if isSamePosition(p.GetPosition(), pos) {
			return p
		}
	}
	return nil
}

func (g *Game) isCurrentPlayerPiece(piece pieces.Piece) bool {
	if g.currentPlayer == 0 {
		return piece.GetColor() == pieces.White
	}
	return piece.GetColor() == pieces.Black
}

func (g *Game) removeActivePiece(piece pieces.Piece) {
	if piece == nil {
		return
	}
	if piece.GetColor() == pieces.White {
		delete(g.whitePieces, piece)
		return
	}
	delete(g.blackPieces, piece)
}

func (g *Game) setActivePiece(piece pieces.Piece) {
	if piece == nil {
		return
	}
	if piece.GetColor() == pieces.White {
		g.whitePieces[piece] = true
		return
	}
	g.blackPieces[piece] = true
}

func isSamePosition(a, b pieces.Position) bool {
	return a.GetRow() == b.GetRow() && a.GetCol() == b.GetCol()
}

func (g *Game) undoMove(undo MoveUndo) {
	from := undo.move.From
	to := undo.move.To

	piece := undo.movedPiece

	// revert board
	g.board.SetPiece(from.GetRow(), from.GetCol(), piece)
	g.board.SetPiece(to.GetRow(), to.GetCol(), undo.capturedPiece)

	// revert position
	piece.SetPosition(from)

	// restore captured piece position (if any)
	if undo.capturedPiece != nil {
		undo.capturedPiece.SetPosition(undo.capturedPosition)
	}
}

func (g *Game) undoBoard(undo MoveUndo) {
	g.undoMove(undo)
	// restore active pieces
	// restore pawn flags
	// restore castling flags
}

func (g *Game) detectMoveType(from, to pieces.Position, piece pieces.Piece) MoveType {
	if g.isCastling(from, to, piece) {
		return Castling
	}

	if g.isEnPassant(from, to, piece) {
		return EnPassant
	}

	return Normal
}

func (g *Game) isCheckmate(player pieces.PieceColor) bool {
	// 1. Must be in check
	if !g.isKingUnderAttack(player) {
		return false
	}

	piecesMap := g.whitePieces
	if player == pieces.Black {
		piecesMap = g.blackPieces
	}

	// 2. Try any pseudo-legal move: capture attacker, block ray, or move king
	for p := range piecesMap {
		from := p.GetPosition()
		if g.board.GetPiece(from.GetRow(), from.GetCol()) != p {
			continue
		}

		candidates := p.GetMoves(g.board)

		for _, to := range candidates {
			moveType := g.detectMoveType(from, to, p)

			if g.isValidMoveForSide(player, g.board, from, to, moveType) {
				return false // escape: capture checker, interpose, or king move
			}
		}
	}

	return true
}

func (g *Game) switchTurn() {
	if g.currentPlayer == 0 {
		g.currentPlayer = 1
	} else {
		g.currentPlayer = 0
	}
}
