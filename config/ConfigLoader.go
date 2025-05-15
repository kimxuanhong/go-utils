package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	LoadConfigFile()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Can not unmarshal config into struct: %v", err)
	}
	return &cfg, nil
}

func ResolveEnvInViper() {
	settings := viper.AllSettings()

	for key, value := range settings {
		// Chỉ xử lý nếu là string
		if strVal, ok := value.(string); ok {
			// Nếu giá trị là tên biến môi trường và tồn tại
			if envVal, exists := os.LookupEnv(strVal); exists {
				viper.Set(key, envVal)
			}
		}
	}
}

func LoadConfigFile() {
	dir, _ := os.Getwd()
	// Load default config
	configPath := filepath.Join(dir, "./resources")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can not load config.yml: %v", filepath.Join(configPath, "config.yml"))
	}

	// Check environment (ví dụ: dev, prod...)
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env != "" {
		viper.SetConfigName("config-" + env)
		err := viper.MergeInConfig()
		if err != nil {
			return
		}
	}
	ResolveEnvInViper()
}
