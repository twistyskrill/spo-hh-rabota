package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"` //Посмотреть про struct теги
	DB         `yaml:"db"`
	HTTPServer `yaml:"http_server"`
	JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeOut time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname" env-default:"handyman"`
	SSLmode  string `yaml:"sslmode" env-default:"disable"`

	// Пул соединений
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"10"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"100"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"1h"`

	// Таймауты
	PingTimeout time.Duration `yaml:"ping_timeout" env-default:"5s"`
}

type JWT struct {
	SecretKey string `yaml:"secret_key" env-default:"key"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	// Если CONFIG_PATH не установлена, используем значение по умолчанию
	if configPath == "" {
		configPath = "./config/local.yaml"
		log.Printf("CONFIG_PATH not set, using default: %s", configPath)
	}

	//проверка существования файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
