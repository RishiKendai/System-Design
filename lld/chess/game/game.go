package game

import (
	"fmt"

	"github.com/RishiKendai/System-Design/lld/chess/board"
	"github.com/RishiKendai/System-Design/lld/chess/common"
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
	"github.com/RishiKendai/System-Design/lld/chess/player"
)

func (g *Game) InitializePieces() {
	for r := range 8 {
		for c := range 8 {
			p := g.board.GetPiece(r, c)
			if p == nil {
				continue
			}

			if p.GetColor() == pieces.White {
				g.whitePieces[p] = true
			} else {
				g.blackPieces[p] = true
			}
		}
	}
}

func NewGame(board *board.Board, players []*player.Player, cPlayer int) *Game {

	g := &Game{
		board:         board,
		players:       players,
		whitePieces:   make(map[pieces.Piece]bool),
		blackPieces:   make(map[pieces.Piece]bool),
		movesHistory:  make([]common.Move, 0),
		currentPlayer: cPlayer,
		gameStatus:    GameStatusInProgress,
	}
	return g
}

func AddNotifier(g *Game, players []*player.Player) {
	for _, p := range players {
		g.notifier = append(g.notifier, p)
	}
}

func (g *Game) notifyPlayers(move common.Move) {
	for _, players := range g.notifier {
		players.Notify(common.GameState{
			LastMove:    &move,
			CurrentTurn: g.players[g.currentPlayer].GetColor(),
			Status:      string(g.gameStatus),
			Board:       g.printBorderedBoard,
		})
	}
}

// MakeMove applies a move. promotion must be non-nil only for pawn moves to the
// last rank; otherwise pass nil. If promotion is required but nil, returns ErrPromotionRequired.
func (g *Game) MakeMove(from, to pieces.Position, promotion *pieces.PieceType) error {
	piece := g.board.GetPiece(from.GetRow(), from.GetCol())
	if piece == nil {
		return fmt.Errorf("no piece")
	}
	mType := g.detectMoveType(from, to, g.board.GetPiece(from.GetRow(), from.GetCol()))
	var capturedPiece pieces.Piece

	if !g.isValidMove(g.board, from, to, mType) {
		return fmt.Errorf("invalid move")
	}

	if promotion != nil && !g.isPromotion(from, to, piece) {
		return fmt.Errorf("unexpected promotion for this move")
	}

	if mType == Normal && g.isPromotion(from, to, piece) && promotion == nil {
		return ErrPromotionRequired
	}

	switch mType {
	case Castling:
		g.handleCastling(from, to, g.board.GetPiece(from.GetRow(), from.GetCol()))
	case EnPassant:
		capturedPiece = g.handleEnPassant(from, to, g.board.GetPiece(from.GetRow(), from.GetCol()))
	case Normal:
		capturedPiece = g.board.Apply(from, to)
		g.removeActivePiece(capturedPiece)
	}

	move := common.Move{
		Piece:    piece,
		From:     from,
		To:       to,
		Captured: capturedPiece,
		IsCheck:  false,
		IsMate:   false,
	}

	if mType == Normal && g.isPromotion(from, to, piece) {
		if err := g.applyPromotionAt(to, piece, *promotion); err != nil {
			return err
		}
	}

	var statePiece pieces.Piece
	switch mType {
	case Normal:
		if g.isPromotion(from, to, piece) {
			statePiece = g.board.GetPiece(to.GetRow(), to.GetCol())
		} else {
			statePiece = piece
		}
	default:
		statePiece = piece
	}
	g.updatePieceState(statePiece)
	g.updateGameStatus(&move)
	g.movesHistory = append(g.movesHistory, move)
	g.notifyPlayers(move)
	if g.gameStatus != GameStatusWhiteWin && g.gameStatus != GameStatusBlackWin {
		g.switchTurn()
	}
	return nil
}

func (g *Game) PrintState() {
	fmt.Println("Game State:")
	fmt.Println("Board:")
	g.printBorderedBoard()
	fmt.Printf("White Pieces: ")
	g.printPieces(pieces.White)
	fmt.Printf("\nBlack Pieces: ")
	g.printPieces(pieces.Black)
	g.printMovesHistory()
	g.printCurrentPlayer()
	g.printGameStatus()
}

func (g *Game) updateGameStatus(move *common.Move) {
	currentPlayer := g.players[g.currentPlayer]

	var opponentPlayer *player.Player
	if currentPlayer.GetColor() == pieces.White {
		opponentPlayer = g.players[1]
	} else {
		opponentPlayer = g.players[0]
	}

	if g.isCheckmate(opponentPlayer.GetColor()) {
		move.IsMate = true
		if currentPlayer.GetColor() == pieces.White {
			g.gameStatus = GameStatusWhiteWin
		} else {
			g.gameStatus = GameStatusBlackWin
		}
		return
	}

	if g.isKingUnderAttack(opponentPlayer.GetColor()) {
		move.IsCheck = true
		g.gameStatus = GameStatusCheck
		return
	}

	g.gameStatus = GameStatusInProgress
}

func (g *Game) updatePieceState(piece pieces.Piece) {
	switch piece.GetPieceType() {
	case pieces.PawnType:
		pawn, ok := piece.(*pieces.Pawn)
		if !ok {
			return
		}
		pawn.MarkMoved()
	case pieces.RookType:
		rook, ok := piece.(*pieces.Rook)
		if !ok {
			return
		}
		rook.MarkMoved()
	case pieces.KingType:
		k, ok := piece.(*pieces.King)
		if !ok {
			return
		}
		k.MarkMoved()
	default:
		return
	}
}

// CurrentPlayerName returns the player to move.
func (g *Game) CurrentPlayerName() string {
	return g.players[g.currentPlayer].GetName()
}

// CurrentTurnColor is WHITE or BLACK for the side to move.
func (g *Game) CurrentTurnColor() pieces.PieceColor {
	return g.players[g.currentPlayer].GetColor()
}

// Status is the game phase after the last completed move (or initial state).
func (g *Game) Status() GameStatus {
	return g.gameStatus
}

// IsFinished reports whether play should stop (win, draw, or stalemate).
func (g *Game) IsFinished() bool {
	switch g.gameStatus {
	case GameStatusWhiteWin, GameStatusBlackWin, GameStatusDraw, GameStatusStalemate:
		return true
	default:
		return false
	}
}
