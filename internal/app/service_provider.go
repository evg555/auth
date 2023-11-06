package app

import (
	"context"
	"log"

	"github.com/evg555/platform-common/pkg/closer"
	"github.com/evg555/platform-common/pkg/db"
	"github.com/evg555/platform-common/pkg/db/pg"
	"github.com/evg555/platform-common/pkg/db/transaction"

	"github.com/evg555/auth/internal/api"
	"github.com/evg555/auth/internal/config"
	"github.com/evg555/auth/internal/repository"
	userRepo "github.com/evg555/auth/internal/repository/user"
	"github.com/evg555/auth/internal/service"
	userService "github.com/evg555/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	db             db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userService service.UserService

	server *api.Server
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PgConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GrpcConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HttpConfig() config.GRPCConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.GRPCConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) Db(ctx context.Context) db.Client {
	if s.db == nil {
		dbc, err := pg.New(ctx, s.PgConfig().DSN())
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

func (s *serviceProvider) UseRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.Db(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UseRepository(ctx), s.TxManager(ctx))
	}

	return s.userService
}

func (s *serviceProvider) Server(ctx context.Context) *api.Server {
	if s.server == nil {
		s.server = api.NewServer(s.UserService(ctx))
	}

	return s.server
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.Db(ctx).DB())
	}

	return s.txManager
}
