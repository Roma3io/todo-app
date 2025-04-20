package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

/*
MustLoad Приставка Must в имени функции обычно говорит,
что функция вместо возврата ошибки
аварийно завершает работу приложения — например,
будет паниковать.
*/
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable not set")
	}
	if !strings.Contains(cfg.StoragePath, "password=") {
		cfg.StoragePath = strings.Replace(
			cfg.StoragePath,
			"postgres://postgres@",
			"postgres://postgres:"+dbPassword+"@",
			1,
		)
	}
	return &cfg
}
