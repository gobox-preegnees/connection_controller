package config

import (
	"sync"
)

type Config struct {

} 

var once sync.Once
var config *Config

func GetGonfig() *Config {

	once.Do(func() {
		config = &Config{}
	})
	return config
}