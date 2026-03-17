package game

import (
	"fmt"
	"strconv"
	"strings"
)

type EnumValueOption struct {
	Value string
	Label string
}

type LaunchParamType int

const (
	LaunchParamFlag LaunchParamType = iota
	LaunchParamString
	LaunchParamInt
	LaunchParamFloat
	LaunchParamEnum
)

type LaunchParamInput struct {
	Type              LaunchParamType
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

type LaunchParam struct {
	paramType         LaunchParamType
	name              string
	description       string
	argument          string
	value             string
	enumValues        []EnumValueOption
	required          bool
	enabled           bool
	argumentSeparator string
	valueSeparator    string
	floatPrecision    *int64
	minInt            *int64
	maxInt            *int64
	minFloat          *float64
	maxFloat          *float64
}

func NewLaunchParam(input LaunchParamInput) (*LaunchParam, error) {
	param := &LaunchParam{
		paramType:         input.Type,
		name:              input.Name,
		description:       input.Description,
		argument:          input.Argument,
		value:             input.Value,
		enumValues:        input.EnumValues,
		required:          input.Required,
		enabled:           input.Enabled,
		argumentSeparator: input.ArgumentSeparator,
		valueSeparator:    input.ValueSeparator,
		floatPrecision:    input.FloatPrecision,
		minInt:            input.MinInt,
		maxInt:            input.MaxInt,
		minFloat:          input.MinFloat,
		maxFloat:          input.MaxFloat,
	}
	return param, nil
}

func (param *LaunchParam) Type() LaunchParamType {
	return param.paramType
}

func (param *LaunchParam) Name() string {
	return param.name
}

func (param *LaunchParam) Description() string {
	return param.description
}

func (param *LaunchParam) Argument() string {
	return param.argument
}

func (param *LaunchParam) Value() string {
	return param.value
}

func (param *LaunchParam) EnumValues() []EnumValueOption {
	return param.enumValues
}

func (param *LaunchParam) Required() bool {
	return param.required
}

func (param *LaunchParam) Enabled() bool {
	return param.enabled
}

func (param *LaunchParam) MinInt() *int64 {
	return param.minInt
}

func (param *LaunchParam) MaxInt() *int64 {
	return param.maxInt
}

func (param *LaunchParam) MinFloat() *float64 {
	return param.minFloat
}

func (param *LaunchParam) MaxFloat() *float64 {
	return param.maxFloat
}

func (param *LaunchParam) SetValue(value string) error {
	param.value = value

	if err := param.Validate(); err != nil {
		return fmt.Errorf("failed to validate param: %w", err)
	}

	return nil
}

func (param *LaunchParam) SetEnumValues(values []EnumValueOption, defaultValue string) error {
	param.enumValues = values

	return param.SetValue(defaultValue)
}

func (param *LaunchParam) Enable() {
	param.enabled = true
}

func (param *LaunchParam) Disable() {
	param.enabled = false
}

func (param *LaunchParam) Validate() error {
	if !param.enabled {
		return nil
	}

	if len(param.name) == 0 {
		return fmt.Errorf("name can not be empty")
	}
	if param.required && !param.enabled {
		return fmt.Errorf("can not be disabled when required")
	}

	switch param.paramType {
	case LaunchParamFlag:
		return nil
	case LaunchParamString:
		if param.required && len(param.value) == 0 {
			return fmt.Errorf("required can not be empty")
		}
		if len(param.value) == 0 {
			return fmt.Errorf("value can not be empty")
		}
		return nil
	case LaunchParamInt:
		value, err := strconv.ParseInt(param.value, 10, 64)
		if err != nil {
			return fmt.Errorf("value is not an integer: %w", err)
		}
		if param.minInt != nil && value < *param.minInt {
			return fmt.Errorf("value is less than the minimum")
		}
		if param.maxInt != nil && value > *param.maxInt {
			return fmt.Errorf("value is greater than the maximum")
		}
		return nil
	case LaunchParamFloat:
		value, err := strconv.ParseFloat(param.value, 64)
		if err != nil {
			return fmt.Errorf("value is not a float: %w", err)
		}
		if param.minFloat != nil && value < *param.minFloat {
			return fmt.Errorf("value is less than the minimum")
		}
		if param.maxFloat != nil && value > *param.maxFloat {
			return fmt.Errorf("value is greater than the maximum")
		}
		return nil
	case LaunchParamEnum:
		if len(param.value) == 0 {
			return fmt.Errorf("value can not be empty")
		}
		if len(param.enumValues) == 0 {
			return fmt.Errorf("enum values can not be empty")
		}
		for _, option := range param.enumValues {
			if len(option.Label) == 0 {
				return fmt.Errorf("enum value label can not be empty")
			}
			if len(option.Value) == 0 {
				return fmt.Errorf("enum value value can not be empty")
			}
		}
		for _, option := range param.enumValues {
			if option.Value == param.value {
				return nil
			}
		}
		return fmt.Errorf("value is not a valid enum value")
	}

	return fmt.Errorf("invalid param type: %d", param.paramType)
}

func (param *LaunchParam) Render() ([]string, error) {
	if !param.enabled {
		return []string{}, nil
	}
	if err := param.Validate(); err != nil {
		return []string{}, fmt.Errorf("failed to validate param: %w", err)
	}

	separator := param.valueSeparator
	if separator == "" {
		separator = "="
	}

	value := param.value
	if param.paramType == LaunchParamFloat && param.floatPrecision != nil {
		v, err := strconv.ParseFloat(param.value, 64)
		if err != nil {
			return []string{}, fmt.Errorf("failed to parse float value: %w", err)
		}
		digits := *param.floatPrecision
		value = fmt.Sprintf("%.*f", digits, v)
	}

	data := struct {
		Arg       []string
		Separator string
		Value     string
	}{
		Arg:       strings.Split(param.argument, " "),
		Separator: separator,
		Value:     value,
	}

	if param.paramType == LaunchParamFlag {
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
