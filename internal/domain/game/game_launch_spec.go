package game

import (
	"fmt"
	"strconv"
	"strings"
)

type GameLaunchSpec struct {
	ExecutablePathRelative string
	Params                 []GameLaunchParam
}

func NewGameLaunchSpec(executablePathRelative string, params []GameLaunchParam) (*GameLaunchSpec, error) {
	return &GameLaunchSpec{
		ExecutablePathRelative: executablePathRelative,
		Params:                 params,
	}, nil
}

func (spec *GameLaunchSpec) UpdateParam(name string, argument string, value string, enabled bool) error {
	for i, p := range spec.Params {
		if p.Name == name && p.Argument == argument {
			spec.Params[i].Value = value
			spec.Params[i].Enabled = enabled
			if err := spec.Params[i].Validate(); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("param with name %s and argument %s not found", name, argument)
}

func (spec *GameLaunchSpec) UpdateParamByName(name string, value string, enabled bool) error {
	for i, p := range spec.Params {
		if p.Name == name {
			spec.Params[i].Value = value
			spec.Params[i].Enabled = enabled
			if err := spec.Params[i].Validate(); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("param with name %s not found", name)
}

func (spec *GameLaunchSpec) Validate() error {
	for _, param := range spec.Params {
		if err := param.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (spec *GameLaunchSpec) Render() ([]string, error) {
	args := []string{}
	remainingSeparator := ""
	nonSpaceSeparator := false

	for i, param := range spec.Params {
		rendered, err := param.Render()
		if err != nil {
			return nil, err
		}
		if len(rendered) == 0 {
			continue
		}

		if len(remainingSeparator) > 0 {
			rendered[0] = remainingSeparator + rendered[0]
			remainingSeparator = ""
		}

		if nonSpaceSeparator {
			args[len(args)-1] = args[len(args)-1] + rendered[0]
			rendered = rendered[1:]
		}

		args = append(args, rendered...)

		separator := param.ArgumentSeparator
		if i < len(spec.Params)-1 && strings.Contains(separator, " ") {
			nonSpaceSeparator = false
			separatorParts := strings.SplitN(separator, " ", 2)
			args[len(args)-1] = args[len(args)-1] + separatorParts[0]
			remainingSeparator = separatorParts[1]
		} else {
			nonSpaceSeparator = true
			args[len(args)-1] = args[len(args)-1] + separator
		}
	}

	return args, nil
}

type GameLaunchParamType int

const (
	GameLaunchParamFlag GameLaunchParamType = iota
	GameLaunchParamString
	GameLaunchParamInt
	GameLaunchParamFloat
	GameLaunchParamEnum
)

type EnumValueOption struct {
	Value string
	Label string
}

type GameLaunchParamInput struct {
	Type              GameLaunchParamType
	Name              string
	Description       string
	Argument          string
	Value             string
	EnumValues        []EnumValueOption
	Required          bool
	Enabled           bool
	ArgumentSeparator string
	ValueSeparator    string
	FloatPrecision    *int64
	MinInt            *int64
	MaxInt            *int64
	MinFloat          *float64
	MaxFloat          *float64
}

type GameLaunchParam struct {
	Type              GameLaunchParamType
	Name              string
	Description       string
	Argument          string
	Value             string
	EnumValues        []EnumValueOption
	Required          bool
	Enabled           bool
	ArgumentSeparator string
	ValueSeparator    string
	FloatPrecision    *int64
	MinInt            *int64
	MaxInt            *int64
	MinFloat          *float64
	MaxFloat          *float64
}

func NewGameLaunchParam(input GameLaunchParamInput) (*GameLaunchParam, error) {
	param := &GameLaunchParam{
		Type:              input.Type,
		Name:              input.Name,
		Description:       input.Description,
		Argument:          input.Argument,
		Value:             input.Value,
		EnumValues:        input.EnumValues,
		Required:          input.Required,
		Enabled:           input.Enabled,
		ArgumentSeparator: input.ArgumentSeparator,
		ValueSeparator:    input.ValueSeparator,
		FloatPrecision:    input.FloatPrecision,
		MinInt:            input.MinInt,
		MaxInt:            input.MaxInt,
		MinFloat:          input.MinFloat,
		MaxFloat:          input.MaxFloat,
	}
	return param, nil
}

func (param *GameLaunchParam) SetValue(value string) error {
	param.Value = value
	if err := param.Validate(); err != nil {
		return err
	}
	return nil
}

func (param *GameLaunchParam) Enable() {
	param.Enabled = true
}

func (param *GameLaunchParam) Disable() {
	param.Enabled = false
}

func (param *GameLaunchParam) Validate() error {
	if !param.Enabled {
		return nil
	}

	if len(param.Name) == 0 {
		return fmt.Errorf("name can not be empty")
	}
	if param.Required && !param.Enabled {
		return fmt.Errorf("can not be disabled when required")
	}

	switch param.Type {
	case GameLaunchParamFlag:
		return nil
	case GameLaunchParamString:
		if param.Required && len(param.Value) == 0 {
			return fmt.Errorf("required can not be empty")
		}
		if len(param.Value) == 0 {
			return fmt.Errorf("value can not be empty")
		}
		return nil
	case GameLaunchParamInt:
		value, err := strconv.ParseInt(param.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("value is not an integer: %w", err)
		}
		if param.MinInt != nil && value < *param.MinInt {
			return fmt.Errorf("value is less than the minimum")
		}
		if param.MaxInt != nil && value > *param.MaxInt {
			return fmt.Errorf("value is greater than the maximum")
		}
		return nil
	case GameLaunchParamFloat:
		value, err := strconv.ParseFloat(param.Value, 64)
		if err != nil {
			return fmt.Errorf("value is not a float: %w", err)
		}
		if param.MinFloat != nil && value < *param.MinFloat {
			return fmt.Errorf("value is less than the minimum")
		}
		if param.MaxFloat != nil && value > *param.MaxFloat {
			return fmt.Errorf("value is greater than the maximum")
		}
		return nil
	case GameLaunchParamEnum:
		if len(param.Value) == 0 {
			return fmt.Errorf("value can not be empty")
		}
		if len(param.EnumValues) == 0 {
			return fmt.Errorf("enum values can not be empty")
		}
		for _, option := range param.EnumValues {
			if len(option.Label) == 0 {
				return fmt.Errorf("enum value label can not be empty")
			}
			if len(option.Value) == 0 {
				return fmt.Errorf("enum value value can not be empty")
			}
		}
		for _, option := range param.EnumValues {
			if option.Value == param.Value {
				return nil
			}
		}
		return fmt.Errorf("value is not a valid enum value")
	}

	return fmt.Errorf("invalid param type: %d", param.Type)
}

func (param *GameLaunchParam) Render() ([]string, error) {
	if !param.Enabled {
		return []string{}, nil
	}
	if err := param.Validate(); err != nil {
		return []string{}, err
	}

	separator := param.ValueSeparator
	if separator == "" {
		separator = "="
	}

	value := param.Value
	if param.Type == GameLaunchParamFloat && param.FloatPrecision != nil {
		v, err := strconv.ParseFloat(param.Value, 64)
		if err != nil {
			return []string{}, err
		}
		digits := *param.FloatPrecision
		value = fmt.Sprintf("%.*f", digits, v)
	}

	data := struct {
		Arg       []string
		Separator string
		Value     string
	}{
		Arg:       strings.Split(param.Argument, " "),
		Separator: separator,
		Value:     value,
	}

	if param.Type == GameLaunchParamFlag {
		return data.Arg, nil
	}

	if strings.Contains(data.Separator, " ") {
		separatorParts := strings.SplitN(data.Separator, " ", 2)
		data.Arg[len(data.Arg)-1] = data.Arg[len(data.Arg)-1] + separatorParts[0]
		data.Arg = append(data.Arg, separatorParts[1]+data.Value)
	} else {
		data.Arg[len(data.Arg)-1] = data.Arg[len(data.Arg)-1] + data.Separator + data.Value
	}

	result := []string{}
	for _, arg := range data.Arg {
		if arg != "" {
			result = append(result, arg)
		}
	}

	return result, nil
}
