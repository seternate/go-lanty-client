package gameviewmodel

import (
	"strconv"

	"fyne.io/fyne/v2/data/binding"
	"github.com/seternate/go-lanty-client/internal/game"
)

type LaunchArgumentTile struct {
	name        string
	description string
	argument    string
	enumValues  map[string]string
	enableTile  bool
	required    bool
	minInt      *int64
	maxInt      *int64
	minFloat    *float64
	maxFloat    *float64

	showList  bool
	showEntry bool

	Enabled             binding.Bool
	Value               binding.String
	ValidationError     binding.String
	ShowValidationError binding.Bool
}

func NewLaunchArgumentTileFromModel(model game.LaunchParam) (*LaunchArgumentTile, error) {
	vm := &LaunchArgumentTile{
		name:                model.Name(),
		description:         model.Description(),
		argument:            model.Argument(),
		enumValues:          make(map[string]string, len(model.EnumValues())),
		enableTile:          !model.Required(),
		required:            model.Required(),
		minInt:              model.MinInt(),
		maxInt:              model.MaxInt(),
		minFloat:            model.MinFloat(),
		maxFloat:            model.MaxFloat(),
		showList:            model.Type() == game.LaunchParamEnum,
		showEntry:           model.Type() == game.LaunchParamString || model.Type() == game.LaunchParamInt || model.Type() == game.LaunchParamFloat,
		Enabled:             binding.NewBool(),
		Value:               binding.NewString(),
		ValidationError:     binding.NewString(),
		ShowValidationError: binding.NewBool(),
	}

	for _, enumValue := range model.EnumValues() {
		vm.enumValues[enumValue.Label] = enumValue.Value
	}

	vm.Enabled.Set(model.Enabled() || model.Required())
	vm.Value.Set(model.Value())

	vm.Value.AddListener(binding.NewDataListener(func() {
		vm.validateNumeric()
	}))
	vm.validateNumeric()

	return vm, nil
}

func (vm *LaunchArgumentTile) AddChangeListener(fn func()) {
	l := binding.NewDataListener(fn)

	vm.Enabled.AddListener(l)
	vm.Value.AddListener(l)
	vm.ValidationError.AddListener(l)
	vm.ShowValidationError.AddListener(l)
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

func (vm *LaunchArgumentTile) GetValidationError() string {
	s, err := vm.ValidationError.Get()
	if err != nil {
		return ""
	}
	return s
}

func (vm *LaunchArgumentTile) GetShowValidationError() bool {
	b, err := vm.ShowValidationError.Get()
	if err != nil {
		return false
	}
	return b
}

func (vm *LaunchArgumentTile) SetValidationError(s string) {
	vm.ValidationError.Set(s)
	vm.ShowValidationError.Set(s != "")
}

func (vm *LaunchArgumentTile) validateNumeric() {
	if !vm.ShowEntry() || (vm.minInt == nil && vm.maxInt == nil && vm.minFloat == nil && vm.maxFloat == nil && !vm.required) {
		vm.SetValidationError("")
		return
	}

	valueStr, err := vm.Value.Get()
	if err != nil {
		vm.SetValidationError("")
		return
	}

	if valueStr == "" {
		if vm.required {
			vm.SetValidationError("Required field can not be empty")
		} else {
			vm.SetValidationError("")
		}
		return
	}

	if vm.minInt != nil || vm.maxInt != nil {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			vm.SetValidationError("Must be a whole number")
			return
		}

		if vm.minInt != nil && vm.maxInt != nil {
			if value < *vm.minInt || value > *vm.maxInt {
				vm.SetValidationError("Must be between " + strconv.FormatInt(*vm.minInt, 10) + " and " + strconv.FormatInt(*vm.maxInt, 10))
				return
			}
		} else if vm.minInt != nil {
			if value < *vm.minInt {
				vm.SetValidationError("Must be at least " + strconv.FormatInt(*vm.minInt, 10))
				return
			}
		} else if vm.maxInt != nil {
			if value > *vm.maxInt {
				vm.SetValidationError("Must be at most " + strconv.FormatInt(*vm.maxInt, 10))
				return
			}
		}

		vm.SetValidationError("")
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		vm.SetValidationError("Must be a number")
		return
	}

	if vm.minFloat != nil && vm.maxFloat != nil {
		if value < *vm.minFloat || value > *vm.maxFloat {
			vm.SetValidationError("Must be between " + strconv.FormatFloat(*vm.minFloat, 'f', -1, 64) + " and " + strconv.FormatFloat(*vm.maxFloat, 'f', -1, 64))
			return
		}
	} else if vm.minFloat != nil {
		if value < *vm.minFloat {
			vm.SetValidationError("Must be at least " + strconv.FormatFloat(*vm.minFloat, 'f', -1, 64))
			return
		}
	} else if vm.maxFloat != nil {
		if value > *vm.maxFloat {
			vm.SetValidationError("Must be at most " + strconv.FormatFloat(*vm.maxFloat, 'f', -1, 64))
			return
		}
	}

	vm.SetValidationError("")
}
