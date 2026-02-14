package gameview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	gamemodel "github.com/seternate/go-lanty-client/internal/ui/model/game"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type GameBrowser struct {
	widget.BaseWidget

	stats *gamemodel.GameStatsModel

	gamelist *GameList
}

func NewGameBrowser(gamelist *GameList, stats *gamemodel.GameStatsModel) *GameBrowser {
	view := &GameBrowser{
		stats:    stats,
		gamelist: gamelist,
	}

	stats.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)
	return view
}

func (view *GameBrowser) CreateRenderer() fyne.WidgetRenderer {
	return newGameRenderer(view)
}

type gameRenderer struct {
	view            *GameBrowser
	headerContainer *fyne.Container
	headerText      *canvas.Text
	statsText       *canvas.Text

	objects []fyne.CanvasObject
}

func newGameRenderer(view *GameBrowser) *gameRenderer {
	renderer := &gameRenderer{
		view: view,
	}

	renderer.headerText = canvas.NewText("Games", color.White)
	renderer.headerText.TextSize = 28
	renderer.headerText.TextStyle = fyne.TextStyle{Bold: true}

	statsText, err := view.stats.StatsText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game stats text")
	}
	renderer.statsText = canvas.NewText(statsText, color.RGBA{180, 180, 180, 255})
	renderer.statsText.TextSize = theme.TextSize()

	headerPadding := theme.InnerPadding() * 2

	paddingTop := canvas.NewRectangle(color.Transparent)
	paddingTop.SetMinSize(fyne.NewSize(0, headerPadding))
	paddingBottom := canvas.NewRectangle(color.Transparent)
	paddingBottom.SetMinSize(fyne.NewSize(0, headerPadding))

	paddingLeft := canvas.NewRectangle(color.Transparent)
	paddingLeft.SetMinSize(fyne.NewSize(headerPadding, 0))
	paddingRight := canvas.NewRectangle(color.Transparent)
	paddingRight.SetMinSize(fyne.NewSize(headerPadding, 0))

	headerLeft := container.NewBorder(paddingTop, paddingBottom, paddingLeft, nil, container.NewWithoutLayout(renderer.headerText))
	headerRight := container.NewBorder(paddingTop, paddingBottom, nil, paddingRight, container.NewWithoutLayout(renderer.statsText))

	renderer.headerContainer = container.NewBorder(nil, nil, headerLeft, headerRight)

	renderer.objects = []fyne.CanvasObject{
		container.NewBorder(renderer.headerContainer, nil, nil, nil, renderer.view.gamelist),
	}

	return renderer
}

func (r *gameRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *gameRenderer) Layout(size fyne.Size) {
	if len(r.objects) > 0 {
		r.objects[0].Resize(size)
		r.objects[0].Move(fyne.NewPos(0, 0))
	}
}

func (r *gameRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, r.headerContainer.MinSize().Height)
}

func (r *gameRenderer) Refresh() {
	statsText, err := r.view.stats.StatsText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game stats text")
	}
	r.statsText.Text = statsText

	r.headerText.Refresh()
	r.statsText.Refresh()
	r.view.gamelist.Refresh()
}

func (r *gameRenderer) Destroy() {}
