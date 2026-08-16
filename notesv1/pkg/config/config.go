package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func LoadConfig[T any](prefix string) (*T, error) {
	var cfg T

	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	return &cfg, nil
}
