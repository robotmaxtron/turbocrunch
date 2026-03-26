// Package main provides a TUI calculator called TurboCrunch with dual backends: SpeedCrunch and Go.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robotmaxtron/turbocrunch/pkg/backend"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// model represents the state of the TUI application.
type model struct {
	textInput textinput.Model
	err       error
	evaluator *backend.EvaluatorWrapper
	config    *backend.Config
	history   []string
}

// initialModel initializes the model with default values.
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter expression..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	config := &backend.Config{Backend: backend.BackendSpeedCrunch}
	return model{
		textInput: ti,
		err:       nil,
		config:    config,
		evaluator: backend.NewEvaluatorWrapper(config),
		history:   []string{},
	}
}

// Init initializes the text input blinking command.
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming messages and updates the model's state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlT:
			if m.config.Backend == backend.BackendSpeedCrunch {
				m.config.Backend = backend.BackendGo
				m.history = append(m.history, "--- Switched to Go backend ---")
			} else {
				m.config.Backend = backend.BackendSpeedCrunch
				m.history = append(m.history, "--- Switched to SpeedCrunch backend ---")
			}
			return m, nil

		case tea.KeyEnter:
			expr := m.textInput.Value()
			if expr != "" {
				res := m.evaluator.Evaluate(expr)
				m.history = append(m.history, fmt.Sprintf("> %s", expr))
				m.history = append(m.history, res)
				m.textInput.SetValue("")
			}
			return m, nil
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

// View renders the application's user interface.
func (m model) View() string {
	backendStr := "SpeedCrunch (Arbitrary Precision)"
	if m.config.Backend == backend.BackendGo {
		backendStr = "Go math/cmplx (High Performance)"
	}
	s := titleStyle.Render("SpeedCrunch TUI") + "  Backend: " + backendStr + "\n\n"

	// Show history
	start := 0
	if len(m.history) > 10 {
		start = len(m.history) - 10
	}
	for _, h := range m.history[start:] {
		if strings.HasPrefix(h, ">") {
			s += h + "\n"
		} else {
			s += "  = " + h + "\n"
		}
	}

	s += "\n" + m.textInput.View() + "\n\n"
	s += historyStyle.Render("ctrl+c to quit • ctrl+t to toggle backend")

	return s
}
