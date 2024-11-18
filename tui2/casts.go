package tui2

import (
	"fmt"
	"strings"

	"github.com/vrypan/fargo/tui2/history"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrypan/fargo/farcaster"
	"github.com/vrypan/fargo/fctools"
	"github.com/vrypan/fargo/tui"
)

type castsBlock struct {
	id     string
	text   string
	height int
}

type castsModel struct {
	cursor    int
	casts     fctools.CastGroup
	width     int
	height    int
	blocks    []castsBlock
	viewStart int
	viewEnd   int
	hashIdx   []fctools.Hash
	history   *history.History
}

func NewCastsModel() *castsModel {
	// test vals 3, "8cbf9d20658bc99b91e38ae77bc5c34cc53da972"
	return &castsModel{
		history: history.New(1024),
	}
}

func (m *castsModel) LoadCasts(fid uint64, hash []byte) *castsModel {
	casts := fctools.NewCastGroup().FromCast(nil, &farcaster.CastId{Fid: fid, Hash: hash}, true)
	m.casts = *casts
	m.blocks = make([]castsBlock, len(casts.Messages))
	m.hashIdx = make([]fctools.Hash, len(casts.Messages))
	m.cursor = 0
	m.viewStart = 0
	m.viewEnd = 0
	m.appendBlocks(nil, 0)
	m.history.Push(history.Path{Type: history.TYPE_THREAD, Fid: fid, Hash: hash})
	return m
}
func (m *castsModel) LoadFid(fid uint64) *castsModel {
	casts := fctools.NewCastGroup().FromFid(nil, fid, 50)
	m.casts = *casts
	m.blocks = make([]castsBlock, len(casts.Messages))
	m.hashIdx = make([]fctools.Hash, len(casts.Messages))
	m.cursor = 0
	m.viewStart = 0
	m.viewEnd = 0
	m.appendBlocks(nil, 0)
	m.history.Push(history.Path{Type: history.TYPE_LIST, Fid: fid})
	return m
}

func (m *castsModel) Init() tea.Cmd {
	return nil
}

func (m *castsModel) initViewport() {
	m.viewStart = 0
	height := 0
	for i, b := range m.blocks {
		if height+b.height+1 > m.height {
			m.cursor = 0
			break
		}
		m.viewEnd = i
		height += b.height + 1
	}
}

func (m *castsModel) appendBlocks(hash *fctools.Hash, padding int) {
	opts := tui.FmtCastOpts{Width: 80}
	if hash == nil && m.casts.Head != (fctools.Hash{}) {
		hash = &m.casts.Head
	}
	if hash != nil { // This is a thread
		m.handleThreadBlocks(hash, padding, opts)
	} else {
		m.handleListBlocks(padding, opts)
	}
}

func (m *castsModel) handleThreadBlocks(hash *fctools.Hash, padding int, opts tui.FmtCastOpts) {
	text := tui.FmtCast(m.casts.Messages[*hash].Message, m.casts.Fnames, padding, padding == 0, &opts)
	m.blocks[m.cursor] = castsBlock{
		id:     hash.String(),
		text:   text,
		height: strings.Count(text, "\n"),
	}
	m.hashIdx[m.cursor] = *hash
	m.cursor++
	for _, reply := range m.casts.Messages[*hash].Replies {
		m.appendBlocks(&reply, padding+4)
	}
}

func (m *castsModel) handleListBlocks(padding int, opts tui.FmtCastOpts) {
	for i, hash := range m.casts.Ordered {
		msg := m.casts.Messages[hash]
		text := tui.FmtCast(msg.Message, m.casts.Fnames, padding, padding == 0, &opts)
		m.hashIdx[i] = hash
		m.blocks[i] = castsBlock{
			id:     hash.String(),
			text:   text,
			height: strings.Count(text, "\n") + 1,
		}
	}
}

func (m *castsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.moveCursorUp()
		case "down", "j":
			m.moveCursorDown()
		case "enter":
			hash := m.hashIdx[m.cursor]
			fid := m.casts.Messages[hash].Message.Data.Fid
			m.LoadCasts(fid, hash.Bytes())
		case "left":
			_, err := m.history.Pop() // Pop the current status
			prev, err := m.history.Peek()
			if err != nil {
				return m, nil
			}
			switch prev.Type {
			case history.TYPE_LIST:
				m.LoadFid(prev.Fid)
			case history.TYPE_THREAD:
				m.LoadCasts(prev.Fid, prev.Hash)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 2
	}
	return m, nil
}

func (m *castsModel) moveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.viewStart {
			m.viewStart = m.cursor
			m.recalculateViewEnd()
		}
	}
}

func (m *castsModel) moveCursorDown() {
	if m.viewEnd == 0 {
		m.initViewport()
	}
	if m.cursor < len(m.blocks)-1 {
		m.cursor++
		if m.cursor > m.viewEnd {
			m.viewEnd = m.cursor
			m.recalculateViewStart()
		}
	}
}

func (m *castsModel) recalculateViewEnd() {
	height := 0
	for i := m.viewStart; i < len(m.blocks); i++ {
		if height+m.blocks[i].height+1 > m.height {
			m.viewEnd = i - 1
			return
		}
		height += m.blocks[i].height + 1
	}
	m.viewEnd = len(m.blocks)
}

func (m *castsModel) recalculateViewStart() {
	height := 0
	i := m.viewEnd
	for ; i >= 0 && height+m.blocks[i].height+1 <= m.height; i-- {
		height += m.blocks[i].height + 1
	}
	m.viewStart = i + 1
}

func (m *castsModel) View() string {
	var s strings.Builder
	if m.viewEnd == 0 {
		m.initViewport()
	}
	height := 0
	for i := m.viewStart; height < m.height && i < len(m.blocks); i++ {
		style := lipgloss.NewStyle().Bold(m.cursor == i)
		if height+m.blocks[i].height+1 < m.height {
			s.WriteString(fmt.Sprintf("%s\n", style.Render(m.blocks[i].text)))
			height += m.blocks[i].height + 1
		} else {
			lines := strings.Split(m.blocks[i].text, "\n")
			for _, line := range lines {
				s.WriteString(fmt.Sprintf("%s\n", style.Render(line)))
				height++
				if height == m.height {
					break
				}
			}
		}
	}
	s.WriteString(strings.Repeat("\n", m.height-height))
	s.WriteString(fmt.Sprintf("\nPress q to quit. %d %d-%d of %d. Height=%d/%d",
		m.cursor, m.viewStart, m.viewEnd, len(m.blocks), height, m.height))
	return s.String()
}

/*
hashBytes, err := hex.DecodeString(hash[2:])
	if err != nil {
		return nil
	}
*/
