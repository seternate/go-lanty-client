package view

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/seternate/go-lanty-client/internal/ui/theme"
	"github.com/seternate/go-lanty-client/internal/ui/viewmodel"
)

type HeaderScrollScreen struct {
	widget.BaseWidget

	vm *viewmodel.HeaderScrollScreen

	scroll       *container.Scroll
	submitButton *widget.Button
	cancelButton *widget.Button
}

func NewHeaderScrollScreen(vm *viewmodel.HeaderScrollScreen, content fyne.CanvasObject) *HeaderScrollScreen {
	view := &HeaderScrollScreen{
		vm:           vm,
		scroll:       container.NewScroll(content),
		submitButton: widget.NewButton("OK", vm.CallOnSubmit),
		cancelButton: widget.NewButton("Cancel", vm.CallOnCancel),
	}

	view.submitButton.Importance = widget.HighImportance
	view.cancelButton.Importance = widget.HighImportance

	view.vm.AddChangeListener(func() {
		view.Refresh()
	})

	view.ExtendBaseWidget(view)

	return view
}

func (view *HeaderScrollScreen) Refresh() {
	if view.vm.HasOnSubmit() {
		view.submitButton.Show()
	} else {
		view.submitButton.Hide()
	}
	if view.vm.HasOnCancel() {
		view.cancelButton.Show()
	} else {
		view.cancelButton.Hide()
	}

	view.scroll.Refresh()
	view.submitButton.Refresh()
	view.cancelButton.Refresh()

	view.BaseWidget.Refresh()
}

func (view *HeaderScrollScreen) CreateRenderer() fyne.WidgetRenderer {
	return newHeaderScrollScreenRenderer(view)
}

type headerScrollScreenRenderer struct {
	view *HeaderScrollScreen

	background       *canvas.Rectangle
	buttonBackground *canvas.Rectangle

	content          *fyne.Container
	headerContainer  *fyne.Container
	buttonsContainer *fyne.Container
	titleText        *canvas.Text
	infoText         *canvas.Text
}

func newHeaderScrollScreenRenderer(view *HeaderScrollScreen) *headerScrollScreenRenderer {
	renderer := &headerScrollScreenRenderer{
		view: view,
	}

	renderer.background = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.background.CornerRadius = theme.CornerRadius()
	renderer.buttonBackground = canvas.NewRectangle(theme.BackgroundColor2())
	renderer.buttonBackground.CornerRadius = theme.CornerRadius()

	renderer.titleText = canvas.NewText(view.vm.GetTitle(), color.White)
	renderer.titleText.TextSize = 28
	renderer.titleText.TextStyle = fyne.TextStyle{Bold: true}

	renderer.infoText = canvas.NewText(view.vm.GetInfo(), color.RGBA{180, 180, 180, 255})
	renderer.infoText.TextSize = theme.TextSize()

	headerPadding := theme.InnerPadding() * 2
	paddingTop := canvas.NewRectangle(color.Transparent)
	paddingTop.SetMinSize(fyne.NewSize(0, headerPadding))
	paddingBottom := canvas.NewRectangle(color.Transparent)
	paddingBottom.SetMinSize(fyne.NewSize(0, headerPadding))
	paddingLeft := canvas.NewRectangle(color.Transparent)
	paddingLeft.SetMinSize(fyne.NewSize(headerPadding, 0))
	paddingRight := canvas.NewRectangle(color.Transparent)
	paddingRight.SetMinSize(fyne.NewSize(headerPadding, 0))
	headerLeft := container.NewBorder(paddingTop, paddingBottom, paddingLeft, nil, container.NewWithoutLayout(renderer.titleText))
	headerRight := container.NewBorder(paddingTop, paddingBottom, nil, paddingRight, container.NewWithoutLayout(renderer.infoText))

	renderer.headerContainer = container.NewBorder(nil, nil, headerLeft, headerRight)

	paddingTopButton := canvas.NewRectangle(color.Transparent)
	paddingTopButton.SetMinSize(fyne.NewSize(0, theme.InnerPadding()*0.5))
	paddingBottomButton := canvas.NewRectangle(color.Transparent)
	paddingBottomButton.SetMinSize(fyne.NewSize(0, theme.InnerPadding()*0.5))
	paddingLeftButton := canvas.NewRectangle(color.Transparent)
	paddingLeftButton.SetMinSize(fyne.NewSize(theme.InnerPadding(), 0))
	paddingRightButton := canvas.NewRectangle(color.Transparent)
	paddingRightButton.SetMinSize(fyne.NewSize(theme.InnerPadding(), 0))

	renderer.buttonsContainer = container.NewBorder(paddingTopButton, paddingBottomButton, paddingLeftButton, paddingRightButton, container.NewHBox(layout.NewSpacer(), renderer.view.submitButton, renderer.view.cancelButton))
	renderer.content = container.NewBorder(renderer.headerContainer, renderer.buttonsContainer, nil, nil, renderer.view.scroll)

	return renderer
}

func (renderer *headerScrollScreenRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		renderer.background,
		renderer.buttonBackground,
		renderer.content,
	}
}

func (renderer *headerScrollScreenRenderer) Layout(size fyne.Size) {
	renderer.content.Resize(size)
	renderer.content.Move(fyne.NewPos(0, 0))

	renderer.background.Resize(fyne.NewSize(renderer.headerContainer.Size().Width, renderer.headerContainer.Size().Height+theme.CornerRadius()))
	renderer.background.Move(fyne.NewPos(0, -theme.CornerRadius()))

	renderer.buttonBackground.Resize(fyne.NewSize(renderer.buttonsContainer.Size().Width, renderer.buttonsContainer.Size().Height+theme.CornerRadius()))
	renderer.buttonBackground.Move(fyne.NewPos(0, size.Height-renderer.buttonsContainer.Size().Height))
}

func (renderer *headerScrollScreenRenderer) MinSize() fyne.Size {
	return renderer.content.MinSize()
}

func (renderer *headerScrollScreenRenderer) Refresh() {
	renderer.titleText.Text = renderer.view.vm.GetTitle()
	renderer.infoText.Text = renderer.view.vm.GetInfo()

	if renderer.view.vm.HasOnSubmit() || renderer.view.vm.HasOnCancel() {
		renderer.buttonsContainer.Show()
	} else {
		renderer.buttonsContainer.Hide()
	}

	renderer.background.Refresh()

	renderer.titleText.Refresh()
	renderer.infoText.Refresh()

	renderer.headerContainer.Refresh()
	renderer.buttonsContainer.Refresh()
	renderer.content.Refresh()
}

func (renderer *headerScrollScreenRenderer) Destroy() {}
