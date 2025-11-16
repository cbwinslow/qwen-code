# egui Demo Template for Golang TUI

This directory contains the egui demo source code saved as a template for creating your own Golang TUI application.

## 📁 Files Saved

### Core Demo Files:

- **`mod.rs`** - Module definitions and traits for demo system
- **`widget_gallery.rs`** - Complete widget gallery with all major UI components
- **`toggle_switch.rs`** - Custom widget implementation tutorial

### Key Concepts Extracted:

#### 1. **Widget System Architecture**

```rust
pub trait View {
    fn ui(&mut self, ui: &mut egui::Ui);
}

pub trait Demo {
    fn is_enabled(&self, _ctx: &egui::Context) -> bool { true }
    fn name(&self) -> &'static str;
    fn show(&mut self, ctx: &egui::Context, open: &mut bool);
}
```

#### 2. **Custom Widget Creation**

The toggle switch demonstrates the 4-step widget creation process:

1. **Size Decision** - Calculate widget dimensions
2. **Space Allocation** - Reserve screen real estate
3. **Interaction Handling** - Process user input
4. **Painting** - Render the widget

#### 3. **UI Component Patterns**

- **Grid Layouts** - Organized component placement
- **Responsive Design** - Adaptive sizing and positioning
- **State Management** - Persistent widget state
- **Animation Support** - Smooth transitions and effects

#### 4. **Widget Gallery Components**

- Labels and text rendering
- Buttons and interactive elements
- Sliders and drag values
- Progress bars with animation
- Color pickers
- Checkboxes and radio buttons
- Separators and visual dividers

## 🚀 Adapting for Golang TUI

### Key Translation Concepts:

#### Rust → Go Mapping:

```rust
// Rust egui pattern
pub trait View {
    fn ui(&mut self, ui: &mut egui::Ui);
}
```

```go
// Go TUI equivalent
type View interface {
    Render(ui *tview.UI) tview.Primitive
}
```

#### Widget Creation Pattern:

```rust
// Rust: 4-step widget creation
let (rect, response) = ui.allocate_exact_size(desired_size, Sense::click());
if response.clicked() { /* handle */ }
if ui.is_rect_visible(rect) { /* paint */ }
```

```go
// Go: TUI widget pattern
func (w *Widget) Draw(screen tcell.Screen) {
    w.Box.DrawForSubscreen(screen, w.Box.GetRect())
    if w.HasFocus() { /* handle input */ }
}
```

### Recommended Go TUI Libraries:

1. **tview** - Feature-rich TUI framework
2. **tcell** - Terminal cell manipulation
3. **bubbletea** - Elm-inspired TUI framework
4. **lipgloss** - Styling and layout

## 📋 Implementation Steps

### 1. **Project Structure**

```
go-tui-app/
├── main.go
├── widgets/
│   ├── button.go
│   ├── toggle.go
│   └── gallery.go
├── views/
│   └── main_view.go
└── state/
    └── app_state.go
```

### 2. **Core Interfaces**

```go
type Widget interface {
    Render() tview.Primitive
    HandleInput(*tcell.EventKey) bool
    SetRect(x, y, width, height int)
    GetRect() (int, int, int, int)
}

type View interface {
    Name() string
    IsEnabled() bool
    Show(app *App, ctx *Context) error
}
```

### 3. **State Management**

```go
type AppState struct {
    CurrentView string
    Widgets map[string]Widget
    Focus    string
    Data     map[string]interface{}
}
```

## 🎯 Next Steps

1. **Choose Go TUI Library** - Based on your requirements
2. **Define Core Interfaces** - Adapt Rust traits to Go
3. **Implement Basic Widgets** - Start with buttons, labels
4. **Create Layout System** - Grid and flexbox equivalents
5. **Add Event Handling** - Keyboard and mouse support
6. **Build Demo App** - Recreate widget gallery in Go

## 📚 Resources

### Go TUI Libraries:

- [tview](https://github.com/rivo/tview) - Full-featured TUI framework
- [bubbletea](https://github.com/charmbracelet/bubbletea) - Modern alternative
- [tcell](https://github.com/gdamore/tcell) - Terminal abstraction
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

### egui Documentation:

- [egui Docs](https://docs.rs/egui/) - Original Rust documentation
- [egui Examples](https://github.com/emilk/egui/tree/master/examples) - More examples

---

**Template Created**: 2025-11-15  
**Source**: emilk/egui @ 01770be13ee4513a960a7db8118b6981e907eb64  
**Purpose**: Golang TUI Development Reference
