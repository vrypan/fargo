package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	db "github.com/vrypan/fargo/localdb"
	"github.com/vrypan/fargo/tui2"
	"github.com/vrypan/fargo/tui2/history"
)

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
	fmt.Println(*t.history)

}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

type tuiModel struct {
	casts     *tui2.CastsModel
	cursor    int
	history   *history.History
	statusBar *tui2.StatusBar
}

func NewTuiModel() *tuiModel {
	statusText := "↑/↓/←/→ navigate • q quit"
	m := tuiModel{
		casts:     tui2.NewCastsModel(),
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
		switch msg.String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "up", "k", "down", "j":
			if current.Type == history.TYPE_LIST || current.Type == history.TYPE_THREAD {
				_, cmd = t.casts.Update(msg)
				cmds = append(cmds, cmd)
			}
		case "enter", "right", "l":
			if current.Type == history.TYPE_LIST {
				cursor, fid, hash := t.casts.Status()
				cmd = func() tea.Msg { return t.statusBar.SetStatus("Loading...") }
				cmds = append(cmds, cmd)
				cmd = func() tea.Msg { return t.casts.LoadCasts(fid, hash) }
				cmds = append(cmds, cmd)
				cmd = func() tea.Msg {
					t.history.SetCursor(cursor)
					t.history.Push(history.Path{
						Type: history.TYPE_THREAD,
						Fid:  fid,
						Hash: hash,
					})
					return nil
				}
				cmds = append(cmds, cmd)
				cmd = func() tea.Msg { return t.statusBar.SetStatus("") }
				cmds = append(cmds, cmd)

				/*
					cmd = func() tea.Msg {
						t.statusBar.SetStatus("Loading...")
						t.casts.LoadCasts(fid, hash)
						t.history.Push(history.Path{
							Type: history.TYPE_THREAD,
							Fid:  fid,
							Hash: hash,
						})
						t.statusBar.SetStatus("")
						t.casts.Update(msg)
						return nil
					}

					cmds = append(cmds, cmd)
				*/
				return t, tea.Sequence(cmds...)
			}
		case "left", "h":
			if current.Type == history.TYPE_THREAD {
				_, err := t.history.Pop() // Pop the current status
				prev, err := t.history.Peek()
				if err != nil {
					return t, nil
				}
				t.casts.LoadFid(prev.Fid)
				t.cursor = prev.Cursor
			}
		}
	case tea.WindowSizeMsg:
		if current.Type == history.TYPE_LIST || current.Type == history.TYPE_THREAD {
			_, cmd := t.casts.Update(tea.WindowSizeMsg{Width: msg.Width, Height: msg.Height - t.statusBar.Height()})
			cmds = append(cmds, cmd)
		}
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
	case history.TYPE_THREAD:
		output.WriteString(t.casts.View())
	}
	output.WriteString(t.statusBar.View())
	return output.String()
}
