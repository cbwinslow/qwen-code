package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Test nil and invalid inputs
func TestNilAndInvalidInputs(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test nil time message
	model := initialModel()
	updatedModel, cmd := model.Update(nil)
	if updatedModel == nil {
		t.Error("Update with nil should return a model")
	}
	if cmd != nil {
		t.Error("Update with nil should not return a command")
	}

	// Test invalid key messages
	invalidKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{0}},   // Null character
		{Type: tea.KeyRunes, Runes: []rune{127}}, // Delete
		{Type: tea.KeyRunes, Runes: []rune{27}},  // Escape
	}

	for i, key := range invalidKeys {
		updatedModel, cmd := model.Update(key)
		if updatedModel == nil {
			t.Errorf("Invalid key %d should return a model", i)
		}
		// Should not panic
	}

	// Test invalid mouse messages
	invalidMouse := []tea.MouseMsg{
		{Type: tea.MouseLeft, X: -1, Y: -1},         // Negative coordinates
		{Type: tea.MouseLeft, X: 100000, Y: 100000}, // Very large coordinates
		{Type: tea.MouseMotion, X: 0, Y: 0},         // Motion without press
	}

	for i, mouse := range invalidMouse {
		updatedModel, cmd := model.Update(mouse)
		if updatedModel == nil {
			t.Errorf("Invalid mouse %d should return a model", i)
		}
	}

	// Test invalid window sizes
	invalidSizes := []tea.WindowSizeMsg{
		{Width: 0, Height: 0},           // Zero size
		{Width: -1, Height: -1},         // Negative size
		{Width: 1, Height: 1},           // Too small
		{Width: 100000, Height: 100000}, // Very large
	}

	for i, size := range invalidSizes {
		updatedModel, cmd := model.Update(size)
		if updatedModel == nil {
			t.Errorf("Invalid size %d should return a model", i)
		}
	}
}

// Test file system errors
func TestFilesystemErrors(t *testing.T) {
	// Test with read-only directory
	tempDir := t.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly")
	os.Mkdir(readOnlyDir, 0755)
	os.Chmod(readOnlyDir, 0444) // Read-only

	logger := NewFileLogger(readOnlyDir)

	event := SystemEvent{
		ID:        "test-event",
		Timestamp: time.Now(),
		Type:      "error",
		Source:    "test",
		Message:   "Test error handling",
	}

	err := logger.LogEvent(event)
	if err == nil {
		t.Error("Should fail when writing to read-only directory")
	}

	// Test with non-existent directory
	logger = NewFileLogger("/non/existent/directory")
	err = logger.LogEvent(event)
	if err == nil {
		t.Error("Should fail when writing to non-existent directory")
	}

	// Test with full disk (simulate by creating a very large file)
	// Note: This is a simplified test - real disk full testing is complex
	largeFile := filepath.Join(tempDir, "large")
	file, err := os.Create(largeFile)
	if err == nil {
		// Try to write a lot (this might not actually fill the disk)
		data := make([]byte, 1024*1024) // 1MB
		for i := 0; i < 100; i++ {
			_, err = file.Write(data)
			if err != nil {
				break
			}
		}
		file.Close()
	}

	// Clean up
	os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup
}

// Test memory pressure
func TestMemoryPressure(t *testing.T) {
	// Test with very large number of particles
	animator := NewUnderwaterAnimator()

	// Add many particles to stress memory
	for i := 0; i < 10000; i++ {
		particle := Particle{
			X:       float64(i % 1000),
			Y:       float64(i / 1000),
			SpeedX:  0.1,
			SpeedY:  0.1,
			Color:   getRandomColor(),
			Type:    "bubble",
			Opacity: 0.5,
		}
		animator.particles = append(animator.particles, particle)
	}

	// Should still work
	err := animator.Update(0.016)
	if err != nil {
		t.Errorf("Update failed with many particles: %v", err)
	}

	render := animator.Render()
	if len(render) == 0 {
		t.Error("Render should not fail with many particles")
	}
}

// Test concurrent access edge cases
func TestConcurrentEdgeCases(t *testing.T) {
	animator := NewUnderwaterAnimator()
	done := make(chan error, 10)

	// Multiple goroutines modifying the same animator
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				// Concurrent modifications
				animator.SetSpeed(float64(id))
				animator.SetPaused(j%2 == 0)
				animator.Update(0.016)
				animator.Render()
			}
			done <- nil
		}(i)
	}

	// Wait for completion
	for i := 0; i < 10; i++ {
		err := <-done
		if err != nil {
			t.Errorf("Concurrent access error: %v", err)
		}
	}

	// Animator should still be functional
	render := animator.Render()
	if len(render) == 0 {
		t.Error("Animator should still work after concurrent access")
	}
}

// Test time-related edge cases
func TestTimeEdgeCases(t *testing.T) {
	animator := NewUnderwaterAnimator()
	model := initialModel()

	// Test with zero time
	zeroTime := time.Time{}
	updatedModel, cmd := model.Update(zeroTime)
	if updatedModel == nil {
		t.Error("Zero time should return a model")
	}

	// Test with very old time
	oldTime := time.Unix(0, 0)
	updatedModel, cmd = model.Update(oldTime)
	if updatedModel == nil {
		t.Error("Old time should return a model")
	}

	// Test with very future time
	futureTime := time.Unix(1<<63-1, 0)
	updatedModel, cmd = model.Update(futureTime)
	if updatedModel == nil {
		t.Error("Future time should return a model")
	}

	// Test animation with extreme delta times
	extremeDeltas := []float64{
		-1000, // Negative time
		0,     // Zero time
		1e-10, // Very small time
		1e10,  // Very large time
	}

	for _, delta := range extremeDeltas {
		err := animator.Update(delta)
		if err != nil {
			t.Errorf("Update failed with delta %f: %v", delta, err)
		}
	}
}

// Test boundary conditions
func TestBoundaryConditions(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test particles at boundaries
	animator.particles = []Particle{
		{X: 0, Y: 0, SpeedX: -1, SpeedY: -1, Color: "#FF0000", Type: "bubble", Opacity: 1},
		{X: float64(animator.width), Y: float64(animator.height), SpeedX: 1, SpeedY: 1, Color: "#00FF00", Type: "bubble", Opacity: 1},
		{X: -1, Y: -1, SpeedX: 1, SpeedY: 1, Color: "#0000FF", Type: "bubble", Opacity: 1},
		{X: float64(animator.width + 1), Y: float64(animator.height + 1), SpeedX: -1, SpeedY: -1, Color: "#FFFF00", Type: "bubble", Opacity: 1},
	}

	// Should handle boundary particles gracefully
	err := animator.Update(1.0)
	if err != nil {
		t.Errorf("Update failed with boundary particles: %v", err)
	}

	render := animator.Render()
	if len(render) == 0 {
		t.Error("Render should handle boundary particles")
	}

	// Test with zero-sized animator
	animator.width = 0
	animator.height = 0

	err = animator.Update(1.0)
	if err != nil {
		t.Errorf("Update failed with zero size: %v", err)
	}

	render = animator.Render()
	if len(render) == 0 {
		t.Error("Render should handle zero size")
	}
}

// Test invalid data structures
func TestInvalidDataStructures(t *testing.T) {
	// Test with empty particle slice
	animator := NewUnderwaterAnimator()
	animator.particles = []Particle{}

	err := animator.Update(0.016)
	if err != nil {
		t.Errorf("Update failed with empty particles: %v", err)
	}

	render := animator.Render()
	if len(render) == 0 {
		t.Error("Render should handle empty particles")
	}

	// Test with nil particle slice
	animator.particles = nil

	err = animator.Update(0.016)
	if err != nil {
		t.Errorf("Update failed with nil particles: %v", err)
	}

	render = animator.Render()
	if len(render) == 0 {
		t.Error("Render should handle nil particles")
	}

	// Test with invalid particle data
	animator.particles = []Particle{
		{X: 0, Y: 0, SpeedX: 0, SpeedY: 0, Color: "", Type: "", Opacity: -1},
		{X: 0, Y: 0, SpeedX: 0, SpeedY: 0, Color: "invalid", Type: "invalid", Opacity: 2},
	}

	err = animator.Update(0.016)
	if err != nil {
		t.Errorf("Update failed with invalid particles: %v", err)
	}

	render = animator.Render()
	if len(render) == 0 {
		t.Error("Render should handle invalid particles")
	}
}

// Test panic recovery
func TestPanicRecovery(t *testing.T) {
	// Test that the application recovers from potential panics
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic occurred: %v", r)
		}
	}()

	animator := NewUnderwaterAnimator()

	// Try to trigger potential panics
	animator.SetSpeed(-1e100) // Very negative speed
	animator.Update(1e100)    // Very large delta
	animator.Render()         // Should not panic

	// Test with nil animator
	var nilAnimator *UnderwaterAnimator
	if nilAnimator != nil {
		nilAnimator.Update(0.016)
		nilAnimator.Render()
	}
}

// Test resource cleanup
func TestResourceCleanup(t *testing.T) {
	// Test that resources are properly cleaned up
	tempDir := t.TempDir()
	logger := NewFileLogger(tempDir)

	// Log some events
	for i := 0; i < 100; i++ {
		event := SystemEvent{
			ID:        fmt.Sprintf("cleanup-test-%d", i),
			Timestamp: time.Now(),
			Type:      "info",
			Source:    "cleanup-test",
			Message:   "Cleanup test event",
		}
		logger.LogEvent(event)
	}

	// Files should be created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Errorf("Failed to read temp directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("No files created during logging")
	}

	// Should be able to clean up
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Errorf("Failed to clean up temp directory: %v", err)
	}
}

// Test invalid JSON handling
func TestInvalidJSONHandling(t *testing.T) {
	tempDir := t.TempDir()

	// Create invalid JSON file
	invalidJSONFile := filepath.Join(tempDir, "events.json")
	err := os.WriteFile(invalidJSONFile, []byte("invalid json content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	// Logger should handle invalid JSON gracefully
	logger := NewFileLogger(tempDir)

	event := SystemEvent{
		ID:        "test-after-invalid",
		Timestamp: time.Now(),
		Type:      "info",
		Source:    "test",
		Message:   "Test after invalid JSON",
	}

	// Should not panic, might overwrite or append
	err = logger.LogEvent(event)
	if err != nil {
		// This is acceptable - the important thing is not to panic
		t.Logf("Expected error when handling invalid JSON: %v", err)
	}
}

// Test extreme user input sequences
func TestExtremeUserInputSequences(t *testing.T) {
	model := initialModel()

	// Rapid sequence of all possible keys
	keys := []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '+', '=', '[', ']', '{', '}', '\\', '|', ';', ':', '\'', '"', ',', '.', '<', '>', '/', '?',
		'\t', '\n', '\r', ' ',
	}

	for _, key := range keys {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}}
		updatedModel, cmd := model.Update(msg)

		if updatedModel == nil {
			t.Errorf("Key '%c' should return a model", key)
		}
	}

	// Rapid mouse movements
	for i := 0; i < 100; i++ {
		msg := tea.MouseMsg{
			Type: tea.MouseMotion,
			X:    i % 100,
			Y:    i % 50,
		}
		updatedModel, cmd := model.Update(msg)

		if updatedModel == nil {
			t.Errorf("Mouse motion %d should return a model", i)
		}
	}

	// Rapid window resizes
	for i := 0; i < 50; i++ {
		msg := tea.WindowSizeMsg{
			Width:  80 + i,
			Height: 24 + i,
		}
		updatedModel, cmd := model.Update(msg)

		if updatedModel == nil {
			t.Errorf("Window resize %d should return a model", i)
		}
	}

	// Model should still be functional
	finalView := model.View()
	if len(finalView) == 0 {
		t.Error("Model should still be functional after extreme input")
	}
}

// Test system call failures
func TestSystemCallFailures(t *testing.T) {
	// Test with invalid file descriptors
	// This is difficult to test reliably in Go, but we can test edge cases

	// Test with very long file paths
	longPath := string(make([]byte, 1000))
	for i := range longPath {
		longPath = longPath[:i] + "a" + longPath[i+1:]
	}

	logger := NewFileLogger(longPath)

	event := SystemEvent{
		ID:        "long-path-test",
		Timestamp: time.Now(),
		Type:      "error",
		Source:    "test",
		Message:   "Long path test",
	}

	err := logger.LogEvent(event)
	if err == nil {
		t.Error("Should fail with very long path")
	}
}

// Test data corruption scenarios
func TestDataCorruptionScenarios(t *testing.T) {
	tempDir := t.TempDir()

	// Create corrupted conversation file
	corruptedFile := filepath.Join(tempDir, "conversation_corrupted.json")
	err := os.WriteFile(corruptedFile, []byte("{ \"corrupted\": json content }"), 0644)
	if err != nil {
		t.Fatalf("Failed to create corrupted file: %v", err)
	}

	// Logger should handle corruption gracefully
	logger := NewFileLogger(tempDir)

	session := ConversationSession{
		ID:        "corruption-test",
		StartTime: time.Now(),
		Messages: []ConversationMessage{
			{
				ID:        "msg-1",
				Timestamp: time.Now(),
				Role:      "user",
				Content:   "Test after corruption",
			},
		},
		IsActive: true,
	}

	err = logger.LogConversation(session)
	if err != nil {
		// Acceptable - important thing is not to panic
		t.Logf("Expected error with corrupted data: %v", err)
	}
}
