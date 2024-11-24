package tui2

import (
	"fmt"
	"strconv"
	"strings"

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

type CastsModel struct {
	cursor      int
	casts       fctools.CastGroup
	width       int
	height      int
	blocks      []castsBlock
	viewStart   int
	viewEnd     int
	hashIdx     []fctools.Hash
	focus       bool
	activeField int
	resultsNum  uint32
}

type View struct {
	Start  int
	End    int
	Cursor int
	Height int
}

func NewCastsModel() *CastsModel {
	return &CastsModel{}
}

func (m *CastsModel) SetResultsCount(count uint32) {
	m.resultsNum = count
}

func (m *CastsModel) LoadCasts(fid uint64, hash []byte) *CastsModel {
	m.prepareCasts(fctools.NewCastGroup().FromCast(nil, &farcaster.CastId{Fid: fid, Hash: hash}, true))
	return m
}

func (m *CastsModel) LoadFid(fid uint64) *CastsModel {
	m.prepareCasts(fctools.NewCastGroup().FromFid(nil, fid, m.resultsNum))
	return m
}

func (m *CastsModel) prepareCasts(casts *fctools.CastGroup) {
	m.focus = false
	m.activeField = 0
	m.casts = *casts
	m.blocks = make([]castsBlock, len(casts.Messages))
	m.hashIdx = make([]fctools.Hash, len(casts.Messages))
	m.cursor = 0
	m.appendBlocks(nil, 0)
	m.cursor = 0
	m.viewStart = 0
	m.viewEnd = 0
}

func (m *CastsModel) Init() tea.Cmd {
	return nil
}

func (m *CastsModel) initViewport() {
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

func (m *CastsModel) appendBlocks(hash *fctools.Hash, padding int) {
	opts := tui.FmtCastOpts{Width: 80}
	if hash == nil && m.casts.Head != (fctools.Hash{}) {
		hash = &m.casts.Head
	}
	if hash != nil {
		m.handleThreadBlocks(hash, padding, opts)
	} else {
		m.handleListBlocks(padding)
	}
}

func (m *CastsModel) handleThreadBlocks(hash *fctools.Hash, padding int, opts tui.FmtCastOpts) {
	idx := m.cursor
	m.hashIdx[idx] = *hash
	text := m.fmtCast(idx, padding)
	m.blocks[idx] = castsBlock{
		id:     hash.String(),
		text:   text,
		height: strings.Count(text, "\n") + 1,
	}
	m.cursor++
	for _, reply := range m.casts.Messages[*hash].Replies {
		m.handleThreadBlocks(&reply, padding+4, opts)
	}
}

func (m *CastsModel) handleListBlocks(padding int) {
	for i, hash := range m.casts.Ordered {
		m.hashIdx[i] = hash
		text := m.fmtCast(i, padding)

		m.blocks[i] = castsBlock{
			id:     hash.String(),
			text:   text,
			height: strings.Count(text, "\n") + 1,
		}
	}
}

func (m *CastsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keys := map[string]func(){
			"ctrl+c": func() { m.Quit() },
			"q":      func() { m.Quit() },
			"up":     m.moveCursorUp,
			"k":      m.moveCursorUp,
			"down":   m.moveCursorDown,
			"j":      m.moveCursorDown,
			"enter":  func() { m.focus = true },
			"right":  func() { m.focus = true },
			"l":      func() { m.focus = true },
		}

		activeKeys := map[string]func(){
			"ctrl+c": func() { m.Quit() },
			"q":      func() { m.Quit() },
			"up": func() {
				if m.activeField > 0 {
					m.activeField--
				}
			},
			"k": func() {
				if m.activeField > 0 {
					m.activeField--
				}
			},
			"down": func() { m.activeField++ },
			"j":    func() { m.activeField++ },
		}

		if m.focus {
			if fn, ok := activeKeys[msg.String()]; ok {
				fn()
			}
		} else if fn, ok := keys[msg.String()]; ok {
			fn()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *CastsModel) Quit() (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m *CastsModel) moveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.viewStart {
			m.viewStart = m.cursor
			m.recalculateViewEnd()
		}
	}
}

func (m *CastsModel) moveCursorDown() {
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

func (m *CastsModel) recalculateViewEnd() {
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

func (m *CastsModel) recalculateViewStart() {
	height := 0
	i := m.viewEnd
	for ; i >= 0 && height+m.blocks[i].height+1 <= m.height; i-- {
		height += m.blocks[i].height + 1
	}
	m.viewStart = i + 1
}

func (m *CastsModel) View() string {
	if m.focus {
		return m.ViewOne()
	}
	return m.ViewAll()
}

func (m *CastsModel) ViewOne() string {
	/*
		out := m.fmtCast(m.cursor, 2)
		height := strings.Count(out, "\n")
		return out + strings.Repeat("\n", m.height-height)
	*/
	var s strings.Builder
	if m.viewEnd == 0 {
		m.initViewport()
	}
	height := 0
	styleActive := lipgloss.NewStyle().Bold(true).BorderForeground(lipgloss.Color("#00aa00"))
	styleInactive := lipgloss.NewStyle().Faint(true)
	for i := m.viewStart; height < m.height && i < len(m.blocks); i++ {
		if height+m.blocks[i].height+1 < m.height {
			if i == m.cursor {
				s.WriteString(fmt.Sprintf("%s\n\n", styleActive.Render(m.fmtCast(m.cursor, 0))))
			} else {
				s.WriteString(fmt.Sprintf("%s\n\n", styleInactive.Render(m.blocks[i].text)))
			}

			height += m.blocks[i].height + 1
		} else {
			for _, line := range strings.Split(m.blocks[i].text, "\n") {
				s.WriteString(fmt.Sprintf("%s\n", styleInactive.Render(line)))
				height++
				if height == m.height {
					break
				}
			}
		}
	}
	s.WriteString(strings.Repeat("\n", m.height-height))
	return s.String()
}

func (m *CastsModel) ViewAll() string {
	var s strings.Builder
	if m.viewEnd == 0 {
		m.initViewport()
	}
	height := 0
	for i := m.viewStart; height < m.height && i < len(m.blocks); i++ {
		style := lipgloss.NewStyle().Bold(m.cursor == i)
		if height+m.blocks[i].height+1 < m.height {
			s.WriteString(fmt.Sprintf("%s\n\n", style.Render(m.blocks[i].text)))
			height += m.blocks[i].height + 1
		} else {
			for _, line := range strings.Split(m.blocks[i].text, "\n") {
				s.WriteString(fmt.Sprintf("%s\n", style.Render(line)))
				height++
				if height == m.height {
					break
				}
			}
		}
	}
	s.WriteString(strings.Repeat("\n", m.height-height))
	return s.String()
}

func (m *CastsModel) GetCast(cursor int) *fctools.Cast {
	return m.casts.Messages[m.hashIdx[m.cursor]]
}

func (m *CastsModel) Status() (uint64, []byte, View) {
	hash := m.hashIdx[m.cursor]
	fid := m.casts.Messages[hash].Message.Data.Fid
	return fid, hash.Bytes(), View{Start: m.viewStart, End: m.viewEnd, Cursor: m.cursor, Height: m.height}
}
func (m *CastsModel) SetView(v View) {
	m.cursor = v.Cursor
	m.viewStart = v.Start
	m.viewEnd = v.End
}

func (m *CastsModel) SetFocus(onoff bool, idx int) {
	m.focus = onoff
	m.activeField = 0
	m.cursor = idx
}
func (m *CastsModel) IsFocus() bool {
	return m.focus
}

func (m *CastsModel) GetItemInFocus() string {
	castHash := m.hashIdx[m.cursor]
	message := m.casts.Messages[castHash].Message
	itemCount := len(message.GetData().GetCastAddBody().Mentions) + len(message.GetData().GetCastAddBody().Embeds) + 1
	items := make([]string, itemCount+1)
	items[1] = "fid:" + strconv.FormatUint(message.Data.Fid, 10)

	i := 1
	for _, fid := range message.GetData().GetCastAddBody().Mentions {
		i++
		items[i] = "fid:" + strconv.FormatUint(fid, 10)
	}
	for _, embed := range message.GetData().GetCastAddBody().Embeds {
		i++
		if embedData := embed.GetCastId(); embedData != nil {
			items[i] = fmt.Sprintf("cst:%d:%x", embedData.Fid, embedData.Hash)
		} else {
			items[i] = fmt.Sprintf("url:%s", embed.GetUrl())
		}
	}
	return items[m.activeField]
}
