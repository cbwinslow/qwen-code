package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Test FileLogger functionality
func TestFileLogger(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	logger := NewFileLogger(tempDir)

	// Test system event logging
	event := SystemEvent{
		ID:        "test-event-1",
		Timestamp: time.Now(),
		Type:      "info",
		Source:    "test",
		Message:   "Test event message",
		Data: map[string]interface{}{
			"test_key": "test_value",
			"number":   42,
		},
	}

	err := logger.LogEvent(event)
	if err != nil {
		t.Fatalf("Failed to log event: %v", err)
	}

	// Verify event file exists and contains correct data
	eventFile := filepath.Join(tempDir, "events.jsonl")
	data, err := os.ReadFile(eventFile)
	if err != nil {
		t.Fatalf("Failed to read event file: %v", err)
	}

	// The logger writes individual JSON objects, one per line
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 1 {
		t.Fatalf("Expected 1 line in events file, got %d", len(lines))
	}

	var loggedEvent SystemEvent
	err = json.Unmarshal([]byte(lines[0]), &loggedEvent)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if loggedEvent.ID != event.ID {
		t.Errorf("Expected event ID %s, got %s", event.ID, loggedEvent.ID)
	}

	if loggedEvent.Message != event.Message {
		t.Errorf("Expected message %s, got %s", event.Message, loggedEvent.Message)
	}

	// Test conversation logging
	session := ConversationSession{
		ID:        "test-session-1",
		StartTime: time.Now(),
		Messages: []ConversationMessage{
			{
				ID:         "msg-1",
				Timestamp:  time.Now(),
				Role:       "user",
				Content:    "Hello AI",
				TokenCount: 2,
				Model:      "test-model",
			},
		},
		IsActive: true,
	}

	err = logger.LogConversation(session)
	if err != nil {
		t.Fatalf("Failed to log conversation: %v", err)
	}

	// Verify conversation file
	convFile := filepath.Join(tempDir, "conversations.jsonl")
	convData, err := os.ReadFile(convFile)
	if err != nil {
		t.Fatalf("Failed to read conversation file: %v", err)
	}

	// The logger writes individual JSON objects, one per line
	lines = strings.Split(strings.TrimSpace(string(convData)), "\n")
	if len(lines) != 1 {
		t.Fatalf("Expected 1 line in conversations file, got %d", len(lines))
	}

	var loggedSession ConversationSession
	err = json.Unmarshal([]byte(lines[0]), &loggedSession)
	if err != nil {
		t.Fatalf("Failed to unmarshal conversation: %v", err)
	}

	if loggedSession.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, loggedSession.ID)
	}

	if len(loggedSession.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(loggedSession.Messages))
	}
}

// Test FileLogger error handling
func TestFileLoggerErrorHandling(t *testing.T) {
	// Test with invalid directory
	logger := NewFileLogger("/invalid/directory/that/does/not/exist")

	event := SystemEvent{
		ID:        "test-event",
		Timestamp: time.Now(),
		Type:      "error",
		Source:    "test",
		Message:   "Test error handling",
	}

	err := logger.LogEvent(event)
	if err == nil {
		t.Error("Expected error when logging to invalid directory")
	}
}

// Test conversation session edge cases
func TestConversationSessionEdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	logger := NewFileLogger(tempDir)

	// Test empty session
	emptySession := ConversationSession{
		ID:        "empty-session",
		StartTime: time.Now(),
		Messages:  []ConversationMessage{},
		IsActive:  false,
	}

	err := logger.LogConversation(emptySession)
	if err != nil {
		t.Errorf("Failed to log empty session: %v", err)
	}

	// Test session with nil metadata
	sessionWithNilMetadata := ConversationSession{
		ID:        "metadata-session",
		StartTime: time.Now(),
		Messages: []ConversationMessage{
			{
				ID:        "msg-1",
				Timestamp: time.Now(),
				Role:      "user",
				Content:   "Test message",
				Metadata:  nil, // Explicitly nil
			},
		},
		IsActive: true,
	}

	err = logger.LogConversation(sessionWithNilMetadata)
	if err != nil {
		t.Errorf("Failed to log session with nil metadata: %v", err)
	}
}

// Test ID generation
func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	if id1 == id2 {
		t.Error("generateID() should produce unique IDs")
	}

	if len(id1) == 0 {
		t.Error("generateID() should not return empty string")
	}

	if len(id1) < 10 {
		t.Error("generateID() should return reasonably long IDs")
	}
}

// Test color functions
func TestColorFunctions(t *testing.T) {
	// Test getRGBFromHex with valid hex colors
	testCases := map[string]string{
		"#FF0000": "255;0;0",
		"#00FF00": "0;255;0",
		"#0000FF": "0;0;255",
		"#FFFFFF": "255;255;255",
		"#000000": "0;0;0",
	}

	for hex, expected := range testCases {
		result := getRGBFromHex(hex)
		if result != expected {
			t.Errorf("getRGBFromHex(%s) = %s, expected %s", hex, result, expected)
		}
	}

	// Test with invalid hex
	invalidHex := "invalid"
	result := getRGBFromHex(invalidHex)
	if result != "255;255;255" {
		t.Errorf("getRGBFromHex(%s) should return default white", invalidHex)
	}

	// Test getRGBFromColor
	color := getRandomColor()
	result = getRGBFromColor(color)
	if result == "" {
		t.Error("getRGBFromColor() should not return empty string")
	}
}

// Test hexToByte function
func TestHexToByte(t *testing.T) {
	testCases := map[string]byte{
		"00": 0,
		"FF": 255,
		"80": 128,
		"0F": 15,
		"F0": 240,
	}

	for hex, expected := range testCases {
		result := hexToByte(hex)
		if result != expected {
			t.Errorf("hexToByte(%s) = %d, expected %d", hex, result, expected)
		}
	}

	// Test with invalid hex
	result := hexToByte("invalid")
	if result != 0 {
		t.Errorf("hexToByte(invalid) should return 0, got %d", result)
	}
}

// Benchmark logging performance
func BenchmarkLogEvent(b *testing.B) {
	tempDir := b.TempDir()
	logger := NewFileLogger(tempDir)

	event := SystemEvent{
		ID:        "bench-event",
		Timestamp: time.Now(),
		Type:      "info",
		Source:    "benchmark",
		Message:   "Benchmark test event",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.ID = fmt.Sprintf("bench-event-%d", i)
		logger.LogEvent(event)
	}
}

// Benchmark conversation logging
func BenchmarkLogConversation(b *testing.B) {
	tempDir := b.TempDir()
	logger := NewFileLogger(tempDir)

	session := ConversationSession{
		ID:        "bench-session",
		StartTime: time.Now(),
		Messages: []ConversationMessage{
			{
				ID:         "msg-1",
				Timestamp:  time.Now(),
				Role:       "user",
				Content:    "Benchmark test message",
				TokenCount: 4,
				Model:      "benchmark-model",
			},
		},
		IsActive: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session.ID = fmt.Sprintf("bench-session-%d", i)
		logger.LogConversation(session)
	}
}
