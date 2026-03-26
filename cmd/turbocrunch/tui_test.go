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

	if !strings.Contains(view, "Backend: SpeedCrunch") {
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
	if !strings.Contains(view, "Backend: Go") {
		t.Errorf("Expected backend to be Go, got %q", view)
	}
	if !strings.Contains(view, "Switched to Go") {
		t.Errorf("Expected history to contain switch message, got %q", view)
	}

	// Toggle back to SpeedCrunch
	newModel, _ = m.Update(msg)
	m = newModel.(model)
	view = m.View()
	if !strings.Contains(view, "Backend: SpeedCrunch") {
		t.Errorf("Expected backend to be SpeedCrunch, got %q", view)
	}
	if !strings.Contains(view, "Switched to SpeedCrunch") {
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
	newModel, cmd := m.Update(msg)
	m = newModel.(model)

	if cmd == nil {
		t.Fatal("Expected command after Enter, got nil")
	}

	// Manually execute the command and send result back
	resultMsg := cmd()
	newModel, _ = m.Update(resultMsg)
	m = newModel.(model)

	// Verify history update in View
	view := m.View()
	if !strings.Contains(view, "2+2") {
		t.Errorf("Expected table to show input, got %q", view)
	}
	if !strings.Contains(view, "4") {
		t.Errorf("Expected table to show result '4', got %q", view)
	}

	// Input should be cleared
	if m.textInput.Value() != "" {
		t.Errorf("Expected input to be cleared after Enter, got %q", m.textInput.Value())
	}
}

func TestTUI_FormulaBook(t *testing.T) {
	m := initialModel()

	// Switch to formula book
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}
	newModel, _ := m.Update(msg)
	m = newModel.(model)

	if m.view != viewFormulaBook {
		t.Errorf("Expected view to be viewFormulaBook, got %v", m.view)
	}

	view := m.View()
	if !strings.Contains(view, "Formula Book") {
		t.Errorf("Expected view to contain 'Formula Book', got %q", view)
	}

	// Select first formula (Circle Area)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = m.Update(enterMsg)
	m = newModel.(model)

	if m.view != viewInput {
		t.Errorf("Expected view to be viewInput after selecting formula, got %v", m.view)
	}

	if !strings.Contains(m.textInput.Value(), "pi * r^2") {
		t.Errorf("Expected input to contain 'pi * r^2', got %q", m.textInput.Value())
	}
}

func TestTUI_ConstantsAndUnits(t *testing.T) {
	m := initialModel()

	// Switch to constants
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, _ := m.Update(msg)
	m = newModel.(model)

	if m.view != viewConstants {
		t.Errorf("Expected view to be viewConstants, got %v", m.view)
	}

	if !strings.Contains(m.View(), "Constants") {
		t.Errorf("Expected view to contain 'Constants', got %q", m.View())
	}

	// Select first constant (e.g., pi or c)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = m.Update(enterMsg)
	m = newModel.(model)

	if m.view != viewInput {
		t.Errorf("Expected view to be viewInput after selecting constant, got %v", m.view)
	}
	if m.textInput.Value() == "" {
		t.Error("Expected input to contain selected constant, got empty string")
	}

	// Switch to units
	msg = tea.KeyMsg{Type: tea.KeyCtrlU}
	newModel, _ = m.Update(msg)
	m = newModel.(model)

	if m.view != viewUnits {
		t.Errorf("Expected view to be viewUnits, got %v", m.view)
	}

	if !strings.Contains(m.View(), "Units") {
		t.Errorf("Expected view to contain 'Units', got %q", m.View())
	}

	// Select first unit
	newModel, _ = m.Update(enterMsg)
	m = newModel.(model)

	if m.view != viewInput {
		t.Errorf("Expected view to be viewInput after selecting unit, got %v", m.view)
	}
}

func TestTUI_Quit(t *testing.T) {
	m := initialModel()

	// Ctrl+Q
	msg := tea.KeyMsg{Type: tea.KeyCtrlQ}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatal("Expected tea.Quit command, got nil")
	}

	// In bubbletea, tea.Quit() returns a tea.QuitMsg which is what p.Run() reacts to.
	// We can't easily check the value of the command without executing it, 
	// but we can check if it returns tea.Quit.
}
