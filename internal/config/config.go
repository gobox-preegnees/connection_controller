package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Debug bool `yaml:"debug" env-default:"false"`
	Redis struct {
		Url string `yaml:"url" env-default:"redis://default:password@localhost:6379/0"`
	}
	Kafka struct {
		Addrs []string `yaml:"addrs" env-default:"localhost:29092"`
		Producer struct {
			SnapshotTopic string `yaml:"snapshot_topic" env-default:"snapshot"`
			Attempts int `yaml:"attemps" env-default:"3"`
			Timeout int `yaml:"timeout" env-default:"10"`
			Sleeptime int `yaml:"sleeptime" env-default:"250"`
		}
		Consumer struct {
			ConsistencyTopic string `yaml:"consistency_topic" env-default:"consistency"`
			GroupId string `yaml:"group_id" env-default:"consistency_group_id"`
			Partition int `yaml:"partition" env-defauld:"0"`
		}
	}
	Http struct {
		Addr string `yaml:"addr" env-default:"localhost:9966"`
		JWTAlg string `yaml:"jwt_alg" env-default:"HS256"`
		Secret string `yaml:"secret" env-default:"secret"`
		CrtPath string `yaml:"crt_path" env-default:"server.crt"`
		KeyPath string `yaml:"key_path" env-default:"server.key"`
	}
} 

var once sync.Once
var cnf Config

// GetGonfig. 
func GetGonfig(path string) *Config {

	once.Do(func() {
		if err := cleanenv.ReadConfig(path, &cnf); err != nil {
			panic(err)
		}
	})
	return &cnf
}