package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty-client/pkg/setting"
)

type Lanty struct {
	widget.BaseWidget

	controller      *controller.Controller
	sidebar         *Sidebar
	gamebrowser     *GameBrowser
	downloadbrowser *DownloadBrowser
	userbrowser     *UserBrowser
}

func NewLanty(controller *controller.Controller) *Lanty {
	lanty := &Lanty{
		controller:      controller,
		sidebar:         NewSidebar(setting.APPLICATION_NAME),
		gamebrowser:     NewGameBrowser(controller),
		downloadbrowser: NewDownloadBrowser(controller),
		userbrowser:     NewUserBrowser(controller),
	}
	lanty.showGameBrowser()

	lanty.sidebar.OnGamesTapped = func() {
		lanty.showGameBrowser()
		lanty.Refresh()
	}
	lanty.sidebar.OnDownloadsTapped = func() {
		lanty.showDownloadBrowser()
		lanty.Refresh()
	}
	lanty.sidebar.OnUsersTapped = func() {
		lanty.showUserBrowser()
		lanty.Refresh()
	}

	lanty.ExtendBaseWidget(lanty)

	return lanty
}

func (widget *Lanty) showGameBrowser() {
	widget.gamebrowser.Show()
	widget.downloadbrowser.Hide()
	widget.userbrowser.Hide()
	widget.Refresh()
}

func (widget *Lanty) showDownloadBrowser() {
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Show()
	widget.userbrowser.Hide()
	widget.Refresh()
}

func (widget *Lanty) showUserBrowser() {
	widget.gamebrowser.Hide()
	widget.downloadbrowser.Hide()
	widget.userbrowser.Show()
	widget.Refresh()
}

func (lanty *Lanty) CreateRenderer() fyne.WidgetRenderer {
	return newLantyRenderer(lanty)
}

type lantyRenderer struct {
	lanty *Lanty
	main  *fyne.Container
}

func newLantyRenderer(lanty *Lanty) *lantyRenderer {
	return &lantyRenderer{
		lanty: lanty,
		main: container.NewMax(
			lanty.gamebrowser,
			lanty.downloadbrowser,
			lanty.userbrowser),
	}
}

func (renderer *lantyRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.lanty.sidebar,
		renderer.main,
	}
}

func (renderer *lantyRenderer) Layout(size fyne.Size) {
	renderer.lanty.sidebar.Move(fyne.NewPos(0, 0))
	renderer.lanty.sidebar.Resize(fyne.NewSize(fyne.Max(size.Width/7, renderer.lanty.sidebar.MinSize().Width), size.Height))
	renderer.main.Resize(fyne.NewSize(size.Width-fyne.Max(size.Width/7, renderer.lanty.sidebar.MinSize().Width), size.Height))
	renderer.main.Move(fyne.NewPos(fyne.Max(size.Width/7, renderer.lanty.sidebar.MinSize().Width), 0))

}

func (renderer *lantyRenderer) MinSize() fyne.Size {
	size := fyne.NewSize(0, 0)
	size.Height = fyne.Max(renderer.lanty.sidebar.MinSize().Height, renderer.main.MinSize().Height)
	size.Width = renderer.main.MinSize().Width + fyne.Max(renderer.main.MinSize().Width/6, renderer.lanty.sidebar.MinSize().Width)
	return size
}

func (renderer *lantyRenderer) Refresh() {

}

func (renderer *lantyRenderer) Destroy() {}
