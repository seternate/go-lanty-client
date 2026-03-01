package userview

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	usermodel "github.com/seternate/go-lanty-client/internal/ui/model/user"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
)

type UserView struct {
	widget.BaseWidget

	users          *usermodel.UserListModel
	tiles          map[string]*UserTile
	tilesContainer *fyne.Container
}

func NewUserView(users *usermodel.UserListModel) *UserView {
	view := &UserView{
		users:          users,
		tiles:          make(map[string]*UserTile),
		tilesContainer: container.NewVBox(),
	}

	view.RefreshUserTiles()

	users.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)
	return view
}

func (view *UserView) RefreshUserTiles() {
	ipsOrdered := view.users.GetUserIPAddressesOrderedByName()
	users := view.users.GetUsers()

	currentIps := make(map[string]bool)
	for _, user := range users {
		currentIps[user.IPAddress] = true
	}

	for ip := range view.tiles {
		if !currentIps[ip] {
			delete(view.tiles, ip)
		}
	}

	for _, user := range users {
		if _, found := view.tiles[user.IPAddress]; found {
			continue
		}
		view.tiles[user.IPAddress] = NewUserTile(user)
	}

	view.tilesContainer.RemoveAll()
	for _, ip := range ipsOrdered {
		tile, found := view.tiles[ip]
		if !found {
			continue
		}
		view.tilesContainer.Add(tile)
	}
}

func (view *UserView) Refresh() {
	view.RefreshUserTiles()
	view.BaseWidget.Refresh()
}

func (view *UserView) CreateRenderer() fyne.WidgetRenderer {
	return newUserRenderer(view)
}

type userRenderer struct {
	view            *UserView
	headerContainer *fyne.Container
	headerText      *canvas.Text
	statsText       *canvas.Text
	scrollContainer *container.Scroll
	tiles           []*UserTile
	objects         []fyne.CanvasObject
}

func newUserRenderer(view *UserView) *userRenderer {
	headerText := canvas.NewText("Users", color.White)
	headerText.TextSize = 28
	headerText.TextStyle = fyne.TextStyle{Bold: true}

	statsText := canvas.NewText("", color.RGBA{180, 180, 180, 255})
	statsText.TextSize = theme.TextSize()

	headerPadding := theme.InnerPadding() * 2
	paddingTop := canvas.NewRectangle(color.Transparent)
	paddingTop.SetMinSize(fyne.NewSize(0, headerPadding))
	paddingBottom := canvas.NewRectangle(color.Transparent)
	paddingBottom.SetMinSize(fyne.NewSize(0, headerPadding))

	paddingLeft := canvas.NewRectangle(color.Transparent)
	paddingLeft.SetMinSize(fyne.NewSize(headerPadding, 0))
	paddingRight := canvas.NewRectangle(color.Transparent)
	paddingRight.SetMinSize(fyne.NewSize(headerPadding, 0))

	headerLeft := container.NewBorder(paddingTop, paddingBottom, paddingLeft, nil, container.NewWithoutLayout(headerText))
	headerRight := container.NewBorder(paddingTop, paddingBottom, nil, paddingRight, container.NewWithoutLayout(statsText))

	headerContainer := container.NewBorder(nil, nil, headerLeft, headerRight)

	scrollContainer := container.NewScroll(view.tilesContainer)
	scrollContainer.SetMinSize(fyne.NewSize(0, 0))

	mainContainer := container.NewBorder(headerContainer, nil, nil, nil, scrollContainer)

	renderer := &userRenderer{
		view:            view,
		headerContainer: headerContainer,
		headerText:      headerText,
		statsText:       statsText,
		scrollContainer: scrollContainer,
		objects: []fyne.CanvasObject{
			mainContainer,
		},
	}

	totalUsers := len(renderer.view.users.GetUsers())
	renderer.statsText.Text = fmt.Sprintf("%d available", totalUsers)

	return renderer
}

func (r *userRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *userRenderer) Layout(size fyne.Size) {
	if len(r.objects) > 0 {
		r.objects[0].Resize(size)
		r.objects[0].Move(fyne.NewPos(0, 0))
	}
}

func (r *userRenderer) MinSize() fyne.Size {
	padding := theme.InnerPadding() * 2
	headerHeight := r.headerText.MinSize().Height + padding*2
	return fyne.NewSize(0, headerHeight)
}

func (r *userRenderer) Refresh() {
	totalUsers := len(r.view.users.GetUsers())
	r.statsText.Text = fmt.Sprintf("%d available", totalUsers)

	r.headerText.Refresh()
	r.statsText.Refresh()

	for _, tile := range r.tiles {
		tile.Refresh()
	}
}

func (r *userRenderer) Destroy() {}
