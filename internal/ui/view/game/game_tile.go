package gameview

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
	gameviewmodel "github.com/seternate/go-lanty-client/internal/ui/viewmodel/game"
)

type GameTile struct {
	widget.BaseWidget

	game *gameviewmodel.GameTile

	startSingleplayerButton *widget.Button
	joinMultiplayerButton   *widget.Button
	hostMultiplayerButton   *widget.Button
	downloadButton          *widget.Button
	openInExplorerButton    *widget.Button
	progressBar             *widget.ProgressBar
	extensionBadge          *widget.Button
}

func NewGameTile(game *gameviewmodel.GameTile) *GameTile {
	view := &GameTile{
		game:                    game,
		startSingleplayerButton: widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		joinMultiplayerButton:   widget.NewButtonWithIcon("", theme.LoginIcon(), nil),
		hostMultiplayerButton:   widget.NewButtonWithIcon("", theme.StorageIcon(), nil),
		downloadButton:          widget.NewButtonWithIcon("", theme.DownloadIcon(), nil),
		openInExplorerButton:    widget.NewButtonWithIcon("", theme.FolderIcon(), nil),
		progressBar:             widget.NewProgressBar(),
		extensionBadge:          widget.NewButtonWithIcon("", theme.ExtensionUpToDateIcon(), nil),
	}

	view.progressBar.TextFormatter = func() string {
		return ""
	}
	view.progressBar.Bind(game.Progress)

	view.startSingleplayerButton.Disable()
	view.joinMultiplayerButton.Disable()
	view.hostMultiplayerButton.Disable()
	view.extensionBadge.Importance = widget.LowImportance

	view.startSingleplayerButton.OnTapped = func() {
		game.StartSingleplayer()
	}
	view.joinMultiplayerButton.OnTapped = func() {
		game.OpenUserSelectionToJoinMultiplayer()
	}
	view.hostMultiplayerButton.OnTapped = func() {
		game.OpenArgumentConfigurationToHostMultiplayer()
	}
	view.downloadButton.OnTapped = func() {
		game.StartDownload()
	}
	view.openInExplorerButton.OnTapped = func() {
		game.OpenDirectoryInExplorer()
	}
	view.extensionBadge.OnTapped = func() {
		// game.OpenExtensions()
		// TODO
	}

	game.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)

	return view
}

func (view *GameTile) Refresh() {
	statusText, err := view.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	supportsJoiningMultiplayer, err := view.game.SupportsJoiningMultiplayer.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile supports joining multiplayer")
	}
	supportsHostingServer, err := view.game.SupportsHostingServer.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile supports hosting server")
	}
	installationDetected, err := view.game.InstallationDetected.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile installation detected")
	}

	if statusText == string(gameviewmodel.StatusInstalled) || (statusText == string(gameviewmodel.StatusCanceled) && installationDetected) {
		view.startSingleplayerButton.Enable()
		if supportsJoiningMultiplayer {
			view.joinMultiplayerButton.Enable()
		} else {
			view.joinMultiplayerButton.Disable()
		}
		if supportsHostingServer {
			view.hostMultiplayerButton.Enable()
		} else {
			view.hostMultiplayerButton.Disable()
		}
	} else {
		view.startSingleplayerButton.Disable()
		view.joinMultiplayerButton.Disable()
		view.hostMultiplayerButton.Disable()

	}

	if installationDetected {
		view.openInExplorerButton.Enable()
	} else {
		view.openInExplorerButton.Disable()
	}

	extensionsVisible, err := view.game.ExtensionsVisible.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile extensions visible")
	}
	hasNewExtensions, err := view.game.HasNewExtensions.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile has new extensions")
	}
	if extensionsVisible {
		view.extensionBadge.Show()
		if hasNewExtensions {
			view.extensionBadge.SetIcon(theme.ExtensionNewIcon())
		} else {
			view.extensionBadge.SetIcon(theme.ExtensionUpToDateIcon())
		}
	} else {
		view.extensionBadge.Hide()
	}

	view.startSingleplayerButton.Refresh()
	view.joinMultiplayerButton.Refresh()
	view.hostMultiplayerButton.Refresh()
	view.openInExplorerButton.Refresh()
	view.downloadButton.Refresh()
	view.progressBar.Refresh()
	view.extensionBadge.Refresh()

	view.BaseWidget.Refresh()
}

func (widget *GameTile) CreateRenderer() fyne.WidgetRenderer {
	return newGameTileRenderer(widget)
}

type gameTileRenderer struct {
	view *GameTile

	nameSizeContainer *fyne.Container

	background *canvas.Rectangle

	icon                  *canvas.Image
	name                  *canvas.Text
	installedFileSizeText *canvas.Text

	statusBackground *canvas.Rectangle
	statusText       *canvas.Text

	progressFrontText *canvas.Text
	progressEndText   *canvas.Text

	iconratio float32

	objects []fyne.CanvasObject
}

func newGameTileRenderer(view *GameTile) *gameTileRenderer {
	renderer := &gameTileRenderer{
		view: view,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	renderer.icon = canvas.NewImageFromImage(view.game.GetIcon())
	renderer.icon.FillMode = canvas.ImageFillContain
	renderer.iconratio = float32(view.game.GetIcon().Bounds().Max.X) / float32(view.game.GetIcon().Bounds().Max.Y)

	renderer.name = canvas.NewText(view.game.GetName(), color.White)
	renderer.name.TextSize = theme.TextSize()

	renderer.installedFileSizeText = canvas.NewText(view.game.GetInstalledFileSize(), color.RGBA{150, 150, 150, 255})
	renderer.installedFileSizeText.TextSize = theme.TextSize() * 0.8

	statusColor, err := view.game.StatusColor.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status color")
	}
	if colorVal, ok := statusColor.(color.Color); ok {
		renderer.statusBackground = canvas.NewRectangle(colorVal)
		renderer.statusBackground.CornerRadius = theme.SmallCornerRadius()
	}

	statusText, err := view.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	renderer.statusText = canvas.NewText(statusText, color.White)
	renderer.statusText.TextSize = theme.TextSize() * 0.9

	progressFrontText, err := view.game.ProgressFrontText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress front text")
	}
	renderer.progressFrontText = canvas.NewText(progressFrontText, color.RGBA{180, 180, 180, 255})
	renderer.progressFrontText.TextSize = theme.TextSize() * 0.85

	progressEndText, err := view.game.ProgressEndText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress end text")
	}
	renderer.progressEndText = canvas.NewText(progressEndText, color.RGBA{180, 180, 180, 255})
	renderer.progressEndText.TextSize = theme.TextSize() * 0.85

	buttonHeight := float32(40)
	buttonWidth := buttonHeight * 2
	buttonSize := fyne.NewSize(buttonWidth, buttonHeight)

	for _, btn := range []*widget.Button{
		renderer.view.startSingleplayerButton,
		renderer.view.joinMultiplayerButton,
		renderer.view.hostMultiplayerButton,
		renderer.view.openInExplorerButton,
		renderer.view.downloadButton,
	} {
		btn.Resize(buttonSize)
	}

	renderer.nameSizeContainer = container.NewWithoutLayout(renderer.name, renderer.installedFileSizeText)

	badgeSize := float32(20)
	renderer.view.extensionBadge.Resize(fyne.NewSize(badgeSize, badgeSize))
	renderer.view.extensionBadge.Hide()

	renderer.objects = []fyne.CanvasObject{
		renderer.background,
		renderer.icon,
		renderer.nameSizeContainer,
		renderer.view.extensionBadge,
		renderer.statusBackground,
		renderer.statusText,
		renderer.progressFrontText,
		renderer.progressEndText,
		renderer.view.startSingleplayerButton,
		renderer.view.joinMultiplayerButton,
		renderer.view.hostMultiplayerButton,
		renderer.view.openInExplorerButton,
		renderer.view.downloadButton,
		renderer.view.progressBar,
	}

	return renderer
}

func (renderer *gameTileRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *gameTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	padding := theme.InnerPadding()
	buttonHeight := float32(32)
	buttonSpacing := padding / 2

	iconHeight := float32(64)
	iconWidth := iconHeight * renderer.iconratio
	renderer.icon.Resize(fyne.NewSize(iconWidth, iconHeight))
	renderer.icon.Move(fyne.NewPos(padding, padding))
	iconRight := padding + iconWidth + padding

	statusText, err := renderer.view.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	statusTextSize := fyne.MeasureText(statusText, renderer.statusText.TextSize, renderer.statusText.TextStyle)
	statusPadding := padding * 0.6
	statusBgWidth := statusTextSize.Width + 2*statusPadding
	statusBgHeight := statusTextSize.Height + 2*statusPadding
	statusX := size.Width - statusBgWidth - padding
	statusY := padding
	nameRowWidth := statusX - iconRight
	renderer.statusBackground.Resize(fyne.NewSize(statusBgWidth, statusBgHeight))
	renderer.statusBackground.Move(fyne.NewPos(statusX, statusY))
	renderer.statusText.Move(fyne.NewPos(statusX+statusPadding, statusY+statusPadding))

	isProgressing, err := renderer.view.game.IsProgressing.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile is progressing")
	}
	progressBarHeight := float32(6)

	iconCenter := padding + iconHeight/2
	iconBottom := padding + iconHeight
	iconTop := padding

	textSpacing := padding * 0.4

	nameHeight := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Height
	blobSizeText, err := renderer.view.game.InstalledFileSize.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile blob size text")
	}
	blobSizeHeight := fyne.MeasureText(blobSizeText, renderer.installedFileSizeText.TextSize, renderer.installedFileSizeText.TextStyle).Height
	containerHeight := nameHeight + blobSizeHeight

	var progressBarColumnX float32
	var progressBarColumnWidth float32
	var progressBarY float32
	var containerY float32

	if isProgressing {
		progressBarColumnX = iconRight
		progressBarColumnWidth = size.Width - progressBarColumnX - padding

		progressBarY = iconBottom - progressBarHeight

		progressTextHeight := fyne.MeasureText("", renderer.progressFrontText.TextSize, renderer.progressFrontText.TextStyle).Height
		progressTextY := progressBarY - progressTextHeight - textSpacing
		renderer.progressFrontText.Move(fyne.NewPos(iconRight, progressTextY))

		containerY = (iconTop + progressTextY - containerHeight) / 2
		extensionsVisible, _ := renderer.view.game.ExtensionsVisible.Get()
		badgeSize := float32(20)
		nameWidth := progressBarColumnWidth
		if extensionsVisible {
			nameTextWidth := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Width
			nameDisplayWidth := min(nameTextWidth, nameRowWidth-badgeSize-padding*0.5)
			nameWidth = nameDisplayWidth
			renderer.view.extensionBadge.Resize(fyne.NewSize(badgeSize, badgeSize))
			renderer.view.extensionBadge.Move(fyne.NewPos(iconRight+nameDisplayWidth+padding*0.5, containerY+(nameHeight-badgeSize)/2))
		}
		renderer.nameSizeContainer.Move(fyne.NewPos(iconRight, containerY))
		renderer.nameSizeContainer.Resize(fyne.NewSize(progressBarColumnWidth, containerHeight))

		renderer.name.Resize(fyne.NewSize(nameWidth, nameHeight))
		renderer.name.Move(fyne.NewPos(0, 0))
		renderer.installedFileSizeText.Resize(fyne.NewSize(progressBarColumnWidth, blobSizeHeight))
		renderer.installedFileSizeText.Move(fyne.NewPos(0, nameHeight))

		endText, _ := renderer.view.game.ProgressEndText.Get()
		endTextSize := fyne.MeasureText(endText, renderer.progressEndText.TextSize, renderer.progressEndText.TextStyle)
		renderer.progressEndText.Move(fyne.NewPos(progressBarColumnX+progressBarColumnWidth-endTextSize.Width, progressTextY))

		renderer.view.progressBar.Resize(fyne.NewSize(progressBarColumnWidth, progressBarHeight))
		renderer.view.progressBar.Move(fyne.NewPos(progressBarColumnX, progressBarY))
	} else {
		containerY = iconCenter - containerHeight/2
		containerWidth := size.Width - iconRight - padding
		renderer.nameSizeContainer.Move(fyne.NewPos(iconRight, containerY))
		renderer.nameSizeContainer.Resize(fyne.NewSize(containerWidth, containerHeight))

		extensionsVisible, _ := renderer.view.game.ExtensionsVisible.Get()
		badgeSize := float32(20)
		nameWidth := containerWidth
		if extensionsVisible {
			nameTextWidth := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Width
			nameDisplayWidth := min(nameTextWidth, nameRowWidth-badgeSize-padding*0.5)
			nameWidth = nameDisplayWidth
			renderer.view.extensionBadge.Resize(fyne.NewSize(badgeSize, badgeSize))
			renderer.view.extensionBadge.Move(fyne.NewPos(iconRight+nameDisplayWidth+padding*0.5, containerY+(nameHeight-badgeSize)/2))
		}

		renderer.name.Resize(fyne.NewSize(nameWidth, nameHeight))
		renderer.name.Move(fyne.NewPos(0, 0))
		renderer.installedFileSizeText.Resize(fyne.NewSize(containerWidth, blobSizeHeight))
		renderer.installedFileSizeText.Move(fyne.NewPos(0, nameHeight))

		renderer.progressFrontText.Hide()
		renderer.progressEndText.Hide()
		renderer.view.progressBar.Hide()
	}

	var buttonY float32
	var buttonStartX float32
	var buttonWidth float32
	buttonWidth = buttonHeight * 2
	buttonY = iconBottom + padding
	buttonStartX = padding

	renderer.view.startSingleplayerButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.view.startSingleplayerButton.Move(fyne.NewPos(buttonStartX, buttonY))

	x := buttonStartX + buttonWidth + buttonSpacing
	renderer.view.joinMultiplayerButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.view.joinMultiplayerButton.Move(fyne.NewPos(x, buttonY))
	x += buttonWidth + buttonSpacing

	renderer.view.hostMultiplayerButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.view.hostMultiplayerButton.Move(fyne.NewPos(x, buttonY))
	x += buttonWidth + buttonSpacing

	renderer.view.downloadButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.view.downloadButton.Move(fyne.NewPos(x, buttonY))
	x += buttonWidth + buttonSpacing

	renderer.view.openInExplorerButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.view.openInExplorerButton.Move(fyne.NewPos(x, buttonY))
}

func (renderer *gameTileRenderer) MinSize() fyne.Size {
	padding := theme.InnerPadding()
	buttonHeight := float32(32)
	buttonSpacing := padding / 2

	iconHeight := float32(64)
	iconWidth := iconHeight * renderer.iconratio
	iconRight := padding + iconWidth + padding

	buttonMinWidth := float32(30)
	buttonRowWidth := 5*buttonMinWidth + 4*buttonSpacing

	minHeight := padding + iconHeight + padding + buttonHeight + padding

	minWidth := iconRight + buttonRowWidth + padding

	return fyne.NewSize(minWidth, minHeight)
}

func (renderer *gameTileRenderer) Refresh() {
	if iconVal, err := renderer.view.game.Icon.Get(); err == nil {
		if img, ok := iconVal.(image.Image); ok && img != nil {
			renderer.icon.Image = img
			renderer.iconratio = float32(img.Bounds().Max.X) / float32(img.Bounds().Max.Y)
		}
	} else {
		log.Error().Err(err).Msg("could not get game tile icon")
	}

	if name, err := renderer.view.game.Name.Get(); err == nil {
		renderer.name.Text = name
	} else {
		log.Error().Err(err).Msg("could not get game tile name")
	}

	statusColor, err := renderer.view.game.StatusColor.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status color")
	}
	if colorVal, ok := statusColor.(color.Color); ok {
		renderer.statusBackground.FillColor = colorVal
		if rgba, ok := colorVal.(color.RGBA); ok && rgba.R == 255 && rgba.G == 255 && rgba.B == 255 {
			renderer.statusText.Color = color.Black
		} else {
			renderer.statusText.Color = color.White
		}
	}

	statusText, err := renderer.view.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	renderer.statusText.Text = statusText

	progressFrontText, err := renderer.view.game.ProgressFrontText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress front text")
	}
	renderer.progressFrontText.Text = progressFrontText

	progressEndText, err := renderer.view.game.ProgressEndText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress end text")
	}
	renderer.progressEndText.Text = progressEndText

	blobSizeText, err := renderer.view.game.InstalledFileSize.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile blob size text")
	}
	renderer.installedFileSizeText.Text = blobSizeText

	isProgressing, err := renderer.view.game.IsProgressing.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile is progressing")
	}

	if isProgressing {
		renderer.view.downloadButton.SetIcon(theme.CancelIcon())
	} else {
		renderer.view.downloadButton.SetIcon(theme.DownloadIcon())
	}

	if isProgressing {
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
	renderer.statusBackground.Refresh()
	renderer.statusText.Refresh()
	renderer.progressFrontText.Refresh()
	renderer.progressEndText.Refresh()
	renderer.installedFileSizeText.Refresh()
	renderer.nameSizeContainer.Refresh()
}

func (renderer *gameTileRenderer) Destroy() {}
