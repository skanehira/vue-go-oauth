package config

import (
	"github.com/jinzhu/configor"
)

// Config config info
type Config struct {
	AppName string `default:"test"`
	Port    string `default:"8080"`
	DBLog   bool
	AppLog  bool
	DB      struct {
		Name     string `default:"test"`
		User     string `default:"test"`
		Password string `default:"test"`
		Port     string `default:"3306"`
		Host     string `default:"localhost"`
	}
	Twitter struct {
		Token            string
		Secret           string
		RequestURI       string
		AuthorizationURI string
		TokenRequestURI  string
		CallbackURI      string
	}
}

// New load config.yaml
func New(file ...string) *Config {
	config := new(Config)

	if len(file) < 1 {
		file = append(file, "./config/config.yaml")
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(config, file[0]); err != nil {
		panic(err)
	}

	return config
}
