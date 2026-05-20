package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Languages map[string]Language
}

type yamlRoot struct {
	Languages []Language `yaml:"languages"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var root yamlRoot
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if len(root.Languages) == 0 {
		return nil, fmt.Errorf("no languages defined in config")
	}

	cfg := &Config{
		Languages: make(map[string]Language, len(root.Languages)),
	}
	for _, lang := range root.Languages {
		if lang.ID == "" {
			return nil, fmt.Errorf("language missing id")
		}
		cfg.Languages[lang.ID] = lang
	}

	return cfg, nil
}

func (c *Config) GetLanguage(id string) (Language, error) {
	lang, ok := c.Languages[id]
	if !ok {
		return Language{}, fmt.Errorf("unknown language: %s", id)
	}
	return lang, nil
}
