package gameview

// import (
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/dialog"
// 	"fyne.io/fyne/v2/storage"
// 	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
// 	model "github.com/seternate/go-lanty-client/internal/ui/model/game"
// 	gameview "github.com/seternate/go-lanty-client/internal/ui/view/game"
// )

// var _ gamecontroller.ExtensionsAdapter = (*ExtensionsAdapter)(nil)

// type ExtensionsAdapter struct {
// 	window fyne.Window
// }

// func NewExtensionsAdapter(window fyne.Window) *ExtensionsAdapter {
// 	return &ExtensionsAdapter{window: window}
// }

// func (a *ExtensionsAdapter) ShowExtensions(
// 	gameName, _, installPath string,
// 	extensions *model.ExtensionListModel,
// 	onUpload func(localPath string),
// 	onDownload func(ext *model.ExtensionModel),
// ) {
// 	onUploadClick := func() {
// 		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
// 			if err != nil || reader == nil {
// 				return
// 			}
// 			path := reader.URI().Path()
// 			_ = reader.Close()
// 			onUpload(path)
// 			dialog.ShowInformation("Upload finished", "Extension uploaded successfully", a.window)
// 		}, a.window)
// 		if installPath != "" {
// 			uri := storage.NewFileURI(installPath)
// 			if listable, err := storage.ListerForURI(uri); err == nil {
// 				fileDialog.SetLocation(listable)
// 			}
// 		}
// 		fileDialog.Resize(fyne.NewSize(800, 600))
// 		fileDialog.Show()
// 	}

// 	content := gameview.NewExtensionDialogContent(gameName, extensions, onUploadClick, onDownload)

// 	d := dialog.NewCustomConfirm("Extensions for "+gameName, "Close", "Cancel", content, func(_ bool) {}, a.window)
// 	d.Resize(fyne.NewSize(500, 400))
// 	d.Show()
// }
