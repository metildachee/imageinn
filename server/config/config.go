package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Elasticsearch struct {
		Index       string `yaml:"index"`
		Url         string `yaml:"url"`
		Sniff       bool   `yaml:"sniff"`
		HealthCheck bool   `yaml:"health-check"`
	} `yaml:"elasticsearch"`
}

func LoadConfig(path string) *Config {
	conf := Config{}

	// Read the YAML file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the YAML into our struct
	if err := yaml.Unmarshal(data, &conf); err != nil {
		log.Fatalf("error: %v", err)
	}

	return &conf
}
