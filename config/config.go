package config

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"
)

var once sync.Once
var instance *Config

// GetConfig get config defined in config.json
func GetConfig() *Config {
	once.Do(func() {
		var config *Config
		pwd, _ := os.Getwd()
		// it would be better to read by env variables or some flags, but here we make it simple
		filepath := path.Join(pwd, "config.json")

		configFile, err := os.Open(filepath)
		defer configFile.Close()
		if err != nil {
			panic(err)
		}

		jsonParser := json.NewDecoder(configFile)
		err = jsonParser.Decode(&config)
		if err != nil {
			panic(err)
		}

		instance = config
	})

	return instance
}

// Config config entry
type Config struct {
	Endpoint string `json:"endpoint"`
	Delay    int    `json:"delay"`
}

// GetDelay get delay in time.duration
func (c *Config) GetDelay() time.Duration {
	if c.Delay == 0 {
		c.Delay = 60
	}
	return time.Duration(c.Delay) * time.Second
}
