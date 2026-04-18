package tui

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/RishiKendai/System-Design/lld/chess/game"
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncViewportSize()
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.phase {
		case phaseGameOver:
			if isQuit(msg.String()) {
				return m, tea.Quit
			}
			return m, nil

		case phasePromotion:
			return m.handlePromotionKey(msg)

		case phasePlay:
			if m.game != nil && scrollHistoryKey(msg.String()) {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
			if isQuit(msg.String()) {
				return m, tea.Quit
			}
			if isHelp(msg.String()) {
				m.pushMsg(helpText)
				return m, nil
			}
			if m.game != nil && m.game.IsFinished() {
				m.phase = phaseGameOver
				return m, nil
			}
			switch msg.Type {
			case tea.KeyEnter:
				return m.submitPlayLine()
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}

		default:
			return m.handleSetupKey(msg)
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func scrollHistoryKey(s string) bool {
	switch strings.ToLower(s) {
	case "up", "down", "pgup", "pgdown":
		return true
	default:
		return false
	}
}

func isQuit(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "quit", "q", "exit":
		return true
	default:
		return false
	}
}

func isHelp(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "help", "?":
		return true
	default:
		return false
	}
}

// afterMoveAck is shown after a legal move when play continues (not mate/resign).
func afterMoveAck(g *game.Game) string {
	s := "OK — " + g.CurrentPlayerName() + " to move."
	if g.Status() == game.GameStatusCheck {
		s += " Check."
	}
	return s
}

func (m *Model) handlePromotionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyEsc {
		m.phase = phasePlay
		m.pushMsg("Promotion cancelled — re-enter the move.")
		return m, nil
	}
	var choice *pieces.PieceType
	switch strings.ToLower(msg.String()) {
	case "q":
		t := pieces.QueenType
		choice = &t
	case "r":
		t := pieces.RookType
		choice = &t
	case "b":
		t := pieces.BishopType
		choice = &t
	case "k":
		t := pieces.KnightType
		choice = &t
	default:
		return m, nil
	}
	err := m.game.MakeMove(m.pendingFrom, m.pendingTo, choice)
	if err != nil {
		m.pushMsg(err.Error())
		m.phase = phasePlay
		return m, nil
	}
	m.phase = phasePlay
	m.textInput.SetValue("")
	if m.game.IsFinished() {
		m.phase = phaseGameOver
		m.pushMsg("Final: " + string(m.game.Status()))
	} else {
		m.pushMsg(afterMoveAck(m.game))
	}
	m.viewport.SetContent(m.moveHistoryText())
	return m, nil
}

func (m *Model) submitPlayLine() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.textInput.Value())
	m.textInput.SetValue("")
	if line == "" {
		return m, nil
	}
	if isQuit(line) {
		return m, tea.Quit
	}
	if isHelp(line) {
		m.pushMsg(helpText)
		return m, nil
	}
	from, to, err := parseMove(line)
	if err != nil {
		m.pushMsg(err.Error())
		return m, nil
	}
	err = m.game.MakeMove(from, to, nil)
	if errors.Is(err, game.ErrPromotionRequired) {
		m.pendingFrom = from
		m.pendingTo = to
		m.phase = phasePromotion
		m.pushMsg("Choose promotion (Q R B K). Esc to cancel.")
		return m, nil
	}
	if err != nil {
		m.pushMsg("Illegal move: " + err.Error())
		return m, nil
	}
	if m.game.IsFinished() {
		m.phase = phaseGameOver
		m.pushMsg("Final: " + string(m.game.Status()))
	} else {
		m.pushMsg(afterMoveAck(m.game))
	}
	m.viewport.SetContent(m.moveHistoryText())
	return m, nil
}

func (m *Model) handleSetupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if isQuit(msg.String()) {
		return m, tea.Quit
	}
	if msg.Type == tea.KeyEnter {
		line := strings.TrimSpace(m.textInput.Value())
		if line == "" {
			return m, nil
		}
		if isQuit(line) {
			return m, tea.Quit
		}
		switch m.phase {
		case phaseName1:
			m.p1name = line
			m.textInput.SetValue("")
			m.textInput.Placeholder = "White or Black"
			m.phase = phaseColor1
			return m, nil
		case phaseColor1:
			c, ok := parsePieceColor(line)
			if !ok {
				m.pushMsg("Invalid color: type White or Black (W or B).")
				return m, nil
			}
			m.p1color = c
			m.textInput.SetValue("")
			m.textInput.Placeholder = "name"
			m.phase = phaseName2
			return m, nil
		case phaseName2:
			m.p2name = line
			m.textInput.SetValue("")
			m.textInput.Placeholder = "White or Black"
			m.phase = phaseColor2
			return m, nil
		case phaseColor2:
			c, ok := parsePieceColor(line)
			if !ok {
				m.pushMsg("Invalid color: type White or Black (W or B).")
				return m, nil
			}
			m.p2color = c
			if m.p1color == m.p2color {
				m.pushMsg("Colors must differ. Start again from player 1.")
				m.phase = phaseName1
				m.textInput.SetValue("")
				m.textInput.Placeholder = "your name"
				return m, nil
			}
			m.startGame()
			return m, nil
		default:
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}
