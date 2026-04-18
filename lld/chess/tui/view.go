package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/RishiKendai/System-Design/lld/chess/common"
	"github.com/RishiKendai/System-Design/lld/chess/pieces"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E0E0E0"))

	leftBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6B8E9B")).
			Padding(0, 1)

	rightTopBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#8B7BB8")).
				Padding(0, 1)

	rightBottomBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7B9B7B")).
				Padding(0, 1)

	lightSquare = lipgloss.Color("#EEEED2")
	darkSquare  = lipgloss.Color("#769656")
	labelMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	whitePieceFG = lipgloss.Color("#F8F8F8")
	blackPieceFG = lipgloss.Color("#141414")
)

func squareStyle(light bool, cellW, cellH int, p pieces.Piece) lipgloss.Style {
	bg := darkSquare
	if light {
		bg = lightSquare
	}
	st := lipgloss.NewStyle().
		Background(bg).
		Width(cellW).
		Height(cellH).
		Align(lipgloss.Center, lipgloss.Center)
	if p == nil {
		return st
	}
	if p.GetColor() == pieces.White {
		return st.Foreground(whitePieceFG).Bold(true)
	}
	return st.Foreground(blackPieceFG).Bold(true)
}

func (m *Model) View() string {
	if m.width < 1 {
		return "Loading…"
	}
	leftW := m.leftWidth()
	rightW := m.rightWidth()

	leftBody := strings.Builder{}
	leftBody.WriteString(titleStyle.Render("Chess") + "\n")
	leftBody.WriteString(labelMuted.Render(m.setupPrompt()) + "\n")
	if m.phase != phasePromotion && m.phase != phaseGameOver {
		leftBody.WriteString(m.textInput.View() + "\n")
	} else if m.phase == phasePromotion {
		leftBody.WriteString(labelMuted.Render("(use Q R B K on keyboard)") + "\n")
	}
	leftBody.WriteString("\n")
	leftBody.WriteString(strings.Join(m.msgs, "\n"))

	leftPane := leftBorder.Width(leftW).Height(m.height - 2).Render(leftBody.String())

	rightTopH := m.rightTopHeight()
	boardBlock := m.renderBoardBlock(rightW - 4)
	var piecesBlock string
	if m.game != nil {
		w := "W " + m.game.ActivePiecesSummary(pieces.White)
		b := "B " + m.game.ActivePiecesSummary(pieces.Black)
		piecesBlock = labelMuted.Render(w) + "\n" + labelMuted.Render(b)
	} else {
		piecesBlock = labelMuted.Render("Pieces appear after setup")
	}
	topInner := boardBlock
	if piecesBlock != "" {
		topInner = boardBlock + "\n" + piecesBlock
	}
	topBox := rightTopBorder.Width(rightW).Height(rightTopH).Render(topInner)

	m.viewport.SetContent(m.moveHistoryText())
	histView := rightBottomBorder.Width(rightW).Height(m.rightBottomHeight()).Render(m.viewport.View())

	rightCol := lipgloss.JoinVertical(lipgloss.Left, topBox, histView)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightCol)
}

func (m *Model) renderBoardBlock(maxInnerW int) string {
	if m.game == nil {
		return labelMuted.Render(strings.Repeat("·", minInt(maxInnerW, 24)))
	}
	const cols = 8
	rankColW := 3
	usable := maxInnerW - rankColW
	if usable < cols*3 {
		usable = cols * 3
	}
	cellW := usable / cols
	if cellW < 3 {
		cellW = 3
	}
	if cellW > 6 {
		cellW = 6
	}
	cellH := 1
	if cellW >= 4 {
		cellH = 2
	}

	fileLabels := func() string {
		var parts []string
		parts = append(parts, strings.Repeat(" ", rankColW))
		for c := 0; c < cols; c++ {
			l := string(rune('a' + c))
			parts = append(parts, labelMuted.Width(cellW).Align(lipgloss.Center).Render(l))
		}
		return strings.Join(parts, "")
	}

	lines := []string{fileLabels()}
	for r := 0; r < 8; r++ {
		var cells []string
		for c := 0; c < 8; c++ {
			light := (r+c)%2 == 0
			p := m.game.PieceAt(pieces.NewPosition(r, c))
			ch := " "
			if p != nil {
				ch = common.GetPieceEmoji(p)
			}
			cells = append(cells, squareStyle(light, cellW, cellH, p).Render(ch))
		}
		// Match rank column height to squares so multi-line cells (cellH>1) do not
		// drop the rank gutter on continuation lines.
		rank := labelMuted.
			Width(rankColW).
			Height(cellH).
			Align(lipgloss.Right, lipgloss.Center).
			Render(fmt.Sprintf("%d", 8-r))
		row := lipgloss.JoinHorizontal(lipgloss.Top, append([]string{rank}, cells...)...)
		lines = append(lines, row)
	}
	lines = append(lines, fileLabels())
	return strings.Join(lines, "\n")
}

func (m *Model) moveHistoryText() string {
	if m.game == nil {
		return ""
	}
	var b strings.Builder
	for _, mv := range m.game.MovesSnapshot() {
		cap := "-"
		if mv.Captured != nil {
			cap = common.GetPieceEmoji(mv.Captured)
		}
		flags := ""
		if mv.IsMate {
			flags += "#"
		} else if mv.IsCheck {
			flags += "+"
		}
		fmt.Fprintf(&b, "%s %s→%s %s%s\n",
			common.GetPieceEmoji(mv.Piece), mv.From.String(), mv.To.String(), cap, flags)
	}
	s := b.String()
	if s == "" {
		return "No moves yet."
	}
	return s
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
