package layout

import (
	"math"

	"fyne.io/fyne/v2"
	"github.com/seternate/go-lanty-client/pkg/theme"
)

type GridScalingLayout struct {
	columns int
}

func NewGridScalingLayout(columns int) *GridScalingLayout {
	return &GridScalingLayout{
		columns: columns,
	}
}

func (layout *GridScalingLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	objectwidth := (size.Width - float32(layout.columns+1)*theme.InnerPadding()) / float32(layout.columns)
	objectheight := float32(0)
	for _, object := range objects {
		objectwidth = fyne.Max(object.MinSize().Width, objectwidth)
		objectheight = fyne.Max(object.MinSize().Height, objectheight)
	}

	for index, object := range objects {
		object.Resize(fyne.NewSize(objectwidth, objectheight))
		object.Move(fyne.NewPos(theme.InnerPadding()+float32(index%layout.columns)*(objectwidth+theme.InnerPadding()), theme.InnerPadding()+float32(index/layout.columns)*(objectheight+theme.InnerPadding())))
	}
}

func (layout *GridScalingLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	objectwidth := float32(0)
	objectheight := float32(0)
	objectsInY := float32(math.Ceil(float64(len(objects)) / float64(layout.columns)))

	for _, object := range objects {
		objectwidth = fyne.Max(object.MinSize().Width, objectwidth)
		objectheight = fyne.Max(object.MinSize().Height, objectheight)
	}
	minwidth := float32(layout.columns)*objectwidth + float32(layout.columns+1)*theme.InnerPadding()
	minheight := objectsInY*objectheight + (objectsInY-1)*theme.InnerPadding() + 2*theme.InnerPadding()

	return fyne.NewSize(minwidth, minheight)
}
