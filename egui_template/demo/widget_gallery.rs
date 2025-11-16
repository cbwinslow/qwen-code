/// Shows off one example of each major type of widget.
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize))]
#[cfg_attr(feature = "serde", serde(default))]
pub struct WidgetGallery {
    enabled: bool,
    visible: bool,
    boolean: bool,
    opacity: f32,
    scalar: f32,
    string: String,
    color: egui::Color32,
    animate_progress_bar: bool,
}

impl Default for WidgetGallery {
    fn default() -> Self {
        Self {
            enabled: true,
            visible: true,
            opacity: 1.0,
            boolean: false,
            scalar: 42.0,
            string: Default::default(),
            color: egui::Color32::LIGHT_BLUE.linear_multiply(0.5),
            animate_progress_bar: false,
        }
    }
}

impl crate::Demo for WidgetGallery {
    fn name(&self) -> &'static str {
        "🗄 Widget Gallery"
    }

    fn show(&mut self, ctx: &egui::Context, open: &mut bool) {
        egui::Window::new(self.name())
            .open(open)
            .resizable([true, false])
            .default_width(280.0)
            .show(ctx, |ui| {
                use crate::View as _;
                self.ui(ui);
            });
    }
}

impl crate::View for WidgetGallery {
    fn ui(&mut self, ui: &mut egui::Ui) {
        let mut ui_builder = egui::UiBuilder::new();
        if !self.enabled {
            ui_builder = ui_builder.disabled();
        }
        if !self.visible {
            ui_builder = ui_builder.invisible();
        }

        ui.scope_builder(ui_builder, |ui| {
            ui.multiply_opacity(self.opacity);

            egui::Grid::new("my_grid")
                .num_columns(2)
                .spacing([40.0, 4.0])
                .striped(true)
                .show(ui, |ui| {
                    self.gallery_grid_contents(ui);
                });
        });

        ui.separator();

        ui.horizontal(|ui| {
            ui.checkbox(&mut self.visible, "Visible")
                .on_hover_text("Uncheck to hide all widgets.");
            if self.visible {
                ui.checkbox(&mut self.enabled, "Interactive")
                    .on_hover_text("Uncheck to inspect how widgets look when disabled.");
                (ui.add(
                    egui::DragValue::new(&mut self.opacity)
                        .speed(0.01)
                        .range(0.0..=1.0),
                ) | ui.label("Opacity"))
                .on_hover_text("Reduce this value to make widgets semi-transparent");
            }
        });
    }
}

impl WidgetGallery {
    fn gallery_grid_contents(&mut self, ui: &mut egui::Ui) {
        let Self {
            enabled: _,
            visible: _,
            opacity: _,
            boolean,
            scalar,
            string,
            color,
            animate_progress_bar,
        } = self;

        ui.label("Label");
        ui.label("Welcome to the widget gallery!");
        ui.end_row();

        ui.label("Button");
        if ui.button("Click me!").clicked() {
            *boolean = !*boolean;
        }
        ui.end_row();

        ui.label("Checkbox");
        ui.checkbox(boolean, "Checkbox");
        ui.end_row();

        ui.label("Slider");
        ui.add(egui::Slider::new(scalar, 0.0..=360.0).suffix("°"));
        ui.end_row();

        ui.label("DragValue");
        ui.add(egui::DragValue::new(scalar).speed(1.0));
        ui.end_row();

        ui.label("ProgressBar");
        let progress = *scalar / 360.0;
        let progress_bar = egui::ProgressBar::new(progress)
            .show_percentage()
            .animate(*animate_progress_bar);
        *animate_progress_bar = ui
            .add(progress_bar)
            .on_hover_text("The progress bar can be animated!")
            .hovered();
        ui.end_row();

        ui.label("Color picker");
        ui.color_edit_button_srgba(color);
        ui.end_row();

        ui.label("Separator");
        ui.separator();
        ui.end_row();
    }
}