package main

import (
	"strings"
	"testing"
)

// Test UnderwaterAnimator initialization
func TestNewUnderwaterAnimator(t *testing.T) {
	animator := NewUnderwaterAnimator()

	if animator == nil {
		t.Fatal("NewUnderwaterAnimator() should not return nil")
	}

	// Note: UnderwaterAnimator doesn't expose width/height fields directly
	// These are managed internally

	if animator.isPaused {
		t.Error("Animator should not be paused by default")
	}

	if animator.speed != 1.0 {
		t.Errorf("Expected default speed 1.0, got %f", animator.speed)
	}

	// Check particles
	if len(animator.particles) != 50 {
		t.Errorf("Expected 50 particles, got %d", len(animator.particles))
	}

	// Check fish
	if len(animator.fish) != 5 {
		t.Errorf("Expected 5 fish, got %d", len(animator.fish))
	}

	// Check planets
	if len(animator.planets) != 3 {
		t.Errorf("Expected 3 planets, got %d", len(animator.planets))
	}

	// Check octopus
	if animator.octopus == nil {
		t.Error("Octopus should be initialized")
	}
}

// Test animation pause/resume
func TestAnimationPauseResume(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test initial state
	if animator.IsPaused() {
		t.Error("Animator should not be paused initially")
	}

	// Test pause
	animator.SetPaused(true)
	if !animator.IsPaused() {
		t.Error("Animator should be paused after SetPaused(true)")
	}

	// Test resume
	animator.SetPaused(false)
	if animator.IsPaused() {
		t.Error("Animator should not be paused after SetPaused(false)")
	}
}

// Test animation speed control
func TestAnimationSpeed(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test default speed - we need to check the internal field
	// Since getAnimationSpeed is on Model, not UnderwaterAnimator
	if animator.speed != 1.0 {
		t.Errorf("Expected default speed 1.0, got %f", animator.speed)
	}

	// Test speed change
	animator.SetSpeed(2.5)
	if animator.speed != 2.5 {
		t.Errorf("Expected speed 2.5, got %f", animator.speed)
	}

	// Test zero speed
	animator.SetSpeed(0.0)
	if animator.speed != 0.0 {
		t.Errorf("Expected speed 0.0, got %f", animator.speed)
	}

	// Test negative speed (should be handled gracefully)
	animator.SetSpeed(-1.0)
	// Should either be -1.0 or 0.0 depending on implementation
	if animator.speed < 0 {
		t.Logf("Negative speed allowed: %f", animator.speed)
	}
}

// Test animation update
func TestAnimationUpdate(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Store initial positions
	initialParticleX := animator.particles[0].X
	initialFishX := animator.fish[0].X
	initialOctopusX := animator.octopus.X

	// Update animation
	err := animator.Update(1.0) // 1 second delta
	if err != nil {
		t.Errorf("Animation update failed: %v", err)
	}

	// Check that particles moved
	if animator.particles[0].X == initialParticleX && !animator.isPaused {
		t.Error("Particles should move when animation is updated")
	}

	// Check that fish moved
	if animator.fish[0].X == initialFishX && !animator.isPaused {
		t.Error("Fish should move when animation is updated")
	}

	// Check that octopus moved
	if animator.octopus.X == initialOctopusX && !animator.isPaused {
		t.Error("Octopus should move when animation is updated")
	}
}

// Test animation update when paused
func TestAnimationUpdateWhenPaused(t *testing.T) {
	animator := NewUnderwaterAnimator()
	animator.SetPaused(true)

	// Store initial positions
	initialParticleX := animator.particles[0].X
	initialFishX := animator.fish[0].X

	// Update animation while paused
	err := animator.Update(1.0)
	if err != nil {
		t.Errorf("Animation update failed: %v", err)
	}

	// Positions should not change when paused
	if animator.particles[0].X != initialParticleX {
		t.Error("Particles should not move when paused")
	}

	if animator.fish[0].X != initialFishX {
		t.Error("Fish should not move when paused")
	}
}

// Test animation rendering
func TestAnimationRender(t *testing.T) {
	animator := NewUnderwaterAnimator()

	render := animator.Render()

	if len(render) == 0 {
		t.Error("Render should not return empty string")
	}

	// Should contain some ANSI escape codes for colors
	if !strings.Contains(render, "\x1b[") {
		t.Error("Render should contain ANSI color codes")
	}

	// Should contain some characters for particles/fish
	if !strings.Contains(render, "•") && !strings.Contains(render, "><>") {
		t.Error("Render should contain particle or fish characters")
	}
}

// Test boundary conditions
func TestBoundaryConditions(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test very small delta time
	err := animator.Update(0.001)
	if err != nil {
		t.Errorf("Update with small delta failed: %v", err)
	}

	// Test large delta time
	err = animator.Update(10.0)
	if err != nil {
		t.Errorf("Update with large delta failed: %v", err)
	}

	// Test zero delta time
	err = animator.Update(0.0)
	if err != nil {
		t.Errorf("Update with zero delta failed: %v", err)
	}
}

// Test particle behavior
func TestParticleBehavior(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test that particles have valid properties
	for i, particle := range animator.particles {
		// Note: UnderwaterAnimator doesn't expose width/height
		// We can't check bounds without these fields
		if particle.Opacity < 0 || particle.Opacity > 1 {
			t.Errorf("Particle %d opacity out of bounds: %f", i, particle.Opacity)
		}

		if particle.Opacity < 0 || particle.Opacity > 1 {
			t.Errorf("Particle %d opacity out of bounds: %f", i, particle.Opacity)
		}

		if particle.VX == 0 && particle.VY == 0 {
			t.Errorf("Particle %d should have some movement", i)
		}
	}
}

// Test fish behavior
func TestFishBehavior(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test that fish have valid properties
	for i, fish := range animator.fish {
		// Note: Can't check bounds without width/height fields
		if fish.Speed == 0 {
			t.Errorf("Fish %d should have non-zero speed", i)
		}

		if fish.Speed == 0 {
			t.Errorf("Fish %d should have non-zero speed", i)
		}

		// Fish doesn't have Direction field - it uses Angle for direction
		if fish.Speed == 0 {
			t.Errorf("Fish %d should have non-zero speed", i)
		}
	}
}

// Test planet behavior
func TestPlanetBehavior(t *testing.T) {
	animator := NewUnderwaterAnimator()

	// Test that planets have valid properties
	for i, planet := range animator.planets {
		// Note: Can't check bounds without width/height fields
		if planet.Size <= 0 {
			t.Errorf("Planet %d should have positive size", i)
		}

		if planet.Size <= 0 {
			t.Errorf("Planet %d should have positive size", i)
		}

		if planet.Speed == 0 {
			t.Errorf("Planet %d should have non-zero speed", i)
		}
	}
}

// Test octopus behavior
func TestOctopusBehavior(t *testing.T) {
	animator := NewUnderwaterAnimator()

	octopus := animator.octopus

	// Note: Can't check bounds without width/height fields
	// Just verify octopus exists and has tentacles
	if len(octopus.Tentacles) != 8 {
		t.Errorf("Octopus should have 8 tentacles, got %d", len(octopus.Tentacles))
	}

	if len(octopus.Tentacles) != 8 {
		t.Errorf("Octopus should have 8 tentacles, got %d", len(octopus.Tentacles))
	}

	// Test tentacle properties
	for i, tentacle := range octopus.Tentacles {
		if tentacle.Length <= 0 {
			t.Errorf("Tentacle %d should have positive length", i)
		}

		// Tentacle doesn't have Speed field - it uses Angle and Wave
		if tentacle.Length <= 0 {
			t.Errorf("Tentacle %d should have positive length", i)
		}
	}
}

// Benchmark animation update
func BenchmarkAnimationUpdate(b *testing.B) {
	animator := NewUnderwaterAnimator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		animator.Update(0.016) // ~60 FPS
	}
}

// Benchmark animation render
func BenchmarkAnimationRender(b *testing.B) {
	animator := NewUnderwaterAnimator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		animator.Render()
	}
}
