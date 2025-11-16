package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Test animation performance under load
func TestAnimationPerformance(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Measure animation update performance
	start := time.Now()
	iterations := 10000

	for i := 0; i < iterations; i++ {
		animator.Update(0.016) // ~60 FPS
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)

	t.Logf("Animation update: %v total, %v average", duration, avgTime)

	// Should be fast enough for 60 FPS
	if avgTime > time.Millisecond {
		t.Errorf("Animation update too slow: %v (should be < 1ms)", avgTime)
	}
}

// Test rendering performance
func TestRenderingPerformance(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Measure rendering performance
	start := time.Now()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		render := animator.Render()
		if len(render) == 0 {
			t.Error("Render should not return empty string")
		}
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)

	t.Logf("Rendering: %v total, %v average", duration, avgTime)

	// Rendering should be reasonably fast
	if avgTime > 10*time.Millisecond {
		t.Errorf("Rendering too slow: %v (should be < 10ms)", avgTime)
	}
}

// Test UI performance
func TestUIPerformance(t *testing.T) {
	model := initialModel()

	// Measure UI update performance
	start := time.Now()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		model.Update(time.Now())
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)

	t.Logf("UI update: %v total, %v average", duration, avgTime)

	// UI updates should be fast
	if avgTime > time.Millisecond {
		t.Errorf("UI update too slow: %v (should be < 1ms)", avgTime)
	}

	// Measure view rendering performance
	start = time.Now()
	for i := 0; i < iterations; i++ {
		view := model.View()
		if len(view) == 0 {
			t.Error("View should not return empty string")
		}
	}

	duration = time.Since(start)
	avgTime = duration / time.Duration(iterations)

	t.Logf("View rendering: %v total, %v average", duration, avgTime)

	// View rendering should be fast
	if avgTime > 5*time.Millisecond {
		t.Errorf("View rendering too slow: %v (should be < 5ms)", avgTime)
	}
}

// Test memory usage during extended operation
func TestMemoryUsageExtended(t *testing.T) {
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	model := initialModel()

	// Simulate extended operation
	iterations := 10000
	for i := 0; i < iterations; i++ {
		// Animation updates
		model.Update(time.Now())

		// View rendering
		model.View()

		// Occasional user interactions
		if i%100 == 0 {
			model.Update(tea.KeyMsg{Type: tea.KeyTab})
		}

		if i%1000 == 0 {
			runtime.GC() // Force garbage collection periodically
		}
	}

	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	memUsed := m2.Alloc - m1.Alloc
	t.Logf("Memory used: %d bytes", memUsed)

	// Memory usage should be reasonable (less than 10MB)
	if memUsed > 10*1024*1024 {
		t.Errorf("Memory usage too high: %d bytes", memUsed)
	}
}

// Test logging performance
func TestLoggingPerformance(t *testing.T) {
	tempDir := t.TempDir()
	logger := NewFileLogger(tempDir)

	event := SystemEvent{
		ID:        "perf-test",
		Timestamp: time.Now(),
		Type:      "info",
		Source:    "performance-test",
		Message:   "Performance test event",
	}

	// Measure logging performance
	start := time.Now()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		event.ID = fmt.Sprintf("perf-test-%d", i)
		logger.LogEvent(event)
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)

	t.Logf("Logging: %v total, %v average", duration, avgTime)

	// Logging should be reasonably fast
	if avgTime > 10*time.Millisecond {
		t.Errorf("Logging too slow: %v (should be < 10ms)", avgTime)
	}
}

// Test stress with many particles
func TestStressManyParticles(t *testing.T) {
	// Create animator with many more particles than default
	animator := NewUnderwaterAnimator()

	// Add many more particles
	for i := 0; i < 500; i++ { // Add 500 more particles
		particle := Particle{
			X:       float64(i % animator.width),
			Y:       float64(i % animator.height),
			SpeedX:  (rand.Float64() - 0.5) * 2,
			SpeedY:  (rand.Float64() - 0.5) * 2,
			Color:   getRandomColor(),
			Type:    "bubble",
			Opacity: rand.Float64(),
		}
		animator.particles = append(animator.particles, particle)
	}

	t.Logf("Total particles: %d", len(animator.particles))

	// Test performance with many particles
	start := time.Now()
	iterations := 100

	for i := 0; i < iterations; i++ {
		animator.Update(0.016)
		animator.Render()
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)

	t.Logf("Stress test (%d particles): %v total, %v average", len(animator.particles), duration, avgTime)

	// Should still be reasonably fast even with many particles
	if avgTime > 50*time.Millisecond {
		t.Errorf("Performance too slow with many particles: %v", avgTime)
	}
}

// Test stress with rapid user interactions
func TestStressRapidInteractions(t *testing.T) {
	model := initialModel()

	// Simulate rapid user interactions
	interactions := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyTab},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
		tea.KeyMsg{Type: tea.KeySpace},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}},
		tea.MouseMsg{Type: tea.MouseLeft, X: 10, Y: 10},
		tea.WindowSizeMsg{Width: 100, Height: 50},
	}

	start := time.Now()
	iterations := 1000

	for i := 0; i < iterations; i++ {
		// Rapid sequence of interactions
		for _, msg := range interactions {
			model.Update(msg)
		}

		// Also include animation updates
		model.Update(time.Now())
		model.View()
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)

	t.Logf("Rapid interactions stress test: %v total, %v average", duration, avgTime)

	// Should handle rapid interactions gracefully
	if avgTime > 100*time.Millisecond {
		t.Errorf("Rapid interactions too slow: %v", avgTime)
	}
}

// Test concurrent stress
func TestConcurrentStress(t *testing.T) {
	model := initialModel()
	done := make(chan bool, 4)
	iterations := 1000

	// Goroutine 1: Animation updates
	go func() {
		for i := 0; i < iterations; i++ {
			model.Update(time.Now())
		}
		done <- true
	}()

	// Goroutine 2: View rendering
	go func() {
		for i := 0; i < iterations; i++ {
			model.View()
		}
		done <- true
	}()

	// Goroutine 3: User interactions
	go func() {
		for i := 0; i < iterations/10; i++ {
			model.Update(tea.KeyMsg{Type: tea.KeyTab})
			model.Update(tea.KeyMsg{Type: tea.KeySpace})
		}
		done <- true
	}()

	// Goroutine 4: Window resizes
	go func() {
		for i := 0; i < iterations/100; i++ {
			width := 80 + (i % 50)
			height := 24 + (i % 20)
			model.Update(tea.WindowSizeMsg{Width: width, Height: height})
		}
		done <- true
	}()

	// Wait for all goroutines
	start := time.Now()
	for i := 0; i < 4; i++ {
		<-done
	}
	duration := time.Since(start)

	t.Logf("Concurrent stress test: %v total", duration)

	// Should complete in reasonable time
	if duration > 10*time.Second {
		t.Errorf("Concurrent operations too slow: %v", duration)
	}

	// Model should still be functional
	finalView := model.View()
	if len(finalView) == 0 {
		t.Error("Model should still be functional after concurrent stress")
	}
}

// Test memory leak detection
func TestMemoryLeakDetection(t *testing.T) {
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Create and destroy many models
	for i := 0; i < 100; i++ {
		model := initialModel()

		// Use the model
		for j := 0; j < 10; j++ {
			model.Update(time.Now())
			model.View()
		}

		// Model should be garbage collected
		model = nil
	}

	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	memUsed := m2.Alloc - m1.Alloc
	t.Logf("Memory after creating/destroying 100 models: %d bytes", memUsed)

	// Should not leak significant memory
	if memUsed > 5*1024*1024 { // 5MB threshold
		t.Errorf("Potential memory leak: %d bytes", memUsed)
	}
}

// Test extreme values stress
func TestExtremeValuesStress(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test with extreme animation speeds
	extremeSpeeds := []float64{0, 0.001, 100, -100, 1e6, -1e6}

	for _, speed := range extremeSpeeds {
		animator.SetSpeed(speed)

		// Should not panic
		err := animator.Update(1.0)
		if err != nil {
			t.Errorf("Update failed with speed %f: %v", speed, err)
		}

		render := animator.Render()
		if len(render) == 0 {
			t.Errorf("Render failed with speed %f", speed)
		}
	}

	// Test with extreme delta times
	extremeDeltas := []float64{0, -1, 1e-6, 1e6, -1e6}

	for _, delta := range extremeDeltas {
		err := animator.Update(delta)
		if err != nil {
			t.Errorf("Update failed with delta %f: %v", delta, err)
		}
	}
}

// Test resource exhaustion
func TestResourceExhaustion(t *testing.T) {
	// Test with very large window sizes
	animator := NewUnderwaterAnimator()

	// Simulate very large window
	animator.width = 10000
	animator.height = 10000

	// Should handle gracefully
	err := animator.Update(0.016)
	if err != nil {
		t.Errorf("Update failed with large window: %v", err)
	}

	// Render might be slow but should not panic
	start := time.Now()
	render := animator.Render()
	duration := time.Since(start)

	t.Logf("Large window render: %v, length: %d", duration, len(render))

	if len(render) == 0 {
		t.Error("Render should not return empty string even for large windows")
	}
}

// Benchmark comparisons
func BenchmarkAnimationVsUI(b *testing.B) {
	animator := NewUnderwaterAnimator()
	model := initialModel()

	b.Run("Animation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			animator.Update(0.016)
		}
	})

	b.Run("UI", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			model.Update(time.Now())
		}
	})

	b.Run("Rendering", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			animator.Render()
		}
	})

	b.Run("View", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			model.View()
		}
	})
}
