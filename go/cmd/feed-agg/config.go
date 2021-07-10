package main

import (
	"fmt"
	"os"
	"time"

	"sigs.k8s.io/yaml"
)

func NewConfig(fp string) (Config, error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return Config{}, err
	}
	for k := range c.Feeds {
		fc := c.Feeds[k]
		fc.refresh, err = time.ParseDuration(c.Feeds[k].Refresh)
		if err != nil {
			return Config{}, err
		}
		c.Feeds[k] = fc
	}

	return c, nil
}

type Config struct {
	Feeds map[string]FeedConfig // id: conf
}

type FeedConfig struct {
	Refresh string
	refresh time.Duration
	URLs    map[string]string // id: url
}
