package table

import (
	"fmt"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/localdb"
	"github.com/vrypan/fargo/tui"
)

type block struct {
	id     string
	text   string
	height int
}

type model struct {
	cursor    int
	hashIdx   []string
	casts     fctools.CastGroup
	width     int
	height    int
	blocks    []block
	viewStart int
	viewEnd   int
}

func (m *model) initViewport() {
	height := 0
	m.viewStart = 0
	for i, b := range m.blocks {
		if height+b.height+1 > m.height {
			m.cursor = 0
			break
		}
		m.viewEnd = i
		height += b.height + 1
	}
}

func (m *model) appendBlocks(hash *fctools.Hash, padding int) {
	if hash == nil && m.casts.Head != (fctools.Hash{}) {
		hash = &m.casts.Head
	}
	opts := tui.FmtCastOpts{Width: 80}
	if hash != nil { // This is a thread
		text := tui.FmtCast(m.casts.Messages[*hash].Message, m.casts.Fnames, padding, padding == 0, &opts)
		m.blocks[m.cursor] = block{
			id:     hash.String(),
			text:   text,
			height: strings.Count(text, "\n"),
		}
		m.cursor++
		for _, reply := range m.casts.Messages[*hash].Replies {
			m.appendBlocks(&reply, padding+4)
		}
	} else { // This is a list
		for i, hash := range m.casts.Ordered {
			msg := m.casts.Messages[hash]
			text := tui.FmtCast(msg.Message, m.casts.Fnames, padding, padding == 0, &opts)
			m.blocks[i] = block{
				id:     hash.String(),
				text:   text,
				height: strings.Count(text, "\n") + 1,
			}
		}
	}
}

func initialModel() model {
	//hashBytes, _ := hex.DecodeString("9d899db71b97f8c92538279946d74b06b529ac8c")
	casts := fctools.NewCastGroup().FromCastFidHash(nil, 3, "8cbf9d20658bc99b91e38ae77bc5c34cc53da972", true)
	// casts := fctools.NewCastGroup().FromFid(nil, 3, 100)
	m := model{
		casts:   *casts,
		hashIdx: make([]string, len(casts.Messages)),
		blocks:  make([]block, len(casts.Messages)),
	}
	m.appendBlocks(nil, 0)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			if m.cursor < m.viewStart {
				m.viewStart = m.cursor
				m.recalculateViewEnd()
			}
		case "down", "j":
			if m.viewEnd == 0 {
				m.initViewport()
			}
			if m.cursor < len(m.blocks)-1 {
				m.cursor++
			}
			if m.cursor > m.viewEnd {
				m.viewEnd = m.cursor
				m.recalculateViewStart()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 2
	}
	return m, nil
}

func (m *model) recalculateViewEnd() {
	height := 0
	for i := m.viewStart; i < len(m.blocks); i++ {
		if height+m.blocks[i].height+1 > m.height {
			break
		}
		m.viewEnd = i
		height += m.blocks[i].height + 1
	}
}

func (m *model) recalculateViewStart() {
	height := 0
	for i := m.viewEnd; i >= 0; i-- {
		blockHeight := m.blocks[i].height + 1
		if height+blockHeight > m.height {
			break
		}
		m.viewStart = i
		height += blockHeight
	}
}

func (m model) View() string {
	var s strings.Builder
	if m.viewEnd == 0 {
		m.initViewport()
	}
	height := 0
	for i := m.viewStart; height < m.height && i < len(m.blocks); i++ {
		style := lipgloss.NewStyle().Bold(m.cursor == i)
		lines := strings.Split(m.blocks[i].text, "\n")
		for _, line := range lines {
			s.WriteString(fmt.Sprintf("%s\n", style.Render(line)))
			height++
			if height == m.height {
				break
			}
		}
	}
	for ; height < m.height; height++ {
		s.WriteString("\n")
	}
	s.WriteString(fmt.Sprintf("\nPress q to quit. %d %d-%d of %d. Height=%d/%d",
		m.cursor, m.viewStart, m.viewEnd, len(m.blocks), height, m.height))
	return s.String()
}

func TestMain(t *testing.T) {
	localdb.Open()
	defer localdb.Close()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	t.Log("Done")
}
