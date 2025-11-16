# 🚀 Multi-Pane TUI System

A feature-rich terminal user interface built with Go and Bubble Tea, featuring resizable panes, progress tracking, and secure secret management.

## ✨ Features

### 🪟 **Resizable Panes**

- **Drag-to-resize** borders for intuitive layout adjustment
- **Multiple panes** with independent content
- **Active pane highlighting** with visual feedback
- **Smart positioning** system

### 📊 **Progress Tracking**

- **Real-time progress bars** with automatic updates
- **Color-coded status** indicators
- **Percentage display** with smooth animations
- **Multiple concurrent** progress tracking

### 🔐 **Secret Management**

- **Secure local storage** in JSON format
- **Add/Edit/Delete** secret operations
- **Tag-based organization** system
- **Timestamp tracking** for audit trail

## 🎮 Controls

### Navigation

- `Tab` - Cycle through panes
- `Shift+Tab` - Reverse cycle through panes
- `Click` - Focus pane or resize border
- `Drag` - Resize panes from borders

### View Modes

- `s` - Switch to Secrets view
- `p` - Switch to Progress view
- `Esc` or `Ctrl+C` - Quit application

### Secret Management

- `n` - Add new secret
- `Enter` - Save secret
- `q` - Cancel editing
- `r` - Refresh secrets list

## 🏗️ Architecture

### Core Components

```go
type Model struct {
    panes       []Pane      // Resizable window panes
    secrets     []Secret     // Encrypted storage
    progress    []Progress   // Task tracking
    activePane  int         // Currently focused pane
    focusedPane int         // Keyboard navigation
    // ... additional state
}
```

### Pane System

```go
type Pane struct {
    ID          string  // Unique identifier
    Title       string  // Display title
    Content     string  // Pane content
    Width       int     // Current width
    Height      int     // Current height
    X           int     // X position
    Y           int     // Y position
    IsActive    bool    // Focus state
    IsResizable bool    // Resize capability
}
```

### Progress Tracking

```go
type Progress struct {
    ID          string  // Task identifier
    Label       string  // Display label
    Percent     float64 // Completion percentage
    IsActive    bool    // Running state
    ShowPercent bool    // Display preference
    Color       string  // Visual theme
}
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Terminal with mouse support
- 80x24 minimum terminal size

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd go-tui-app

# Install dependencies
go mod tidy

# Run the application
go run main.go
```

### Build

```bash
# Build for current platform
go build -o tui-app main.go

# Build for multiple platforms
GOOS=linux go build -o tui-app-linux main.go
GOOS=windows go build -o tui-app.exe main.go
GOOS=darwin go build -o tui-app-macos main.go
```

## 📁 File Structure

```
go-tui-app/
├── main.go              # Main application entry point
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── .go-tui-secrets.json # Encrypted secrets storage
└── README.md            # This documentation
```

## 🎨 Customization

### Adding New Panes

```go
newPane := Pane{
    ID:          "custom",
    Title:       "Custom Pane",
    Content:     "Your content here",
    Width:       50,
    Height:      15,
    X:           10,
    Y:           10,
    IsActive:    false,
    IsResizable: true,
}

model.panes = append(model.panes, newPane)
```

### Custom Progress Bars

```go
newProgress := Progress{
    ID:          "custom-task",
    Label:       "Custom Task",
    Percent:     0.0,
    IsActive:    true,
    ShowPercent: true,
    Color:       "purple",
}

model.progress = append(model.progress, newProgress)
```

### Theme Customization

```go
// Modify color schemes in the styles section
var (
    titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF79C6"))
    activeStyle   = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Foreground(lipgloss.Color("#86E1FC"))
    // ... more styles
)
```

## 🔧 Configuration

### Environment Variables

- `TUI_SECRET_FILE` - Custom secrets file path
- `TUI_CONFIG_DIR` - Configuration directory
- `TUI_THEME` - Color theme (light/dark)

### Secret Storage Format

```json
[
  {
    "id": "unique-id",
    "name": "API Key",
    "value": "secret-value",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "tags": ["api", "production"]
  }
]
```

## 🐛 Troubleshooting

### Common Issues

**Terminal too small**

```
Error: Terminal too small! Please resize to at least 80x24
Solution: Increase terminal window size
```

**Mouse not working**

```
Solution: Ensure terminal supports mouse events
- For tmux: `setw -g mouse on`
- For screen: `termcapinfo`
```

**Colors not displaying**

```
Solution: Check terminal color support
export TERM=xterm-256color
```

### Debug Mode

Run with verbose logging:

```bash
DEBUG=1 go run main.go
```

## 🔄 Development

### Adding Features

1. **Update Model struct** - Add new state fields
2. **Implement Update logic** - Handle new message types
3. **Add View rendering** - Create display functions
4. **Update controls** - Add keyboard shortcuts

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Benchmark
go test -bench=. ./...
```

## 📚 Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [lipgloss/table](https://github.com/charmbracelet/lipgloss/table) - Table rendering
- [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) - Terminal utilities

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make your changes
4. Add tests
5. Submit pull request

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- Built with [Charm Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Styled with [Charm Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Inspired by [egui](https://github.com/emilk/egui) widget patterns

---

**Built with ❤️ using Go and Bubble Tea**
