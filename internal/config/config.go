package config

import (
	"flag"
	"os"
	"time"

	"github.com/PritOriginal/problem-map-server/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  logger.Environment `yaml:"env" env:"ENV" env-default:"local"`
	REST Server             `yaml:"server"`
}

type Server struct {
	Host    string `yaml:"host" env:"SERVER_HOST"`
	Port    int    `yaml:"port" env:"SEFVER_PORT"`
	Timeout struct {
		Write time.Duration `yaml:"write" env:"SERVER_TIMEOUT_WRITE"`
		Read  time.Duration `yaml:"read" env:"SERVER_TIMEOUT_READ"`
		Idle  time.Duration `yaml:"idle" env:"SERVER_TIMEOUT_IDLE"`
	} `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
