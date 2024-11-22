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
	"github.com/vrypan/fargo/tui2/history"
)

var exploreCmd = &cobra.Command{
	Use:   "explore [@username]/casts",
	Short: "Interactive Farcaster explorer",
	Long: `It only supports "@username/casts" for now.
Ex.: fargo explore @dwr/casts
`,
	Run: exploreRun,
}

func exploreRun(cmd *cobra.Command, args []string) {
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
	t := NewTuiModel()
	t.casts.SetResultsCount(countFlag)

	switch {
	case len(parts) == 1 && parts[0] == "casts":
		t.casts.LoadFid(user.Fid)
		t.history.Push(history.Path{
			Type: history.TYPE_LIST,
			Fid:  280,
		})
		p := tea.NewProgram(t)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.AddCommand(exploreCmd)
	exploreCmd.Flags().IntP("count", "c", 0, "Number of casts to show when getting @user/casts")
}

type tuiModel struct {
	casts     *tui2.CastsModel
	cursor    int
	history   *history.History
	statusBar *tui2.StatusBar
}

func NewTuiModel() *tuiModel {
	statusText := "↑/↓/←/→ navigate • q quit"
	m := &tuiModel{
		casts:     tui2.NewCastsModel(),
		history:   history.New(1024),
		statusBar: tui2.NewStatusBar().SetText(statusText).SetHeight(1),
	}
	return m
}

func (t *tuiModel) Init() tea.Cmd {
	return nil
}

func (t *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	current, err := t.history.Peek()
	if err != nil {
		log.Fatal(err)
	}
	handleKeyMsg := func(msg tea.KeyMsg) {
		switch current.Type {
		case history.TYPE_LIST, history.TYPE_THREAD:
			switch msg.String() {
			case "ctrl+c", "q":
				cmds = append(cmds, tea.Quit)
			case "up", "k", "down", "j":
				_, cmd := t.casts.Update(msg)
				cmds = append(cmds, cmd)
			case "enter", "right", "l":
				cursor, fid, hash := t.casts.Status()
				if current.Type == history.TYPE_LIST {
					cmd := func() tea.Msg {
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
				} else if current.Type == history.TYPE_THREAD {
					t.history.SetCursor(cursor)
					t.history.Push(history.Path{
						Type: history.TYPE_CAST,
						Fid:  fid,
						Hash: hash,
					})
					t.casts.SetFocus(true, cursor)
				}
			case "left", "h":
				if _, err := t.history.Pop(); err == nil {
					if prev, err := t.history.Peek(); err == nil {
						switch prev.Type {
						case history.TYPE_LIST:
							t.casts.LoadFid(prev.Fid)
						case history.TYPE_THREAD:
							t.casts.LoadCasts(prev.Fid, prev.Hash)
						case history.TYPE_CAST:
							t.casts.LoadCasts(prev.Fid, prev.Hash)
							t.casts.SetFocus(true, prev.Cursor)
						}
						t.cursor = prev.Cursor
					}
				}
			}
		case history.TYPE_CAST:
			switch msg.String() {
			case "esc", "q":
				cmds = append(cmds, tea.Quit)
			case "enter":
				action := t.casts.GetItemInFocus()
				if action == "" {
					return
				}
				switch action[0:3] {
				case "fid":
					nextFid, _ := strconv.ParseUint(action[4:], 10, 64)
					cursor, _, _ := t.casts.Status()
					cmd := func() tea.Msg {
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
					cmd := func() tea.Msg {
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
				}
			case "left":
				if _, err := t.history.Pop(); err == nil {
					if prev, err := t.history.Peek(); err == nil {
						if prev.Hash == nil {
							cmd := func() tea.Msg {
								t.statusBar.SetStatus("Loading...")
								t.casts.LoadFid(prev.Fid)
								return t.statusBar.SetStatus("")
							}
							cmds = append(cmds, cmd)
						} else {
							cmd := func() tea.Msg {
								t.statusBar.SetStatus("Loading...")
								t.casts.LoadCasts(prev.Fid, prev.Hash)
								return t.statusBar.SetStatus("")
							}
							cmds = append(cmds, cmd)
						}
						t.casts.SetFocus(false, 0)
						t.cursor = prev.Cursor
					}
				}
			default:
				_, cmd := t.casts.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		adjustSize := func(model tea.Model, h int) {
			_, cmd := (model).Update(tea.WindowSizeMsg{Width: msg.Width, Height: msg.Height - h})
			cmds = append(cmds, cmd)
		}
		h := t.statusBar.Height()
		adjustSize(t.casts, h)
		adjustSize(t.statusBar, 0)
	}
	return t, tea.Batch(cmds...)
}

func (t *tuiModel) View() string {
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
