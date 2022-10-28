package tui

import (
	"fmt"
	"strings"

	"pixel/tui/constants"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"maunium.net/go/mautrix/id"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Implement different tea messages sent by the client.
	// I.e., constants.Message message data sent in a Matrix room.
	case constants.Message:
		m.displayMessages()
	case constants.Room:
		m.list.InsertItem(-1, Room(msg.Name))
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - msg.Width/4
		m.viewport.Height = msg.Height - msg.Height/4
		m.displayMessages()
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Tab):
			m.nextElement()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	switch m.Focus {
	case Input:
		return m.UpdateInput(msg)
	case Feed:
		return m.UpdateFeed(msg)
	default:
		return m.UpdateList(msg)
	}
}

func (m *Model) UpdateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var liCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.list.CursorDown()
			m.displayMessages()
			m.list.CursorUp()
		case "up":
			m.list.CursorUp()
			m.displayMessages()
			m.list.CursorDown()
		}
	}

	m.list, liCmd = m.list.Update(msg)
	return m, liCmd
}

func (m *Model) UpdateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Enter):
			if m.textarea.Focused() {
				if m.textarea.Value() != "" {
					// TODO: send text, there's other options too for later (i.e., images)
					message := m.textarea.Value()
					m.textarea.SetValue("")
					return m, m.SendTextMessage(message)
				}
			}
		}
	}

	m.textarea, tiCmd = m.textarea.Update(msg)
	return m, tiCmd
}

func (m *Model) UpdateFeed(msg tea.Msg) (tea.Model, tea.Cmd) {
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, vpCmd
}

// SetContent performs text wrapping before setting the content in the viewport
func (m *Model) SetContent(text string) {
	wrap := lipgloss.NewStyle().Width(m.viewport.Width)
	// TODO: fix this with reflow
	m.viewport.SetContent(wrap.Render(text))
}

/* HELPERS */

func (m *Model) SendTextMessage(msg string) tea.Cmd {
	return func() tea.Msg {
		room, _ := m.list.SelectedItem().(Room)
		resp, err := m.client.SendText(id.RoomID(m.rooms[string(room)]), msg)
		if err != nil {
			return errMsg(err)
		}
		return resp
	}
}

// nextElement toggles between the message entry and room list
func (m *Model) nextElement() {
	if m.Focus == Feed {
		m.Focus = List
	} else {
		m.Focus++
	}
	m.handleInput()
}

func (m *Model) handleInput() {
	if m.Focus != Input {
		m.textarea.Blur()
	} else {
		m.textarea.Focus()
	}
}

// displayMessages sets the displayed messages based on which room is selected.
func (m *Model) displayMessages() {
	if len(m.list.Items()) > 0 {

		// get the current position of the cursuor and use that to access the message map
		idx := m.list.Cursor()
		rooms := m.list.Items()
		id := rooms[idx].(Room)
		roomId := m.rooms[string(id)]

		// set content based on selected room
		m.SetContent(strings.Join(m.msgMap[roomId], "\n"))
		m.viewport.GotoBottom()
	}
}
