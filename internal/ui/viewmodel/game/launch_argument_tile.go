package gameviewmodel

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/seternate/go-lanty-client/internal/game"
)

type LaunchArgumentTile struct {
	name        string
	description string
	argument    string
	enumValues  map[string]string
	enableTile  bool
	// minInt      *int64
	// maxInt      *int64
	// minFloat    *float64
	// maxFloat    *float64
	showList  bool
	showEntry bool

	Enabled binding.Bool
	Value   binding.String
}

func NewLaunchArgumentTileFromModel(model game.LaunchParam) (*LaunchArgumentTile, error) {
	vm := &LaunchArgumentTile{
		name:        model.Name(),
		description: model.Description(),
		argument:    model.Argument(),
		enumValues:  make(map[string]string, len(model.EnumValues())),
		enableTile:  !model.Required(),
		// minInt:      model.MinInt(),
		// maxInt:      model.MaxInt(),
		// minFloat:    model.MinFloat(),
		// maxFloat:    model.MaxFloat(),
		showList:  model.Type() == game.LaunchParamEnum,
		showEntry: model.Type() == game.LaunchParamString || model.Type() == game.LaunchParamInt || model.Type() == game.LaunchParamFloat,
		Enabled:   binding.NewBool(),
		Value:     binding.NewString(),
	}

	for _, enumValue := range model.EnumValues() {
		vm.enumValues[enumValue.Label] = enumValue.Value
	}

	vm.Enabled.Set(model.Enabled() || model.Required())
	vm.Value.Set(model.Value())

	return vm, nil
}

func (vm *LaunchArgumentTile) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)

	vm.Enabled.AddListener(l)
	vm.Value.AddListener(l)
}

func (vm *LaunchArgumentTile) GetName() string {
	return vm.name
}

func (vm *LaunchArgumentTile) GetDescription() string {
	return vm.description
}

func (vm *LaunchArgumentTile) GetArgument() string {
	return vm.argument
}

func (vm *LaunchArgumentTile) GetEnumValues() map[string]string {
	return vm.enumValues
}

func (vm *LaunchArgumentTile) ShowList() bool {
	return vm.showList
}

func (vm *LaunchArgumentTile) ShowEntry() bool {
	return vm.showEntry
}

func (vm *LaunchArgumentTile) GetEnabled() bool {
	enabled, err := vm.Enabled.Get()
	if err != nil {
		return false
	}
	return enabled
}

func (vm *LaunchArgumentTile) GetValueForList() string {
	value, err := vm.Value.Get()
	if err != nil {
		return ""
	}

	for label, v := range vm.enumValues {
		if v == value {
			return label
		}
	}

	return ""
}

func (vm *LaunchArgumentTile) SetValueFromList(value string) {
	vm.Value.Set(vm.enumValues[value])
}

func (vm *LaunchArgumentTile) TileEnabled() bool {
	return vm.enableTile
}
