package game

import (
	"fmt"
	"strings"
)

type LaunchSpecMode string

const (
	LaunchSpecModePlay LaunchSpecMode = "play"
	LaunchSpecModeJoin LaunchSpecMode = "join"
	LaunchSpecModeHost LaunchSpecMode = "host"
)

type LaunchSpec struct {
	ExecutablePathRelative string
	params                 []LaunchParam
}

func NewLaunchSpec(executablePathRelative string, params []LaunchParam) (*LaunchSpec, error) {
	return &LaunchSpec{
		ExecutablePathRelative: executablePathRelative,
		params:                 params,
	}, nil
}

func (spec *LaunchSpec) Params() []LaunchParam {
	return spec.params
}

func (spec *LaunchSpec) UpdateParam(name string, argument string, value string, enabled bool) error {
	for i, p := range spec.params {
		if p.name == name && p.argument == argument {
			if enabled {
				spec.params[i].Enable()
			} else {
				spec.params[i].Disable()
			}

			err := spec.params[i].SetValue(value)
			if err != nil {
				return fmt.Errorf("failed to set value for param: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("param with name %q and argument %q not found", name, argument)
}

func (spec *LaunchSpec) UpdateParamByName(name string, value string, enabled bool) error {
	for i, p := range spec.params {
		if p.name == name {
			if enabled {
				spec.params[i].Enable()
			} else {
				spec.params[i].Disable()
			}

			err := spec.params[i].SetValue(value)
			if err != nil {
				return fmt.Errorf("failed to set value for param: %w", err)
			}

			return nil
		}
	}
	return fmt.Errorf("param with name %q not found", name)
}

func (spec *LaunchSpec) Validate() error {
	for _, param := range spec.params {
		err := param.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate param %q: %w", param.name, err)
		}
	}
	return nil
}

func (spec *LaunchSpec) Render() ([]string, error) {
	args := []string{}
	remainingSeparator := ""
	nonSpaceSeparator := false

	for i, param := range spec.params {
		rendered, err := param.Render()
		if err != nil {
			return nil, fmt.Errorf("failed to render param %q: %w", param.name, err)
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

		separator := param.argumentSeparator
		if i < len(spec.params)-1 && strings.Contains(separator, " ") {
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
