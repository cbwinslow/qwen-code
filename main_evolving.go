package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== ANIMATION SYSTEM ====================

type AnimationState struct {
	Time      float64
	Frame     int
	Speed     float64
	Direction int // 1 or -1
}

type Particle struct {
	X           float64
	Y           float64
	VX          float64
	VY          float64
	Size        float64
	Color       string
	Opacity     float64
	Lifetime    float64
	MaxLifetime float64
}

type Star struct {
	X       float64
	Y       float64
	Size    float64
	Bright  float64
	Twinkle float64
}

type Planet struct {
	X     float64
	Y     float64
	Orbit float64
	Size  float64
	Color string
	Speed float64
	Angle float64
}

type Octopus struct {
	X         float64
	Y         float64
	Angle     float64
	Tentacles []Tentacle
	Color     string
	Speed     float64
}

type Tentacle struct {
	Angle  float64
	Length float64
	Wave   float64
}

type Fish struct {
	X         float64
	Y         float64
	Angle     float64
	Speed     float64
	Size      float64
	Color     string
	WavePhase float64
}

// ==================== MODEL ====================

type Model struct {
	width   int
	height  int
	focused bool

	// Animation state
	anim      AnimationState
	particles []Particle
	stars     []Star
	planets   []Planet
	octopus   *Octopus
	fish      []Fish

	// Background gradient
	gradientPos float64

	// UI panes
	panes      []Pane
	activePane int

	// Time tracking
	startTime time.Time
}

type Pane struct {
	ID       string
	Title    string
	Content  string
	Width    int
	Height   int
	X        int
	Y        int
	IsActive bool
	Opacity  float64
}

// ==================== INITIALIZE ====================

func initialModel() Model {
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
		{
			X: 20, Y: 10, Orbit: 15, Size: 2,
			Color: "#FF6B6B", Speed: 0.02, Angle: 0,
		},
		{
			X: 70, Y: 20, Orbit: 8, Size: 1.5,
			Color: "#4ECDC4", Speed: 0.03, Angle: math.Pi,
		},
		{
			X: 50, Y: 5, Orbit: 12, Size: 1,
			Color: "#95E1D3", Speed: 0.015, Angle: math.Pi / 2,
		},
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

	// Create UI panes
	panes := []Pane{
		{
			ID:       "main",
			Title:    "🌌 Living Workspace",
			Content:  "Welcome to the evolving TUI!\n\nWatch the underwater world come alive\n\nFeatures:\n• Dynamic particle system\n• Orbiting planets\n• Swimming octopus\n• Moving fish schools\n• Twinkling stars\n\nPress 'q' to exit",
			Width:    40,
			Height:   12,
			X:        30,
			Y:        28,
			IsActive: true,
			Opacity:  0.9,
		},
		{
			ID:    "stats",
			Title: "📊 System Stats",
			Content: fmt.Sprintf("Particles: %d\nStars: %d\nPlanets: %d\nFish: %d\nOctopus: %s\n\nRuntime: %v",
				len(particles), len(stars), len(planets), len(fish), "Active",
				time.Since(m.startTime).Round(time.Second)),
			Width:    30,
			Height:   10,
			X:        75,
			Y:        28,
			IsActive: false,
			Opacity:  0.8,
		},
		{
			ID:       "controls",
			Title:    "🎮 Controls",
			Content:  "[Tab] Switch Pane\n[Space] Pause/Resume\n[+/-] Speed Control\n[R] Reset Animation\n[Q] Quit",
			Width:    25,
			Height:   8,
			X:        5,
			Y:        28,
			IsActive: false,
			Opacity:  0.8,
		},
	}

	return Model{
		width:       100,
		height:      40,
		focused:     true,
		anim:        AnimationState{Time: 0, Frame: 0, Speed: 1.0, Direction: 1},
		particles:   particles,
		stars:       stars,
		planets:     planets,
		octopus:     octopus,
		fish:        fish,
		gradientPos: 0,
		panes:       panes,
		activePane:  0,
		startTime:   time.Now(),
	}
}

// ==================== UPDATE ====================

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		tea.Tick(time.Second/60), // 60 FPS
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = int(msg.Width), int(msg.Height)
		return m, nil

	case tea.TickMsg:
		return m.updateAnimation(msg)

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	default:
		return m, nil
	}
}

func (m *Model) updateAnimation(msg tea.TickMsg) (tea.Model, tea.Cmd) {
	m.anim.Time += float64(msg) * m.anim.Speed

	// Update particles
	for i := range m.particles {
		p := &m.particles[i]
		p.X += p.VX
		p.Y += p.VY
		p.Lifetime++

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
	for i := range m.stars {
		star := &m.stars[i]
		star.Twinkle += 0.1
		star.Bright = 0.5 + 0.5*math.Sin(star.Twinkle)
	}

	// Update planets (orbiting)
	for i := range m.planets {
		planet := &m.planets[i]
		planet.Angle += planet.Speed
		planet.X = 50 + math.Cos(planet.Angle)*planet.Orbit
		planet.Y = 15 + math.Sin(planet.Angle)*planet.Orbit*0.5
	}

	// Update octopus
	if m.octopus != nil {
		m.octopus.Angle += m.octopus.Speed
		m.octopus.X = 50 + math.Cos(m.octopus.Angle)*5
		m.octopus.Y = 25 + math.Sin(m.octopus.Angle)*2

		for i := range m.octopus.Tentacles {
			tentacle := &m.octopus.Tentacles[i]
			tentacle.Wave += 0.05
			waveOffset := math.Sin(tentacle.Wave) * 0.3
			tentacle.Angle = float64(i)*(math.Pi*2/8) + waveOffset
		}
	}

	// Update fish (swimming)
	for i := range m.fish {
		fish := &m.fish[i]
		fish.X += math.Cos(fish.Angle) * fish.Speed
		fish.Y += math.Sin(fish.Angle) * fish.Speed * 0.3
		fish.Angle += 0.02
		fish.WavePhase += 0.05

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
	m.gradientPos += 0.005
	if m.gradientPos > 1 {
		m.gradientPos = 0
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc, tea.Key('q'):
		return m, tea.Quit

	case tea.KeyTab:
		m.activePane = (m.activePane + 1) % len(m.panes)
		for i := range m.panes {
			m.panes[i].IsActive = (i == m.activePane)
		}
		return m, nil

	case tea.Key(' '):
		m.anim.Speed *= 0.5
		if m.anim.Speed < 0.1 {
			m.anim.Speed = 0.1
		}
		return m, nil

	case tea.Key('+'):
		m.anim.Speed *= 2
		if m.anim.Speed > 5 {
			m.anim.Speed = 5
		}
		return m, nil

	case tea.Key('r'):
		*m = initialModel()
		return m, nil

	case tea.Key('p'):
		m.focused = !m.focused
		return m, nil
	}

	return m, nil
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
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

// ==================== RENDER ====================

func (m Model) View() string {
	if m.width < 80 || m.height < 40 {
		return "Terminal too small! Please resize to at least 80x40"
	}

	var content strings.Builder

	// Render animated background
	content.WriteString(m.renderBackground())

	// Render UI panes on top
	content.WriteString(m.renderPanes())

	return content.String()
}

func (m Model) renderBackground() string {
	var bg strings.Builder

	// Create gradient background
	for y := 0; y < m.height; y++ {
		var line strings.Builder

		for x := 0; x < m.width; x++ {
			// Calculate gradient color
			intensity := (math.Sin((float64(x)/10+m.gradientPos)*math.Pi) + 1) / 2
			depth := float64(y) / float64(m.height)

			// Ocean gradient from deep blue to lighter blue
			r := int(10 + depth*20 + intensity*10)
			g := int(30 + depth*30 + intensity*20)
			b := int(60 + depth*40 + intensity*30)

			// Check if there's a UI pane at this position
			hasPane := false
			for _, pane := range m.panes {
				if x >= pane.X && x < pane.X+pane.Width &&
					y >= pane.Y && y < pane.Y+pane.Height {
					hasPane = true
					break
				}
			}

			if hasPane {
				// Darker background under panes
				r = int(float64(r) * 0.3)
				g = int(float64(g) * 0.3)
				b = int(float64(b) * 0.3)
			}

			color := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
			line.WriteString(color)
			line.WriteString("•")
			line.WriteString("\x1b[0m")
		}

		bg.WriteString(line.String() + "\n")
	}

	return bg.String()
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
				Opacity(float64(pane.Opacity)).
				Padding(0, 1).
				Width(pane.Width).
				Height(pane.Height)

			paneContent := fmt.Sprintf("%s\n\n%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#86E1FC")).Bold(true).Render(pane.Title),
				pane.Content)

			renderedPane := style.Render(paneContent)
			content.WriteString(lipgloss.Place(m.height, m.width,
				lipgloss.Left, lipgloss.Top,
				renderedPane))
		}
	}

	// Runtime info
	runtime := time.Since(m.startTime)
	runtimeText := fmt.Sprintf("Runtime: %v | FPS: %.0f | Particles: %d",
		runtime.Round(time.Second),
		60.0*m.anim.Speed,
		len(m.particles))

	content.WriteString(lipgloss.Place(m.height, m.width,
		lipgloss.Left, lipgloss.Bottom,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#1a1a2e")).
			Padding(0, 1).
			Render(runtimeText)))

	return content.String()
}

// ==================== HELPERS ====================

func getRandomColor() string {
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8"}
	return colors[rand.Intn(len(colors))]
}

func getRandomFishColor() string {
	colors := []string{"#FF69B4", "#FFB347", "#87CEEB", "#98FB98", "#DDA0DD", "#F0E68C"}
	return colors[rand.Intn(len(colors))]
}

// ==================== MAIN ====================

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
