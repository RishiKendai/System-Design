package tui

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

func inGrid(r, c int) bool {
	return r >= 0 && r < 8 && c >= 0 && c < 8
}

func parseMove(line string) (pieces.Position, pieces.Position, error) {
	fields := strings.Fields(line)
	switch len(fields) {
	case 4:
		r0, err0 := strconv.Atoi(fields[0])
		c0, err1 := strconv.Atoi(fields[1])
		r1, err2 := strconv.Atoi(fields[2])
		c1, err3 := strconv.Atoi(fields[3])
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil {
			return pieces.Position{}, pieces.Position{}, fmt.Errorf("use four integers 0–7, or two squares like e2 e4")
		}
		if !inGrid(r0, c0) || !inGrid(r1, c1) {
			return pieces.Position{}, pieces.Position{}, fmt.Errorf("coordinates must be in 0..7")
		}
		return pieces.NewPosition(r0, c0), pieces.NewPosition(r1, c1), nil
	case 2:
		from, err := parseSquare(fields[0])
		if err != nil {
			return pieces.Position{}, pieces.Position{}, err
		}
		to, err := parseSquare(fields[1])
		if err != nil {
			return pieces.Position{}, pieces.Position{}, err
		}
		return from, to, nil
	default:
		return pieces.Position{}, pieces.Position{}, fmt.Errorf("use four numbers (r c r c) or two squares (e2 e4); help for more")
	}
}

func parseSquare(s string) (pieces.Position, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if len(s) < 2 {
		return pieces.Position{}, fmt.Errorf("invalid square %q", s)
	}
	file := rune(s[0])
	if file < 'a' || file > 'h' {
		return pieces.Position{}, fmt.Errorf("file must be a–h in %q", s)
	}
	col := int(file - 'a')
	rankStr := s[1:]
	for _, ch := range rankStr {
		if !unicode.IsDigit(ch) {
			return pieces.Position{}, fmt.Errorf("invalid rank in %q", s)
		}
	}
	rank, err := strconv.Atoi(rankStr)
	if err != nil || rank < 1 || rank > 8 {
		return pieces.Position{}, fmt.Errorf("rank must be 1–8 in %q", s)
	}
	row := 8 - rank
	return pieces.NewPosition(row, col), nil
}

func parsePieceColor(s string) (pieces.PieceColor, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "white", "w":
		return pieces.White, true
	case "black", "b":
		return pieces.Black, true
	default:
		return "", false
	}
}
