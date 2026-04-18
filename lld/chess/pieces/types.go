package pieces

type Position struct {
	row int
	col int
}

type PieceColor string

type PieceType string

const (
	PawnType   PieceType = "PAWN"
	RookType   PieceType = "ROOK"
	KnightType PieceType = "KNIGHT"
	BishopType PieceType = "BISHOP"
	QueenType  PieceType = "QUEEN"
	KingType   PieceType = "KING"
)

const (
	White PieceColor = "WHITE"
	Black PieceColor = "BLACK"
)

type BoardView interface {
	GetPiece(r, c int) Piece
	IsInsideBoard(pos Position) bool
	IsEmpty(pos Position) bool
}

type Piece interface {
	CanMove(board BoardView, from, to Position) bool
	GetPosition() Position
	SetPosition(p Position)
	GetPieceType() PieceType
	GetColor() PieceColor
	CanAttack(board BoardView, from, to Position) bool
	GetMoves(board BoardView) []Position
}

type BasePiece struct {
	color    PieceColor
	position Position
	ptype    PieceType
}

type Pawn struct {
	BasePiece
	hasMoved   bool
	isPromoted bool
}

type Rook struct {
	BasePiece
	hasMoved   bool
	hasCastled bool
}

type Knight struct {
	BasePiece
}

type Bishop struct {
	BasePiece
}

type Queen struct {
	BasePiece
}

type King struct {
	BasePiece
	hasMoved   bool
	hasCastled bool
}
