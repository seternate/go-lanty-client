package gamebrowser

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/seternate/go-lanty-client/pkg/controller"
	"github.com/seternate/go-lanty/pkg/game"
)

type Gametile struct {
	Container *fyne.Container
	info      *Info
	control   *Control
	icon      *canvas.Image

	controller *controller.GameController
	game       game.Game
}

func NewGametile(controller *controller.GameController, game game.Game) (tile *Gametile) {
	tile = &Gametile{
		controller: controller,
		game:       game,
		info:       NewInfo(controller, game),
		control:    NewControl(controller, game),
	}
	tile.setIcon(controller.GetIcon(game))

	tile.Container = container.NewMax(canvas.NewRectangle(color.RGBA{126, 126, 126, 255}),
		container.NewBorder(nil,
			nil,
			tile.icon,
			nil,
			container.NewVBox(
				tile.info.Container,
				layout.NewSpacer(),
				container.NewPadded(
					tile.control.Container,
				),
			),
		),
	)

	return
}

func (tile *Gametile) setIcon(icon image.Image) {
	image := canvas.NewImageFromImage(icon)
	image.SetMinSize(fyne.NewSize(130, 130))
	tile.icon = image
}
