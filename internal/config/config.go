package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var config *Config

func Load() (*Config, error) {
	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "/etc/config.json"
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	config.setupDatabases()

	return config, nil
}

func (c *Config) setupDatabases() {
	var countMainDB int

	if len(c.Databases) == 0 {
		panic("no database configs found")
	}

	for i := range c.Databases {
		if countMainDB > 1 {
			panic("multiple main databases found")
		}

		if c.Databases[i].ReadOnly {
			c.Databases[i].Nickname = fmt.Sprintf("%s_readonly", config.Databases[i].Nickname)
		}
		if c.Databases[i].IsDefault {
			countMainDB++
		}
	}

}

func Get() *Config {
	if config == nil {
		panic("config not loaded")
	}
	return config
}

func (c *Config) GetMainDatabaseNickname() string {
	for i := range c.Databases {
		if c.Databases[i].IsDefault {
			return c.Databases[i].Nickname
		}
	}

	return ""
}

func (c *Config) GetDatabase(nickname string) *DatabaseConfig {
	for i := range c.Databases {
		if c.Databases[i].Nickname == nickname {
			return &c.Databases[i]
		}
	}

	return nil
}
