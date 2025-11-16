package main

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Test initialModel creation
func TestInitialModel(t *testing.T) {
	model := initialModel()

	if model.animator == nil {
		t.Error("Model should have an animator")
	}

	if model.logger == nil {
		t.Error("Model should have a logger")
	}

	if len(model.panes) != 3 {
		t.Errorf("Expected 3 panes, got %d", len(model.panes))
	}

	if model.activePane != 0 {
		t.Errorf("Expected active pane 0, got %d", model.activePane)
	}

	if model.isRecording {
		t.Error("Model should not be recording initially")
	}

	// Check pane properties
	paneIDs := []string{"main", "conversation", "monitoring"}
	for i, expectedID := range paneIDs {
		if model.panes[i].ID != expectedID {
			t.Errorf("Pane %d should have ID %s, got %s", i, expectedID, model.panes[i].ID)
		}

		if model.panes[i].Width <= 0 {
			t.Errorf("Pane %d should have positive width", i)
		}

		if model.panes[i].Height <= 0 {
			t.Errorf("Pane %d should have positive height", i)
		}

		if i == 0 && !model.panes[i].IsActive {
			t.Error("First pane should be active")
		}
	}
}

// Test model initialization
func TestModelInit(t *testing.T) {
	model := initialModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Init should return a command")
	}
}

// Test model update with key messages
func TestModelUpdateKeyMsg(t *testing.T) {
	model := initialModel()

	// Test Tab key
	tabMsg := tea.KeyMsg{Type: tea.KeyTab, Runes: []rune{'\t'}}
	updatedModel, cmd := model.Update(tabMsg)

	if updatedModel == nil {
		t.Error("Update should return a model")
	}

	if cmd == nil {
		t.Error("Update should return a command")
	}

	// Check that active pane changed
	modelPtr, ok := updatedModel.(*Model)
	if !ok {
		t.Fatal("Failed to type assert updated model")
	}
	if modelPtr.activePane == model.activePane {
		t.Error("Tab key should change active pane")
	}

	// Test 's' key for recording toggle
	sMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updatedModel, _ = model.Update(sMsg)

	modelPtr, ok = updatedModel.(*Model)
	if !ok {
		t.Fatal("Failed to type assert updated model")
	}
	if modelPtr.isRecording == model.isRecording {
		t.Error("'s' key should toggle recording state")
	}

	// Test 'c' key for clearing conversation
	cMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	updatedModel, _ = model.Update(cMsg)

	// Should not panic and should return a valid model
	if updatedModel == nil {
		t.Error("'c' key update should return a model")
	}

	// Test 'm' key for monitoring
	mMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}
	updatedModel, _ = model.Update(mMsg)

	if updatedModel == nil {
		t.Error("'m' key update should return a model")
	}

	// Test space key for pause/resume
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	updatedModel, _ = model.Update(spaceMsg)

	if updatedModel == nil {
		t.Error("Space key update should return a model")
	}

	// Test '+' key for speed increase
	plusMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}}
	updatedModel, _ = model.Update(plusMsg)

	if updatedModel == nil {
		t.Error("'+' key update should return a model")
	}

	// Test '-' key for speed decrease
	minusMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}}
	updatedModel, _ = model.Update(minusMsg)

	if updatedModel == nil {
		t.Error("'-' key update should return a model")
	}

	// Test 'r' key for reset
	rMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updatedModel, _ = model.Update(rMsg)

	if updatedModel == nil {
		t.Error("'r' key update should return a model")
	}

	// Test Ctrl+C
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd = model.Update(ctrlCMsg)

	if updatedModel == nil {
		t.Error("Ctrl+C update should return a model")
	}

	if cmd == nil {
		t.Error("Ctrl+C should return a quit command")
	}
}

// Test model update with time messages (animation ticks)
func TestModelUpdateTimeMsg(t *testing.T) {
	model := initialModel()

	// Test time message for animation
	timeMsg := time.Now()
	updatedModel, cmd := model.Update(timeMsg)

	if updatedModel == nil {
		t.Error("Time update should return a model")
	}

	if cmd == nil {
		t.Error("Time update should return a command")
	}
}

// Test model update with window size messages
func TestModelUpdateWindowSizeMsg(t *testing.T) {
	model := initialModel()

	// Test window size message
	sizeMsg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	updatedModel, cmd := model.Update(sizeMsg)

	if updatedModel == nil {
		t.Error("Window size update should return a model")
	}

	if cmd == nil {
		t.Error("Window size update should return a command")
	}

	// Check that animator was updated
	updatedAnimator := updatedModel.(*Model).animator
	if updatedAnimator == nil {
		t.Error("Animator should still exist after window size update")
	}
}

// Test model update with mouse messages
func TestModelUpdateMouseMsg(t *testing.T) {
	model := initialModel()

	// Test mouse click
	mouseMsg := tea.MouseMsg{
		Type: tea.MouseLeft,
		X:    10,
		Y:    5,
	}

	updatedModel, cmd := model.Update(mouseMsg)

	if updatedModel == nil {
		t.Error("Mouse update should return a model")
	}

	if cmd == nil {
		t.Error("Mouse update should return a command")
	}
}

// Test model view rendering
func TestModelView(t *testing.T) {
	model := initialModel()

	view := model.View()

	if len(view) == 0 {
		t.Error("View should not return empty string")
	}

	// Should contain pane titles
	if !strings.Contains(view, "AI Workspace") {
		t.Error("View should contain main pane title")
	}

	if !strings.Contains(view, "Conversation") {
		t.Error("View should contain conversation pane title")
	}

	if !strings.Contains(view, "Monitor") {
		t.Error("View should contain monitoring pane title")
	}

	if !strings.Contains(view, "Conversation") {
		t.Error("View should contain conversation pane title")
	}

	if !strings.Contains(view, "Monitor") {
		t.Error("View should contain monitoring pane title")
	}
}

// Test pane switching
func TestPaneSwitching(t *testing.T) {
	model := initialModel()

	// Test switching to next pane
	initialPane := model.activePane
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(*Model)

	if model.activePane == initialPane {
		t.Error("Tab should switch to next pane")
	}

	// Test switching through all panes
	for i := 0; i < len(model.panes)*2; i++ {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model = updatedModel.(*Model)

		if model.activePane < 0 || model.activePane >= len(model.panes) {
			t.Errorf("Active pane %d is out of bounds", model.activePane)
		}
	}
}

// Test recording functionality
func TestRecordingFunctionality(t *testing.T) {
	model := initialModel()

	// Test starting recording
	if model.isRecording {
		t.Error("Should not be recording initially")
	}

	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model = updatedModel.(*Model)

	if !model.isRecording {
		t.Error("Should be recording after toggle")
	}

	// Test stopping recording
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model = updatedModel.(*Model)

	if model.isRecording {
		t.Error("Should not be recording after second toggle")
	}
}

// Test conversation clearing
func TestClearConversation(t *testing.T) {
	model := initialModel()

	// Start recording first
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model = updatedModel.(*Model)

	// Add some content to conversation
	model.panes[1].Content = "Test conversation content"

	// Clear conversation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	model = updatedModel.(*Model)

	// Check that conversation was cleared
	expectedContent := "Conversation cleared.\n\nPress 's' to start a new recording session."
	if model.panes[1].Content != expectedContent {
		t.Errorf("Expected cleared conversation content, got: %s", model.panes[1].Content)
	}
}

// Test monitoring display
func TestMonitoringDisplay(t *testing.T) {
	model := initialModel()

	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(*Model)

	// Check that monitoring pane was updated
	monitoringContent := model.panes[2].Content

	if len(monitoringContent) == 0 {
		t.Error("Monitoring content should not be empty")
	}

	// Should contain some monitoring information
	if !strings.Contains(monitoringContent, "System") && !strings.Contains(monitoringContent, "Status") {
		t.Error("Monitoring content should contain system information")
	}
}

// Test conversation formatting
func TestConversationFormatting(t *testing.T) {
	model := initialModel()

	// Test with empty conversation
	formatted := model.formatConversationDisplay()

	if len(formatted) == 0 {
		t.Error("Formatted conversation should not be empty")
	}

	// Should contain default message
	if !strings.Contains(formatted, "No conversation") {
		t.Error("Should contain no conversation message")
	}
}

// Test monitoring formatting
func TestMonitoringFormatting(t *testing.T) {
	model := initialModel()

	formatted := model.formatMonitoringDisplay()

	if len(formatted) == 0 {
		t.Error("Formatted monitoring should not be empty")
	}

	// Should contain system status
	if !strings.Contains(formatted, "System") {
		t.Error("Should contain system information")
	}
}

// Test animation speed control
func TestAnimationSpeedControl(t *testing.T) {
	model := initialModel()

	initialSpeed := model.getAnimationSpeed()

	// Test speed increase
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	model = updatedModel.(*Model)
	newSpeed := model.getAnimationSpeed()

	if newSpeed <= initialSpeed {
		t.Error("'+' should increase animation speed")
	}

	// Test speed decrease
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	model = updatedModel.(*Model)
	decreasedSpeed := model.getAnimationSpeed()

	if decreasedSpeed >= newSpeed {
		t.Error("'-' should decrease animation speed")
	}
}

// Test animation reset
func TestAnimationReset(t *testing.T) {
	model := initialModel()

	// Change animation state
	model.animator.SetSpeed(2.0)
	model.animator.SetPaused(true)

	// Reset animation
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	model = updatedModel.(*Model)

	// Check that animation was reset
	if model.getAnimationSpeed() != 1.0 {
		t.Error("Animation speed should be reset to 1.0")
	}

	if model.animator.IsPaused() {
		t.Error("Animation should not be paused after reset")
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	model := initialModel()

	// Test unknown key
	unknownMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedModel, cmd := model.Update(unknownMsg)

	if updatedModel == nil {
		t.Error("Unknown key should still return a model")
	}

	if cmd != nil {
		t.Error("Unknown key should not return a command")
	}

	// Test mouse click outside all panes
	outsideMouseMsg := tea.MouseMsg{
		Type: tea.MouseLeft,
		X:    1000,
		Y:    1000,
	}

	updatedModel, cmd = model.Update(outsideMouseMsg)

	if updatedModel == nil {
		t.Error("Outside mouse click should still return a model")
	}

	// Test empty time message
	emptyTimeMsg := time.Time{}
	updatedModel, cmd = model.Update(emptyTimeMsg)

	if updatedModel == nil {
		t.Error("Empty time message should still return a model")
	}
}

// Benchmark model update
func BenchmarkModelUpdate(b *testing.B) {
	model := initialModel()
	timeMsg := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.Update(timeMsg)
	}
}

// Benchmark model view
func BenchmarkModelView(b *testing.B) {
	model := initialModel()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.View()
	}
}
