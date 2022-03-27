package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Proxy Proxy `yaml:"proxy"`
}

type Proxy struct {
	Port           int      `yaml:"port"`
	AllowedDomains []string `yaml:"allowed_domains,omitempty"`
}

func Load(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
