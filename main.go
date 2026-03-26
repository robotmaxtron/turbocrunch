package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	textInput textinput.Model
	err       error
	evaluator *EvaluatorWrapper
	config    *Config
	history   []string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter expression..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	config := &Config{Backend: BackendSpeedCrunch}
	return model{
		textInput: ti,
		err:       nil,
		config:    config,
		evaluator: NewEvaluatorWrapper(config),
		history:   []string{},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlT:
			if m.config.Backend == BackendSpeedCrunch {
				m.config.Backend = BackendGo
				m.history = append(m.history, "--- Switched to Go backend ---")
			} else {
				m.config.Backend = BackendSpeedCrunch
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

func (m model) View() string {
	backend := "SpeedCrunch (Arbitrary Precision)"
	if m.config.Backend == BackendGo {
		backend = "Go math/cmplx (High Performance)"
	}
	s := titleStyle.Render("SpeedCrunch TUI") + "  Backend: " + backend + "\n\n"

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
