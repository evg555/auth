package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	service "github.com/evg555/auth/pkg/user_v1"
)

const grpcPort = 8000

type server struct {
	service.UnimplementedUserV1Server
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	service.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", listener.Addr())

	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Create(_ context.Context, req *service.CreateRequest) (*service.CreateResponse, error) {
	log.Println("creating user...")
	log.Printf("req: %+v", req)

	id := rand.Int63()

	return &service.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(_ context.Context, req *service.GetRequest) (*service.GetResponse, error) {
	log.Printf("getting user with id: %d\n", req.GetId())

	return &service.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      service.Role_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) Update(_ context.Context, req *service.UpdateRequest) (*emptypb.Empty, error) {
	log.Println("updating user...")
	log.Printf("req: %+v", req)

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *service.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("deleting user with id: %d\n", req.GetId())

	return &emptypb.Empty{}, nil
}
