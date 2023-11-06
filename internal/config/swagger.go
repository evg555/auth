package config

import (
	"errors"
	"net"
	"os"
)

const (
	swaggerHostEnvName = "SWAGGER_HOST"
	swaggerPortEnvName = "SWAGGER_PORT"
)

type SwaggerConfig interface {
	Address() string
}

type swaggerConfig struct {
	host string
	port string
}

func (g swaggerConfig) Address() string {
	return net.JoinHostPort(g.host, g.port)
}

func NewSwaggerConfig() (SwaggerConfig, error) {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("HTTP host not found")
	}

	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("HTTP port not found")
	}

	return &swaggerConfig{
		host: host,
		port: port,
	}, nil
}
