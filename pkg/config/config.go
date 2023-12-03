package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func MustLoad[Config any]() *Config {
	configPath := fetchPath()
	if configPath == "" {
		panic("Path to config is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config on this path not found: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Error in config reading: " + err.Error())
	}

	return &cfg
}

func fetchPath() string {
	return os.Getenv("CONFIG_PATH")
}
