package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local" env-required:"true" env:"ENV"`
	StoragePath string `yaml:"storage_path" env-required:"true" env:"STORAGE_PATH"`
	HTTPServer  `yaml:"http_server" env:"HTTP_SERVER"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080" env:"ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s" env:"TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s" env:"IDLE_TIMEOUT"`
	User        string        `yaml:"user" env-required:"true" env:"USER"`
	Password    string        `yaml:"password" env-required:"true" env:"PASSWORD"`
}

func MustLoad() *Config {
	basePath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	configPath := basePath + "/.env"

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf(err.Error())
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
