package config

import (
	"github.com/joho/godotenv"
	"time"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type HTTPConfig interface {
	Address() string
}

type SwaggerConfig interface {
	Address() string
}

type TokenConfig interface {
	AccessTokenSecretKey() string
	AccessTokenExpiration() time.Duration
	RefreshTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AuthPrefix() string
}

type KafkaProducerConfig interface {
	Brokers() ([]string, error)
}
