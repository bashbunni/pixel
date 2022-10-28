package tui

import (
	"fmt"
	"io"

	"pixel/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Room is the name of a Matrix room
type Room string
type ItemDelegate struct{}

func (r Room) FilterValue() string { return "" }

func (d ItemDelegate) Height() int                               { return 1 }
func (d ItemDelegate) Spacing() int                              { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Room)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return styles.SelectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

func CreateList() list.Model {
	list := list.New([]list.Item{}, ItemDelegate{}, styles.DefaultWidth, styles.ListHeight)
	list.SetFilteringEnabled(false)
	list.DisableQuitKeybindings()
	list.KeyMap.CursorUp.SetKeys("up")
	list.KeyMap.CursorDown.SetKeys("down")
	// disable these keys ("g" and "G") while the list is inactive - it interferes with typing otherwise
	list.KeyMap.GoToStart.SetEnabled(false)
	list.KeyMap.GoToEnd.SetEnabled(false)
	list.KeyMap.Filter.SetEnabled(false)

	list.Title = "Rooms"
	list.SetStatusBarItemName("Room", "Rooms")
	return list
}
