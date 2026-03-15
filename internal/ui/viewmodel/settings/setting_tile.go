package settingsviewmodel

import (
	"fyne.io/fyne/v2/data/binding"
)

type FolderPicker interface {
	ShowFolderPickerDialog(location string, callback func(path string))
}

type SettingTile struct {
	label       string
	description string

	Value            binding.String
	showFolderPicker bool

	folderPicker FolderPicker
	onSave       func()
}

func NewSettingTile(label string, description string, value binding.String, showFolderPicker bool, folderPicker FolderPicker, onSave func()) *SettingTile {
	return &SettingTile{
		label:            label,
		description:      description,
		Value:            value,
		showFolderPicker: showFolderPicker,
		folderPicker:     folderPicker,
		onSave:           onSave,
	}
}

func (vm *SettingTile) OpenFolderPicker() {
	if vm.showFolderPicker && vm.folderPicker != nil {
		vm.folderPicker.ShowFolderPickerDialog(vm.GetValue(), func(path string) {
			vm.Value.Set(path)
			vm.Save()
		})
	}
}

func (vm *SettingTile) Save() {
	if vm.onSave != nil {
		vm.onSave()
	}
}

func (vm *SettingTile) GetLabel() string {
	return vm.label
}

func (vm *SettingTile) GetDescription() string {
	return vm.description
}

func (vm *SettingTile) GetValue() string {
	value, err := vm.Value.Get()
	if err != nil {
		return ""
	}

	return value
}

func (vm *SettingTile) GetShowFolderPicker() bool {
	return vm.showFolderPicker
}
