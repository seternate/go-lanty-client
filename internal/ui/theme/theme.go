package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	fynetheme "fyne.io/fyne/v2/theme"
)

func BackgroundColor() color.Color {
	return color.RGBA{126, 126, 126, 255}
}

func BackgroundColor2() color.Color {
	return fynetheme.Color(fynetheme.ColorNameInputBorder)
}

func ForegroundColor() color.Color {
	return fynetheme.Color(fynetheme.ColorNameForeground)
}

func TextSize() float32 {
	return 14
}

func InnerPadding() float32 {
	return fynetheme.InnerPadding()
}

func CornerRadius() float32 {
	return 12 // Modern, rounded corner radius
}

func SmallCornerRadius() float32 {
	return 6 // Smaller corner radius for small components like status badges
}

func MediaPlayIcon() fyne.Resource {
	return fynetheme.MediaPlayIcon()
}

func LoginIcon() fyne.Resource {
	return fynetheme.LoginIcon()
}

func FolderIcon() fyne.Resource {
	return fynetheme.FolderIcon()
}

func DownloadIcon() fyne.Resource {
	return fynetheme.DownloadIcon()
}

func UploadIcon() fyne.Resource {
	return fynetheme.UploadIcon()
}

func StorageIcon() fyne.Resource {
	return fynetheme.StorageIcon()
}

func CancelIcon() fyne.Resource {
	return fynetheme.CancelIcon()
}

func ContentCopyIcon() fyne.Resource {
	return fynetheme.ContentCopyIcon()
}

func ExtensionUpToDateIcon() fyne.Resource {
	return fynetheme.ConfirmIcon()
}

func ExtensionNewIcon() fyne.Resource {
	return fynetheme.ViewRefreshIcon()
}

const (
	StatusDownloading = "downloading"
	StatusExtracting  = "extracting"
	StatusError       = "error"
	StatusReady       = "ready"
	StatusNotReady    = "not ready"
	StatusCanceled    = "canceled"
)

func StatusColor(status string) color.Color {
	switch status {
	case StatusNotReady:
		return color.RGBA{97, 97, 97, 255} // Material Design Grey-600
	case StatusDownloading:
		return color.RGBA{33, 150, 243, 255} // Material Design Blue-500
	case StatusExtracting:
		return color.RGBA{156, 39, 176, 255} // Material Design Purple-500
	case StatusError:
		return color.RGBA{244, 67, 54, 255} // Material Design Red-500
	case StatusReady:
		return color.RGBA{76, 175, 80, 255} // Material Design Green-500
	case StatusCanceled:
		return color.RGBA{255, 152, 0, 255} // Material Design Orange-500
	default:
		return color.RGBA{97, 97, 97, 255} // Material Design Grey-600 (default)
	}
}
