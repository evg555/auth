package app

import (
	"context"
	"github.com/evg555/auth/internal/client/db"
	"github.com/evg555/auth/internal/client/db/pg"
	"github.com/evg555/auth/internal/closer"
	"github.com/evg555/auth/internal/config"
	"github.com/evg555/auth/internal/repository"
	"github.com/evg555/auth/internal/service"
	"log"

	"github.com/evg555/auth/internal/api"
	userRepo "github.com/evg555/auth/internal/repository/user"
	userService "github.com/evg555/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	db             db.Client
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

func (s *serviceProvider) GetDB() db.Client {
	if s.db == nil {
		ctx := context.Background()

		dbc, err := pg.New(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("%v", err)
		}

		if err = dbc.DB().Ping(ctx); err != nil {
			log.Fatalf("ping error: %v", err)
		}

		s.db = dbc

		closer.Add(func() error {
			err = dbc.Close()
			if err != nil {
				return err
			}
			return nil
		})
	}

	return s.db
}

func (s *serviceProvider) GetUseRepository() repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.GetDB())
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
