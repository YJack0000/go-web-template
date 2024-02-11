package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

type (
	Config struct {
		TestEnv string `env:"Test_ENV"`

		TWCC struct {
			APIKey string `env:"TWCC_API_KEY,required"`
		}

		HTTP struct {
			Port string `env:"HTTP_PORT,required" envDefault:"8080"`
		}

		Log struct {
			Level string `env:"LOG_LEVEL,required" envDefault:"debug"`
		}
	}
)

func NewConfig() (*Config, error) {
	var cfg Config
	env.Parse(&cfg)

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)

	return &cfg, nil
}
