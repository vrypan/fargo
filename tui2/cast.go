package tui2

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrypan/fargo/fctools"
)

type Field int

const (
	Fid Field = iota
	InReplyTo
	EmbedCast1
	EmbedCast2
	EmbedUrl1
	EmbedUrl2
)

type CastModel struct {
	cast        *fctools.Cast
	activeField Field
	width       int
	height      int
}

func NewCastModel() *CastModel {
	m := CastModel{}
	return &m
}

func (m *CastModel) SetCast(cast *fctools.Cast) *CastModel {
	m.cast = cast
	return m
}

func (m *CastModel) Init() tea.Cmd {
	return nil
}

func (m *CastModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			return m, tea.Quit
		case "down", "j":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *CastModel) View() string {
	mainColor := lipgloss.AdaptiveColor{Light: "#333", Dark: "#aaa"}
	//strongColor := lipgloss.AdaptiveColor{Light: "#000", Dark: "#fff"}
	var style = lipgloss.NewStyle().
		Foreground(mainColor).
		PaddingLeft(1).
		Width(m.width).Height(m.height).AlignVertical(lipgloss.Top)

	return style.Render(fmt.Sprintf("%s/0x%s\n", m.cast.Fid(), m.cast.Hash()))
}
