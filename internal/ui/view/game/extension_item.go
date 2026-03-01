package gameview

// import (
// 	"image/color"
// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/canvas"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// 	model "github.com/seternate/go-lanty-client/internal/ui/model/game"
// 	"github.com/seternate/go-lanty-client/internal/ui/theme"
// )

// const doubleTapWindow = 400 * time.Millisecond

// type ExtensionItem struct {
// 	widget.BaseWidget

// 	ext            *model.ExtensionModel
// 	onDoubleTapped func(ext *model.ExtensionModel)

// 	lastTapTime time.Time

// 	nameLabel   *canvas.Text
// 	fileLabel   *canvas.Text
// 	sizeLabel   *canvas.Text
// 	statusIcon  *widget.Icon
// 	content     *fyne.Container
// }

// func NewExtensionItem(ext *model.ExtensionModel, onDoubleTapped func(ext *model.ExtensionModel)) *ExtensionItem {
// 	item := newExtensionItemWithoutModel(onDoubleTapped)
// 	item.SetExtension(ext)
// 	return item
// }

// func NewExtensionItemWithoutModel(onDoubleTapped func(ext *model.ExtensionModel)) *ExtensionItem {
// 	return newExtensionItemWithoutModel(onDoubleTapped)
// }

// func newExtensionItemWithoutModel(onDoubleTapped func(ext *model.ExtensionModel)) *ExtensionItem {
// 	item := &ExtensionItem{
// 		onDoubleTapped: onDoubleTapped,
// 	}

// 	item.nameLabel = canvas.NewText("", color.White)
// 	item.nameLabel.TextSize = theme.TextSize()
// 	item.nameLabel.TextStyle = fyne.TextStyle{Bold: true}

// 	item.fileLabel = canvas.NewText("", color.RGBA{150, 150, 150, 255})
// 	item.fileLabel.TextSize = theme.TextSize() * 0.85

// 	item.sizeLabel = canvas.NewText("", color.RGBA{150, 150, 150, 255})
// 	item.sizeLabel.TextSize = theme.TextSize() * 0.8

// 	item.statusIcon = widget.NewIcon(theme.ExtensionUpToDateIcon())
// 	item.statusIcon.Hide()

// 	item.content = container.NewBorder(
// 		nil,
// 		nil,
// 		nil,
// 		item.statusIcon,
// 		container.NewVBox(
// 			item.nameLabel,
// 			container.NewHBox(item.fileLabel, item.sizeLabel),
// 		),
// 	)

// 	item.ExtendBaseWidget(item)
// 	return item
// }

// func (item *ExtensionItem) SetExtension(ext *model.ExtensionModel) {
// 	if item.ext != nil {
// 		// Remove old listener - we'd need to track it; for simplicity, we rely on
// 		// the same ext being passed and listener already registered.
// 	}
// 	item.ext = ext
// 	if ext != nil {
// 		ext.AddChangeListener(func() {
// 			item.Refresh()
// 		})
// 	}
// 	item.Refresh()
// }

// func (item *ExtensionItem) Tapped(*fyne.PointEvent) {
// 	now := time.Now()
// 	if !item.lastTapTime.IsZero() && now.Sub(item.lastTapTime) < doubleTapWindow && item.onDoubleTapped != nil && item.ext != nil {
// 		item.onDoubleTapped(item.ext)
// 		item.lastTapTime = time.Time{}
// 		return
// 	}
// 	item.lastTapTime = now
// }

// func (item *ExtensionItem) CreateRenderer() fyne.WidgetRenderer {
// 	background := canvas.NewRectangle(theme.BackgroundColor2())
// 	background.CornerRadius = theme.SmallCornerRadius()

// 	return widget.NewSimpleRenderer(container.NewStack(background, item.content))
// }

// func (item *ExtensionItem) Refresh() {
// 	if item.ext == nil {
// 		return
// 	}
// 	name, _ := item.ext.Name.Get()
// 	size, _ := item.ext.Size.Get()
// 	downloaded, _ := item.ext.Downloaded.Get()

// 	item.nameLabel.Text = name
// 	item.fileLabel.Text = item.ext.Filename
// 	if size != "" {
// 		item.sizeLabel.Text = "(" + size + ")"
// 		item.sizeLabel.Show()
// 	} else {
// 		item.sizeLabel.Text = ""
// 		item.sizeLabel.Hide()
// 	}

// 	if downloaded {
// 		item.statusIcon.SetResource(theme.ExtensionUpToDateIcon())
// 		item.statusIcon.Show()
// 	} else {
// 		item.statusIcon.Hide()
// 	}

// 	item.nameLabel.Refresh()
// 	item.fileLabel.Refresh()
// 	item.sizeLabel.Refresh()
// 	item.statusIcon.Refresh()
// }
