package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

type Config struct {
	LogLevel string
}

func ParseConfig(args ...string) (*Config, *flag.FlagSet, error) {
	config := &Config{}
	flagset := flag.NewFlagSet("lanty", flag.ContinueOnError)

	flagset.Usage = func() {
		fmt.Printf("Usage: %s [options]\nOptions:\n", filepath.Base(args[0]))
		flagset.PrintDefaults()
	}

	flagset.StringVar(&config.LogLevel, "loglevel", "warning", "Log level of the application [disable, trace, debug, info, warning, error, panic, fatal]")

	err := flagset.Parse(args[1:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse into config: %w", err)
	}

	return config, flagset, nil
}
