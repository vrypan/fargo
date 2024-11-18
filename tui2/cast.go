package tui2

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

type castModel struct {
	cast        *fctools.Cast
	activeField Field
}

func NewCastModel() *castModel {
	m := castModel{}
	return &m
}

func (m *castModel) SetCast(cast *fctools.Cast) *castModel {
	m.cast = cast
	return m
}

func (m *castModel) Init() tea.Cmd {
	return nil
}

func (m *castModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}
	return m, nil
}

func (m *castModel) View() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%s/0x%s", m.cast.Fid(), m.cast.Hash()))
	return s.String()
}
