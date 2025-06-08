package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string   `yaml:"env" env-required:"true"`
	Database   Database `yaml:"database" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type Database struct {
	Host           string `yaml:"host" env-default:"localhost"`
	Port           string `yaml:"port" env-default:"5432"`
	User           string `yaml:"user" env-required:"true"`
	Password       string `yaml:"password" env-required:"true"`
	DbName         string `yaml:"dbname" env-required:"true"`
	SslMode        string `yaml:"ssl_mode" env-default:"require"`
	MigrationsPath string `yaml:"migrations_path" env-default:"/internal/storage/migrations"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	return &cfg
}

func GetDBConnectionString(config *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DbName,
		config.Database.SslMode)
}
