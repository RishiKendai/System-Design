package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	gamemanager "github.com/lld/chess/gameManager"
	"github.com/lld/chess/pieces"
	"github.com/lld/chess/player"
)

const banner = `
--------------------------------
  Chess CLI
--------------------------------
Each player: name + color (White or Black).
You need exactly one White and one Black; White moves first.
Commands: quit (or q) — exit anytime
          help — move syntax
--------------------------------
`

const helpText = `
Move formats (choose one):
  • Four numbers: fromRow fromCol toRow toCol  (each 0–7)
    Rows: 0 = top (Black's back rank), 7 = bottom (White's back rank)
    Cols: 0 = a-file … 7 = h-file
    Example: 6 4 5 4  (pawn one square forward from starting e-file)
  • Chess squares: from to  (e.g. e2 e4)

Other: quit | q  — leave the game
`

func main() {
	fmt.Print(banner)
	sc := bufio.NewScanner(os.Stdin)

	var players []*player.Player
	for {
		p1 := readPlayerDraft(sc, 1)
		p2 := readPlayerDraft(sc, 2)
		if p1.color == p2.color {
			fmt.Println("Invalid setup: colors must differ (one White, one Black). Try again from player 1.")
			fmt.Println()
			continue
		}
		if p1.color == pieces.White {
			players = []*player.Player{
				player.NewPlayer(uuid.New().String(), p1.name, pieces.White, true),
				player.NewPlayer(uuid.New().String(), p2.name, pieces.Black, true),
			}
		} else {
			players = []*player.Player{
				player.NewPlayer(uuid.New().String(), p2.name, pieces.White, true),
				player.NewPlayer(uuid.New().String(), p1.name, pieces.Black, true),
			}
		}
		break
	}

	gm := gamemanager.NewGameManager()
	g := gm.CreateGame(players)

	fmt.Println("\nGame started. White to play.")
	g.PrintState()

	for !g.IsFinished() {
		cur := g.CurrentPlayerName()
		side := strings.ToLower(string(g.CurrentTurnColor()))
		fmt.Printf("\n%s's turn (%s). Enter move, or 'quit' / 'help':\n", cur, side)

		line := readLine(sc)
		if handleMetaCommand(line) {
			continue
		}
		if shouldQuit(line) {
			fmt.Println("Goodbye.")
			return
		}

		from, to, err := parseMove(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := g.MakeMove(from, to); err != nil {
			fmt.Println("Illegal move:", err)
			continue
		}

		g.PrintState()
	}

	fmt.Printf("\nFinal result: %s\n", g.Status())
	fmt.Println("Game over.")
}

func readLine(sc *bufio.Scanner) string {
	for {
		s := readRawLine(sc)
		if s != "" {
			return s
		}
		fmt.Println("Please enter a non-empty line, or 'quit'.")
	}
}

func readRawLine(sc *bufio.Scanner) string {
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println("Goodbye.")
		os.Exit(0)
	}
	return strings.TrimSpace(sc.Text())
}

type playerDraft struct {
	name  string
	color pieces.PieceColor
}

func readPlayerDraft(sc *bufio.Scanner, n int) playerDraft {
	var d playerDraft
	for {
		fmt.Printf("Player %d name: ", n)
		name := readRawLine(sc)
		if shouldQuit(name) {
			fmt.Println("Goodbye.")
			os.Exit(0)
		}
		if name != "" {
			d.name = name
			break
		}
		fmt.Println("Name cannot be empty. Try again.")
	}
	for {
		fmt.Printf("Player %d color (White / Black): ", n)
		line := readRawLine(sc)
		if shouldQuit(line) {
			fmt.Println("Goodbye.")
			os.Exit(0)
		}
		c, ok := parsePieceColor(line)
		if ok {
			d.color = c
			return d
		}
		fmt.Println("Invalid color: type White or Black (W or B ok). Try again.")
	}
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

func shouldQuit(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "quit", "q", "exit":
		return true
	default:
		return false
	}
}

func handleMetaCommand(line string) bool {
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "help", "h", "?":
		fmt.Print(helpText)
		return true
	default:
		return false
	}
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
			return pieces.Position{}, pieces.Position{}, fmt.Errorf("move: expected four integers 0–7, or two squares like e2 e4")
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
		return pieces.Position{}, pieces.Position{}, fmt.Errorf("move: use four numbers (fromRow fromCol toRow toCol) or two squares (e2 e4); type 'help' for examples")
	}
}

func inGrid(r, c int) bool {
	return r >= 0 && r < 8 && c >= 0 && c < 8
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
