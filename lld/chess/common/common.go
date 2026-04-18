package common

import "github.com/RishiKendai/System-Design/lld/chess/pieces"

func GetPieceEmoji(p pieces.Piece) string {
	switch p.GetPieceType() {
	case pieces.PawnType:
		if p.GetColor() == pieces.White {
			return "\u2659"
		}
		return "\u265f"
	case pieces.RookType:
		if p.GetColor() == pieces.White {
			return "\u2656"
		}
		return "\u265c"
	case pieces.KnightType:
		if p.GetColor() == pieces.White {
			return "\u2658"
		}
		return "\u265e"
	case pieces.BishopType:
		if p.GetColor() == pieces.White {
			return "\u2657"
		}
		return "\u265d"
	case pieces.QueenType:
		if p.GetColor() == pieces.White {
			return "\u2655"
		}
		return "\u265b"
	case pieces.KingType:
		if p.GetColor() == pieces.White {
			return "\u2654"
		}
		return "\u265a"
	default:
		return " "
	}
}
