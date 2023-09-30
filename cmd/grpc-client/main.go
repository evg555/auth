package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"

	service "github.com/evg555/auth/pkg/user_v1"
)

const (
	address = "localhost:8000"
	userId  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := service.NewUserV1Client(conn)

	resp1, err := c.Get(ctx, &service.GetRequest{Id: userId})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	log.Printf(color.RedString("Get user:\n"), color.GreenString("%+v", resp1))

	password := gofakeit.Password(true, true, true, true, true, 10)

	resp2, err := c.Create(ctx, &service.CreateRequest{
		Name:            gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        password,
		PasswordConfirm: password,
		Role:            service.Role_USER,
	})
	if err != nil {
		log.Fatal("failed to create user", err)
	}

	log.Printf(color.RedString("Create user:\n"), color.GreenString("user id: %d", resp2.GetId()))

	resp3, err := c.Update(ctx, &service.UpdateRequest{
		Id:    userId,
		Name:  wrapperspb.String(gofakeit.Name()),
		Email: nil,
	})
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf(color.RedString("Update user:\n"), color.GreenString("%+v", resp3))

	resp4, err := c.Delete(ctx, &service.DeleteRequest{Id: userId})
	if err != nil {
		log.Fatalf("failed to delete user by id: %v", err)
	}

	log.Printf(color.RedString("Delete user:\n"), color.GreenString("%+v", resp4))
}
