package game

import (
	"fmt"
	"slices"

	"github.com/lld/chess/common"
	"github.com/lld/chess/pieces"
)

// piecePrintOrder: king first, then queen, then descending material (rook, bishop, knight, pawn).
func piecePrintOrder(t pieces.PieceType) int {
	switch t {
	case pieces.KingType:
		return 0
	case pieces.QueenType:
		return 1
	case pieces.RookType:
		return 2
	case pieces.BishopType:
		return 3
	case pieces.KnightType:
		return 4
	case pieces.PawnType:
		return 5
	default:
		return 6
	}
}

func (g *Game) printBorderedBoard() {
	const (
		h3   = "───"
		topL = "┌"
		topM = "┬"
		topR = "┐"
		midL = "├"
		midM = "┼"
		midR = "┤"
		botL = "└"
		botM = "┴"
		botR = "┘"
		pipe = "│"
	)

	hLine := func(left, mid, right string) {
		fmt.Print("   ")
		fmt.Print(left)
		for c := range 8 {
			if c > 0 {
				fmt.Print(mid)
			}
			fmt.Print(h3)
		}
		fmt.Println(right)
	}

	fmt.Print("     ")
	for c := range 8 {
		fmt.Printf(" %c  ", 'a'+rune(c))
	}
	fmt.Println()

	hLine(topL, topM, topR)
	for r := range 8 {
		if r > 0 {
			hLine(midL, midM, midR)
		}
		fmt.Printf(" %d ", 8-r)
		for c := range 8 {
			cell := " "
			if p := g.getPieceAtPosition(pieces.NewPosition(r, c)); p != nil {
				cell = common.GetPieceEmoji(p)
			}
			fmt.Printf("%s %s ", pipe, cell)
		}
		fmt.Println(pipe)
	}
	hLine(botL, botM, botR)

	fmt.Print("     ")
	for c := range 8 {
		fmt.Printf(" %c  ", 'a'+rune(c))
	}
	fmt.Println()
}

func (g *Game) printPieces(color pieces.PieceColor) {
	pieceSet := g.whitePieces
	if color == pieces.Black {
		pieceSet = g.blackPieces
	}
	list := make([]pieces.Piece, 0, len(pieceSet))
	for p := range pieceSet {
		list = append(list, p)
	}
	slices.SortFunc(list, func(a, b pieces.Piece) int {
		oa, ob := piecePrintOrder(a.GetPieceType()), piecePrintOrder(b.GetPieceType())
		if oa != ob {
			return oa - ob
		}
		pa, pb := a.GetPosition(), b.GetPosition()
		if d := pa.GetRow() - pb.GetRow(); d != 0 {
			return d
		}
		return pa.GetCol() - pb.GetCol()
	})
	for _, p := range list {
		fmt.Printf("%s ", common.GetPieceEmoji(p))
	}
}

func (g *Game) printMovesHistory() {
	fmt.Println("\nMoves History:")
	for _, m := range g.movesHistory {
		captured := "-"
		if m.Captured != nil {
			captured = common.GetPieceEmoji(m.Captured)
		}
		fmt.Println(common.GetPieceEmoji(m.Piece), " ", m.From, " ", m.To, " ", captured, " ", m.IsCheck, " ", m.IsMate)
	}
}

func (g *Game) printCurrentPlayer() {
	fmt.Printf("Current Player: %s\n", g.players[g.currentPlayer].GetName())
}

func (g *Game) printGameStatus() {
	fmt.Printf("Game Status: %s\n", g.gameStatus)
}
