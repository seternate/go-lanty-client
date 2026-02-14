package gameview

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog/log"
	model "github.com/seternate/go-lanty-client/internal/ui/model/game"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type GameTile struct {
	widget.BaseWidget

	game *model.GameTileModel

	playButton     *widget.Button
	joinButton     *widget.Button
	serverButton   *widget.Button
	downloadButton *widget.Button
	openButton     *widget.Button
	progressBar    *widget.ProgressBar

	OnPlayButtonPressed     func(slug string)
	OnJoinButtonPressed     func(slug string)
	OnServerButtonPressed   func(slug string)
	OnDownloadButtonPressed func(slug string)
	OnOpenButtonPressed     func(slug string)
}

func NewGameTile(game *model.GameTileModel) *GameTile {
	widget := &GameTile{
		game:           game,
		playButton:     widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		joinButton:     widget.NewButtonWithIcon("", theme.LoginIcon(), nil),
		serverButton:   widget.NewButtonWithIcon("", theme.StorageIcon(), nil),
		downloadButton: widget.NewButtonWithIcon("", theme.DownloadIcon(), nil),
		openButton:     widget.NewButtonWithIcon("", theme.FolderIcon(), nil),
		progressBar:    widget.NewProgressBar(),
	}

	widget.progressBar.TextFormatter = func() string {
		return ""
	}
	widget.progressBar.Bind(game.Progress)

	widget.playButton.Disable()
	widget.joinButton.Disable()
	widget.serverButton.Disable()

	widget.playButton.OnTapped = func() {
		if widget.OnPlayButtonPressed != nil {
			widget.OnPlayButtonPressed(widget.game.Slug)
		}
	}
	widget.joinButton.OnTapped = func() {
		if widget.OnJoinButtonPressed != nil {
			widget.OnJoinButtonPressed(widget.game.Slug)
		}
	}
	widget.serverButton.OnTapped = func() {
		if widget.OnServerButtonPressed != nil {
			widget.OnServerButtonPressed(widget.game.Slug)
		}
	}
	widget.downloadButton.OnTapped = func() {
		if widget.OnDownloadButtonPressed != nil {
			widget.OnDownloadButtonPressed(widget.game.Slug)
		}
	}
	widget.openButton.OnTapped = func() {
		if widget.OnOpenButtonPressed != nil {
			widget.OnOpenButtonPressed(widget.game.Slug)
		}
	}

	game.AddChangeListener(func() {
		widget.Refresh()
	})

	widget.ExtendBaseWidget(widget)

	return widget
}

func (widget *GameTile) Refresh() {
	statusText, err := widget.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	supportsJoiningMultiplayer, err := widget.game.SupportsJoiningMultiplayer.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile supports joining multiplayer")
	}
	supportsHostingServer, err := widget.game.SupportsHostingServer.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile supports hosting server")
	}

	installationDetected, err := widget.game.InstallationDetected.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile installation detected")
	}

	if statusText == string(model.StatusInstalled) || (statusText == string(model.StatusCanceled) && installationDetected) {
		widget.playButton.Enable()
		if supportsJoiningMultiplayer {
			widget.joinButton.Enable()
		} else {
			widget.joinButton.Disable()
		}
		if supportsHostingServer {
			widget.serverButton.Enable()
		} else {
			widget.serverButton.Disable()
		}
	} else {
		widget.playButton.Disable()
		widget.joinButton.Disable()
		widget.serverButton.Disable()

	}

	if installationDetected {
		widget.openButton.Enable()
	} else {
		widget.openButton.Disable()
	}

	widget.playButton.Refresh()
	widget.joinButton.Refresh()
	widget.serverButton.Refresh()
	widget.openButton.Refresh()
	widget.downloadButton.Refresh()
	widget.progressBar.Refresh()

	widget.BaseWidget.Refresh()
}

func (widget *GameTile) CreateRenderer() fyne.WidgetRenderer {
	return newGameTileRenderer(widget)
}

type gameTileRenderer struct {
	widget *GameTile

	background        *canvas.Rectangle
	icon              *canvas.Image
	name              *canvas.Text
	statusBackground  *canvas.Rectangle
	statusText        *canvas.Text
	progressFrontText *canvas.Text
	progressEndText   *canvas.Text
	blobSizeText      *canvas.Text

	nameSizeContainer *fyne.Container

	iconratio float32

	objects []fyne.CanvasObject
}

func newGameTileRenderer(tile *GameTile) *gameTileRenderer {
	renderer := &gameTileRenderer{
		widget:    tile,
		iconratio: 1.0,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()

	iv, err := tile.game.Icon.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile icon")
	}
	if img, ok := iv.(image.Image); ok {
		renderer.icon = canvas.NewImageFromImage(img)
		renderer.icon.FillMode = canvas.ImageFillContain
		renderer.iconratio = float32(img.Bounds().Max.X) / float32(img.Bounds().Max.Y)
	}

	nameText, err := tile.game.Name.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile name")
	}
	renderer.name = canvas.NewText(nameText, color.White)
	renderer.name.TextSize = theme.TextSize()

	statusColor, err := tile.game.StatusColor.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status color")
	}
	if colorVal, ok := statusColor.(color.Color); ok {
		renderer.statusBackground = canvas.NewRectangle(colorVal)
		renderer.statusBackground.CornerRadius = theme.SmallCornerRadius()
	}

	statusText, err := tile.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	renderer.statusText = canvas.NewText(statusText, color.White)
	renderer.statusText.TextSize = theme.TextSize() * 0.9

	progressFrontText, err := tile.game.ProgressFrontText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress front text")
	}
	renderer.progressFrontText = canvas.NewText(progressFrontText, color.RGBA{180, 180, 180, 255})
	renderer.progressFrontText.TextSize = theme.TextSize() * 0.85

	progressEndText, err := tile.game.ProgressEndText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress end text")
	}
	renderer.progressEndText = canvas.NewText(progressEndText, color.RGBA{180, 180, 180, 255})
	renderer.progressEndText.TextSize = theme.TextSize() * 0.85

	blobSizeText, err := tile.game.BlobSizeText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile blob size text")
	}
	renderer.blobSizeText = canvas.NewText(blobSizeText, color.RGBA{150, 150, 150, 255})
	renderer.blobSizeText.TextSize = theme.TextSize() * 0.8

	buttonHeight := float32(40)
	buttonWidth := buttonHeight * 2
	buttonSize := fyne.NewSize(buttonWidth, buttonHeight)

	for _, btn := range []*widget.Button{
		renderer.widget.playButton,
		renderer.widget.joinButton,
		renderer.widget.serverButton,
		renderer.widget.openButton,
		renderer.widget.downloadButton,
	} {
		btn.Resize(buttonSize)
	}

	renderer.nameSizeContainer = container.NewWithoutLayout(renderer.name, renderer.blobSizeText)

	renderer.objects = []fyne.CanvasObject{
		renderer.background,
		renderer.icon,
		renderer.nameSizeContainer,
		renderer.statusBackground,
		renderer.statusText,
		renderer.progressFrontText,
		renderer.progressEndText,
		renderer.widget.playButton,
		renderer.widget.joinButton,
		renderer.widget.serverButton,
		renderer.widget.openButton,
		renderer.widget.downloadButton,
		renderer.widget.progressBar,
	}

	return renderer
}

func (renderer *gameTileRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *gameTileRenderer) Layout(size fyne.Size) {
	renderer.background.Resize(size)

	padding := theme.InnerPadding()
	buttonHeight := float32(40)
	buttonSpacing := padding / 2

	iconHeight := float32(64)
	iconWidth := iconHeight * renderer.iconratio
	renderer.icon.Resize(fyne.NewSize(iconWidth, iconHeight))
	renderer.icon.Move(fyne.NewPos(padding, padding))
	iconRight := padding + iconWidth + padding

	statusText, err := renderer.widget.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	statusTextSize := fyne.MeasureText(statusText, renderer.statusText.TextSize, renderer.statusText.TextStyle)
	statusPadding := padding * 0.6
	statusBgWidth := statusTextSize.Width + 2*statusPadding
	statusBgHeight := statusTextSize.Height + 2*statusPadding
	statusX := size.Width - statusBgWidth - padding
	statusY := padding
	renderer.statusBackground.Resize(fyne.NewSize(statusBgWidth, statusBgHeight))
	renderer.statusBackground.Move(fyne.NewPos(statusX, statusY))
	renderer.statusText.Move(fyne.NewPos(statusX+statusPadding, statusY+statusPadding))

	isProgressing, err := renderer.widget.game.IsProgressing.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile is progressing")
	}
	progressBarHeight := float32(6)

	iconCenter := padding + iconHeight/2
	iconBottom := padding + iconHeight
	iconTop := padding

	textSpacing := padding * 0.4

	nameHeight := fyne.MeasureText(renderer.name.Text, renderer.name.TextSize, renderer.name.TextStyle).Height
	blobSizeText, err := renderer.widget.game.BlobSizeText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile blob size text")
	}
	blobSizeHeight := fyne.MeasureText(blobSizeText, renderer.blobSizeText.TextSize, renderer.blobSizeText.TextStyle).Height
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
		renderer.nameSizeContainer.Move(fyne.NewPos(iconRight, containerY))
		renderer.nameSizeContainer.Resize(fyne.NewSize(progressBarColumnWidth, containerHeight))

		renderer.name.Resize(fyne.NewSize(progressBarColumnWidth, nameHeight))
		renderer.name.Move(fyne.NewPos(0, 0))
		renderer.blobSizeText.Resize(fyne.NewSize(progressBarColumnWidth, blobSizeHeight))
		renderer.blobSizeText.Move(fyne.NewPos(0, nameHeight))

		endText, _ := renderer.widget.game.ProgressEndText.Get()
		endTextSize := fyne.MeasureText(endText, renderer.progressEndText.TextSize, renderer.progressEndText.TextStyle)
		renderer.progressEndText.Move(fyne.NewPos(progressBarColumnX+progressBarColumnWidth-endTextSize.Width, progressTextY))

		renderer.widget.progressBar.Resize(fyne.NewSize(progressBarColumnWidth, progressBarHeight))
		renderer.widget.progressBar.Move(fyne.NewPos(progressBarColumnX, progressBarY))
	} else {
		containerY = iconCenter - containerHeight/2
		containerWidth := size.Width - iconRight - padding
		renderer.nameSizeContainer.Move(fyne.NewPos(iconRight, containerY))
		renderer.nameSizeContainer.Resize(fyne.NewSize(containerWidth, containerHeight))

		renderer.name.Resize(fyne.NewSize(containerWidth, nameHeight))
		renderer.name.Move(fyne.NewPos(0, 0))
		renderer.blobSizeText.Resize(fyne.NewSize(containerWidth, blobSizeHeight))
		renderer.blobSizeText.Move(fyne.NewPos(0, nameHeight))

		renderer.progressFrontText.Hide()
		renderer.progressEndText.Hide()
		renderer.widget.progressBar.Hide()
	}

	var buttonY float32
	var buttonStartX float32
	var buttonWidth float32
	buttonWidth = buttonHeight * 2
	buttonY = iconBottom + padding
	buttonStartX = padding

	renderer.widget.playButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.widget.playButton.Move(fyne.NewPos(buttonStartX, buttonY))

	x := buttonStartX + buttonWidth + buttonSpacing
	renderer.widget.joinButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.widget.joinButton.Move(fyne.NewPos(x, buttonY))
	x += buttonWidth + buttonSpacing

	renderer.widget.serverButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.widget.serverButton.Move(fyne.NewPos(x, buttonY))
	x += buttonWidth + buttonSpacing

	renderer.widget.downloadButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.widget.downloadButton.Move(fyne.NewPos(x, buttonY))
	x += buttonWidth + buttonSpacing

	renderer.widget.openButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))
	renderer.widget.openButton.Move(fyne.NewPos(x, buttonY))
}

func (renderer *gameTileRenderer) MinSize() fyne.Size {
	padding := theme.InnerPadding()
	buttonHeight := float32(40)
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
	if iconVal, err := renderer.widget.game.Icon.Get(); err == nil {
		if img, ok := iconVal.(image.Image); ok && img != nil {
			renderer.icon.Image = img
			renderer.iconratio = float32(img.Bounds().Max.X) / float32(img.Bounds().Max.Y)
		}
	} else {
		log.Error().Err(err).Msg("could not get game tile icon")
	}

	if name, err := renderer.widget.game.Name.Get(); err == nil {
		renderer.name.Text = name
	} else {
		log.Error().Err(err).Msg("could not get game tile name")
	}

	statusColor, err := renderer.widget.game.StatusColor.Get()
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

	statusText, err := renderer.widget.game.StatusText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile status text")
	}
	renderer.statusText.Text = statusText

	progressFrontText, err := renderer.widget.game.ProgressFrontText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress front text")
	}
	renderer.progressFrontText.Text = progressFrontText

	progressEndText, err := renderer.widget.game.ProgressEndText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile progress end text")
	}
	renderer.progressEndText.Text = progressEndText

	blobSizeText, err := renderer.widget.game.BlobSizeText.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile blob size text")
	}
	renderer.blobSizeText.Text = blobSizeText

	isProgressing, err := renderer.widget.game.IsProgressing.Get()
	if err != nil {
		log.Error().Err(err).Msg("could not get game tile is progressing")
	}

	if isProgressing {
		renderer.widget.downloadButton.SetIcon(theme.CancelIcon())
	} else {
		renderer.widget.downloadButton.SetIcon(theme.DownloadIcon())
	}

	if isProgressing {
		renderer.progressFrontText.Show()
		renderer.widget.progressBar.Show()
		renderer.progressEndText.Show()
	} else {
		renderer.progressFrontText.Hide()
		renderer.progressEndText.Hide()
		renderer.widget.progressBar.Hide()
	}

	renderer.background.Refresh()
	renderer.icon.Refresh()
	renderer.name.Refresh()
	renderer.statusBackground.Refresh()
	renderer.statusText.Refresh()
	renderer.progressFrontText.Refresh()
	renderer.progressEndText.Refresh()
	renderer.blobSizeText.Refresh()
	renderer.nameSizeContainer.Refresh()
}

func (renderer *gameTileRenderer) Destroy() {}
