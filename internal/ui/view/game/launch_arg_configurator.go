package gameview

// import (
// 	"fmt"
// 	"maps"
// 	"slices"
// 	"strconv"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/dialog"
// 	"fyne.io/fyne/v2/layout"
// 	"fyne.io/fyne/v2/widget"
// 	gameapp "github.com/seternate/go-lanty-client/internal/application/game"
// 	gamecontroller "github.com/seternate/go-lanty-client/internal/ui/controller/game"
// 	model "github.com/seternate/go-lanty-client/internal/ui/model/game"
// )

// var _ gamecontroller.HostConfigAdapter = (*HostConfigAdapter)(nil)

// type HostConfigAdapter struct {
// 	window fyne.Window
// }

// func NewHostConfigAdapter(window fyne.Window) *HostConfigAdapter {
// 	return &HostConfigAdapter{window: window}
// }

// func (adapter *HostConfigAdapter) ShowHostConfig(gameName string, configFields []model.HostConfigField, onSubmit func(values []gameapp.ArgInput)) {
// 	var widgets []*widget.FormItem
// 	fieldWidgets := make([]func() (gameapp.ArgInput, bool, error), 0)

// 	for i := range configFields {
// 		f := &configFields[i]

// 		switch f.Type {
// 		case model.FieldFlag:
// 			enabledCheck := widget.NewCheck("", nil)
// 			enabledCheck.SetChecked(f.Enabled || f.Required)
// 			if f.Required {
// 				enabledCheck.Disable()
// 			}
// 			fieldWidgets = append(fieldWidgets, func() (gameapp.ArgInput, bool, error) {
// 				return gameapp.ArgInput{
// 					Name:     f.Name,
// 					Argument: f.Argument,
// 					Enabled:  enabledCheck.Checked,
// 					Value:    "",
// 				}, true, nil
// 			})
// 			widgets = append(widgets, widget.NewFormItem(f.Name, container.NewBorder(nil, nil, nil, nil, enabledCheck)))
// 		case model.FieldString:
// 			entry := widget.NewEntry()
// 			entry.Disable()
// 			enabledCheck := widget.NewCheck("", nil)
// 			enabledCheck.SetChecked(f.Enabled || f.Required)
// 			if f.Enabled || f.Required {
// 				entry.Enable()
// 			}
// 			if f.Required {
// 				enabledCheck.Disable()
// 			}
// 			enabledCheck.OnChanged = func(checked bool) {
// 				if checked {
// 					entry.Enable()
// 				} else {
// 					entry.Disable()
// 				}
// 			}
// 			entry.SetText(f.Value)
// 			fieldWidgets = append(fieldWidgets, func() (gameapp.ArgInput, bool, error) {
// 				if entry.Text == "" {
// 					return gameapp.ArgInput{
// 						Name:     f.Name,
// 						Argument: f.Argument,
// 						Enabled:  enabledCheck.Checked,
// 						Value:    "",
// 					}, entry.Validate() == nil, nil
// 				}
// 				return gameapp.ArgInput{
// 					Name:     f.Name,
// 					Argument: f.Argument,
// 					Enabled:  enabledCheck.Checked,
// 					Value:    entry.Text,
// 				}, entry.Validate() == nil, nil
// 			})
// 			widgets = append(widgets, widget.NewFormItem(f.Name, container.NewBorder(nil, nil, enabledCheck, nil, entry)))

// 		case model.FieldEnum:
// 			enumValues := slices.Collect(maps.Keys(f.EnumValues))
// 			selectBox := widget.NewSelect(enumValues, nil)
// 			selectBox.Disable()
// 			enabledCheck := widget.NewCheck("", nil)
// 			enabledCheck.SetChecked(f.Enabled || f.Required)
// 			if f.Enabled || f.Required {
// 				selectBox.Enable()
// 			}
// 			if f.Required {
// 				enabledCheck.Disable()
// 			}
// 			enabledCheck.OnChanged = func(checked bool) {
// 				if checked {
// 					selectBox.Enable()
// 				} else {
// 					selectBox.Disable()
// 				}
// 			}
// 			value := ""
// 			for l, v := range f.EnumValues {
// 				if f.Value == v {
// 					value = l
// 					break
// 				}
// 			}
// 			selectBox.SetSelected(value)
// 			fieldWidgets = append(fieldWidgets, func() (gameapp.ArgInput, bool, error) {
// 				return gameapp.ArgInput{
// 					Name:     f.Name,
// 					Argument: f.Argument,
// 					Enabled:  enabledCheck.Checked,
// 					Value:    f.EnumValues[selectBox.Selected],
// 				}, true, nil
// 			})
// 			widgets = append(widgets, widget.NewFormItem(f.Name, container.NewBorder(nil, nil, enabledCheck, nil, selectBox)))

// 		case model.FieldInt:
// 			entry := widget.NewEntry()
// 			entry.SetPlaceHolder("Enter number")
// 			entry.Disable()

// 			enabledCheck := widget.NewCheck("", nil)
// 			enabledCheck.SetChecked(f.Enabled || f.Required)
// 			if f.Enabled || f.Required {
// 				entry.Enable()
// 			}
// 			if f.Required {
// 				enabledCheck.Disable()
// 			}
// 			enabledCheck.OnChanged = func(checked bool) {
// 				if checked {
// 					entry.Enable()
// 				} else {
// 					entry.Disable()
// 				}
// 			}
// 			entry.SetText(f.Value)

// 			entry.Validator = func(s string) error {
// 				if !enabledCheck.Checked {
// 					return nil
// 				}

// 				if s == "" {
// 					return fmt.Errorf("required")
// 				}

// 				val, err := strconv.ParseInt(s, 10, 64)
// 				if err != nil {
// 					return fmt.Errorf("invalid number")
// 				}

// 				if f.MinInt != nil && val < *f.MinInt {
// 					return fmt.Errorf("min %d", *f.MinInt)
// 				}
// 				if f.MaxInt != nil && val > *f.MaxInt {
// 					return fmt.Errorf("max %d", *f.MaxInt)
// 				}

// 				return nil
// 			}

// 			fieldWidgets = append(fieldWidgets, func() (gameapp.ArgInput, bool, error) {
// 				if entry.Text == "" {
// 					return gameapp.ArgInput{
// 						Name:     f.Name,
// 						Argument: f.Argument,
// 						Enabled:  enabledCheck.Checked,
// 						Value:    "",
// 					}, entry.Validate() == nil, nil
// 				}
// 				return gameapp.ArgInput{
// 					Name:     f.Name,
// 					Argument: f.Argument,
// 					Enabled:  enabledCheck.Checked,
// 					Value:    entry.Text,
// 				}, entry.Validate() == nil, nil
// 			})
// 			widgets = append(widgets, widget.NewFormItem(f.Name, container.NewBorder(nil, nil, enabledCheck, nil, entry)))

// 		case model.FieldFloat:
// 			entry := widget.NewEntry()
// 			entry.SetPlaceHolder("Enter number")
// 			entry.Disable()

// 			enabledCheck := widget.NewCheck("", nil)
// 			enabledCheck.SetChecked(f.Enabled || f.Required)
// 			if f.Enabled || f.Required {
// 				entry.Enable()
// 			}
// 			enabledCheck.OnChanged = func(checked bool) {
// 				if checked {
// 					entry.Enable()
// 				} else {
// 					entry.Disable()
// 				}
// 			}

// 			entry.SetText(f.Value)

// 			entry.Validator = func(s string) error {
// 				if !enabledCheck.Checked {
// 					return nil
// 				}

// 				if s == "" {
// 					return fmt.Errorf("required")
// 				}

// 				val, err := strconv.ParseFloat(s, 64)
// 				if err != nil {
// 					return fmt.Errorf("invalid number")
// 				}

// 				if f.MinFloat != nil && val < *f.MinFloat {
// 					return fmt.Errorf("min %.2f", *f.MinFloat)
// 				}
// 				if f.MaxFloat != nil && val > *f.MaxFloat {
// 					return fmt.Errorf("max %.2f", *f.MaxFloat)
// 				}

// 				return nil
// 			}

// 			fieldWidgets = append(fieldWidgets, func() (gameapp.ArgInput, bool, error) {
// 				if entry.Text == "" {
// 					return gameapp.ArgInput{
// 						Name:     f.Name,
// 						Argument: f.Argument,
// 						Enabled:  enabledCheck.Checked,
// 						Value:    "",
// 					}, entry.Validate() == nil, nil
// 				}
// 				return gameapp.ArgInput{
// 					Name:     f.Name,
// 					Argument: f.Argument,
// 					Enabled:  enabledCheck.Checked,
// 					Value:    entry.Text,
// 				}, entry.Validate() == nil, nil
// 			})
// 			widgets = append(widgets, widget.NewFormItem(f.Name, container.NewBorder(nil, nil, enabledCheck, nil, entry)))
// 		}
// 	}

// 	form := widget.NewForm(widgets...)

// 	cancelBtn := widget.NewButton("Cancel", nil)

// 	submitBtn := widget.NewButton("Start Server", nil)

// 	buttonBar := container.NewHBox(
// 		layout.NewSpacer(),
// 		cancelBtn,
// 		submitBtn,
// 	)

// 	content := container.NewVBox(
// 		form,
// 		buttonBar,
// 	)

// 	d := dialog.NewCustomWithoutButtons(
// 		fmt.Sprintf("Configure hosting for '%s'", gameName),
// 		content,
// 		adapter.window,
// 	)
// 	cancelBtn.OnTapped = func() {
// 		d.Hide()
// 	}
// 	submitBtn.OnTapped = func() {
// 		values := []gameapp.ArgInput{}
// 		for _, getter := range fieldWidgets {
// 			argInput, valid, err := getter()
// 			if !valid {
// 				return
// 			}
// 			if err != nil {
// 				continue
// 			}
// 			values = append(values, argInput)
// 		}
// 		onSubmit(values)
// 		d.Hide()
// 	}

// 	d.Resize(fyne.NewSize(800, 600))
// 	d.Show()
// }
