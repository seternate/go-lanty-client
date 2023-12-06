package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/setting"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type Sidebar struct {
	widget.BaseWidget

	text    string
	buttons map[string]*widget.Button

	OnGamesTapped     func()
	OnDownloadsTapped func()
	OnUsersTapped     func()
	OnSettingsTapped  func()
}

func NewSidebar(text string) *Sidebar {
	sidebar := &Sidebar{
		text: text,
		buttons: map[string]*widget.Button{
			"games":     widget.NewButtonWithIcon("Games", fynetheme.MediaPlayIcon(), nil),
			"downloads": widget.NewButtonWithIcon("Downloads", fynetheme.DownloadIcon(), nil),
			"users":     widget.NewButtonWithIcon("Users", fynetheme.AccountIcon(), nil),
			"settings":  widget.NewButtonWithIcon("Settings", fynetheme.SettingsIcon(), nil),
		},
	}
	sidebar.ExtendBaseWidget(sidebar)

	sidebar.buttons["games"].OnTapped = func() {
		if sidebar.OnGamesTapped != nil {
			sidebar.OnGamesTapped()
		}
	}
	sidebar.buttons["downloads"].OnTapped = func() {
		if sidebar.OnDownloadsTapped != nil {
			sidebar.OnDownloadsTapped()
		}
	}
	sidebar.buttons["users"].OnTapped = func() {
		if sidebar.OnUsersTapped != nil {
			sidebar.OnUsersTapped()
		}
	}
	sidebar.buttons["settings"].OnTapped = func() {
		if sidebar.OnSettingsTapped != nil {
			sidebar.OnSettingsTapped()
		}
	}
	for _, button := range sidebar.buttons {
		button.Alignment = widget.ButtonAlignLeading
	}

	return sidebar
}

func (widget *Sidebar) CreateRenderer() fyne.WidgetRenderer {
	return newSidebarRenderer(widget)
}

type sidebarRenderer struct {
	widget     *Sidebar
	background *canvas.Rectangle
	text       *canvas.Text
	version    *canvas.Text
	objects    []fyne.CanvasObject
}

func newSidebarRenderer(widget *Sidebar) *sidebarRenderer {
	renderer := &sidebarRenderer{
		widget:     widget,
		background: canvas.NewRectangle(fynetheme.InputBorderColor()),
		text:       canvas.NewText(widget.text, theme.ForegroundColor()),
		version:    canvas.NewText(setting.VERSION, theme.ForegroundColor()),
	}
	renderer.objects = []fyne.CanvasObject{
		renderer.background,
		renderer.text,
		renderer.version,
	}
	for _, button := range renderer.widget.buttons {
		renderer.objects = append(renderer.objects, button)
	}
	renderer.text.TextSize = 32
	renderer.text.TextStyle.Bold = true

	return renderer
}

func (renderer *sidebarRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *sidebarRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	textsize := fyne.MeasureText(renderer.text.Text, renderer.text.TextSize, renderer.text.TextStyle)
	renderer.text.Move(fyne.NewPos((size.Width-textsize.Width)/2, theme.InnerPadding()))

	index := 0
	for _, button := range []*widget.Button{
		renderer.widget.buttons["games"],
		renderer.widget.buttons["downloads"],
		renderer.widget.buttons["users"],
		renderer.widget.buttons["settings"],
	} {
		button.Resize(fyne.NewSize(size.Width-2*theme.InnerPadding(), button.MinSize().Height))
		button.Move(fyne.NewPos(theme.InnerPadding(), textsize.Height*1.5+(button.MinSize().Height+theme.InnerPadding())*float32(index)))
		index++
	}

	textsize = fyne.MeasureText(renderer.version.Text, renderer.version.TextSize, renderer.version.TextStyle)
	renderer.version.Move(fyne.NewPos(theme.InnerPadding(), size.Height-theme.InnerPadding()-textsize.Height))
}

func (renderer *sidebarRenderer) MinSize() fyne.Size {
	textsize := fyne.MeasureText(renderer.text.Text, renderer.text.TextSize, renderer.text.TextStyle)
	minHeight := textsize.Height*1.5 + theme.InnerPadding()
	minWidth := textsize.Width + 2*theme.InnerPadding()

	for _, button := range renderer.widget.buttons {
		minHeight += button.MinSize().Height + theme.InnerPadding()
		minWidth = fyne.Max(minWidth, button.MinSize().Width+2*theme.InnerPadding())
	}
	textsize = fyne.MeasureText(renderer.version.Text, renderer.version.TextSize, renderer.version.TextStyle)
	minHeight += textsize.Height + theme.InnerPadding()

	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *sidebarRenderer) Refresh() {
	renderer.text.Text = renderer.widget.text
	renderer.background.Refresh()
	renderer.text.Refresh()
	renderer.version.Refresh()
	for _, button := range renderer.widget.buttons {
		button.Refresh()
	}
}

func (renderer *sidebarRenderer) Destroy() {}
