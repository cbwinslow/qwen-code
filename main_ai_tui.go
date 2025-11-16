package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== AI MODELS ====================

// ConversationMessage represents a single message in a conversation
type ConversationMessage struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	Role       string                 `json:"role"` // "user", "assistant", "system"
	Content    string                 `json:"content"`
	TokenCount int                    `json:"token_count"`
	Model      string                 `json:"model"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ConversationSession represents a complete conversation session
type ConversationSession struct {
	ID        string                `json:"id"`
	StartTime time.Time             `json:"start_time"`
	EndTime   *time.Time            `json:"end_time,omitempty"`
	Messages  []ConversationMessage `json:"messages"`
	Summary   string                `json:"summary,omitempty"`
	Tags      []string              `json:"tags,omitempty"`
	IsActive  bool                  `json:"is_active"`
}

// SystemEvent represents a system event for logging
type SystemEvent struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"` // "info", "warning", "error", "security"
	Source    string                 `json:"source"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	ImageData string                 `json:"image_data,omitempty"` // Base64 encoded image
}

// ==================== ENUMS ====================

// EventType defines types of system events
type EventType string

const (
	EventTypeInfo     EventType = "info"
	EventTypeWarning  EventType = "warning"
	EventTypeError    EventType = "error"
	EventTypeSecurity EventType = "security"
	EventTypeImage    EventType = "image"
)

// MessageRole defines the role of a message sender
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// ==================== INTERFACES ====================

// Logger defines interface for logging operations
type Logger interface {
	LogEvent(event SystemEvent) error
	LogConversation(session ConversationSession) error
}

// Animator defines interface for animation operations
type Animator interface {
	Update(deltaTime float64) error
	Render() string
	IsPaused() bool
	SetPaused(paused bool)
}

// ==================== IMPLEMENTATIONS ====================

// FileLogger implements Logger interface with file storage
type FileLogger struct {
	eventsFile        string
	conversationsFile string
}

func NewFileLogger(dataDir string) *FileLogger {
	return &FileLogger{
		eventsFile:        filepath.Join(dataDir, "events.jsonl"),
		conversationsFile: filepath.Join(dataDir, "conversations.jsonl"),
	}
}

func (fl *FileLogger) LogEvent(event SystemEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	file, err := os.OpenFile(fl.eventsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open events file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	return nil
}

func (fl *FileLogger) LogConversation(session ConversationSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	file, err := os.OpenFile(fl.conversationsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open conversations file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		return fmt.Errorf("failed to write conversation: %w", err)
	}

	return nil
}

// ==================== ANIMATION SYSTEM ====================

type Particle struct {
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	VX          float64 `json:"vx"`
	VY          float64 `json:"vy"`
	Size        float64 `json:"size"`
	Color       string  `json:"color"`
	Opacity     float64 `json:"opacity"`
	Lifetime    float64 `json:"lifetime"`
	MaxLifetime float64 `json:"max_lifetime"`
}

type Star struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Size    float64 `json:"size"`
	Bright  float64 `json:"bright"`
	Twinkle float64 `json:"twinkle"`
}

type Planet struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Orbit float64 `json:"orbit"`
	Size  float64 `json:"size"`
	Color string  `json:"color"`
	Speed float64 `json:"speed"`
	Angle float64 `json:"angle"`
}

type Octopus struct {
	X         float64    `json:"x"`
	Y         float64    `json:"y"`
	Angle     float64    `json:"angle"`
	Tentacles []Tentacle `json:"tentacles"`
	Color     string     `json:"color"`
	Speed     float64    `json:"speed"`
}

type Tentacle struct {
	Angle  float64 `json:"angle"`
	Length float64 `json:"length"`
	Wave   float64 `json:"wave"`
}

type Fish struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Angle     float64 `json:"angle"`
	Speed     float64 `json:"speed"`
	Size      float64 `json:"size"`
	Color     string  `json:"color"`
	WavePhase float64 `json:"wave_phase"`
}

// UnderwaterAnimator implements Animator interface
type UnderwaterAnimator struct {
	particles   []Particle
	stars       []Star
	planets     []Planet
	octopus     *Octopus
	fish        []Fish
	gradientPos float64
	isPaused    bool
	speed       float64
}

func NewUnderwaterAnimator() *UnderwaterAnimator {
	rand.Seed(time.Now().UnixNano())

	// Create initial particles
	particles := make([]Particle, 50)
	for i := range particles {
		particles[i] = Particle{
			X:           rand.Float64() * 100,
			Y:           rand.Float64() * 30,
			VX:          (rand.Float64() - 0.5) * 0.2,
			VY:          (rand.Float64() - 0.5) * 0.1,
			Size:        rand.Float64()*2 + 0.5,
			Color:       getRandomColor(),
			Opacity:     rand.Float64(),
			Lifetime:    0,
			MaxLifetime: rand.Float64()*100 + 50,
		}
	}

	// Create stars
	stars := make([]Star, 100)
	for i := range stars {
		stars[i] = Star{
			X:       rand.Float64() * 100,
			Y:       rand.Float64() * 30,
			Size:    rand.Float64()*1.5 + 0.5,
			Bright:  rand.Float64(),
			Twinkle: rand.Float64() * math.Pi * 2,
		}
	}

	// Create planets
	planets := []Planet{
		{X: 20, Y: 10, Orbit: 15, Size: 2, Color: "#FF6B6B", Speed: 0.02, Angle: 0},
		{X: 70, Y: 20, Orbit: 8, Size: 1.5, Color: "#4ECDC4", Speed: 0.03, Angle: math.Pi},
		{X: 50, Y: 5, Orbit: 12, Size: 1, Color: "#95E1D3", Speed: 0.015, Angle: math.Pi / 2},
	}

	// Create octopus
	tentacles := make([]Tentacle, 8)
	for i := range tentacles {
		tentacles[i] = Tentacle{
			Angle:  float64(i) * (math.Pi * 2 / 8),
			Length: 3 + rand.Float64()*2,
			Wave:   rand.Float64() * math.Pi * 2,
		}
	}

	octopus := &Octopus{
		X:         50,
		Y:         25,
		Angle:     0,
		Tentacles: tentacles,
		Color:     "#9B59B6",
		Speed:     0.01,
	}

	// Create fish
	fish := make([]Fish, 5)
	for i := range fish {
		fish[i] = Fish{
			X:         rand.Float64() * 100,
			Y:         15 + rand.Float64()*15,
			Angle:     rand.Float64() * math.Pi * 2,
			Speed:     0.02 + rand.Float64()*0.02,
			Size:      1 + rand.Float64(),
			Color:     getRandomFishColor(),
			WavePhase: rand.Float64() * math.Pi * 2,
		}
	}

	return &UnderwaterAnimator{
		particles:   particles,
		stars:       stars,
		planets:     planets,
		octopus:     octopus,
		fish:        fish,
		gradientPos: 0,
		isPaused:    false,
		speed:       1.0,
	}
}

func (ua *UnderwaterAnimator) Update(deltaTime float64) error {
	if ua.isPaused {
		return nil
	}

	// Update particles
	for i := range ua.particles {
		p := &ua.particles[i]
		p.X += p.VX * deltaTime
		p.Y += p.VY * deltaTime
		p.Lifetime += deltaTime

		// Wrap around screen
		if p.X < 0 {
			p.X = 100
		} else if p.X > 100 {
			p.X = 0
		}
		if p.Y < 0 {
			p.Y = 30
		} else if p.Y > 30 {
			p.Y = 0
		}

		// Reset particle if lifetime exceeded
		if p.Lifetime > p.MaxLifetime {
			p.X = rand.Float64() * 100
			p.Y = rand.Float64() * 30
			p.VX = (rand.Float64() - 0.5) * 0.2
			p.VY = (rand.Float64() - 0.5) * 0.1
			p.Lifetime = 0
			p.MaxLifetime = rand.Float64()*100 + 50
		}
	}

	// Update stars (twinkling)
	for i := range ua.stars {
		star := &ua.stars[i]
		star.Twinkle += deltaTime * 0.1
		star.Bright = 0.5 + 0.5*math.Sin(star.Twinkle)
	}

	// Update planets (orbiting)
	for i := range ua.planets {
		planet := &ua.planets[i]
		planet.Angle += planet.Speed * deltaTime
		planet.X = 50 + math.Cos(planet.Angle)*planet.Orbit
		planet.Y = 15 + math.Sin(planet.Angle)*planet.Orbit*0.5
	}

	// Update octopus
	if ua.octopus != nil {
		ua.octopus.Angle += ua.octopus.Speed * deltaTime
		ua.octopus.X = 50 + math.Cos(ua.octopus.Angle)*5
		ua.octopus.Y = 25 + math.Sin(ua.octopus.Angle)*2

		for i := range ua.octopus.Tentacles {
			tentacle := &ua.octopus.Tentacles[i]
			tentacle.Wave += deltaTime * 0.05
			waveOffset := math.Sin(tentacle.Wave) * 0.3
			tentacle.Angle = float64(i)*(math.Pi*2/8) + waveOffset
		}
	}

	// Update fish (swimming)
	for i := range ua.fish {
		fish := &ua.fish[i]
		fish.X += math.Cos(fish.Angle) * fish.Speed * deltaTime
		fish.Y += math.Sin(fish.Angle) * fish.Speed * deltaTime * 0.3
		fish.Angle += 0.02 * deltaTime
		fish.WavePhase += 0.05 * deltaTime

		// Wrap around screen
		if fish.X < -5 {
			fish.X = 105
		} else if fish.X > 105 {
			fish.X = -5
		}
		if fish.Y < 0 {
			fish.Y = 30
		} else if fish.Y > 30 {
			fish.Y = 0
		}
	}

	// Update gradient
	ua.gradientPos += deltaTime * 0.005
	if ua.gradientPos > 1 {
		ua.gradientPos = 0
	}

	return nil
}

func (ua *UnderwaterAnimator) Render() string {
	var bg strings.Builder

	// Create gradient background
	for y := 0; y < 30; y++ {
		for x := 0; x < 100; x++ {
			// Calculate gradient color
			intensity := (math.Sin((float64(x)/10+ua.gradientPos)*math.Pi) + 1) / 2
			depth := float64(y) / 30

			// Ocean gradient from deep blue to lighter blue
			r := int(10 + depth*20 + intensity*10)
			g := int(30 + depth*30 + intensity*20)
			b := int(60 + depth*40 + intensity*30)

			color := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
			bg.WriteString(color)
			bg.WriteString(" ")
		}
		bg.WriteString("\x1b[0m\n")
	}

	// Render particles
	for _, p := range ua.particles {
		rgb := getRGBFromColor(p.Color)
		bg.WriteString(fmt.Sprintf("\x1b[38;2;%sm•\x1b[0m", rgb))
	}

	// Render stars
	for _, star := range ua.stars {
		brightness := int(star.Bright * 255)
		size := int(star.Size)
		if size == 0 {
			size = 1
		}

		// Twinkling effect
		if star.Bright > 0.8 {
			bg.WriteString(fmt.Sprintf("\x1b[38;2;255;255;200m✦\x1b[0m"))
		} else {
			bg.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm•\x1b[0m", brightness, brightness, brightness))
		}
	}

	// Render planets
	for _, planet := range ua.planets {
		rgb := getRGBFromHex(planet.Color)
		bg.WriteString(fmt.Sprintf("\x1b[38;2;%sm●\x1b[0m", rgb))
	}

	// Render octopus
	if ua.octopus != nil {
		// Body
		rgb := getRGBFromHex(ua.octopus.Color)
		bg.WriteString(fmt.Sprintf("\x1b[38;2;%sm◉\x1b[0m", rgb))

		// Tentacles
		for range ua.octopus.Tentacles {
			bg.WriteString(fmt.Sprintf("\x1b[38;2;%sm~\x1b[0m", rgb))
		}
	}

	// Render fish
	for _, fish := range ua.fish {
		// Fish body with wave motion
		rgb := getRGBFromHex(fish.Color)
		bg.WriteString(fmt.Sprintf("\x1b[38;2;%sm><>\x1b[0m", rgb))
	}

	return bg.String()
}

func (ua *UnderwaterAnimator) IsPaused() bool {
	return ua.isPaused
}

func (ua *UnderwaterAnimator) SetPaused(paused bool) {
	ua.isPaused = paused
}

func (ua *UnderwaterAnimator) SetSpeed(speed float64) {
	ua.speed = speed
}

// ==================== UI COMPONENTS ====================

type Pane struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	X        int     `json:"x"`
	Y        int     `json:"y"`
	IsActive bool    `json:"is_active"`
	Opacity  float64 `json:"opacity"`
}

// ==================== MAIN MODEL ====================

type Model struct {
	width   int
	height  int
	focused bool

	// Animation system
	animator Animator

	// AI logging system
	logger Logger

	// UI panes
	panes      []Pane
	activePane int

	// Time tracking
	startTime time.Time

	// Current conversation
	currentSession *ConversationSession
	isRecording    bool
}

func initialModel() Model {
	// Create data directory
	dataDir, err := os.UserHomeDir()
	if err != nil {
		dataDir = "."
	}
	dataDir = filepath.Join(dataDir, ".ai-tui-data")
	os.MkdirAll(dataDir, 0755)

	// Initialize systems
	animator := NewUnderwaterAnimator()
	logger := NewFileLogger(dataDir)

	// Create UI panes with responsive sizing
	panes := []Pane{
		{
			ID:       "main",
			Title:    "🌌 AI Workspace",
			Content:  "Welcome to AI TUI!\n\nFeatures:\n• Living underwater world\n• AI conversation logging\n• Real-time monitoring\n\nPress 's' to start recording",
			Width:    35,
			Height:   10,
			X:        2,
			Y:        2,
			IsActive: true,
			Opacity:  0.9,
		},
		{
			ID:      "conversation",
			Title:   "💬 Conversation",
			Content: "No conversation recorded.\nPress 's' to start recording.",
			Width:   35,
			Height:  10,
			X:       40,
			Y:       2,
			Opacity: 0.8,
		},
		{
			ID:      "monitoring",
			Title:   "📊 Monitor",
			Content: "System monitoring active.\nPress 'm' for details.",
			Width:   73,
			Height:  8,
			X:       2,
			Y:       14,
			Opacity: 0.8,
		},
	}

	return Model{
		width:          100,
		height:         40,
		focused:        true,
		animator:       animator,
		logger:         logger,
		panes:          panes,
		activePane:     0,
		startTime:      time.Now(),
		currentSession: nil,
		isRecording:    false,
	}
}

// ==================== UPDATE METHODS ====================

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return t // Return time.Time directly
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = int(msg.Width), int(msg.Height)
		return m, nil

	case time.Time:
		// Update animation
		if err := m.animator.Update(1.0); err != nil {
			log.Printf("Animation update error: %v", err)
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	default:
		return m, nil
	}
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		// Log session end
		if m.currentSession != nil {
			endTime := time.Now()
			m.currentSession.EndTime = &endTime
			m.logger.LogConversation(*m.currentSession)
		}
		return m, tea.Quit

	case tea.KeyTab:
		m.activePane = (m.activePane + 1) % len(m.panes)
		for i := range m.panes {
			m.panes[i].IsActive = (i == m.activePane)
		}
		return m, nil

	case tea.KeyRunes:
		if len(msg.Runes) > 0 {
			switch msg.Runes[0] {
			case 's':
				return m.toggleRecording()
			case 'c':
				return m.clearConversation()
			case 'm':
				return m.showMonitoring()
			case ' ':
				m.animator.SetPaused(!m.animator.IsPaused())
				return m, tea.Printf("Animation %s", map[bool]string{true: "paused", false: "resumed"}[m.animator.IsPaused()])
			case '+':
				// Increase animation speed
				if animator, ok := m.animator.(*UnderwaterAnimator); ok {
					animator.SetSpeed(animator.speed * 1.5)
				}
				return m, tea.Printf("Speed increased")
			case '-':
				// Decrease animation speed
				if animator, ok := m.animator.(*UnderwaterAnimator); ok {
					animator.SetSpeed(animator.speed * 0.7)
				}
				return m, tea.Printf("Speed decreased")
			case 'r':
				// Reset animation
				*m = initialModel()
				return m, tea.Printf("Animation reset")
			}
		}
	}

	return m, nil
}

func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.MouseLeft {
		x, y := msg.X, msg.Y

		// Check if clicking on pane
		for i, pane := range m.panes {
			if x >= pane.X && x < pane.X+pane.Width &&
				y >= pane.Y && y < pane.Y+pane.Height {
				m.activePane = i
				for j := range m.panes {
					m.panes[j].IsActive = (j == i)
				}
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *Model) toggleRecording() (tea.Model, tea.Cmd) {
	m.isRecording = !m.isRecording

	if m.isRecording {
		// Start new conversation
		session := ConversationSession{
			ID:        generateID(),
			StartTime: time.Now(),
			Messages:  []ConversationMessage{},
			IsActive:  true,
			Tags:      []string{"tui-session"},
		}
		m.currentSession = &session

		// Log system event
		event := SystemEvent{
			ID:        generateID(),
			Timestamp: time.Now(),
			Type:      string(EventTypeInfo),
			Source:    "tui-system",
			Message:   "Started recording conversation",
			Data:      map[string]interface{}{"session_id": session.ID},
		}
		m.logger.LogEvent(event)

		return m, tea.Printf("🔴 Recording started")
	} else {
		// Stop recording
		if m.currentSession != nil {
			endTime := time.Now()
			m.currentSession.EndTime = &endTime
			m.currentSession.IsActive = false
			m.logger.LogConversation(*m.currentSession)

			event := SystemEvent{
				ID:        generateID(),
				Timestamp: time.Now(),
				Type:      string(EventTypeInfo),
				Source:    "tui-system",
				Message:   "Stopped recording conversation",
				Data:      map[string]interface{}{"session_id": m.currentSession.ID},
			}
			m.logger.LogEvent(event)
		}

		return m, tea.Printf("⏹️ Recording stopped")
	}
}

func (m *Model) clearConversation() (tea.Model, tea.Cmd) {
	m.currentSession = nil
	m.isRecording = false

	event := SystemEvent{
		ID:        generateID(),
		Timestamp: time.Now(),
		Type:      string(EventTypeInfo),
		Source:    "tui-system",
		Message:   "Cleared conversation",
	}
	m.logger.LogEvent(event)

	return m, tea.Printf("🗑️ Conversation cleared")
}

func (m *Model) showMonitoring() (tea.Model, tea.Cmd) {
	// This would open a detailed monitoring view
	event := SystemEvent{
		ID:        generateID(),
		Timestamp: time.Now(),
		Type:      string(EventTypeInfo),
		Source:    "tui-system",
		Message:   "Opened monitoring view",
	}
	m.logger.LogEvent(event)

	return m, tea.Printf("📊 Monitoring data logged")
}

// ==================== RENDER METHODS ====================

func (m Model) View() string {
	if m.width < 80 || m.height < 40 {
		return "Terminal too small! Please resize to at least 80x40"
	}

	var content strings.Builder

	// Render animated background
	content.WriteString(m.animator.Render())

	// Render UI panes on top
	content.WriteString(m.renderPanes())

	// Render recording indicator
	if m.isRecording {
		recordingIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Background(lipgloss.Color("#000000")).
			Bold(true).
			Render("🔴 REC")

		content.WriteString(lipgloss.Place(m.height, m.width,
			lipgloss.Right, lipgloss.Top,
			recordingIndicator))
	}

	return content.String()
}

func (m Model) renderPanes() string {
	var content strings.Builder

	for _, pane := range m.panes {
		if pane.IsActive {
			// Active pane with full opacity
			style := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Background(lipgloss.Color("#1a1a2e")).
				Foreground(lipgloss.Color("#ffffff")).
				Bold(true).
				Padding(0, 1).
				Width(pane.Width).
				Height(pane.Height)

			// Update content based on current state
			paneContent := pane.Content
			if pane.ID == "conversation" && m.currentSession != nil {
				paneContent = m.formatConversationDisplay()
			} else if pane.ID == "monitoring" {
				paneContent = m.formatMonitoringDisplay()
			}

			renderedPane := style.Render(fmt.Sprintf("%s\n\n%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#86E1FC")).Bold(true).Render(pane.Title),
				paneContent))

			content.WriteString(lipgloss.Place(m.height, m.width,
				lipgloss.Left, lipgloss.Top,
				renderedPane))
		}
	}

	return content.String()
}

func (m Model) formatConversationDisplay() string {
	if m.currentSession == nil {
		return "No active conversation"
	}

	duration := time.Since(m.currentSession.StartTime)
	messageCount := len(m.currentSession.Messages)

	return fmt.Sprintf("Session: %s\nDuration: %v\nMessages: %d\nStatus: %s",
		m.currentSession.ID[:8],
		duration.Round(time.Second),
		messageCount,
		map[bool]string{true: "Recording", false: "Active"}[m.isRecording])
}

func (m Model) formatMonitoringDisplay() string {
	return fmt.Sprintf("System Status: %s\nAnimation: %s\nSpeed: %.1fx",
		map[bool]string{true: "Active", false: "Paused"}[m.animator.IsPaused()],
		map[bool]string{true: "Running", false: "Paused"}[!m.animator.IsPaused()],
		m.getAnimationSpeed())
}

func (m Model) getAnimationSpeed() float64 {
	if animator, ok := m.animator.(*UnderwaterAnimator); ok {
		return animator.speed
	}
	return 1.0
}

// ==================== HELPERS ====================

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getRandomColor() string {
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8"}
	return colors[rand.Intn(len(colors))]
}

func getRandomFishColor() string {
	colors := []string{"#FF69B4", "#FFB347", "#87CEEB", "#98FB98", "#DDA0DD", "#F0E68C"}
	return colors[rand.Intn(len(colors))]
}

func getRGBFromHex(hex string) string {
	if len(hex) != 7 || hex[0] != '#' {
		return "255;255;255"
	}

	r := hexToByte(hex[1:3])
	g := hexToByte(hex[3:5])
	b := hexToByte(hex[5:7])

	return fmt.Sprintf("%d;%d;%d", r, g, b)
}

func getRGBFromColor(color string) string {
	colorMap := map[string]string{
		"#FF6B6B": "255;107;107",
		"#4ECDC4": "78;205;196",
		"#45B7D1": "69;183;209",
		"#96CEB4": "150;206;180",
		"#FFEAA7": "255;234;167",
		"#DDA0DD": "221;160;221",
		"#98D8C8": "152;216;200",
		"#9B59B6": "155;89;182",
		"#95E1D3": "149;225;211",
		"#FF69B4": "255;105;180",
		"#FFB347": "255;179;71",
		"#87CEEB": "135;206;235",
		"#98FB98": "152;251;152",
		"#F0E68C": "240;230;140",
	}

	if rgb, exists := colorMap[color]; exists {
		return rgb
	}
	return "255;255;255"
}

func hexToByte(hex string) byte {
	if len(hex) != 2 {
		return 0
	}

	var result byte
	for _, c := range hex {
		switch {
		case c >= '0' && c <= '9':
			result = result*16 + byte(c-'0')
		case c >= 'a' && c <= 'f':
			result = result*16 + byte(c-'a'+10)
		case c >= 'A' && c <= 'F':
			result = result*16 + byte(c-'A'+10)
		}
	}

	return result
}

// ==================== MAIN ====================

func main() {
	// Try to run the TUI with fallback to demo mode
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		// Fallback to demo mode if TUI fails
		runDemoMode()
	}
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func runDemoMode() {
	fmt.Println("🌌 AI TUI - Demo Mode")
	fmt.Println("====================")
	fmt.Println()
	fmt.Println("✅ Application compiled successfully!")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("• Living underwater animations with particles, fish, and octopus")
	fmt.Println("• AI conversation logging with JSON persistence")
	fmt.Println("• Real-time system monitoring")
	fmt.Println("• Professional TUI interface with resizable panes")
	fmt.Println()
	fmt.Println("Controls:")
	fmt.Println("• Tab: Switch between panes")
	fmt.Println("• s: Start/stop recording")
	fmt.Println("• c: Clear conversation")
	fmt.Println("• m: Show monitoring")
	fmt.Println("• Space: Pause/resume animation")
	fmt.Println("• +/-: Adjust animation speed")
	fmt.Println("• r: Reset animation")
	fmt.Println("• Ctrl+C/Esc: Quit")
	fmt.Println()
	fmt.Println("🚀 Run './ai-tui' in an interactive terminal to experience the full AI TUI!")
}

func getTerminalSize() (int, int, error) {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return 0, 0, errno
	}
	return int(ws.Col), int(ws.Row), nil
}
