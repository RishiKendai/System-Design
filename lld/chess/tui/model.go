package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/RishiKendai/System-Design/lld/chess/game"
	gamemanager "github.com/RishiKendai/System-Design/lld/chess/gameManager"
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
	"github.com/RishiKendai/System-Design/lld/chess/player"
)

type phase int

const (
	phaseName1 phase = iota
	phaseColor1
	phaseName2
	phaseColor2
	phasePlay
	phasePromotion
	phaseGameOver
)

const helpText = `Moves: e2 e4  or  6 4 5 4
Rows 0–7 (0=black back rank). Cols 0–7 (a–h).
quit / q — exit   help / ? — this hint
Promotion: Q R B K (K=knight)`

// Model is the Bubble Tea root model for the chess CLI.
type Model struct {
	phase   phase
	width   int
	height  int
	msgs    []string

	textInput textinput.Model
	viewport  viewport.Model

	// setup drafts
	p1name  string
	p1color pieces.PieceColor
	p2name  string
	p2color pieces.PieceColor

	gm   *gamemanager.GameManager
	game *game.Game

	pendingFrom pieces.Position
	pendingTo   pieces.Position
}

// NewModel returns the initial TUI model (player setup then game).
func NewModel() *Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 72
	ti.Width = 28
	ti.Placeholder = "your name"

	vp := viewport.New(0, 0)
	vp.SetContent("Moves will appear here.")

	return &Model{
		phase:     phaseName1,
		textInput: ti,
		viewport:  vp,
		gm:        gamemanager.NewGameManager(),
		msgs:      nil,
	}
}

func (m *Model) pushMsg(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	m.msgs = append(m.msgs, s)
	if len(m.msgs) > 18 {
		m.msgs = m.msgs[len(m.msgs)-18:]
	}
}

func (m *Model) leftWidth() int {
	if m.width < 1 {
		return 24
	}
	w := m.width / 4
	if w < 20 {
		w = 20
	}
	if w > 30 {
		w = 30
	}
	if w >= m.width-20 {
		w = max(18, (m.width*9)/40)
	}
	return w
}

func (m *Model) rightWidth() int {
	return max(20, m.width-m.leftWidth()-3)
}

func (m *Model) rightColumnInnerHeight() int {
	h := m.height - 2
	if h < 8 {
		h = 8
	}
	return h
}

func (m *Model) rightTopHeight() int {
	h := m.rightColumnInnerHeight()
	const histMin = 4
	top := (h * 4) / 5
	if top+histMin > h-1 {
		top = h - histMin - 1
	}
	return max(8, top)
}

func (m *Model) rightBottomHeight() int {
	return max(4, m.rightColumnInnerHeight()-m.rightTopHeight()-1)
}

func (m *Model) syncViewportSize() {
	m.textInput.Width = max(8, m.leftWidth()-4)
	inner := m.rightWidth() - 2
	if inner < 8 {
		inner = 8
	}
	m.viewport.Width = inner
	m.viewport.Height = max(2, m.rightBottomHeight()-2)
}

func (m *Model) startGame() {
	players := m.orderedPlayers()
	m.game = m.gm.CreateGame(players)
	m.phase = phasePlay
	m.textInput.SetValue("")
	m.textInput.Placeholder = "e2 e4"
	m.pushMsg("Game started. White to play.")
	m.viewport.SetContent(m.moveHistoryText())
}

func (m *Model) orderedPlayers() []*player.Player {
	if m.p1color == pieces.White {
		return []*player.Player{
			player.NewPlayer(uuid.New().String(), m.p1name, pieces.White, true),
			player.NewPlayer(uuid.New().String(), m.p2name, pieces.Black, true),
		}
	}
	return []*player.Player{
		player.NewPlayer(uuid.New().String(), m.p2name, pieces.White, true),
		player.NewPlayer(uuid.New().String(), m.p1name, pieces.Black, true),
	}
}

func (m *Model) setupPrompt() string {
	switch m.phase {
	case phaseName1:
		return "Player 1 — name"
	case phaseColor1:
		return "Player 1 — color (White / Black)"
	case phaseName2:
		return "Player 2 — name"
	case phaseColor2:
		return "Player 2 — color (White / Black)"
	case phasePromotion:
		return "Promotion — press Q R B or K (knight)"
	case phaseGameOver:
		return "Game over — q to quit"
	default:
		cur := ""
		if m.game != nil {
			cur = m.game.CurrentPlayerName()
			side := strings.ToLower(string(m.game.CurrentTurnColor()))
			line := cur + "'s turn (" + side + ") — move"
			if m.game.Status() == game.GameStatusCheck {
				line += " (check)"
			}
			return line
		}
		return "Move"
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}
