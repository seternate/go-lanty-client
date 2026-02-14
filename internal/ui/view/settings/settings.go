package settingsview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	model "github.com/seternate/go-lanty-client/internal/ui/model/settings"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type SettingsView struct {
	widget.BaseWidget

	settings    *model.SettingsModel
	tiles       []*SettingTile
	saveButton  *widget.Button
	resetButton *widget.Button

	OnSavePressed  func()
	OnResetPressed func()
}

func NewSettingsView(window fyne.Window, settings *model.SettingsModel) *SettingsView {
	usernameTile := NewSettingTile(
		"Username",
		"Enter your username",
		settings.Username,
		window,
		false,
	)

	serverURLTile := NewSettingTile(
		"Server URL",
		"Server address (e.g., http://localhost:8080)",
		settings.ServerURL,
		window,
		false,
	)

	gameDirectoryTile := NewSettingTile(
		"Game Directory",
		"Directory where games are installed",
		settings.GameDirectory,
		window,
		true,
	)

	view := &SettingsView{
		settings:    settings,
		tiles:       []*SettingTile{usernameTile, serverURLTile, gameDirectoryTile},
		resetButton: widget.NewButton("Reset", nil),
		saveButton:  widget.NewButton("Save", nil),
	}
	view.saveButton.Importance = widget.HighImportance

	view.saveButton.OnTapped = func() {
		if view.OnSavePressed != nil {
			view.OnSavePressed()
			view.Refresh()
		}
	}
	view.resetButton.OnTapped = func() {
		if view.OnResetPressed != nil {
			view.OnResetPressed()
			view.Refresh()
		}
	}

	settings.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)
	return view
}

func (view *SettingsView) Refresh() {
	if view.settings.IsDirty() {
		view.saveButton.Enable()
		view.resetButton.Enable()
	} else {
		view.saveButton.Disable()
		view.resetButton.Disable()
	}
	view.BaseWidget.Refresh()
}

func (view *SettingsView) CreateRenderer() fyne.WidgetRenderer {
	return newSettingsRenderer(view)
}

type settingsRenderer struct {
	view            *SettingsView
	headerContainer *fyne.Container
	headerText      *canvas.Text
	scrollContainer *container.Scroll
	formContainer   *fyne.Container
	buttonContainer *fyne.Container
	objects         []fyne.CanvasObject
}

func newSettingsRenderer(view *SettingsView) *settingsRenderer {
	headerText := canvas.NewText("Settings", color.White)
	headerText.TextSize = 28
	headerText.TextStyle = fyne.TextStyle{Bold: true}

	headerPadding := theme.InnerPadding() * 2
	paddingTop := canvas.NewRectangle(color.Transparent)
	paddingTop.SetMinSize(fyne.NewSize(0, headerPadding))
	paddingBottom := canvas.NewRectangle(color.Transparent)
	paddingBottom.SetMinSize(fyne.NewSize(0, headerPadding))

	paddingLeft := canvas.NewRectangle(color.Transparent)
	paddingLeft.SetMinSize(fyne.NewSize(headerPadding, 0))
	paddingRight := canvas.NewRectangle(color.Transparent)
	paddingRight.SetMinSize(fyne.NewSize(headerPadding, 0))

	headerLeft := container.NewBorder(paddingTop, paddingBottom, paddingLeft, paddingRight, container.NewWithoutLayout(headerText))
	headerContainer := container.NewBorder(nil, nil, headerLeft, nil)

	var tileObjects []fyne.CanvasObject
	for _, tile := range view.tiles {
		tileObjects = append(tileObjects, tile)
	}

	formContainer := container.NewVBox(tileObjects...)

	scrollContainer := container.NewScroll(formContainer)
	scrollContainer.SetMinSize(fyne.NewSize(0, 0))

	buttonGrid := container.NewGridWithColumns(2, view.resetButton, view.saveButton)
	buttonRow := container.NewHBox(layout.NewSpacer(), buttonGrid)
	buttonContainer := container.NewPadded(buttonRow)

	mainContainer := container.NewBorder(headerContainer, buttonContainer, nil, nil, scrollContainer)

	renderer := &settingsRenderer{
		view:            view,
		headerContainer: headerContainer,
		headerText:      headerText,
		scrollContainer: scrollContainer,
		formContainer:   formContainer,
		buttonContainer: buttonContainer,
		objects: []fyne.CanvasObject{
			mainContainer,
		},
	}

	return renderer
}

func (r *settingsRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *settingsRenderer) Layout(size fyne.Size) {
	if len(r.objects) > 0 {
		r.objects[0].Resize(size)
		r.objects[0].Move(fyne.NewPos(0, 0))
	}
}

func (r *settingsRenderer) MinSize() fyne.Size {
	padding := theme.InnerPadding() * 2
	headerHeight := r.headerText.MinSize().Height + padding*2
	buttonHeight := r.buttonContainer.MinSize().Height
	return fyne.NewSize(0, headerHeight+buttonHeight)
}

func (r *settingsRenderer) Refresh() {
	r.headerText.Refresh()

	for _, tile := range r.view.tiles {
		tile.Refresh()
	}

	r.buttonContainer.Refresh()
}

func (r *settingsRenderer) Destroy() {}
