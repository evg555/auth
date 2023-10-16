package app

import (
	"context"
	"github.com/evg555/auth/internal/closer"
	"github.com/evg555/auth/internal/config"
	"github.com/evg555/auth/internal/repository"
	"github.com/evg555/auth/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"

	"github.com/evg555/auth/internal/api"
	userRepo "github.com/evg555/auth/internal/repository/user"
	userService "github.com/evg555/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool         *pgxpool.Pool
	userRepository repository.UserRepository

	userService service.UserService

	server *api.Server
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) GetPGPool() *pgxpool.Pool {
	if s.pgPool == nil {
		ctx := context.Background()

		pool, err := pgxpool.Connect(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		if err = pool.Ping(ctx); err != nil {
			log.Fatalf("ping error: %v", err)
		}

		s.pgPool = pool

		closer.Add(func() error {
			pool.Close()
			return nil
		})
	}

	return s.pgPool
}

func (s *serviceProvider) GetUseRepository() repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.GetPGPool())
	}

	return s.userRepository
}

func (s *serviceProvider) GetUserService() service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.GetUseRepository())
	}

	return s.userService
}

func (s *serviceProvider) GetServer() *api.Server {
	if s.server == nil {
		s.server = api.NewServer(s.GetUserService())
	}

	return s.server
}
