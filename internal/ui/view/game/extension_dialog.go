package gameview

// import (
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// 	model "github.com/seternate/go-lanty-client/internal/ui/model/game"
// 	"github.com/seternate/go-lanty-client/internal/ui/theme"
// )

// func NewExtensionDialogContent(
// 	gameName string,
// 	extensions *model.ExtensionListModel,
// 	onUpload func(),
// 	onDownload func(ext *model.ExtensionModel),
// ) fyne.CanvasObject {
// 	title := widget.NewLabel("Extensions for " + gameName)
// 	title.TextStyle = fyne.TextStyle{Bold: true}
// 	title.TextStyle.Monospace = false

// 	list := widget.NewList(
// 		func() int {
// 			return len(extensions.GetExtensions())
// 		},
// 		func() fyne.CanvasObject {
// 			return NewExtensionItemWithoutModel(onDownload)
// 		},
// 		func(id widget.ListItemID, o fyne.CanvasObject) {
// 			extList := extensions.GetExtensions()
// 			if id >= len(extList) {
// 				return
// 			}
// 			item := o.(*ExtensionItem)
// 			item.SetExtension(extList[id])
// 		},
// 	)

// 	extensions.AddChangeListener(func() {
// 		list.Refresh()
// 	})

// 	uploadBtn := widget.NewButtonWithIcon("Upload", theme.UploadIcon(), onUpload)
// 	uploadBtn.Importance = widget.HighImportance

// 	scroll := container.NewScroll(list)
// 	scroll.SetMinSize(fyne.NewSize(400, 300))

// 	content := container.NewBorder(
// 		title,
// 		container.NewHBox(
// 			uploadBtn,
// 		),
// 		nil,
// 		nil,
// 		scroll,
// 	)

// 	return content
// }
