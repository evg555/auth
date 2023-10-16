package main

import (
	"context"
	"flag"
	"github.com/evg555/auth/internal/api"
	"github.com/evg555/auth/internal/config"
	userRepo "github.com/evg555/auth/internal/repository/user"
	userService "github.com/evg555/auth/internal/service/user"
	proto "github.com/evg555/auth/pkg/user_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed load grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed load pg config: %v", err)
	}

	conn, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	repo := userRepo.NewRepository(conn)
	srv := userService.NewService(repo)
	server := api.NewServer(srv)

	listener, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	proto.RegisterUserV1Server(s, server)

	log.Printf("server listening at %v", listener.Addr())

	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
