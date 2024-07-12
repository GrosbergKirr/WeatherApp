package internal

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel       string `yaml:"log_level"`
	DatabaseConfig `yaml:"storage"`
	HttpConfig     `yaml:"http_server_weather"`
	SideApiUrl     `yaml:"side_api_url"`
}
type DatabaseConfig struct {
	DBUsername string `yaml:"db_user" env:"DB_USER"`
	DBPassword string `yaml:"db_password" env:"DB_PASS"`
	DBAddress  string `yaml:"db_address" env:"DB_ADDRESS"`
	DBName     string `yaml:"db_name" env:"DB_NAME"`
	DBMode     string `yaml:"db_ssl_mode" env:"DB_SSL_MODE"`
}

type HttpConfig struct {
	ServerAddress     string        `yaml:"server_address" env:"HTTP_SERVER_ADDRESS"`
	ServerTimeout     time.Duration `yaml:"server_timeout" env:"HTTP_SERVER_TIMEOUT"`
	ServerIdleTimeout time.Duration `yaml:"server_idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT"`
}

type SideApiUrl struct {
	CitiesUrl  string `yaml:"get_cities_url" env:"CITES_URL"`
	WeatherUrl string `yaml:"get_weather_url" env:"WEATHER_URL"`
	ApiKey     string `yaml:"api_key" env:"API_KEY"`
}

func SetupConfig(configPath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal("Cant setup config")
	}
	return &cfg
}
