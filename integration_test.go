package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Test full TUI integration
func TestTUIIntegration(t *testing.T) {
	// Create temporary directory for test data
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Initialize model
	model := initialModel()

	// Test initialization
	cmd := model.Init()
	if cmd == nil {
		t.Error("Model initialization should return a command")
	}

	// Test a sequence of user interactions
	interactions := []tea.Msg{
		// Start recording
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
		// Switch to conversation pane
		tea.KeyMsg{Type: tea.KeyTab},
		// Switch to monitoring pane
		tea.KeyMsg{Type: tea.KeyTab},
		// Pause animation
		tea.KeyMsg{Type: tea.KeySpace},
		// Increase speed
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}},
		// Decrease speed
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}},
		// Reset animation
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}},
		// Clear conversation
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}},
		// Stop recording
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
	}

	for i, msg := range interactions {
		updatedModel, cmd := model.Update(msg)
		if updatedModel == nil {
			t.Errorf("Interaction %d should return a model", i)
		}

		model = updatedModel.(*Model)

		// Verify model is still valid after each interaction
		if model.animator == nil {
			t.Errorf("Interaction %d: Animator should not be nil", i)
		}

		if model.logger == nil {
			t.Errorf("Interaction %d: Logger should not be nil", i)
		}

		if len(model.panes) != 3 {
			t.Errorf("Interaction %d: Should have 3 panes", i)
		}
	}

	// Test final state
	if model.isRecording {
		t.Error("Should not be recording at the end")
	}

	// Verify that data files were created
	dataDir := filepath.Join(tempDir, ".ai-tui-data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Error("Data directory should be created")
	}
}

// Test recording workflow integration
func TestRecordingWorkflowIntegration(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	model := initialModel()

	// Start recording
	model, _ = model.toggleRecording()
	if !model.isRecording {
		t.Error("Should be recording")
	}

	// Simulate some animation updates
	for i := 0; i < 10; i++ {
		model.Update(time.Now())
	}

	// Stop recording
	model, _ = model.toggleRecording()
	if model.isRecording {
		t.Error("Should not be recording")
	}

	// Check that conversation was saved
	dataDir := filepath.Join(tempDir, ".ai-tui-data")
	files, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("Failed to read data directory: %v", err)
	}

	conversationFiles := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "conversation_") {
			conversationFiles++
		}
	}

	if conversationFiles == 0 {
		t.Error("At least one conversation file should be created")
	}
}

// Test animation and UI integration
func TestAnimationUIIntegration(t *testing.T) {
	model := initialModel()

	// Test that animation updates affect UI
	initialView := model.View()

	// Update animation
	model.Update(time.Now())

	updatedView := model.View()

	// Views should be different due to animation changes
	if initialView == updatedView {
		t.Error("Animation updates should change the view")
	}

	// Test pause/resume affects animation
	model.animator.SetPaused(true)
	pausedView := model.View()

	model.animator.SetPaused(false)
	resumedView := model.View()

	// Views might be different, but both should be valid
	if len(pausedView) == 0 {
		t.Error("Paused view should not be empty")
	}

	if len(resumedView) == 0 {
		t.Error("Resumed view should not be empty")
	}
}

// Test window resize integration
func TestWindowResizeIntegration(t *testing.T) {
	model := initialModel()

	// Test various window sizes
	sizes := []tea.WindowSizeMsg{
		{Width: 80, Height: 24},  // Minimum
		{Width: 100, Height: 40}, // Standard
		{Width: 200, Height: 60}, // Large
	}

	for i, size := range sizes {
		updatedModel, cmd := model.Update(size)
		if updatedModel == nil {
			t.Errorf("Resize %d should return a model", i)
		}

		if cmd == nil {
			t.Errorf("Resize %d should return a command", i)
		}

		model = updatedModel.(*Model)

		// Verify animator was updated
		if model.animator.width != size.Width {
			t.Errorf("Resize %d: Animator width should be updated", i)
		}

		if model.animator.height != size.Height {
			t.Errorf("Resize %d: Animator height should be updated", i)
		}
	}
}

// Test mouse interaction integration
func TestMouseInteractionIntegration(t *testing.T) {
	model := initialModel()

	// Test clicking on each pane
	for paneIndex := 0; paneIndex < len(model.panes); paneIndex++ {
		pane := &model.panes[paneIndex]

		// Click in the middle of the pane
		mouseMsg := tea.MouseMsg{
			Type: tea.MouseLeft,
			X:    pane.X + pane.Width/2,
			Y:    pane.Y + pane.Height/2,
		}

		updatedModel, cmd := model.Update(mouseMsg)
		if updatedModel == nil {
			t.Errorf("Click on pane %d should return a model", paneIndex)
		}

		model = updatedModel.(*Model)

		// Verify active pane changed
		if model.activePane != paneIndex {
			t.Errorf("Click on pane %d should make it active", paneIndex)
		}
	}
}

// Test error handling integration
func TestErrorHandlingIntegration(t *testing.T) {
	// Test with invalid home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", "/invalid/directory")
	defer os.Setenv("HOME", originalHome)

	// Model should still initialize, even with invalid directory
	model := initialModel()

	if model == nil {
		t.Error("Model should initialize even with invalid home directory")
	}

	// Recording should handle errors gracefully
	model, _ = model.toggleRecording()

	// Should not panic
	if model == nil {
		t.Error("Model should still exist after recording toggle with invalid directory")
	}
}

// Test concurrent access
func TestConcurrentAccess(t *testing.T) {
	model := initialModel()

	// Test concurrent updates
	done := make(chan bool, 2)

	// Goroutine 1: Animation updates
	go func() {
		for i := 0; i < 100; i++ {
			model.Update(time.Now())
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 2: UI updates
	go func() {
		for i := 0; i < 50; i++ {
			model.View()
			time.Sleep(time.Millisecond * 2)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Model should still be valid
	if model.animator == nil {
		t.Error("Animator should still exist after concurrent access")
	}
}

// Test memory usage
func TestMemoryUsage(t *testing.T) {
	model := initialModel()

	// Simulate extended usage
	for i := 0; i < 1000; i++ {
		// Animation updates
		model.Update(time.Now())

		// UI updates
		model.View()

		// Occasional user interactions
		if i%100 == 0 {
			model.Update(tea.KeyMsg{Type: tea.KeyTab})
		}

		if i%200 == 0 {
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
		}
	}

	// Model should still be functional
	finalView := model.View()
	if len(finalView) == 0 {
		t.Error("Model should still produce valid view after extended use")
	}
}

// Test data persistence integration
func TestDataPersistenceIntegration(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create first model and record data
	model1 := initialModel()
	model1, _ = model1.toggleRecording()

	// Simulate some activity
	for i := 0; i < 5; i++ {
		model1.Update(time.Now())
		time.Sleep(time.Millisecond * 10)
	}

	model1, _ = model1.toggleRecording()

	// Create second model and verify data persistence
	model2 := initialModel()

	// Both models should have valid loggers
	if model1.logger == nil || model2.logger == nil {
		t.Error("Both models should have loggers")
	}

	// Data directory should exist
	dataDir := filepath.Join(tempDir, ".ai-tui-data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Error("Data directory should exist")
	}
}

// Test full workflow simulation
func TestFullWorkflowSimulation(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	model := initialModel()

	// Simulate typical user session
	workflow := []struct {
		name string
		msg  tea.Msg
	}{
		{"Start recording", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}},
		{"Switch to conversation", tea.KeyMsg{Type: tea.KeyTab}},
		{"Switch to monitoring", tea.KeyMsg{Type: tea.KeyTab}},
		{"Back to main", tea.KeyMsg{Type: tea.KeyTab}},
		{"Pause animation", tea.KeyMsg{Type: tea.KeySpace}},
		{"Resume animation", tea.KeyMsg{Type: tea.KeySpace}},
		{"Increase speed", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}}},
		{"Decrease speed", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}}},
		{"Reset animation", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}},
		{"Stop recording", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}},
	}

	for _, step := range workflow {
		updatedModel, cmd := model.Update(step.msg)
		if updatedModel == nil {
			t.Errorf("Step '%s' should return a model", step.name)
		}

		model = updatedModel.(*Model)

		// Verify model is still in valid state
		if model.animator == nil {
			t.Errorf("Step '%s': Animator should not be nil", step.name)
		}

		if len(model.panes) != 3 {
			t.Errorf("Step '%s': Should have 3 panes", step.name)
		}

		// Generate view to ensure no panics
		view := model.View()
		if len(view) == 0 {
			t.Errorf("Step '%s': View should not be empty", step.name)
		}
	}

	// Final verification
	if model.isRecording {
		t.Error("Should not be recording at end of workflow")
	}
}
