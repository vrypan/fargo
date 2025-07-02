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
	"github.com/vrypan/fargo/config"
	db "github.com/vrypan/fargo/localdb"
	"github.com/vrypan/fargo/tui2"
	history "github.com/vrypan/fargo/tui2/history2"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive [@username]/casts",
	Short: "Interactive Farcaster explorer",
	Long: `**Experimental and buggy feature**
It only supports "@username/casts" for now.
Ex.: fargo explore @dwr/casts
`,
	Run: interactiveRun,
}

func interactiveRun(cmd *cobra.Command, args []string) {
	countFlag := uint32(config.GetInt("get.count"))

	db.Open()
	defer db.Close()
	user, parts := parse_url(args)
	if user == nil {
		log.Fatal("User not found")
	}
	if c, _ := cmd.Flags().GetInt("count"); c > 0 {
		countFlag = uint32(c)
	}
	t := NewTuiModel2()
	t.casts.SetResultsCount(countFlag)

	switch {
	case len(parts) == 1 && parts[0] == "casts":
		t.casts.LoadFid(user.Fid)
		status := t.casts.GetStatus()
		t.history.Push(status)
		p := tea.NewProgram(t)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
	interactiveCmd.Flags().IntP("count", "c", 0, "Number of casts to show when getting @user/casts")
}

type tuiModel2 struct {
	casts   *tui2.CastsModel
	cursor  int
	history *history.History
}

func NewTuiModel2() *tuiModel2 {
	m := &tuiModel2{
		casts:   tui2.NewCastsModel(),
		history: history.New(1024),
	}
	return m
}

func (t *tuiModel2) Init() tea.Cmd {
	return nil
}

func cmdSequence(m ...tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	for _, msg := range m {
		cmds = append(cmds, func() tea.Msg { return msg })
	}
	return cmds
}
func (t *tuiModel2) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	status := t.casts.GetStatus()
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "enter", "right":
			if !status.Focus {
				t.history.Update(status)
				switch status.View {
				case tui2.VIEW_LIST:
					cast := t.casts.GetCast(status.Cursor)
					fid := cast.Message.Data.Fid
					hash := cast.Message.Hash
					t.history.Push(tui2.CastsStatus{Fid: fid, Hash: hash, View: tui2.VIEW_THREAD})
					cmds = cmdSequence(
						tui2.UpdateStatusBar{Text: "Loading..."},
						tui2.LoadCastId{Fid: fid, Hash: hash},
						tui2.UpdateStatusBar{Text: ""},
					)
					return t, tea.Sequence(cmds...)
				case tui2.VIEW_THREAD:
					t.casts.SetFocus(true, status.Cursor)
				}
			} else {
				action := t.casts.GetItemInFocus()
				if len(action) < 3 {
					return t, tea.Sequence(cmds...)
				}
				actionType := action[0:3]
				switch actionType {
				case "fid":
					parts := strings.Split(action, ":")
					nextFid, _ := strconv.ParseUint(parts[1], 10, 64)
					t.history.Update(status)
					t.history.Push(tui2.CastsStatus{Fid: nextFid, View: tui2.VIEW_LIST})
					cmds = cmdSequence(
						tui2.UpdateStatusBar{Text: "Loading..."},
						tui2.LoadFid{Fid: nextFid},
						tui2.UpdateStatusBar{Text: ""},
					)
					return t, tea.Sequence(cmds...)
				case "cst":
					parts := strings.Split(action, ":")
					nextFid, _ := strconv.ParseUint(parts[1], 10, 64)
					nextHash, _ := hex.DecodeString(parts[2])
					t.history.Update(status)
					t.history.Push(tui2.CastsStatus{Fid: nextFid, Hash: nextHash, View: tui2.VIEW_THREAD})
					cmds = cmdSequence(
						tui2.UpdateStatusBar{Text: "Loading..."},
						tui2.LoadCastId{Fid: nextFid, Hash: nextHash},
						tui2.UpdateStatusBar{Text: ""},
					)
					return t, tea.Sequence(cmds...)
				case "url":
					url := action[4:]
					OpenUrl(url)
					return t, nil
				}
			}
		case "left":
			if status.Focus {
				t.casts.SetFocus(false, status.Cursor)
				return t, nil
			}
			t.history.Pop()
			last, err := t.history.Peek()
			if err != nil {
				return t, nil
			}
			switch last.View {
			case tui2.VIEW_LIST:
				cmds = cmdSequence(
					tui2.UpdateStatusBar{Text: "Loading..."},
					tui2.LoadFid{Fid: last.Fid},
					tui2.MsgUpdateView{Start: last.ViewStart, End: last.ViewEnd, Cursor: last.Cursor, Height: last.Height},
					tui2.UpdateStatusBar{Text: ""},
				)
				return t, tea.Sequence(cmds...)
			case tui2.VIEW_THREAD:
				cmds = cmdSequence(
					tui2.UpdateStatusBar{Text: "Loading..."},
					tui2.LoadCastId{Fid: last.Fid, Hash: last.Hash},
					tui2.MsgUpdateView{Start: last.ViewStart, End: last.ViewEnd, Cursor: last.Cursor, Height: last.Height},
					tui2.UpdateStatusBar{Text: ""},
				)
				return t, tea.Sequence(cmds...)
			}
		default:
			t.casts.Update(msg)
		}
	case tea.WindowSizeMsg:
		t.casts.Update(msg)
	default:
		t.casts.Update(msg)
	}
	return t, nil
}

func (t *tuiModel2) View() string {
	return t.casts.View()
}
