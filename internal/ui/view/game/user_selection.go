package gameview

// import (
// 	"fmt"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/dialog"
// 	"fyne.io/fyne/v2/widget"
// 	"github.com/seternate/go-lanty-client/internal/domain/user"
// 	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
// 	model "github.com/seternate/go-lanty-client/internal/ui/model/user"
// 	userview "github.com/seternate/go-lanty-client/internal/ui/view/user"
// )

// type UserProvider interface {
// 	ListUsers() ([]user.UserCatalogItem, error)
// }

// var _ gamecontroller.JoinUserAdapter = (*JoinUserAdapter)(nil)

// type JoinUserAdapter struct {
// 	window       fyne.Window
// 	userProvider UserProvider
// }

// func NewJoinUserAdapter(window fyne.Window, userProvider UserProvider) *JoinUserAdapter {
// 	return &JoinUserAdapter{window: window, userProvider: userProvider}
// }

// func (adapter *JoinUserAdapter) OpenUserSelection(gameName string, onUserSelected func(ipAddress string)) {
// 	users, err := adapter.userProvider.ListUsers()
// 	if err != nil {
// 		return
// 	}

// 	userList := model.NewUserListModel()
// 	for _, user := range users {
// 		tile := model.NewUserTileModel(user.IP)
// 		tile.SetName(user.Name)
// 		userList.AddUser(tile)
// 	}

// 	sortedUsers := userList.GetUserNamesOrdered()

// 	selectedIndex := -1

// 	list := widget.NewList(
// 		func() int {
// 			return len(sortedUsers)
// 		},
// 		func() fyne.CanvasObject {
// 			canvasObject := userview.NewUserTileWithoutModel("", "")
// 			canvasObject.HideCopyButton()
// 			return canvasObject
// 			// return widget.NewLabel("")
// 		},
// 		func(i widget.ListItemID, o fyne.CanvasObject) {
// 			canvasObject := o.(*userview.UserTile)
// 			canvasObject.SetIPAddress(userList.GetUserByName(sortedUsers[i]).IPAddress)
// 			canvasObject.SetName(sortedUsers[i])
// 			// o.(*widget.Label).SetText(sortedUsers[i])
// 		},
// 	)

// 	list.OnSelected = func(id widget.ListItemID) {
// 		selectedIndex = id
// 	}

// 	content := container.NewBorder(
// 		nil,
// 		nil,
// 		nil,
// 		nil,
// 		list,
// 	)

// 	d := dialog.NewCustomConfirm(
// 		fmt.Sprintf("Select a user to join '%s'", gameName),
// 		"Join",
// 		"Cancel",
// 		content,
// 		func(confirmed bool) {
// 			if !confirmed {
// 				return
// 			}

// 			if selectedIndex < 0 {
// 				dialog.ShowInformation(
// 					"No selection",
// 					"Please select a user before confirming.",
// 					adapter.window,
// 				)
// 				return
// 			}

// 			selectedUsername := sortedUsers[selectedIndex]
// 			user := userList.GetUserByName(selectedUsername)
// 			onUserSelected(user.IPAddress)
// 		},
// 		adapter.window,
// 	)
// 	d.Resize(fyne.NewSize(800, 600))
// 	d.Show()
// }
