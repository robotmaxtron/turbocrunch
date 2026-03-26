// Package main provides a TUI calculator called TurboCrunch with dual backends: SpeedCrunch and Go.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/robotmaxtron/turbocrunch/pkg/backend"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// keyMap defines the keybindings for the application.
type keyMap struct {
	Up            key.Binding
	Down          key.Binding
	Enter         key.Binding
	Help          key.Binding
	Quit          key.Binding
	ToggleBackend key.Binding
	NextTheme     key.Binding
	ClearHistory  key.Binding
	CycleFormat   key.Binding
	FormulaBook   key.Binding
	Browser       key.Binding
	Variables     key.Binding
	Unset         key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.ToggleBackend, k.CycleFormat, k.FormulaBook, k.Browser, k.Variables}
}

// FullHelp returns keybindings to be shown in the full help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.ToggleBackend, k.CycleFormat, k.NextTheme, k.ClearHistory, k.FormulaBook, k.Browser, k.Variables},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "previous"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "evaluate"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "ctrl+h"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q", "quit"),
	),
	ToggleBackend: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "toggle backend"),
	),
	NextTheme: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "next theme"),
	),
	ClearHistory: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "clear history"),
	),
	CycleFormat: key.NewBinding(
		key.WithKeys("ctrl+f"),
		key.WithHelp("ctrl+f", "cycle format"),
	),
	FormulaBook: key.NewBinding(
		key.WithKeys("ctrl+b"),
		key.WithHelp("ctrl+b", "formula book"),
	),
	Browser: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "browser"),
	),
	Variables: key.NewBinding(
		key.WithKeys("ctrl+v"),
		key.WithHelp("ctrl+v", "variables"),
	),
	Unset: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "unset"),
	),
}

// model represents the state of the TUI application.
type model struct {
	textInput  textinput.Model
	help       help.Model
	keys       keyMap
	err        error
	evaluator  *backend.EvaluatorWrapper
	config     *backend.Config
	history    []string
	historyIdx int
	showHelp   bool
	themeIdx   int
	helpText   string
	showFormulas bool
	formulaIdx   int
	showBrowser  bool
	browserIdx   int
	browserMode  int // 0: Constants, 1: Units, 2: Functions
	showVars     bool
	varIdx       int
	tabCompletions []string
	tabIdx         int
	lastTabWord    string
	showFilter     bool
	filterInput    textinput.Model
}

// initialModel initializes the model with default values.
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter expression..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	ti.Prompt = "> "

	fi := textinput.New()
	fi.Placeholder = "Filter..."
	fi.CharLimit = 156
	fi.Width = 40
	fi.Prompt = "Filter: "

	config := &backend.Config{
		Backend: backend.BackendSpeedCrunch,
		Theme:   backend.DefaultTheme,
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	helpText, _ := renderer.Render(getHelpContent())

	return model{
		textInput:  ti,
		help:       help.New(),
		keys:       keys,
		err:        nil,
		config:     config,
		evaluator:  backend.NewEvaluatorWrapper(config),
		history:    []string{},
		historyIdx: -1,
		showHelp:   false,
		themeIdx:   0,
		helpText:   helpText,
		showFormulas: false,
		formulaIdx:   0,
		showBrowser:  false,
		browserIdx:   0,
		browserMode:  0,
		showVars:     false,
		varIdx:       0,
		tabCompletions: []string{},
		tabIdx:         -1,
		lastTabWord:    "",
		showFilter:     false,
		filterInput:    fi,
	}
}

func getHelpContent() string {
	formulaBook := "\n#### Formula Book\n"
	for _, f := range backend.CommonFormulas {
		formulaBook += fmt.Sprintf("- **%s**: `%s` (%s)\n", f.Name, f.Template, f.Description)
	}
	return helpContent + formulaBook
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			if m.showFilter {
				m.showFilter = false
				m.filterInput.SetValue("")
				return m, nil
			}
			if m.showFormulas || m.showBrowser || m.showVars || m.showHelp {
				m.showFormulas = false
				m.showBrowser = false
				m.showVars = false
				m.showHelp = false
				m.showFilter = false
				m.filterInput.SetValue("")
				return m, nil
			}
			return m, tea.Quit

		case msg.String() == "/":
			if (m.showFormulas || m.showBrowser || m.showVars) && !m.showFilter {
				m.showFilter = true
				m.filterInput.Focus()
				return m, nil
			}

		case msg.String() == "tab":
			if m.showBrowser {
				m.browserMode = (m.browserMode + 1) % 3
				m.browserIdx = 0
				return m, nil
			}
			if !m.showFormulas && !m.showVars && !m.showHelp {
				// Tab completion
				curr := m.textInput.Value()
				pos := m.textInput.Position()
				if pos > 0 {
					// Word before cursor
					start := pos - 1
					for start >= 0 {
						c := curr[start]
						if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
							start--
						} else {
							break
						}
					}
					start++
					word := curr[start:pos]

					if word != "" {
						if m.tabIdx == -1 || word != m.lastTabWord {
							// New completion
							m.tabCompletions = []string{}
							m.lastTabWord = word

							// Get all possible completions
							consts := m.evaluator.GetConstants()
							for _, c := range consts {
								if strings.HasPrefix(c.Name, word) {
									m.tabCompletions = append(m.tabCompletions, c.Name)
								}
							}
							units := m.evaluator.GetUnits()
							for _, u := range units {
								if strings.HasPrefix(u, word) {
									m.tabCompletions = append(m.tabCompletions, u)
								}
							}
							funcs := m.evaluator.GetFunctions()
							for _, f := range funcs {
								if strings.HasPrefix(f.Identifier, word) {
									m.tabCompletions = append(m.tabCompletions, f.Identifier)
								}
							}
							m.tabIdx = 0
						} else {
							// Cycle
							m.tabIdx = (m.tabIdx + 1) % len(m.tabCompletions)
						}

						if len(m.tabCompletions) > 0 {
							completion := m.tabCompletions[m.tabIdx]
							newVal := curr[:start] + completion + curr[pos:]
							m.textInput.SetValue(newVal)
							m.textInput.SetCursor(start + len(completion))
							m.lastTabWord = completion // Update so next tab cycles correctly
						}
					}
				}
				return m, nil
			}

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, m.keys.ToggleBackend):
			if m.config.Backend == backend.BackendSpeedCrunch {
				m.config.Backend = backend.BackendGo
				m.history = append(m.history, "--- Switched to Go backend ---")
			} else {
				m.config.Backend = backend.BackendSpeedCrunch
				m.history = append(m.history, "--- Switched to SpeedCrunch backend ---")
			}
			return m, nil

		case key.Matches(msg, m.keys.NextTheme):
			m.themeIdx = (m.themeIdx + 1) % len(backend.Themes)
			m.config.Theme = backend.Themes[m.themeIdx]
			return m, nil

		case key.Matches(msg, m.keys.ClearHistory):
			m.history = []string{}
			m.historyIdx = -1
			return m, nil

		case key.Matches(msg, m.keys.CycleFormat):
			current := m.evaluator.GetResultFormat()
			var next byte
			var formatStr string
			switch current {
			case 'd':
				next = 'h'
				formatStr = "Hexadecimal"
			case 'h':
				next = 'b'
				formatStr = "Binary"
			case 'b':
				next = 'o'
				formatStr = "Octal"
			default:
				next = 'd'
				formatStr = "Decimal"
			}
			m.evaluator.SetResultFormat(next)
			m.history = append(m.history, fmt.Sprintf("--- Switched to %s format ---", formatStr))
			return m, nil

		case key.Matches(msg, m.keys.FormulaBook):
			m.showFormulas = !m.showFormulas
			if m.showFormulas {
				m.showHelp = false
				m.showBrowser = false
				m.showVars = false
			}
			return m, nil

		case key.Matches(msg, m.keys.Browser):
			m.showBrowser = !m.showBrowser
			if m.showBrowser {
				m.showHelp = false
				m.showFormulas = false
				m.showVars = false
				m.browserIdx = 0
			}
			return m, nil

		case key.Matches(msg, m.keys.Variables):
			m.showVars = !m.showVars
			if m.showVars {
				m.showHelp = false
				m.showFormulas = false
				m.showBrowser = false
				m.varIdx = 0
			}
			return m, nil

		case key.Matches(msg, m.keys.Unset):
			if m.showVars {
				vars := m.evaluator.GetVariables()
				userFuncs := m.evaluator.GetUserFunctions()
				if m.varIdx < len(vars) {
					v := vars[m.varIdx]
					m.evaluator.UnsetVariable(v.Name)
					m.history = append(m.history, fmt.Sprintf("--- Variable '%s' unset ---", v.Name))
				} else if m.varIdx < len(vars)+len(userFuncs) {
					f := userFuncs[m.varIdx-len(vars)]
					m.evaluator.UnsetUserFunction(f.Name)
					m.history = append(m.history, fmt.Sprintf("--- Function '%s' unset ---", f.Name))
				}
				if m.varIdx > 0 && m.varIdx >= len(vars)+len(userFuncs)-1 {
					m.varIdx--
				}
				return m, nil
			}

		case key.Matches(msg, m.keys.Up):
			if m.showFormulas {
				count := 0
				filter := strings.ToLower(m.filterInput.Value())
				for _, f := range backend.CommonFormulas {
					if filter == "" || strings.Contains(strings.ToLower(f.Name), filter) || strings.Contains(strings.ToLower(f.Description), filter) || strings.Contains(strings.ToLower(f.Template), filter) {
						count++
					}
				}
				if m.formulaIdx > 0 {
					m.formulaIdx--
				}
				return m, nil
			}
			if m.showBrowser {
				if m.browserIdx > 0 {
					m.browserIdx--
				}
				return m, nil
			}
			if m.showVars {
				if m.varIdx > 0 {
					m.varIdx--
				}
				return m, nil
			}
			if len(m.history) > 0 {
				startIdx := m.historyIdx
				if startIdx == -1 {
					startIdx = len(m.history) - 1
				} else if startIdx > 0 {
					startIdx--
				} else {
					// We are at the first entry, nowhere else to go up.
					return m, nil
				}

				// Find the previous input line, starting backwards from startIdx.
				// This prevents infinite cycling and out-of-bounds access.
				found := false
				for i := startIdx; i >= 0; i-- {
					if strings.HasPrefix(m.history[i], "> ") {
						m.historyIdx = i
						m.textInput.SetValue(strings.TrimPrefix(m.history[i], "> "))
						m.textInput.CursorEnd()
						found = true
						break
					}
				}
				if !found && m.historyIdx == -1 {
					// No valid input lines found at all.
					m.historyIdx = -1
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Down):
			if m.showFormulas {
				count := 0
				filter := strings.ToLower(m.filterInput.Value())
				for _, f := range backend.CommonFormulas {
					if filter == "" || strings.Contains(strings.ToLower(f.Name), filter) || strings.Contains(strings.ToLower(f.Description), filter) || strings.Contains(strings.ToLower(f.Template), filter) {
						count++
					}
				}
				if m.formulaIdx < count-1 {
					m.formulaIdx++
				}
				return m, nil
			}
			if m.showBrowser {
				var count int
				filter := strings.ToLower(m.filterInput.Value())
				switch m.browserMode {
				case 0:
					consts := m.evaluator.GetConstants()
					for _, c := range consts {
						if filter == "" || strings.Contains(strings.ToLower(c.Name), filter) || strings.Contains(strings.ToLower(c.Category), filter) {
							count++
						}
					}
				case 1:
					units := m.evaluator.GetUnits()
					for _, u := range units {
						if filter == "" || strings.Contains(strings.ToLower(u), filter) {
							count++
						}
					}
				case 2:
					funcs := m.evaluator.GetFunctions()
					for _, f := range funcs {
						if filter == "" || strings.Contains(strings.ToLower(f.Identifier), filter) || strings.Contains(strings.ToLower(f.Name), filter) {
							count++
						}
					}
				}
				if m.browserIdx < count-1 {
					m.browserIdx++
				}
				return m, nil
			}
			if m.showVars {
				count := 0
				filter := strings.ToLower(m.filterInput.Value())
				vars := m.evaluator.GetVariables()
				for _, v := range vars {
					if filter == "" || strings.Contains(strings.ToLower(v.Name), filter) {
						count++
					}
				}
				userFuncs := m.evaluator.GetUserFunctions()
				for _, f := range userFuncs {
					if filter == "" || strings.Contains(strings.ToLower(f.Name), filter) {
						count++
					}
				}
				if m.varIdx < count-1 {
					m.varIdx++
				}
				return m, nil
			}
			if m.historyIdx != -1 {
				// Find the next input line, starting forwards from the current index.
				found := false
				for i := m.historyIdx + 1; i < len(m.history); i++ {
					if strings.HasPrefix(m.history[i], "> ") {
						m.historyIdx = i
						m.textInput.SetValue(strings.TrimPrefix(m.history[i], "> "))
						m.textInput.CursorEnd()
						found = true
						break
					}
				}

				if !found {
					// No more valid inputs, return to the empty prompt.
					m.historyIdx = -1
					m.textInput.SetValue("")
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			if m.showFormulas {
				filter := strings.ToLower(m.filterInput.Value())
				filtered := []backend.Formula{}
				for _, f := range backend.CommonFormulas {
					if filter == "" || strings.Contains(strings.ToLower(f.Name), filter) || strings.Contains(strings.ToLower(f.Description), filter) || strings.Contains(strings.ToLower(f.Template), filter) {
						filtered = append(filtered, f)
					}
				}
				if m.formulaIdx < len(filtered) {
					formula := filtered[m.formulaIdx]
					m.textInput.SetValue(formula.Template)
					m.textInput.CursorEnd()
				}
				m.showFormulas = false
				m.showFilter = false
				m.filterInput.SetValue("")
				return m, nil
			}
			if m.showBrowser {
				filter := strings.ToLower(m.filterInput.Value())
				var value string
				switch m.browserMode {
				case 0:
					consts := m.evaluator.GetConstants()
					filtered := []backend.Constant{}
					for _, c := range consts {
						if filter == "" || strings.Contains(strings.ToLower(c.Name), filter) || strings.Contains(strings.ToLower(c.Category), filter) {
							filtered = append(filtered, c)
						}
					}
					if m.browserIdx < len(filtered) {
						value = filtered[m.browserIdx].Name
					}
				case 1:
					units := m.evaluator.GetUnits()
					filtered := []string{}
					for _, u := range units {
						if filter == "" || strings.Contains(strings.ToLower(u), filter) {
							filtered = append(filtered, u)
						}
					}
					if m.browserIdx < len(filtered) {
						value = filtered[m.browserIdx]
					}
				case 2:
					funcs := m.evaluator.GetFunctions()
					filtered := []backend.Function{}
					for _, f := range funcs {
						if filter == "" || strings.Contains(strings.ToLower(f.Identifier), filter) || strings.Contains(strings.ToLower(f.Name), filter) {
							filtered = append(filtered, f)
						}
					}
					if m.browserIdx < len(filtered) {
						value = filtered[m.browserIdx].Identifier + "()"
					}
				}
				if value != "" {
					curr := m.textInput.Value()
					m.textInput.SetValue(curr + value)
					m.textInput.CursorEnd()
				}
				m.showBrowser = false
				m.showFilter = false
				m.filterInput.SetValue("")
				return m, nil
			}
			if m.showVars {
				filter := strings.ToLower(m.filterInput.Value())
				vars := m.evaluator.GetVariables()
				userFuncs := m.evaluator.GetUserFunctions()

				filteredVars := []backend.Variable{}
				for _, v := range vars {
					if filter == "" || strings.Contains(strings.ToLower(v.Name), filter) {
						filteredVars = append(filteredVars, v)
					}
				}
				filteredFuncs := []backend.UserFunction{}
				for _, f := range userFuncs {
					if filter == "" || strings.Contains(strings.ToLower(f.Name), filter) {
						filteredFuncs = append(filteredFuncs, f)
					}
				}

				var value string
				if m.varIdx < len(filteredVars) {
					value = filteredVars[m.varIdx].Name
				} else if m.varIdx < len(filteredVars)+len(filteredFuncs) {
					value = filteredFuncs[m.varIdx-len(filteredVars)].Name + "()"
				}
				if value != "" {
					curr := m.textInput.Value()
					m.textInput.SetValue(curr + value)
					m.textInput.CursorEnd()
				}
				m.showVars = false
				m.showFilter = false
				m.filterInput.SetValue("")
				return m, nil
			}
			expr := m.textInput.Value()
			if expr != "" {
				res := m.evaluator.Evaluate(expr)
				m.history = append(m.history, fmt.Sprintf("> %s", expr))
				m.history = append(m.history, res)
				m.textInput.SetValue("")
				m.historyIdx = -1
			}
			return m, nil
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	}

	if m.showFilter {
		m.filterInput, cmd = m.filterInput.Update(msg)
		// Reset selected index if filter changed to avoid out of bounds
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() != "up" && msg.String() != "down" && msg.String() != "enter" {
			m.formulaIdx = 0
			m.browserIdx = 0
			m.varIdx = 0
		}
		return m, cmd
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the application's user interface.
func (m model) View() string {
	theme := m.config.Theme

	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.TitleBG)).
		Foreground(lipgloss.Color(theme.TitleFG)).
		Padding(0, 1).
		Bold(true)

	historyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.History))

	resultStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Result))

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Operator))

	if m.showHelp {
		return m.renderHelp()
	}

	if m.showFormulas {
		return m.renderFormulas()
	}

	if m.showBrowser {
		return m.renderBrowser()
	}

	if m.showVars {
		return m.renderVariables()
	}

	backendStr := "SpeedCrunch (Arbitrary Precision)"
	if m.config.Backend == backend.BackendGo {
		backendStr = "Go (High Precision & Complex numbers)"
	}
	s := titleStyle.Render("TurboCrunch") + "  Backend: " + backendStr + "  Theme: " + theme.Name + "\n\n"

	// Show history
	start := 0
	if len(m.history) > 10 {
		start = len(m.history) - 10
	}
	for _, h := range m.history[start:] {
		if strings.HasPrefix(h, ">") {
			s += historyStyle.Render(h) + "\n"
		} else if strings.HasPrefix(h, "---") {
			s += historyStyle.Italic(true).Render(h) + "\n"
		} else {
			s += "  = " + resultStyle.Render(h) + "\n"
		}
	}

	m.textInput.PromptStyle = promptStyle
	m.textInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Cursor))

	s += "\n" + m.textInput.View() + "\n\n"
	s += m.help.View(m.keys)

	return s
}

func (m model) renderHelp() string {
	theme := m.config.Theme
	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.TitleBG)).
		Foreground(lipgloss.Color(theme.TitleFG)).
		Padding(0, 1).
		Bold(true).
		MarginBottom(1)

	return titleStyle.Render("HELP") + "\n" + m.helpText + "\n" + m.help.View(m.keys)
}

func (m model) renderFormulas() string {
	theme := m.config.Theme
	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.TitleBG)).
		Foreground(lipgloss.Color(theme.TitleFG)).
		Padding(0, 1).
		Bold(true).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Result)).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.History))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.HelpDesc)).
		Italic(true)

	s := titleStyle.Render("FORMULA BOOK") + "\n\n"
	if m.showFilter {
		s += m.filterInput.View() + "\n\n"
	} else {
		s += "Select a formula to insert into the input prompt (Press / to filter):\n\n"
	}

	filter := strings.ToLower(m.filterInput.Value())
	filtered := []backend.Formula{}
	for _, f := range backend.CommonFormulas {
		if filter == "" || strings.Contains(strings.ToLower(f.Name), filter) || strings.Contains(strings.ToLower(f.Description), filter) || strings.Contains(strings.ToLower(f.Template), filter) {
			filtered = append(filtered, f)
		}
	}

	if len(filtered) == 0 {
		s += "  No formulas match the filter.\n\n"
	}

	for i, f := range filtered {
		if i == m.formulaIdx {
			s += selectedStyle.Render(fmt.Sprintf("> %s", f.Name)) + "  " + descStyle.Render(f.Description) + "\n"
			s += "  " + selectedStyle.Render(f.Template) + "\n\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("  %s", f.Name)) + "\n\n"
		}
	}

	s += "\n" + m.help.View(m.keys)
	return s
}

func (m model) renderBrowser() string {
	theme := m.config.Theme
	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.TitleBG)).
		Foreground(lipgloss.Color(theme.TitleFG)).
		Padding(0, 1).
		Bold(true).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Result)).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.History))

	tabStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginRight(1)

	activeTabStyle := tabStyle.Copy().
		Background(lipgloss.Color(theme.ActiveTitleBG)).
		Foreground(lipgloss.Color(theme.ActiveTitleFG)).
		Bold(true)

	inactiveTabStyle := tabStyle.Copy().
		Background(lipgloss.Color(theme.InactiveTitleBG)).
		Foreground(lipgloss.Color(theme.InactiveTitleFG))

	s := titleStyle.Render("BROWSER") + "\n\n"
	if m.showFilter {
		s += m.filterInput.View() + "\n\n"
	}

	tabs := []string{"Constants", "Units", "Functions"}
	for i, t := range tabs {
		if i == m.browserMode {
			s += activeTabStyle.Render(t)
		} else {
			s += inactiveTabStyle.Render(t)
		}
	}
	s += "\n\n"

	var items []string
	filter := strings.ToLower(m.filterInput.Value())
	switch m.browserMode {
	case 0:
		consts := m.evaluator.GetConstants()
		for _, c := range consts {
			if filter == "" || strings.Contains(strings.ToLower(c.Name), filter) || strings.Contains(strings.ToLower(c.Category), filter) {
				items = append(items, fmt.Sprintf("%-15s = %-20s (%s)", c.Name, c.Value, c.Category))
			}
		}
	case 1:
		units := m.evaluator.GetUnits()
		for _, u := range units {
			if filter == "" || strings.Contains(strings.ToLower(u), filter) {
				items = append(items, u)
			}
		}
	case 2:
		funcs := m.evaluator.GetFunctions()
		for _, f := range funcs {
			if filter == "" || strings.Contains(strings.ToLower(f.Identifier), filter) || strings.Contains(strings.ToLower(f.Name), filter) {
				items = append(items, fmt.Sprintf("%-15s - %s", f.Identifier, f.Name))
			}
		}
	}

	if len(items) == 0 {
		s += "  No items match the filter.\n\n"
	}

	// Simple pagination/scrolling
	height := 15
	start := 0
	if m.browserIdx >= height {
		start = m.browserIdx - height + 1
	}

	for i := start; i < len(items) && i < start+height; i++ {
		if i == m.browserIdx {
			s += selectedStyle.Render(fmt.Sprintf("> %s", items[i])) + "\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("  %s", items[i])) + "\n"
		}
	}

	s += "\n(Tab: switch category, ↑/↓: navigate, Enter: insert)\n"
	s += "\n" + m.help.View(m.keys)
	return s
}

func (m model) renderVariables() string {
	theme := m.config.Theme
	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.TitleBG)).
		Foreground(lipgloss.Color(theme.TitleFG)).
		Padding(0, 1).
		Bold(true).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Result)).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.History))

	s := titleStyle.Render("VARIABLES & USER FUNCTIONS") + "\n\n"
	if m.showFilter {
		s += m.filterInput.View() + "\n\n"
	}

	vars := m.evaluator.GetVariables()
	userFuncs := m.evaluator.GetUserFunctions()

	filter := strings.ToLower(m.filterInput.Value())
	filteredVars := []backend.Variable{}
	for _, v := range vars {
		if filter == "" || strings.Contains(strings.ToLower(v.Name), filter) {
			filteredVars = append(filteredVars, v)
		}
	}
	filteredFuncs := []backend.UserFunction{}
	for _, f := range userFuncs {
		if filter == "" || strings.Contains(strings.ToLower(f.Name), filter) {
			filteredFuncs = append(filteredFuncs, f)
		}
	}

	if len(filteredVars) == 0 && len(filteredFuncs) == 0 {
		if filter != "" {
			s += "No variables or functions match the filter.\n"
		} else {
			s += "No user-defined variables or functions found.\n"
		}
	} else {
		idx := 0
		for _, v := range filteredVars {
			line := fmt.Sprintf("%-15s = %s", v.Name, v.Value)
			if idx == m.varIdx {
				s += selectedStyle.Render(fmt.Sprintf("> %s", line)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("  %s", line)) + "\n"
			}
			idx++
		}
		for _, f := range filteredFuncs {
			line := fmt.Sprintf("%s(%s) = %s", f.Name, strings.Join(f.Arguments, ";"), f.Expression)
			if idx == m.varIdx {
				s += selectedStyle.Render(fmt.Sprintf("> %s", line)) + "\n"
			} else {
				s += normalStyle.Render(fmt.Sprintf("  %s", line)) + "\n"
			}
			idx++
		}
	}

	s += "\n(↑/↓: navigate, Enter: insert name, Ctrl+D: unset)\n"
	s += "\n" + m.help.View(m.keys)
	return s
}

const helpContent = `
### TurboCrunch Help

TurboCrunch is a dual-backend calculator with SpeedCrunch and Go support.

#### Navigation
- **↑ / ↓**: History navigation / Browser selection
- **Enter**: Evaluate expression / Insert selection
- **Tab**: Autocomplete (Main screen) / Switch categories (Browser)
- **/**: Filter (Formula Book, Browser, Variables)
- **Ctrl+D**: Unset variable/function (Variables screen)
- **Ctrl+L**: Clear history
- **Ctrl+B**: Formula Book
- **Ctrl+U**: Browser (Constants, Units, Functions)
- **Ctrl+V**: Variables & User Functions
- **? / Ctrl+H**: Toggle this help menu
- **Q / Ctrl+C / Esc**: Quit / Close current screen

#### Configuration
- **Ctrl+T**: Toggle backend (SpeedCrunch <-> Go)
- **Ctrl+N**: Cycle through themes
- **Ctrl+F**: Cycle result format (Dec -> Hex -> Bin -> Oct)

#### Backends
- **SpeedCrunch**: Provides arbitrary precision, complex numbers, and a vast library of functions and constants.
- **Go**: High-precision calculations using the ` + "`apd`" + ` package, plus support for complex numbers with ` + "`math/cmplx`" + `.

#### Examples
- ` + "`sin(pi/4)`" + `
- ` + "`(3+4j)*(1-j)`" + `
- ` + "`light * 1 second`" + ` (SpeedCrunch only)
- ` + "`f(x) = x^2`" + ` then ` + "`f(5)`" + ` (SpeedCrunch only)
- ` + "`area(w; h) = w * h`" + ` then ` + "`area(10; 20)`" + ` (SpeedCrunch only)
`
