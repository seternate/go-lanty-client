package gameview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
)

type GameTile struct {
	widget.BaseWidget

	vm *gameviewmodel.GameTile

	startSingleplayerButton *widget.Button
	joinMultiplayerButton   *widget.Button
	hostMultiplayerButton   *widget.Button
	installationButton      *widget.Button
	openInExplorerButton    *widget.Button
	progressBar             *widget.ProgressBar
}

func NewGameTile(vm *gameviewmodel.GameTile) *GameTile {
	view := &GameTile{
		vm:                      vm,
		startSingleplayerButton: widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		joinMultiplayerButton:   widget.NewButtonWithIcon("", theme.LoginIcon(), nil),
		hostMultiplayerButton:   widget.NewButtonWithIcon("", theme.StorageIcon(), nil),
		installationButton:      widget.NewButtonWithIcon("", theme.DownloadIcon(), nil),
		openInExplorerButton:    widget.NewButtonWithIcon("", theme.FolderIcon(), nil),
		progressBar:             widget.NewProgressBar(),
	}

	view.progressBar.TextFormatter = func() string {
		return ""
	}
	view.progressBar.Bind(vm.Progress)

	view.startSingleplayerButton.OnTapped = func() {
		vm.StartSingleplayer()
	}
	view.joinMultiplayerButton.OnTapped = func() {
		vm.OpenUserSelectionToJoinMultiplayer()
	}
	view.hostMultiplayerButton.OnTapped = func() {
		vm.OpenArgumentConfigurationToHostMultiplayer()
	}
	view.installationButton.OnTapped = func() {
		vm.ToggleInstallation()
	}
	view.openInExplorerButton.OnTapped = func() {
		vm.OpenDirectoryInExplorer()
	}

	vm.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)

	return view
}

func (view *GameTile) Refresh() {
	if view.vm.EnableSingleplayerButton() {
		view.startSingleplayerButton.Enable()
	} else {
		view.startSingleplayerButton.Disable()
	}

	if view.vm.EnableJoinMultiplayerButton() {
		view.joinMultiplayerButton.Enable()
	} else {
		view.joinMultiplayerButton.Disable()
	}

	if view.vm.EnableHostMultiplayerButton() {
		view.hostMultiplayerButton.Enable()
	} else {
		view.hostMultiplayerButton.Disable()
	}

	if view.vm.ShowInstallationStopIcon() {
		view.installationButton.SetIcon(theme.CancelIcon())
	} else {
		view.installationButton.SetIcon(theme.DownloadIcon())
	}

	if view.vm.EnableOpenInExplorerButton() {
		view.openInExplorerButton.Enable()
	} else {
		view.openInExplorerButton.Disable()
	}

	view.startSingleplayerButton.Refresh()
	view.joinMultiplayerButton.Refresh()
	view.hostMultiplayerButton.Refresh()
	view.installationButton.Refresh()
	view.openInExplorerButton.Refresh()
	view.progressBar.Refresh()

	view.BaseWidget.Refresh()
}

func (widget *GameTile) CreateRenderer() fyne.WidgetRenderer {
	return newGameTileRenderer(widget)
}

type gameTileRenderer struct {
	view *GameTile

	background *canvas.Rectangle

	icon                  *canvas.Image
	iconratio             float32
	name                  *canvas.Text
	installedFileSizeText *canvas.Text

	statusBackground *canvas.Rectangle
	statusText       *canvas.Text

	progressFrontText *canvas.Text
	progressEndText   *canvas.Text

	objects []fyne.CanvasObject
}

func newGameTileRenderer(view *GameTile) *gameTileRenderer {
	renderer := &gameTileRenderer{
		view: view,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	renderer.icon = canvas.NewImageFromImage(view.vm.GetIcon())
	renderer.icon.FillMode = canvas.ImageFillContain
	renderer.iconratio = float32(view.vm.GetIcon().Bounds().Max.X) / float32(view.vm.GetIcon().Bounds().Max.Y)
	renderer.name = canvas.NewText(view.vm.GetName(), color.White)
	renderer.name.TextSize = theme.TextSize()
	renderer.installedFileSizeText = canvas.NewText(view.vm.GetInstalledFileSize(), color.RGBA{150, 150, 150, 255})
	renderer.installedFileSizeText.TextSize = theme.TextSize() * 0.8

	renderer.statusBackground = canvas.NewRectangle(view.vm.GetStatusColor())
	renderer.statusBackground.CornerRadius = theme.SmallCornerRadius()
	renderer.statusText = canvas.NewText(view.vm.GetStatusText(), color.White)
	renderer.statusText.TextSize = theme.TextSize() * 0.9

	renderer.progressFrontText = canvas.NewText(view.vm.GetProgressFrontText(), color.RGBA{180, 180, 180, 255})
	renderer.progressFrontText.TextSize = theme.TextSize() * 0.85
	renderer.progressEndText = canvas.NewText(view.vm.GetProgressEndText(), color.RGBA{180, 180, 180, 255})
	renderer.progressEndText.TextSize = theme.TextSize() * 0.85

	buttonHeight := float32(32)
	buttonWidth := buttonHeight * 2
	buttonSize := fyne.NewSize(buttonWidth, buttonHeight)
	for _, button := range []*widget.Button{
		renderer.view.startSingleplayerButton,
		renderer.view.joinMultiplayerButton,
		renderer.view.hostMultiplayerButton,
		renderer.view.openInExplorerButton,
		renderer.view.installationButton,
	} {
		button.Resize(buttonSize)
	}

	renderer.objects = []fyne.CanvasObject{
		renderer.background,
		renderer.icon,
		renderer.name,
		renderer.installedFileSizeText,
		renderer.statusBackground,
		renderer.statusText,
		renderer.progressFrontText,
		renderer.progressEndText,
		renderer.view.startSingleplayerButton,
		renderer.view.joinMultiplayerButton,
		renderer.view.hostMultiplayerButton,
		renderer.view.openInExplorerButton,
		renderer.view.installationButton,
		renderer.view.progressBar,
	}

	return renderer
}

func (renderer *gameTileRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *gameTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	iconHeight := float32(64)
	iconWidth := iconHeight * renderer.iconratio
	renderer.icon.Resize(fyne.NewSize(iconWidth, iconHeight))
	renderer.icon.Move(fyne.NewPos(theme.InnerPadding(), theme.InnerPadding()))
	iconRight := renderer.icon.Position().X + iconWidth + theme.InnerPadding()
	iconTop := renderer.icon.Position().Y
	iconLeft := renderer.icon.Position().X
	iconBottom := renderer.icon.Position().Y + iconHeight + theme.InnerPadding()

	progressBarSize := fyne.NewSize(size.Width-iconRight-theme.InnerPadding(), float32(6))
	renderer.view.progressBar.Resize(progressBarSize)
	renderer.view.progressBar.Move(fyne.NewPos(iconRight, iconBottom-theme.InnerPadding()-progressBarSize.Height))
	progressFrontTextSize := fyne.MeasureText(renderer.view.vm.GetProgressFrontText(), renderer.progressFrontText.TextSize, renderer.progressFrontText.TextStyle)
	renderer.progressFrontText.Move(fyne.NewPos(iconRight, renderer.view.progressBar.Position().Y-progressFrontTextSize.Height-theme.InnerPadding()/2))
	progressEndTextSize := fyne.MeasureText(renderer.view.vm.GetProgressEndText(), renderer.progressEndText.TextSize, renderer.progressEndText.TextStyle)
	renderer.progressEndText.Move(fyne.NewPos(size.Width-progressEndTextSize.Width-theme.InnerPadding(), renderer.view.progressBar.Position().Y-progressEndTextSize.Height-theme.InnerPadding()/2))

	nameTextSize := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle)
	installedFileSizeTextSize := fyne.MeasureText(renderer.installedFileSizeText.Text, renderer.installedFileSizeText.TextSize, renderer.installedFileSizeText.TextStyle)
	if renderer.view.vm.ShowProgressIndicator() {
		renderer.name.Move(fyne.NewPos(iconRight, iconTop))
		renderer.installedFileSizeText.Move(fyne.NewPos(iconRight, renderer.name.Position().Y+nameTextSize.Height))
	} else {
		renderer.name.Move(fyne.NewPos(iconRight, iconTop+iconHeight/2-(nameTextSize.Height+installedFileSizeTextSize.Height)/2))
		renderer.installedFileSizeText.Move(fyne.NewPos(iconRight, renderer.name.Position().Y+nameTextSize.Height))
	}

	statusTextSize := fyne.MeasureText(renderer.view.vm.GetStatusText(), renderer.statusText.TextSize, renderer.statusText.TextStyle)
	renderer.statusBackground.Resize(fyne.NewSize(statusTextSize.Width+1.2*theme.InnerPadding(), statusTextSize.Height+1.2*theme.InnerPadding()))
	renderer.statusBackground.Move(fyne.NewPos(size.Width-renderer.statusBackground.Size().Width-theme.InnerPadding(), iconTop))
	renderer.statusText.Move(fyne.NewPos(renderer.statusBackground.Position().X+0.6*theme.InnerPadding(), renderer.statusBackground.Position().Y+0.6*theme.InnerPadding()))

	renderer.view.startSingleplayerButton.Move(fyne.NewPos(iconLeft, iconBottom))
	x := renderer.view.startSingleplayerButton.Position().X + renderer.view.startSingleplayerButton.Size().Width + theme.InnerPadding()/2
	renderer.view.joinMultiplayerButton.Move(fyne.NewPos(x, iconBottom))
	x += renderer.view.joinMultiplayerButton.Size().Width + theme.InnerPadding()/2
	renderer.view.hostMultiplayerButton.Move(fyne.NewPos(x, iconBottom))
	x += renderer.view.hostMultiplayerButton.Size().Width + theme.InnerPadding()/2
	renderer.view.installationButton.Move(fyne.NewPos(x, iconBottom))
	x += renderer.view.installationButton.Size().Width + theme.InnerPadding()/2
	renderer.view.openInExplorerButton.Move(fyne.NewPos(x, iconBottom))
}

func (renderer *gameTileRenderer) MinSize() fyne.Size {
	minHeight := renderer.view.startSingleplayerButton.Position().Y + renderer.view.startSingleplayerButton.Size().Height + theme.InnerPadding()
	minWidth := renderer.view.openInExplorerButton.Position().X + renderer.view.openInExplorerButton.Size().Width + theme.InnerPadding()

	return fyne.NewSize(float32(minWidth), float32(minHeight))
}

func (renderer *gameTileRenderer) Refresh() {
	renderer.icon.Image = renderer.view.vm.GetIcon()
	renderer.iconratio = float32(renderer.view.vm.GetIcon().Bounds().Max.X) / float32(renderer.view.vm.GetIcon().Bounds().Max.Y)
	renderer.name.Text = renderer.view.vm.GetName()
	renderer.installedFileSizeText.Text = renderer.view.vm.GetInstalledFileSize()
	renderer.statusBackground.FillColor = renderer.view.vm.GetStatusColor()
	renderer.statusText.Text = renderer.view.vm.GetStatusText()
	renderer.progressFrontText.Text = renderer.view.vm.GetProgressFrontText()
	renderer.progressEndText.Text = renderer.view.vm.GetProgressEndText()

	if renderer.view.vm.ShowProgressIndicator() {
		renderer.progressFrontText.Show()
		renderer.view.progressBar.Show()
		renderer.progressEndText.Show()
	} else {
		renderer.progressFrontText.Hide()
		renderer.progressEndText.Hide()
		renderer.view.progressBar.Hide()
	}

	renderer.background.Refresh()
	renderer.icon.Refresh()
	renderer.name.Refresh()
	renderer.installedFileSizeText.Refresh()
	renderer.statusBackground.Refresh()
	renderer.statusText.Refresh()
	renderer.progressFrontText.Refresh()
	renderer.progressEndText.Refresh()
}

func (renderer *gameTileRenderer) Destroy() {}
