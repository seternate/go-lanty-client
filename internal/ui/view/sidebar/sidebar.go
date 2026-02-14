package sidebar

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type Sidebar struct {
	widget.BaseWidget

	icon              fyne.Resource
	version           string
	header            fyne.CanvasObject
	headerIcon        *canvas.Image
	iconTextContainer fyne.CanvasObject
	navContainer      fyne.CanvasObject
	footer            fyne.CanvasObject
	gamesBtn          *widget.Button
	usersBtn          *widget.Button
	settingsBtn       *widget.Button

	OnGamesButtonPressed    func()
	OnUsersButtonPressed    func()
	OnSettingsButtonPressed func()
}

func NewSidebar(icon fyne.Resource, version string) *Sidebar {
	sidebar := &Sidebar{
		icon:    icon,
		version: version,
	}

	sidebar.gamesBtn = widget.NewButtonWithIcon("Games", fynetheme.FolderIcon(), nil)
	sidebar.usersBtn = widget.NewButtonWithIcon("Users", fynetheme.AccountIcon(), nil)
	sidebar.settingsBtn = widget.NewButtonWithIcon("Settings", fynetheme.SettingsIcon(), nil)

	sidebar.gamesBtn.OnTapped = func() {
		if sidebar.OnGamesButtonPressed != nil {
			sidebar.OnGamesButtonPressed()
		}
	}
	sidebar.usersBtn.OnTapped = func() {
		if sidebar.OnUsersButtonPressed != nil {
			sidebar.OnUsersButtonPressed()
		}
	}
	sidebar.settingsBtn.OnTapped = func() {
		if sidebar.OnSettingsButtonPressed != nil {
			sidebar.OnSettingsButtonPressed()
		}
	}

	sidebar.styleNavButton(sidebar.gamesBtn)
	sidebar.styleNavButton(sidebar.usersBtn)
	sidebar.styleNavButton(sidebar.settingsBtn)

	headerIcon := canvas.NewImageFromResource(icon)
	headerIcon.FillMode = canvas.ImageFillContain
	headerIcon.SetMinSize(fyne.NewSize(56, 56))
	sidebar.headerIcon = headerIcon

	headerText := canvas.NewText("Lanty", theme.ForegroundColor())
	headerText.TextSize = 22
	headerText.TextStyle = fyne.TextStyle{Bold: true}

	subtitleColor := color.RGBA{150, 150, 150, 255}
	subtitleText := canvas.NewText("For LAN Party", subtitleColor)
	subtitleText.TextSize = 12

	textContainer := container.NewVBox(
		headerText,
		subtitleText,
	)

	iconTextContainer := container.NewHBox(
		headerIcon,
		container.NewHBox(),
		textContainer,
	)
	sidebar.iconTextContainer = iconTextContainer

	sidebar.header = container.NewWithoutLayout(iconTextContainer)

	sidebar.navContainer = container.NewVBox(
		sidebar.gamesBtn,
		container.NewVBox(),
		sidebar.usersBtn,
		container.NewVBox(),
		sidebar.settingsBtn,
	)

	footerColor := color.RGBA{150, 150, 150, 255}
	versionText := canvas.NewText(version, footerColor)
	versionText.TextSize = 13
	versionText.Alignment = fyne.TextAlignTrailing

	sidebar.footer = container.NewWithoutLayout(versionText)

	sidebar.ExtendBaseWidget(sidebar)
	return sidebar
}

func (s *Sidebar) styleNavButton(btn *widget.Button) {
	btn.Importance = widget.LowImportance
	btn.Alignment = widget.ButtonAlignLeading
}

func (s *Sidebar) CreateRenderer() fyne.WidgetRenderer {
	return newSidebarRenderer(s)
}

type sidebarRenderer struct {
	sidebar           *Sidebar
	background        *canvas.Rectangle
	header            fyne.CanvasObject
	headerIcon        *canvas.Image
	iconTextContainer fyne.CanvasObject
	divider           *canvas.Rectangle
	navContainer      fyne.CanvasObject
	footer            fyne.CanvasObject
	objects           []fyne.CanvasObject
}

func newSidebarRenderer(sidebar *Sidebar) *sidebarRenderer {
	background := canvas.NewRectangle(theme.BackgroundColor2())
	background.CornerRadius = 0

	dividerColor := color.RGBA{100, 100, 100, 100}
	divider := canvas.NewRectangle(dividerColor)
	divider.CornerRadius = 0

	renderer := &sidebarRenderer{
		sidebar:           sidebar,
		background:        background,
		header:            sidebar.header,
		headerIcon:        sidebar.headerIcon,
		iconTextContainer: sidebar.iconTextContainer,
		divider:           divider,
		navContainer:      sidebar.navContainer,
		footer:            sidebar.footer,
		objects: []fyne.CanvasObject{
			background,
			sidebar.header,
			divider,
			sidebar.navContainer,
			sidebar.footer,
		},
	}

	return renderer
}

func (r *sidebarRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *sidebarRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.background.Move(fyne.NewPos(0, 0))

	headerPadding := theme.InnerPadding() * 1.5
	headerHeight := float32(80.0)

	iconSize := float32(56.0)
	if r.headerIcon != nil {
		r.headerIcon.Resize(fyne.NewSize(iconSize, iconSize))
	}

	dividerHeight := float32(1.0)
	verticalPadding := headerPadding * 1.2

	dividerTop := verticalPadding + headerHeight + verticalPadding

	midpoint := dividerTop / 2
	headerTop := midpoint - (headerHeight / 2)

	r.header.Resize(fyne.NewSize(size.Width-2*headerPadding, headerHeight))
	r.header.Move(fyne.NewPos(headerPadding, headerTop))

	if r.iconTextContainer != nil {
		contentSize := r.iconTextContainer.MinSize()
		contentY := (headerHeight - contentSize.Height) / 2
		r.iconTextContainer.Resize(contentSize)
		r.iconTextContainer.Move(fyne.NewPos(0, contentY))
	}

	dividerPadding := headerPadding * 0.5
	dividerWidth := size.Width - 2*dividerPadding
	r.divider.Resize(fyne.NewSize(dividerWidth, dividerHeight))
	r.divider.Move(fyne.NewPos(dividerPadding, dividerTop))

	footerHeight := float32(40.0)
	footerPadding := headerPadding
	footerTop := size.Height - footerHeight

	navTop := dividerTop + dividerHeight + 8
	navWidth := size.Width - 2*headerPadding
	navHeight := footerTop - navTop - headerPadding
	r.navContainer.Resize(fyne.NewSize(navWidth, navHeight))
	r.navContainer.Move(fyne.NewPos(headerPadding, navTop))

	footerWidth := size.Width - 2*footerPadding
	r.footer.Resize(fyne.NewSize(footerWidth, footerHeight))
	r.footer.Move(fyne.NewPos(footerPadding, footerTop))

	if footerContainer, ok := r.footer.(*fyne.Container); ok && len(footerContainer.Objects) > 0 {
		versionText := footerContainer.Objects[0]
		textSize := versionText.MinSize()
		textX := footerWidth - textSize.Width
		textY := (footerHeight - textSize.Height) / 2
		versionText.Resize(textSize)
		versionText.Move(fyne.NewPos(textX, textY))
	}
}

func (r *sidebarRenderer) MinSize() fyne.Size {
	headerPadding := theme.InnerPadding() * 1.5

	sidebarWidth := float32(220.0)
	buttonHeight := float32(40.0)
	navHeight := buttonHeight*3 + float32(16)
	footerHeight := float32(40.0)

	headerHeight := float32(80.0)
	totalHeight := headerHeight + navHeight + footerHeight + 3*headerPadding + float32(8)

	return fyne.NewSize(sidebarWidth, totalHeight)
}

func (r *sidebarRenderer) Refresh() {
	r.background.Refresh()
}

func (r *sidebarRenderer) Destroy() {}
