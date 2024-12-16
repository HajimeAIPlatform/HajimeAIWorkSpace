package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Config represents the configuration for an agent
type Config struct {
	Name        string
	Properties  map[string]string
	PrivateData map[string]string
	mu          sync.Mutex
}

func NewConfig(name string) *Config {
	return &Config{
		Name:        name,
		Properties:  make(map[string]string),
		PrivateData: make(map[string]string),
	}
}

func (c *Config) SetProperty(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Properties[key] = value
}

func (c *Config) GetProperty(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Properties[key]
}

func (c *Config) SetPropertyFromJsonFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	c.Properties = data
	return nil
}

func (c *Config) SetPrivateDataFromJsonFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	c.PrivateData = data
	return nil
}
