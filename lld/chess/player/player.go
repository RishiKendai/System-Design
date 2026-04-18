package player

import (
	"fmt"

	"github.com/RishiKendai/System-Design/lld/chess/common"
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

func NewPlayer(id, name string, color pieces.PieceColor, isHuman bool) *Player {
	return &Player{
		id:      id,
		name:    name,
		color:   color,
		isHuman: isHuman,
	}
}

func (p *Player) GetID() string {
	return p.id
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetColor() pieces.PieceColor {
	return p.color
}

func (p *Player) IsHuman() bool {
	return p.isHuman
}

// Notify implements [game.Notifier].
func (p *Player) Notify(gs common.GameState) {
	if gs.CurrentTurn != p.color {
		// Board UI
		gs.Board()
		fmt.Println("--------------------------------")
		fmt.Printf("Last move: ")
		fmt.Printf("%s %s -> %s", common.GetPieceEmoji(gs.LastMove.Piece), gs.LastMove.From.String(), gs.LastMove.To.String())
		if gs.LastMove.Captured != nil {
			fmt.Printf(" (captured: %s)", common.GetPieceEmoji(gs.LastMove.Captured))
		}
		if gs.LastMove.IsCheck {
			fmt.Printf(" (check)")
		}
		if gs.LastMove.IsMate {
			fmt.Printf(" (mate)")
		}
		fmt.Println("--------------------------------")
		return
	}
}
