package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	db "github.com/vrypan/fargo/localdb"
	"github.com/vrypan/fargo/tui2"
	"github.com/vrypan/fargo/tui2/history"
)

var DEBUG string

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "A brief description of your command",
	Run:   tuiRun,
}

func tuiRun(cmd *cobra.Command, args []string) {
	db.Open()
	defer db.Close()
	t := NewTuiModel()

	t.casts.LoadFid(280)
	t.history.Push(history.Path{
		Type: history.TYPE_LIST,
		Fid:  280,
	})
	p := tea.NewProgram(t)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	fmt.Println(t.history)

}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

type tuiModel struct {
	casts     *tui2.CastsModel
	cast      *tui2.CastModel
	cursor    int
	history   *history.History
	statusBar *tui2.StatusBar
}

func NewTuiModel() *tuiModel {
	statusText := "↑/↓/←/→ navigate • q quit"
	m := tuiModel{
		casts:     tui2.NewCastsModel(),
		cast:      tui2.NewCastModel(),
		history:   history.New(1024),
		statusBar: tui2.NewStatusBar().SetText(statusText).SetHeight(1),
	}
	return &m
}

func (t tuiModel) Init() tea.Cmd {
	return nil
}

func (t tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	current, err := t.history.Peek()
	if err != nil {
		log.Fatal(err)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch current.Type {
		case history.TYPE_LIST, history.TYPE_THREAD:
			switch msg.String() {
			case "ctrl+c", "q":
				return t, tea.Quit
			case "up", "k", "down", "j":
				_, cmd = t.casts.Update(msg)
				cmds = append(cmds, cmd)
			case "enter", "right", "l":
				cursor, fid, hash := t.casts.Status()
				if current.Type == history.TYPE_LIST {
					cmd = func() tea.Msg {
						t.statusBar.SetStatus("Loading...")
						t.casts.LoadCasts(fid, hash)
						t.history.SetCursor(cursor)
						t.history.Push(history.Path{
							Type: history.TYPE_THREAD,
							Fid:  fid,
							Hash: hash,
						})
						return t.statusBar.SetStatus("")
					}
					cmds = append(cmds, cmd)
				}
				if current.Type == history.TYPE_THREAD {
					t.history.SetCursor(cursor)
					t.history.Push(history.Path{
						Type: history.TYPE_CAST,
						Fid:  fid,
						Hash: hash,
					})
					t.casts.SetFocus(true)
				}
			case "left", "h":
				_, err := t.history.Pop() // Pop the current status
				prev, err := t.history.Peek()
				if err != nil {
					return t, nil
				}
				switch prev.Type {
				case history.TYPE_LIST:
					t.casts.LoadFid(prev.Fid)
				case history.TYPE_THREAD:
					t.casts.LoadCasts(prev.Fid, prev.Hash)
				}
				t.cursor = prev.Cursor
			}
		case history.TYPE_CAST:
			switch msg.String() {
			case "esc", "q":
				return t, tea.Quit
			case "enter":
				action := t.casts.GetItemInFocus()
				if action == "" {
					return t, nil
				}
				switch action[0:3] {
				case "fid":
					nextFid, _ := strconv.ParseUint(action[4:], 10, 64)
					cursor, _, _ := t.casts.Status()
					cmd = func() tea.Msg {
						t.statusBar.SetStatus("Loading...")
						t.history.SetCursor(cursor)
						t.casts.LoadFid(nextFid)
						t.history.Push(history.Path{
							Type: history.TYPE_LIST,
							Fid:  nextFid,
						})
						return t.statusBar.SetStatus("")
					}
					cmds = append(cmds, cmd)
				case "cst":
					parts := strings.Split(action, ":")
					nextFid, _ := strconv.ParseUint(parts[1], 10, 64)
					nextHash, _ := hex.DecodeString(parts[2])
					cursor, fid, hash := t.casts.Status()
					cmd = func() tea.Msg {
						t.statusBar.SetStatus("Loading...")
						t.casts.LoadCasts(nextFid, nextHash)
						t.history.SetCursor(cursor)
						t.history.Push(history.Path{
							Type: history.TYPE_THREAD,
							Fid:  fid,
							Hash: hash,
						})
						return t.statusBar.SetStatus("")
					}
					cmds = append(cmds, cmd)
				case "url":
				}
			case "left":
				t.history.Pop()
				prev, err := t.history.Peek()
				if err != nil {
					return t, nil
				}
				if prev.Hash == nil {
					t.casts.LoadFid(prev.Fid)
				} else {
					t.casts.LoadCasts(prev.Fid, prev.Hash)
				}
				t.casts.SetFocus(false)
				t.cursor = prev.Cursor
				// _, cmd = t.casts.Update(msg)
				//cmds = append(cmds, cmd)
			default:
				_, cmd = t.casts.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		_, cmd := t.casts.Update(tea.WindowSizeMsg{Width: msg.Width, Height: msg.Height - t.statusBar.Height()})
		cmds = append(cmds, cmd)
		_, cmd = t.cast.Update(tea.WindowSizeMsg{Width: msg.Width, Height: msg.Height - t.statusBar.Height()})
		cmds = append(cmds, cmd)
		_, cmd = t.statusBar.Update(msg)
		cmds = append(cmds, cmd)
	}
	return t, tea.Batch(cmds...)
}

func (t tuiModel) View() string {
	current, _ := t.history.Peek()
	var output strings.Builder
	switch current.Type {
	case history.TYPE_LIST:
		output.WriteString(t.casts.View())
	case history.TYPE_THREAD, history.TYPE_CAST:
		output.WriteString(t.casts.View())
	}
	output.WriteString(t.statusBar.View())
	return output.String()
}
