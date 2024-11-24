package tui2

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
	text   string
	status string // Status is displayed in addition to text.
	width  int
	height int
}

func NewStatusBar() *StatusBar {
	return &StatusBar{}
}

func (s *StatusBar) SetText(t string) *StatusBar {
	s.text = t
	return s
}
func (s *StatusBar) SetStatus(status string) *StatusBar {
	s.status = status
	return s
}
func (s *StatusBar) GetStatus() string {
	return s.status
}
func (s *StatusBar) SetHeight(h int) *StatusBar {
	s.height = h
	return s
}
func (s *StatusBar) Height() int {
	return s.height + 2
}

func (s *StatusBar) Init() tea.Cmd {
	return nil
}

func (s *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
	default:

	}
	return s, nil
}

func (s *StatusBar) View() string {
	mainColor := lipgloss.AdaptiveColor{Light: "#aaa", Dark: "#777"}
	strongColor := lipgloss.AdaptiveColor{Light: "#000", Dark: "#fff"}
	var style = lipgloss.NewStyle().
		Foreground(mainColor).
		PaddingLeft(1).
		PaddingBottom(1).
		Width(s.width).Height(s.height).AlignVertical(lipgloss.Top).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(mainColor)
	statusText := ""
	if s.status != "" {
		statusText = lipgloss.NewStyle().Blink(true).Foreground(strongColor).Render(s.status)
	}
	mainText := lipgloss.NewStyle().
		Width(s.width - len(s.status) - 3).
		Foreground(mainColor).
		Render(s.text)
	return style.Render(mainText + " " + statusText)

}
