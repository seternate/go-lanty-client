package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
)

type DefaultUsernameBrowser struct {
	widget.BaseWidget

	controller *controller.Controller

	settingsbrowser *SettingsBrowser
}

func NewDefaultUsernameBrowser(controller *controller.Controller, window fyne.Window) (browser *DefaultUsernameBrowser) {
	browser = &DefaultUsernameBrowser{
		controller:      controller,
		settingsbrowser: NewSettingsBrowser(controller, window),
	}
	browser.ExtendBaseWidget(browser)

	return
}

func (widget *DefaultUsernameBrowser) SetOnSubmit(onSubmit func()) {
	widget.settingsbrowser.SetOnSubmit(onSubmit)
}

func (widget *DefaultUsernameBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newDefaulusernamebrowserRenderer(widget)
}

type defaultusernamebrowserRenderer struct {
	widget     *DefaultUsernameBrowser
	background *canvas.Rectangle
}

func newDefaulusernamebrowserRenderer(widget *DefaultUsernameBrowser) *defaultusernamebrowserRenderer {
	renderer := &defaultusernamebrowserRenderer{
		widget:     widget,
		background: canvas.NewRectangle(fynetheme.BackgroundColor()),
	}

	return renderer
}

func (renderer *defaultusernamebrowserRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		renderer.background,
		renderer.widget.settingsbrowser,
	}

	return objects
}

func (renderer *defaultusernamebrowserRenderer) Layout(size fyne.Size) {
	renderer.background.Move(fyne.NewPos(0, 0))
	renderer.background.Resize(size)
	renderer.widget.settingsbrowser.Move(fyne.NewPos(0, 0))
	renderer.widget.settingsbrowser.Resize(size)
}

func (renderer *defaultusernamebrowserRenderer) MinSize() fyne.Size {
	return renderer.widget.settingsbrowser.MinSize()
}

func (renderer *defaultusernamebrowserRenderer) Refresh() {
	renderer.background.Refresh()
	renderer.widget.settingsbrowser.Refresh()
}

func (renderer *defaultusernamebrowserRenderer) Destroy() {}
