package tui2

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vrypan/fargo/farcaster"
	"github.com/vrypan/fargo/fctools"
)

// Message types
type UpdateStatusBar struct {
	Text string
}
type LoadFid struct {
	Fid uint64
}
type LoadCastId struct {
	Fid  uint64
	Hash []byte
}
type MsgUpdateView = View

type ViewType int

const (
	VIEW_PROFILE ViewType = iota
	VIEW_LIST
	VIEW_THREAD
	VIEW_CAST
)

type CastsModel struct {
	casts   fctools.CastGroup
	blocks  []castsBlock
	hashIdx []fctools.Hash

	width  int
	height int

	view        ViewType
	viewStart   int
	viewEnd     int
	cursor      int
	focus       bool
	activeField int

	resultsNum uint32
	statusBar  *StatusBar
}

type View struct {
	Start  int
	End    int
	Cursor int
	Height int
}

func NewCastsModel() *CastsModel {
	m := CastsModel{}
	statusText := "↑/↓/←/→ navigate • q quit"
	m.statusBar = NewStatusBar().SetText(statusText).SetHeight(1)
	return &m
}

func (m *CastsModel) SetResultsCount(count uint32) {
	m.resultsNum = count
}

func (m *CastsModel) LoadCasts(fid uint64, hash []byte) *CastsModel {
	m.view = VIEW_THREAD
	m.prepareModel(fctools.NewCastGroup().FromCast(nil, &farcaster.CastId{Fid: fid, Hash: hash}, true))
	return m
}

func (m *CastsModel) LoadFid(fid uint64) *CastsModel {
	m.view = VIEW_LIST
	m.prepareModel(fctools.NewCastGroup().FromFid(nil, fid, m.resultsNum))
	return m
}

func (m *CastsModel) prepareModel(casts *fctools.CastGroup) {
	m.focus = false
	m.activeField = 0
	m.casts = *casts
	m.blocks = make([]castsBlock, len(casts.Messages))
	m.hashIdx = make([]fctools.Hash, len(casts.Messages))
	m.cursor = 0
	m.renderBlocks(nil, 0)
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
		if height+b.height+1 > m.WindowHeight() {
			m.cursor = 0
			break
		}
		m.viewEnd = i
		height += b.height + 1
	}
}

func (m *CastsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keys := map[string]func(){
			"ctrl+c": func() { m.Quit() },
			"q":      func() { m.Quit() },
			"up":     m.moveCursorUp,
			"down":   m.moveCursorDown,
			"enter":  func() { m.focus = true },
			"right":  func() { m.focus = true },
		}

		activeKeys := map[string]func(){
			"ctrl+c": func() { m.Quit() },
			"q":      func() { m.Quit() },
			"up": func() {
				if m.activeField > 0 {
					m.activeField--
				}
			},
			"down": func() { m.activeField++ },
			"left": func() { m.focus = false },
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
		m.statusBar.Update(msg)
	case UpdateStatusBar:
		m.statusBar.SetStatus(msg.Text)
	case LoadFid:
		m.LoadFid(msg.Fid)
	case LoadCastId:
		m.LoadCasts(msg.Fid, msg.Hash)
	case MsgUpdateView:
		m.SetView(msg)
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
		if height+m.blocks[i].height+1 > m.WindowHeight() {
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
	for ; i >= 0 && height+m.blocks[i].height+1 <= m.WindowHeight(); i-- {
		height += m.blocks[i].height + 1
	}
	m.viewStart = i + 1
}

func (m *CastsModel) WindowHeight() int {
	return max(0, m.height-m.statusBar.Height())
}

func (m *CastsModel) View() string {
	ret := m.ViewCasts()
	ret += m.statusBar.View()
	return ret
}

func (m *CastsModel) ViewCasts() string {
	var s strings.Builder
	if m.viewEnd == 0 {
		m.initViewport()
	}
	height := 0
	activeStyle := lipgloss.NewStyle().Bold(true).BorderForeground(lipgloss.Color("#00aa00")).Border(lipgloss.NormalBorder(), false, false, false, true)
	inactiveStyle := lipgloss.NewStyle().Faint(true)
	for i := m.viewStart; height < m.WindowHeight() && i < len(m.blocks); i++ {
		var style lipgloss.Style
		if height+m.blocks[i].height+1 <= m.WindowHeight() {
			if m.focus {
				style = inactiveStyle.PaddingLeft(m.blocks[i].padding + 1)
				if i == m.cursor {
					style = activeStyle.PaddingLeft(m.blocks[i].padding + 1)
					s.WriteString(fmt.Sprintf("%s\n\n", style.Render(m.fmtCast(m.cursor, 0))))
				} else {
					style = style.Border(lipgloss.HiddenBorder(), false, false, false, true)
					s.WriteString(fmt.Sprintf("%s\n\n", style.Render(m.blocks[i].text)))
				}
			} else {
				if m.cursor == i {
					style = lipgloss.NewStyle().Bold(true).Border(lipgloss.NormalBorder(), false, false, false, true).PaddingLeft(m.blocks[i].padding + 1)
					s.WriteString(fmt.Sprintf("%s\n\n", style.Render(m.blocks[i].text)))
				} else {
					style = lipgloss.NewStyle().Bold(false).Border(lipgloss.HiddenBorder(), false, false, false, true).PaddingLeft(m.blocks[i].padding + 1)
					s.WriteString(fmt.Sprintf("%s\n\n", style.Render(m.blocks[i].text)))
				}

			}
			height += m.blocks[i].height + 1
		} else {
			for _, line := range strings.Split(m.blocks[i].text, "\n") {
				s.WriteString(fmt.Sprintf("%s\n", style.PaddingLeft(m.blocks[i].padding+1).Render(line)))
				height++
				if height == m.WindowHeight() {
					break
				}
			}
		}
	}
	s.WriteString(strings.Repeat("\n", m.WindowHeight()-height))
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

// Gets the model's status
// Can be used by a parent model to decide on actions.
//

type CastsStatus struct {
	Fid         uint64
	Hash        []byte
	Width       int
	Height      int
	View        ViewType
	ViewStart   int
	ViewEnd     int
	Cursor      int
	Focus       bool
	ActiveField int
}

func (s CastsStatus) String() string {
	return fmt.Sprintf("@%d/0x%x [%d-%d-%d] Focus=%t", s.Fid, s.Hash, s.ViewStart, s.Cursor, s.ViewEnd, s.Focus)
}

func (m *CastsModel) GetStatus() CastsStatus {
	hash := m.hashIdx[m.cursor]
	fid := m.casts.Messages[hash].Message.Data.Fid
	return CastsStatus{
		Fid:         fid,
		Hash:        hash.Bytes(),
		Width:       m.width,
		Height:      m.height,
		View:        m.view,
		ViewStart:   m.viewStart,
		ViewEnd:     m.viewEnd,
		Cursor:      m.cursor,
		Focus:       m.focus,
		ActiveField: m.activeField,
	}
}
