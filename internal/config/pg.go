package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	pgHostEnvName     = "POSTGRES_HOST"
	pgPortEnvName     = "POSTGRES_PORT"
	pgDBEnvNAme       = "POSTGRES_DB"
	pgUserEnvNAme     = "POSTGRES_USER"
	pgPasswordEnvName = "POSTGRES_PASSWORD"
)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

func (p pgConfig) DSN() string {
	return p.dsn
}

func NewPGConfig() (PGConfig, error) {
	host := os.Getenv(pgHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("pg host not found")
	}

	port := os.Getenv(pgPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("pg port not found")
	}

	db := os.Getenv(pgDBEnvNAme)
	if len(db) == 0 {
		return nil, errors.New("pg db not found")
	}

	user := os.Getenv(pgUserEnvNAme)
	if len(user) == 0 {
		return nil, errors.New("pg user not found")
	}

	pass := os.Getenv(pgPasswordEnvName)
	if len(pass) == 0 {
		return nil, errors.New("pg pass not found")
	}

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, db, user, pass,
	)

	return &pgConfig{dsn: dsn}, nil
}
