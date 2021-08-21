package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

func LoadConfig() (*koanf.Koanf, error) {
	var k = koanf.New(".")
	if err := k.Load(file.Provider("./config.yaml"), yaml.Parser()); err != nil {
		return nil, err
	}
	return k, nil
}
