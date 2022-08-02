package main

// A simple login manager.

import (
	"fmt"
	"os"
	"strings"

	term "golang.org/x/term"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/msteinert/pam"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#005577"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	wrongStyle			= lipgloss.NewStyle().Foreground(lipgloss.Color("#CC3333"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Login ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))

	passWrong = false
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "username"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+shift+c":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {  // if focus is on Login
				result := checkPassword(m.inputs[0].Value(), m.inputs[1].Value())
				if result {
					return m, tea.Quit
				} else {
					passWrong = true
					m.inputs[1].SetValue("")
				}
			} else if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	if passWrong && m.inputs[1].View() == "" {
		m.inputs[1].TextStyle = wrongStyle
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	width, height, _ := term.GetSize(int(os.Stdin.Fd()))

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, b.String())
}

func checkPassword(username string, password string) bool {
	trans, _ := pam.StartFunc("login", username, func(style pam.Style, msg string) (string, error) {
		return password, nil
	})

	err := trans.Authenticate(0)
	return err == nil
}

func main() {
	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
