package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	ListenAddr string `yaml:"listen_addr" env-required:"true"`
	// ....
}

// MustLoad loads the configuration from the specified file and returns a pointer to a Config struct.
// It will panic if the file path is not specified, if the file does not exist, or if the configuration cannot be read.
//
// Returns:
// - A pointer to the populated Config struct.
// - Panics if any error occurs during file reading or parsing.
func MustLoad() *Config {
	configPath := fetchConfig()

	if configPath == "" {
		panic("Config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path doesn't exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

// fetchConfig retrieves the path to the configuration file.
// It first looks for a `-config` flag in the command line arguments.
// If the flag is not provided, it tries to get the value from the `CONFIG` environment variable.
//
// Returns:
// - The path to the configuration file as a string.
// - If neither the command line flag nor the environment variable is set, returns an empty string.
func fetchConfig() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG")
	}

	return res
}
