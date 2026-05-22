package config

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Languages         map[string]Language
	MaxConcurrentJobs int
	QueueTimeoutS     int
}

type yamlRoot struct {
	Languages         []Language `yaml:"languages"`
	MaxConcurrentJobs int        `yaml:"max_concurrent_jobs"`
	QueueTimeoutS     int        `yaml:"queue_timeout_s"`
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

	maxJobs := root.MaxConcurrentJobs
	if maxJobs <= 0 {
		maxJobs = runtime.NumCPU()
	}

	queueTimeout := root.QueueTimeoutS
	if queueTimeout <= 0 {
		queueTimeout = 30
	}

	cfg := &Config{
		Languages:         make(map[string]Language, len(root.Languages)),
		MaxConcurrentJobs: maxJobs,
		QueueTimeoutS:     queueTimeout,
	}
	for _, lang := range root.Languages {
		if lang.ID == "" {
			return nil, fmt.Errorf("language missing id")
		}

		// Validate limits
		if err := validateLimits(lang.Run.Limits); err != nil {
			return nil, fmt.Errorf("language %s: run limits: %w", lang.ID, err)
		}
		if lang.Build != nil {
			if err := validateLimits(lang.Build.Limits); err != nil {
				return nil, fmt.Errorf("language %s: build limits: %w", lang.ID, err)
			}
		}

		cfg.Languages[lang.ID] = lang
	}

	return cfg, nil
}

func validateLimits(l Limits) error {
	if l.WallTimeS <= 0 {
		return fmt.Errorf("wall_time_s must be > 0")
	}
	if l.MemoryKB <= 0 {
		return fmt.Errorf("memory_kb must be > 0")
	}
	if l.MaxProcesses <= 0 {
		return fmt.Errorf("max_processes must be > 0")
	}
	return nil
}

func (c *Config) GetLanguage(id string) (Language, error) {
	lang, ok := c.Languages[id]
	if !ok {
		return Language{}, fmt.Errorf("unknown language: %s", id)
	}
	return lang, nil
}
