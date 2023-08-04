package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HookPath   string
	HookPort   int
	HookSecret string
	Whitelist  []string
}

func LoadConfig(c string) (Config, error) {
	file, err := os.Open(c)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
