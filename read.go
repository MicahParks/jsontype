package jsontype

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	// EnvVarConfigJSON is the environment variable that can be used to provide the JSON configuration for the Read
	// function.
	EnvVarConfigJSON = "CONFIG_JSON"
	// EnvVarConfigPath is the environment variable that can be used to provide the path to the JSON configuration file
	// for the Read function.
	EnvVarConfigPath = "CONFIG_PATH"
)

// Config is any data structure that can unmarshalled from JSON.
type Config[T any] interface {
	// ApplyDefaults applies default values to the configuration.
	//
	// For example, if a zero value is left for a *jsontype.JSONType[time.Duration], the default value can be set here.
	ApplyDefaults() T
}

// Read is a convenience function to read JSON configuration. It will first check the environment variable in the
// EnvVarConfigJSON for raw JSON, then it will check the environment variable in the EnvVarConfigPath for the path to a
// JSON file. If neither are set, it will attempt to read "config.json" in the current working directory. If that file
// does not exist, it will return an os.ErrNotExist error.
func Read[T Config[T]]() (T, error) {
	var (
		config T
		data   []byte
		err    error
		source string
	)
	configJSON := os.Getenv(EnvVarConfigJSON)
	configPath := os.Getenv(EnvVarConfigPath)
	if configPath == "" {
		configPath = "config.json"
	}

	if configJSON != "" {
		source = "environment variable"
		data = []byte(configJSON)
	} else {
		source = "file"
		data, err = os.ReadFile(configPath)
		if err != nil {
			return config, fmt.Errorf("failed to read config file at path: %q: %w", configPath, err)
		}
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration from %s: %w", source, err)
	}

	config = config.ApplyDefaults()

	return config, nil
}
