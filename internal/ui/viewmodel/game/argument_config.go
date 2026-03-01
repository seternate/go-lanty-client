package gameviewmodel

// import (
// 	"fmt"

// 	"fyne.io/fyne/v2/data/binding"
// 	"github.com/seternate/go-lanty-client/internal/game"
// )

// type ArgumentConfigForm struct {
// 	Title    string
// 	Fields   []*ArgumentField
// 	OnSubmit func()
// }

// func NewArgumentConfigFormFromLaunchParams(title string, params []game.LaunchParam) (*ArgumentConfigForm, error) {
// 	vm := &ArgumentConfigForm{
// 		Title:    title,
// 		Fields:   make([]*ArgumentField, len(params)),
// 		OnSubmit: nil,
// 	}

// 	for _, param := range params {
// 		field, err := NewArgumentFieldFromLaunchParam(param)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create argument field: %w", err)
// 		}
// 		vm.Fields = append(vm.Fields, field)
// 	}

// 	return vm, nil
// }

// type FieldType int

// const (
// 	FieldFlag FieldType = iota
// 	FieldString
// 	FieldInt
// 	FieldFloat
// 	FieldEnum
// )

// type ArgumentField struct {
// 	Type        FieldType
// 	Name        string
// 	Description string
// 	Argument    string
// 	Required    bool
// 	EnumValues  map[string]string
// 	MinInt      *int64
// 	MaxInt      *int64
// 	MinFloat    *float64
// 	MaxFloat    *float64

// 	Enabled binding.Bool
// 	Value   binding.String
// }

// func NewArgumentFieldFromLaunchParam(param game.LaunchParam) (*ArgumentField, error) {
// 	var fieldType FieldType

// 	switch param.Type() {
// 	case game.LaunchParamFlag:
// 		fieldType = FieldFlag
// 	case game.LaunchParamString:
// 		fieldType = FieldString
// 	case game.LaunchParamInt:
// 		fieldType = FieldInt
// 	case game.LaunchParamFloat:
// 		fieldType = FieldFloat
// 	case game.LaunchParamEnum:
// 		fieldType = FieldEnum
// 	default:
// 		return nil, fmt.Errorf("invalid param type: %d", param.Type())
// 	}

// 	enumValues := map[string]string{}
// 	for _, option := range param.EnumValues() {
// 		enumValues[option.Label] = option.Value
// 	}

// 	enabled := binding.NewBool()
// 	enabled.Set(param.Enabled())

// 	value := binding.NewString()
// 	value.Set(param.Value())

// 	return &ArgumentField{
// 		Name:        param.Name(),
// 		Description: param.Description(),
// 		Argument:    param.Argument(),
// 		Required:    param.Required(),
// 		Type:        fieldType,
// 		EnumValues:  enumValues,
// 		MinInt:      param.MinInt(),
// 		MaxInt:      param.MaxInt(),
// 		MinFloat:    param.MinFloat(),
// 		MaxFloat:    param.MaxFloat(),
// 		Enabled:     enabled,
// 		Value:       value,
// 	}, nil
// }
