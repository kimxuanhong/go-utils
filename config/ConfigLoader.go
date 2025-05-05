package app

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

// LoadConfig loads configuration from YAML files and overrides values from environment variables.
// Usage: cfg, err := LoadConfig[MyConfig]()
// Priority:
//  1. Load 'resources/config.yml' as the base config
//  2. If ENV["CONFIG"] is set (e.g., config-dev.yml), it loads and overrides values from that file
//  3. For every string field, if the value looks like an ENV key (e.g., "SERVER_PORT"),
//     and the environment has a value for it, then override it.
func LoadConfig[T any]() (*T, error) {
	var cfg T

	// Kiểm tra: cfg phải là struct (dù đã dùng generic)
	v := reflect.ValueOf(&cfg).Elem()
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("LoadConfig only works with struct types, got %s", v.Kind())
	}

	basePath := "resources"
	mainConfig := filepath.Join(basePath, "config.yml")

	// Step 1: Load base config
	if err := loadYAML(mainConfig, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load main config: %w", err)
	}

	// Step 2: Optional override with CONFIG env
	if overrideFile := os.Getenv("CONFIG"); overrideFile != "" {
		envConfig := filepath.Join(basePath, overrideFile)
		if _, err := os.Stat(envConfig); err == nil {
			if err := loadYAML(envConfig, &cfg); err != nil {
				return nil, fmt.Errorf("failed to load override config: %w", err)
			}
		} else {
			log.Printf("Override config file %s not found, skipping\n", envConfig)
		}
	}

	// Step 3: Override string fields with ENV values
	resolveEnvValues(&cfg)

	return &cfg, nil
}

// loadYAML unmarshals YAML file into cfg
func loadYAML(file string, cfg interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close() // Safe to ignore close error
	}(f)

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(cfg)
}

// resolveEnvValues replaces string fields with matching ENV values
func resolveEnvValues(cfg interface{}) {
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)

		// Skip if field cannot be set (unexported/private)
		if !fieldVal.CanSet() {
			continue
		}

		switch fieldVal.Kind() {
		case reflect.Struct:
			resolveEnvValues(fieldVal.Addr().Interface())
		case reflect.String:
			envKey := fieldVal.String()
			if envVal := os.Getenv(envKey); envVal != "" {
				fieldVal.SetString(envVal)
			}
		default:
		}
	}
}
