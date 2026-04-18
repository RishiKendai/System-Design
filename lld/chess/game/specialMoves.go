package game

import (
	"fmt"

	"github.com/RishiKendai/System-Design/lld/chess/pieces"
	"github.com/RishiKendai/System-Design/lld/chess/utils"
)

func (g *Game) applyPromotionAt(to pieces.Position, pawn pieces.Piece, typ pieces.PieceType) error {
	var promoted pieces.Piece
	switch typ {
	case pieces.QueenType:
		promoted = pieces.NewQueen(pawn.GetColor(), to)
	case pieces.RookType:
		promoted = pieces.NewRook(pawn.GetColor(), to)
	case pieces.BishopType:
		promoted = pieces.NewBishop(pawn.GetColor(), to)
	case pieces.KnightType:
		promoted = pieces.NewKnight(pawn.GetColor(), to)
	default:
		return fmt.Errorf("invalid promotion: use Q, R, B, or K (knight)")
	}

	g.removeActivePiece(pawn)
	g.board.SetPiece(to.GetRow(), to.GetCol(), promoted)
	promoted.SetPosition(to)
	g.setActivePiece(promoted)
	return nil
}

func (g *Game) isCastling(from, to pieces.Position, piece pieces.Piece) bool {
	_, ok := piece.(*pieces.King)
	return ok && utils.Abs(to.GetCol()-from.GetCol()) == 2 && to.GetRow() == from.GetRow()
}

func (g *Game) validateCastling(from, to pieces.Position, player pieces.PieceColor) bool {
	if g.isKingUnderAttack(player) {
		return false
	}

	kPiece := g.board.GetPiece(from.GetRow(), from.GetCol())
	king, ok := kPiece.(*pieces.King)
	if !ok || king.GetColor() != player {
		return false
	}

	rowDiff := to.GetRow() - from.GetRow()
	colDiff := to.GetCol() - from.GetCol()
	if rowDiff != 0 || utils.Abs(colDiff) != 2 {
		return false
	}
	step := colDiff / 2 // +1 kingside (toward h-file), -1 queenside

	rookFromCol := 0
	if step > 0 {
		rookFromCol = 7
	}
	rPiece := g.board.GetPiece(from.GetRow(), rookFromCol)
	rook, ok := rPiece.(*pieces.Rook)
	if !ok || rook.GetColor() != player {
		return false
	}

	if king.HasMoved() || rook.HasMoved() {
		return false
	}

	kCol, rCol := from.GetCol(), rookFromCol
	startCol := min(kCol, rCol)
	endCol := max(kCol, rCol)
	for c := startCol + 1; c < endCol; c++ {
		if g.board.GetPiece(from.GetRow(), c) != nil {
			return false
		}
	}

	// King may not castle through or into check (starting square already not in check).
	for _, c := range []int{from.GetCol(), from.GetCol() + step, to.GetCol()} {
		pos := pieces.NewPosition(from.GetRow(), c)
		if g.isSquareUnderAttack(pos, player) {
			return false
		}
	}

	return true
}

func (g *Game) handleCastling(from, to pieces.Position, piece pieces.Piece) error {
	king, ok := piece.(*pieces.King)
	if !ok {
		return fmt.Errorf("castling: not a king")
	}

	row := from.GetRow()
	colDiff := to.GetCol() - from.GetCol()
	if utils.Abs(colDiff) != 2 || to.GetRow() != row {
		return fmt.Errorf("castling: invalid king destination")
	}
	step := colDiff / 2

	rookFromCol := 0
	if step > 0 {
		rookFromCol = 7
	}
	rPiece := g.board.GetPiece(row, rookFromCol)
	rook, ok := rPiece.(*pieces.Rook)
	if !ok {
		return fmt.Errorf("castling: rook missing")
	}

	rookToCol := to.GetCol() - step

	g.board.Clear(from.GetRow(), from.GetCol())
	g.board.SetPiece(to.GetRow(), to.GetCol(), king)
	king.SetPosition(to)

	g.board.Clear(row, rookFromCol)
	g.board.SetPiece(row, rookToCol, rook)
	rook.SetPosition(pieces.NewPosition(row, rookToCol))

	rook.MarkMoved()

	return nil
}

func (g *Game) validateEnPassant(from, to pieces.Position, piece pieces.Piece) bool {
	pawn, ok := piece.(*pieces.Pawn)
	if !ok {
		return false
	}

	// Must move diagonally by 1
	rowDiff := to.GetRow() - from.GetRow()
	colDiff := to.GetCol() - from.GetCol()

	dir := pawn.GetDirection()

	if rowDiff != dir || utils.Abs(colDiff) != 1 {
		return false
	}

	// Destination must be EMPTY
	if g.board.GetPiece(to.GetRow(), to.GetCol()) != nil {
		return false
	}

	// Must have at least one prev move
	if len(g.movesHistory) == 0 {
		return false
	}

	lastMove := g.movesHistory[len(g.movesHistory)-1]

	// last move piece must be a pawn
	lastPiece := g.board.GetPiece(lastMove.To.GetRow(), lastMove.To.GetCol())
	lastPawn, ok := lastPiece.(*pieces.Pawn)
	if !ok {
		return false
	}

	if lastPawn.GetColor() == pawn.GetColor() {
		return false
	}

	// Last move must be 2-step forward
	lastRowDiff := lastMove.To.GetRow() - lastMove.From.GetRow()
	if utils.Abs(lastRowDiff) != 2 {
		return false
	}

	// The pawn must be adjacent to current pawn
	if lastMove.To.GetRow() != from.GetRow() {
		return false
	}

	if lastMove.To.GetCol() != to.GetCol() {
		return false
	}

	return true
}

func (g *Game) isEnPassant(from, to pieces.Position, piece pieces.Piece) bool {
	pawn, ok := piece.(*pieces.Pawn)
	if !ok {
		return false
	}

	rowDiff := to.GetRow() - from.GetRow()
	colDiff := to.GetCol() - from.GetCol()

	if rowDiff != pawn.GetDirection() || utils.Abs(colDiff) != 1 {
		return false
	}

	// Same geometry as a diagonal capture, but en passant always lands on an empty square.
	if g.board.GetPiece(to.GetRow(), to.GetCol()) != nil {
		return false
	}

	return true
}

func (g *Game) isSafeMoveEnPassant(from, to pieces.Position, piece pieces.Piece, player pieces.PieceColor) bool {
	// Pawn captured en passant sits on the same rank as 'from', on the file of 'to'.
	capturedPos := pieces.NewPosition(from.GetRow(), to.GetCol())
	captured := g.board.GetPiece(capturedPos.GetRow(), capturedPos.GetCol())

	// ---- SIMULATE ----
	g.board.SetPiece(to.GetRow(), to.GetCol(), piece)
	g.board.Clear(from.GetRow(), from.GetCol())
	piece.SetPosition(to)
	g.board.Clear(capturedPos.GetRow(), capturedPos.GetCol())

	inCheck := g.isKingUnderAttack(player)

	// ---- UNDO ----
	g.board.SetPiece(from.GetRow(), from.GetCol(), piece)
	g.board.Clear(to.GetRow(), to.GetCol())
	piece.SetPosition(from)
	if captured != nil {
		g.board.SetPiece(capturedPos.GetRow(), capturedPos.GetCol(), captured)
	}

	return !inCheck
}

func (g *Game) handleEnPassant(from, to pieces.Position, piece pieces.Piece) pieces.Piece {
	// move pawn
	g.board.SetPiece(to.GetRow(), to.GetCol(), piece)
	g.board.Clear(from.GetRow(), from.GetCol())
	piece.SetPosition(to)

	// remove captured pawn
	capturedPos := pieces.NewPosition(from.GetRow(), to.GetCol())
	capturedPiece := g.board.GetPiece(capturedPos.GetRow(), capturedPos.GetCol())
	g.removeActivePiece(capturedPiece)
	g.board.Clear(capturedPos.GetRow(), capturedPos.GetCol())
	return capturedPiece
}

func (g *Game) isPromotion(from, to pieces.Position, piece pieces.Piece) bool {
	pawn, ok := piece.(*pieces.Pawn)
	if !ok {
		return false
	}

	if pawn.GetColor() == pieces.White && to.GetRow() == 0 {
		return true
	}
	if pawn.GetColor() == pieces.Black && to.GetRow() == 7 {
		return true
	}

	return false
}

