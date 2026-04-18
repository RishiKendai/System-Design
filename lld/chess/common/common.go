package common

import "github.com/lld/chess/pieces"

func GetPieceEmoji(p pieces.Piece) string {
	switch p.GetPieceType() {
	case pieces.PawnType:
		if p.GetColor() == pieces.White {
			return "\u265f"
		}
		return "\u2659"
	case pieces.RookType:
		if p.GetColor() == pieces.White {
			return "\u265c"
		}
		return "\u2656"
	case pieces.KnightType:
		if p.GetColor() == pieces.White {
			return "\u265e"
		}
		return "\u2658"
	case pieces.BishopType:
		if p.GetColor() == pieces.White {
			return "\u265d"
		}
		return "\u2657"
	case pieces.QueenType:
		if p.GetColor() == pieces.White {
			return "\u265b"
		}
		return "\u2655"
	case pieces.KingType:
		if p.GetColor() == pieces.White {
			return "\u265a"
		}
		return "\u2654"
	default:
		return " "
	}
}
