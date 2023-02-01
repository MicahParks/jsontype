package jsontype_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/MicahParks/jsontype"
)

type myConfig struct {
	MyString string `json:"myString"`
}

func (c myConfig) DefaultsAndValidate() (myConfig, error) {
	if c.MyString == "" {
		c.MyString = "default"
	}
	return c, nil
}

type errorConfig struct {
	MyString string `json:"myString"`
}

func (c errorConfig) DefaultsAndValidate() (errorConfig, error) {
	return errorConfig{}, jsontype.ErrDefaultsAndValidate
}

func TestReadError(t *testing.T) {
	err := os.Setenv(jsontype.EnvVarConfigJSON, "{}")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	_, err = jsontype.Read[errorConfig]()
	if !errors.Is(err, jsontype.ErrDefaultsAndValidate) {
		t.Fatalf("Invalid error: %v", err)
	}
	err = os.Unsetenv(jsontype.EnvVarConfigJSON)
	if err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}
}

func TestReadEnv(t *testing.T) {
	err := os.Setenv(jsontype.EnvVarConfigJSON, "{}")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	config, err := jsontype.Read[myConfig]()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if config.MyString != "default" {
		t.Fatalf("Invalid config read.")
	}
	err = os.Unsetenv(jsontype.EnvVarConfigJSON)
	if err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}
}

func TestReadEnvBadJSON(t *testing.T) {
	err := os.Setenv(jsontype.EnvVarConfigJSON, "{")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	_, err = jsontype.Read[myConfig]()
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Fatalf("Invalid error: %v", err)
	}
	err = os.Unsetenv(jsontype.EnvVarConfigJSON)
	if err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}
}

func TestReadFile(t *testing.T) {
	file, err := os.CreateTemp("", "jsontypetest")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer os.Remove(file.Name())

	_, err = file.WriteString("{}")
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	err = os.Setenv(jsontype.EnvVarConfigPath, file.Name())
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	config, err := jsontype.Read[myConfig]()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	if config.MyString != "default" {
		t.Fatalf("Invalid config read")
	}
	err = os.Unsetenv(jsontype.EnvVarConfigPath)
	if err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}
}

func TestReadNone(t *testing.T) {
	_, err := jsontype.Read[myConfig]()
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Invalid error: %v", err)
	}
}
