package config

import (
	"github.com/jinzhu/configor"
)

// Config config info
type Config struct {
	AppName string `default:"zmemo"`
	Port    string `default:"8080"`
	DBLog   bool
	DB      struct {
		Name     string `default:"pgw"`
		User     string `default:"pgw"`
		Password string `default:"pgw"`
		Port     string `default:"3306"`
		Host     string `default:"localhost"`
	}
}

// New load config.yaml
func New(file ...string) *Config {
	config := new(Config)

	if len(file) < 1 {
		file = append(file, "./config/config.yaml")
	}

	if err := configor.Load(config, file[0]); err != nil {
		panic(err)
	}

	return config
}
