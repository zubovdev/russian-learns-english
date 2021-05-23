package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

var configPath string

type config struct {
	YandexDictApiKey string `yaml:"yandex_dict_api_key"`
}

func getConfig() (*config, error) {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	return cfg, yaml.Unmarshal(b, cfg)
}
