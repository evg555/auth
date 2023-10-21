package app

import (
	"context"
	"flag"
	"github.com/evg555/auth/internal/closer"
	"github.com/evg555/auth/internal/config"
	proto "github.com/evg555/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	grpsServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.RunGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()
	if err := config.Load(configPath); err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = NewServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpsServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpsServer)
	proto.RegisterUserV1Server(a.grpsServer, a.serviceProvider.Server(ctx))

	return nil
}

func (a *App) RunGRPCServer() error {
	listener, err := net.Listen("tcp", a.serviceProvider.GrpcConfig().Address())
	if err != nil {
		return err
	}

	log.Printf("server listening at %v", listener.Addr())

	if err = a.grpsServer.Serve(listener); err != nil {
		return err
	}

	return nil
}
