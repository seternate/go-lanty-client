package theme

import (
	"image/color"

	fynetheme "fyne.io/fyne/v2/theme"
)

func BackgroundColor() color.Color {
	return color.RGBA{126, 126, 126, 255}
}

func TextSize() float32 {
	return 14
}

func InnerPadding() float32 {
	return fynetheme.InnerPadding()
}

func ForegroundColor() color.Color {
	return fynetheme.ForegroundColor()
}
