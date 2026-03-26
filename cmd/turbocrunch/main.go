// Package main provides a TUI calculator called TurboCrunch with dual backends: SpeedCrunch and Go.
package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robotmaxtron/turbocrunch/pkg/backend"
)

var (
	scientificFunctions = []string{
		"abs", "absdev", "arccos", "arcosh", "arcsin", "arctan", "arctan2", "arsinh", "artanh",
		"average", "bin", "binomcdf", "binommean", "binompmf", "binomvar", "cart", "cbrt",
		"ceil", "conj", "cos", "cosh", "cot", "csc", "datetime", "dec", "degrees", "erf",
		"erfc", "exp", "floor", "frac", "gamma", "gcd", "geomean", "gradians", "hex",
		"hypercdf", "hypermean", "hyperpmf", "hypervar", "idiv", "imag", "int", "lb",
		"lg", "ln", "lngamma", "log", "mask", "max", "median", "min", "mod", "ncr",
		"npr", "oct", "phase", "poicdf", "poimean", "poipmf", "poivar", "polar",
		"product", "radians", "real", "round", "sec", "sgn", "shl", "shr", "sin",
		"sinh", "sqrt", "stddev", "sum", "tan", "tanh", "trunc", "unmask", "variance",
	}

	functionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))  // Blue
	numberStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange
	resultStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Greenish

	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

var (
	helpStyle = help.Styles{
		ShortKey:       lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true),
		ShortDesc:      lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		FullKey:        lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true),
		FullDesc:       lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		ShortSeparator: lipgloss.NewStyle().Foreground(lipgloss.Color("238")),
		FullSeparator:  lipgloss.NewStyle().Foreground(lipgloss.Color("238")),
	}

	activeTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 1).
				Bold(true)

	inactiveTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("238")).
				Foreground(lipgloss.Color("245")).
				Padding(0, 1)

	activeBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62"))

	inactiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// viewType represents the current view of the TUI application.
type viewType int

const (
	viewInput viewType = iota
	viewFormulaBook
	viewConstants
	viewUnits
	viewHistory
)

// keyMap defines a set of keybindings. To display help and screen prompts, use the `help` bubble.
type keyMap struct {
	Up            key.Binding
	Down          key.Binding
	Enter         key.Binding
	Help          key.Binding
	Quit          key.Binding
	ToggleBackend key.Binding
	ToggleAngle   key.Binding
	FormulaBook   key.Binding
	Constants     key.Binding
	Units         key.Binding
	History       key.Binding
	Back          key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the help.KeyMap interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// help.KeyMap interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.ToggleBackend, k.ToggleAngle, k.FormulaBook, k.Constants, k.Units, k.Back},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+q"),
		key.WithHelp("ctrl+q", "quit"),
	),
	ToggleBackend: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "toggle backend"),
	),
	ToggleAngle: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "toggle angle mode"),
	),
	FormulaBook: key.NewBinding(
		key.WithKeys("f", "ctrl+f"),
		key.WithHelp("f", "formula book"),
	),
	Constants: key.NewBinding(
		key.WithKeys("c", "ctrl+c"),
		key.WithHelp("c", "constants"),
	),
	Units: key.NewBinding(
		key.WithKeys("u", "ctrl+u"),
		key.WithHelp("u", "units"),
	),
	History: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "history"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "f", "ctrl+f"),
		key.WithHelp("esc", "back"),
	),
}

// formulaItem implements list.Item interface.
type formulaItem struct {
	formula backend.Formula
}

func (i formulaItem) Title() string       { return i.formula.Name }
func (i formulaItem) Description() string { return i.formula.Description }
func (i formulaItem) FilterValue() string { return i.formula.Name }

// constantItem implements list.Item interface.
type constantItem struct {
	constant backend.Constant
}

func (i constantItem) Title() string       { return i.constant.Name }
func (i constantItem) Description() string { return fmt.Sprintf("%s [%s]", i.constant.Value, i.constant.Category) }
func (i constantItem) FilterValue() string { return i.constant.Name + " " + i.constant.Category }

// unitItem implements list.Item interface.
type unitItem struct {
	name string
}

func (i unitItem) Title() string       { return i.name }
func (i unitItem) Description() string { return "" }
func (i unitItem) FilterValue() string { return i.name }

// model represents the state of the TUI application.
type model struct {
	textInput     textinput.Model
	err           error
	evaluator     *backend.EvaluatorWrapper
	config        *backend.Config
	history       []string
	table         table.Model
	list          list.Model
	constantsList list.Model
	unitsList     list.Model
	spinner       spinner.Model
	help          help.Model
	keys          keyMap
	view          viewType
	loading       bool
}

// initialModel initializes the model with default values.
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter expression..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	ti.SetSuggestions(scientificFunctions)
	ti.ShowSuggestions = true
	ti.Prompt = "> "

	// Setup Table
	columns := []table.Column{
		{Title: "Expression", Width: 30},
		{Title: "Result", Width: 30},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Setup List
	items := make([]list.Item, len(backend.CommonFormulas))
	for i, f := range backend.CommonFormulas {
		items[i] = formulaItem{formula: f}
	}
	l := list.New(items, list.NewDefaultDelegate(), 60, 20)
	l.Title = "Formula Book"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = activeTitleStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

	config := &backend.Config{Backend: backend.BackendSpeedCrunch}
	evaluator := backend.NewEvaluatorWrapper(config)

	// Setup Constants List
	consts := evaluator.GetConstants()
	constItems := make([]list.Item, len(consts))
	for i, c := range consts {
		constItems[i] = constantItem{constant: c}
	}
	cl := list.New(constItems, list.NewDefaultDelegate(), 60, 20)
	cl.Title = "Constants"
	cl.SetShowStatusBar(false)
	cl.SetFilteringEnabled(true)
	cl.Styles.Title = activeTitleStyle
	cl.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	cl.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

	// Load session
	_ = evaluator.LoadSession("session.json")

	// Setup Units List
	units := evaluator.GetUnits()
	unitItems := make([]list.Item, len(units))
	for i, u := range units {
		unitItems[i] = unitItem{name: u}
	}
	ul := list.New(unitItems, list.NewDefaultDelegate(), 60, 20)
	ul.Title = "Units"
	ul.SetShowStatusBar(false)
	ul.SetFilteringEnabled(true)
	ul.Styles.Title = activeTitleStyle
	ul.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	ul.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

	// Setup Spinner
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	h := help.New()
	h.Styles = helpStyle

	return model{
		textInput:     ti,
		table:         t,
		list:          l,
		constantsList: cl,
		unitsList:     ul,
		spinner:       spin,
		help:          help.New(),
		keys:          keys,
		err:           nil,
		config:        config,
		evaluator:     evaluator,
		history:       []string{},
		view:          viewInput,
	}
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

// Update handles incoming messages and updates the model's state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.constantsList.SetSize(msg.Width-h, msg.Height-v)
		m.unitsList.SetSize(msg.Width-h, msg.Height-v)
		m.table.SetWidth(msg.Width - h)
		m.table.SetHeight(msg.Height - v - 15) // Adjust for header and input
		m.help.Width = msg.Width
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			if m.view == viewInput {
				_ = m.evaluator.SaveSession("session.json")
				return m, tea.Quit
			}
		}

		if m.view == viewFormulaBook {
			return m.updateFormulaBook(msg)
		}
		if m.view == viewConstants {
			return m.updateConstants(msg)
		}
		if m.view == viewUnits {
			return m.updateUnits(msg)
		}
		if m.view == viewHistory {
			return m.updateHistory(msg)
		}
		return m.updateInput(msg)

	case evaluationResultMsg:
		m.loading = false
		m.table.Focus() // Re-focus table after evaluation
		m.history = append(m.history, msg.expr)
		rows := m.table.Rows()
		rows = append(rows, table.Row{msg.expr, resultStyle.Render(msg.res)})
		m.table.SetRows(rows)
		m.table.GotoBottom()
		return m, nil

	case error:
		m.err = msg
		m.loading = false
		return m, nil
	}

	return m, nil
}

type evaluationResultMsg struct {
	expr string
	res  string
}

func (m model) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(keyMsg, m.keys.ToggleAngle):
			if m.config.Backend == backend.BackendSpeedCrunch {
				current := m.evaluator.GetAngleMode()
				var next byte
				var modeName string
				switch current {
				case 'r':
					next = 'd'
					modeName = "Degrees"
				case 'd':
					next = 'g'
					modeName = "Gradians"
				default:
					next = 'r'
					modeName = "Radians"
				}
				m.evaluator.SetAngleMode(next)
				rows := m.table.Rows()
				rows = append(rows, table.Row{"--- Angle Mode ---", fmt.Sprintf("Switched to %s", modeName)})
				m.table.SetRows(rows)
				m.table.GotoBottom()
			}
			return m, nil

		case key.Matches(keyMsg, m.keys.ToggleBackend):
			if m.config.Backend == backend.BackendSpeedCrunch {
				m.config.Backend = backend.BackendGo
			} else {
				m.config.Backend = backend.BackendSpeedCrunch
			}
			rows := m.table.Rows()
			backendName := "SpeedCrunch"
			if m.config.Backend == backend.BackendGo {
				backendName = "Go"
			}
			rows = append(rows, table.Row{"--- Backend ---", fmt.Sprintf("Switched to %s", backendName)})
			m.table.SetRows(rows)
			m.table.GotoBottom()
			return m, nil

		case key.Matches(keyMsg, m.keys.FormulaBook):
			if m.textInput.Value() == "" || keyMsg.String() == "ctrl+f" {
				m.view = viewFormulaBook
				return m, nil
			}

		case key.Matches(keyMsg, m.keys.Constants):
			if m.textInput.Value() == "" || keyMsg.String() == "ctrl+c" {
				m.view = viewConstants
				return m, nil
			}

		case key.Matches(keyMsg, m.keys.Units):
			if m.textInput.Value() == "" || keyMsg.String() == "ctrl+u" {
				m.view = viewUnits
				return m, nil
			}

		case key.Matches(keyMsg, m.keys.History):
			m.view = viewHistory
			m.table.Focus()
			return m, nil

		case key.Matches(keyMsg, m.keys.Enter):
			expr := m.textInput.Value()
			if expr != "" {
				m.loading = true
				m.textInput.SetValue("")
				m.table.Blur() // De-focus table during evaluation
				return m, func() tea.Msg {
					// Simulate some delay for spinner to be visible
					time.Sleep(100 * time.Millisecond)
					res := m.evaluator.Evaluate(expr)
					return evaluationResultMsg{expr: expr, res: res}
				}
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	// Syntax highlighting
	val := m.textInput.Value()
	if funcRegex.MatchString(val) {
		m.textInput.TextStyle = functionStyle
	} else if numberRegex.MatchString(val) {
		m.textInput.TextStyle = numberStyle
	} else {
		m.textInput.TextStyle = lipgloss.NewStyle()
	}

	return m, cmd
}

func (m model) updateFormulaBook(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) && !m.list.SettingFilter() {
			m.view = viewInput
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) && !m.list.SettingFilter() {
			i, ok := m.list.SelectedItem().(formulaItem)
			if ok {
				m.textInput.SetValue(m.textInput.Value() + i.formula.Template)
				m.view = viewInput
				m.textInput.Focus()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateConstants(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) && !m.constantsList.SettingFilter() {
			m.view = viewInput
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) && !m.constantsList.SettingFilter() {
			i, ok := m.constantsList.SelectedItem().(constantItem)
			if ok {
				m.textInput.SetValue(m.textInput.Value() + i.constant.Name)
				m.view = viewInput
				m.textInput.Focus()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.constantsList, cmd = m.constantsList.Update(msg)
	return m, cmd
}

func (m model) updateUnits(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) && !m.unitsList.SettingFilter() {
			m.view = viewInput
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) && !m.unitsList.SettingFilter() {
			i, ok := m.unitsList.SelectedItem().(unitItem)
			if ok {
				m.textInput.SetValue(m.textInput.Value() + " " + i.name)
				m.view = viewInput
				m.textInput.Focus()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.unitsList, cmd = m.unitsList.Update(msg)
	return m, cmd
}

func (m model) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Back) || key.Matches(keyMsg, m.keys.History) || keyMsg.String() == "esc" {
			m.view = viewInput
			m.table.Blur()
			m.textInput.Focus()
			return m, nil
		}
		if key.Matches(keyMsg, m.keys.Enter) {
			row := m.table.SelectedRow()
			if len(row) > 0 {
				m.textInput.SetValue(row[0])
			}
			m.view = viewInput
			m.table.Blur()
			m.textInput.Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

var (
	funcRegex   = regexp.MustCompile(`\b(` + strings.Join(scientificFunctions, "|") + `)\b`)
	numberRegex = regexp.MustCompile(`\b\d+(\.\d*)?\b`)
)

// View renders the application's user interface.
func (m model) View() string {
	if m.view == viewFormulaBook {
		return docStyle.Render(m.list.View())
	}
	if m.view == viewConstants {
		return docStyle.Render(m.constantsList.View())
	}
	if m.view == viewUnits {
		return docStyle.Render(m.unitsList.View())
	}
	if m.view == viewHistory {
		// In history view, we still show the table but it's focused
		m.table.Focus()
	} else {
		m.table.Blur()
	}

	backendStr := "SpeedCrunch"
	if m.config.Backend == backend.BackendGo {
		backendStr = "Go"
	}

	angleMode := ""
	if m.config.Backend == backend.BackendSpeedCrunch {
		mode := m.evaluator.GetAngleMode()
		switch mode {
		case 'r':
			angleMode = " (Rad)"
		case 'd':
			angleMode = " (Deg)"
		case 'g':
			angleMode = " (Grad)"
		}
	}

	header := activeTitleStyle.Render("SpeedCrunch TUI") + "  Backend: " + backendStr + angleMode

	if m.loading {
		header += " " + m.spinner.View()
	}

	inputView := m.textInput.View()
	tableView := m.table.View()

	// Apply focus highlighting
	if m.view == viewHistory {
		tableView = activeBorderStyle.Render(tableView)
		inputView = inactiveBorderStyle.Render(inputView)
	} else {
		inputView = activeBorderStyle.Render(inputView)
		tableView = inactiveBorderStyle.Render(tableView)
	}

	s := header + "\n\n"
	s += tableView + "\n"
	s += "\n" + inputView + "\n\n"
	s += m.help.View(m.keys)

	return docStyle.Render(s)
}
