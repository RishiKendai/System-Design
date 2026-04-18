package game

import "errors"

// ErrPromotionRequired is returned when a pawn reaches the last rank and the
// caller must supply a promotion piece type on the next MakeMove attempt.
var ErrPromotionRequired = errors.New("promotion required: choose Q, R, B, or K (knight)")
