package game

import (
	"github.com/RishiKendai/System-Design/lld/chess/board"
	"github.com/RishiKendai/System-Design/lld/chess/common"
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
	"github.com/RishiKendai/System-Design/lld/chess/player"
)

type GameStatus string

const (
	GameStatusInProgress GameStatus = "IN_PROGRESS"
	GameStatusWhiteWin   GameStatus = "WHITE_WIN"
	GameStatusBlackWin   GameStatus = "BLACK_WIN"
	GameStatusDraw       GameStatus = "DRAW"
	GameStatusStalemate  GameStatus = "STALEMATE"
	GameStatusCheck      GameStatus = "CHECK"
)

type Notifier interface {
	Notify(gameState common.GameState)
}

type Game struct {
	board         *board.Board
	players       []*player.Player
	whitePieces   map[pieces.Piece]bool // active white pieces
	blackPieces   map[pieces.Piece]bool // active black pieces
	movesHistory  []common.Move         // moves history
	currentPlayer int                   // 0 for white, 1 for black
	gameStatus    GameStatus            // game status
	notifier      []Notifier            // notifier
}

type MoveUndo struct {
	move struct {
		From pieces.Position
		To   pieces.Position
	}
	movedPiece       pieces.Piece
	capturedPiece    pieces.Piece
	prevPosition     pieces.Position
	capturedPosition pieces.Position
}

type MoveType int

const (
	Normal MoveType = iota
	Castling
	EnPassant
)
