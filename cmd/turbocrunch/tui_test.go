package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTUI_InitialState(t *testing.T) {
	m := initialModel()
	view := m.View()

	if !strings.Contains(view, "SpeedCrunch TUI") {
		t.Errorf("Expected view to contain 'SpeedCrunch TUI', got %q", view)
	}

	if !strings.Contains(view, "Backend: SpeedCrunch (Arbitrary Precision)") {
		t.Errorf("Expected initial backend to be SpeedCrunch, got %q", view)
	}

	if !strings.Contains(view, "Enter expression...") {
		t.Errorf("Expected view to contain placeholder, got %q", view)
	}
}

func TestTUI_ToggleBackend(t *testing.T) {
	m := initialModel()

	// Toggle to Go backend
	msg := tea.KeyMsg{Type: tea.KeyCtrlT}
	newModel, _ := m.Update(msg)
	m = newModel.(model)

	view := m.View()
	if !strings.Contains(view, "Backend: Go math/cmplx (High Performance)") {
		t.Errorf("Expected backend to be Go, got %q", view)
	}
	if !strings.Contains(view, "--- Switched to Go backend ---") {
		t.Errorf("Expected history to contain switch message, got %q", view)
	}

	// Toggle back to SpeedCrunch
	newModel, _ = m.Update(msg)
	m = newModel.(model)
	view = m.View()
	if !strings.Contains(view, "Backend: SpeedCrunch (Arbitrary Precision)") {
		t.Errorf("Expected backend to be SpeedCrunch, got %q", view)
	}
	if !strings.Contains(view, "--- Switched to SpeedCrunch backend ---") {
		t.Errorf("Expected history to contain switch back message, got %q", view)
	}
}

func TestTUI_EvaluateExpression(t *testing.T) {
	m := initialModel()

	// Simulate typing "2+2"
	keys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune("2")},
		{Type: tea.KeyRunes, Runes: []rune("+")},
		{Type: tea.KeyRunes, Runes: []rune("2")},
	}

	for _, k := range keys {
		newModel, _ := m.Update(k)
		m = newModel.(model)
	}

	// Verify input value
	if m.textInput.Value() != "2+2" {
		t.Errorf("Expected input value to be '2+2', got %q", m.textInput.Value())
	}

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	m = newModel.(model)

	// Verify history update in View
	view := m.View()
	if !strings.Contains(view, "> 2+2") {
		t.Errorf("Expected history to show input, got %q", view)
	}
	if !strings.Contains(view, "= 4") {
		t.Errorf("Expected history to show result '4', got %q", view)
	}

	// Input should be cleared
	if m.textInput.Value() != "" {
		t.Errorf("Expected input to be cleared after Enter, got %q", m.textInput.Value())
	}
}

func TestTUI_Quit(t *testing.T) {
	m := initialModel()

	// Ctrl+C
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatal("Expected tea.Quit command, got nil")
	}

	// In bubbletea, tea.Quit() returns a tea.QuitMsg which is what p.Run() reacts to.
	// We can't easily check the value of the command without executing it, 
	// but we can check if it returns tea.Quit.
}
