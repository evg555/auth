package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	service "github.com/evg555/auth/pkg/user_v1"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = 8000
	dbDSN    = "host=localhost port=5432 dbname=auth user=auth password=auth sslmode=disable"
	table    = "users"
)

type server struct {
	service.UnimplementedUserV1Server
	db *pgx.Conn
}

func main() {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	service.RegisterUserV1Server(s, &server{db: con})

	log.Printf("server listening at %v", listener.Addr())

	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Create(ctx context.Context, req *service.CreateRequest) (*service.CreateResponse, error) {
	log.Println("creating user...")
	log.Printf("req: %+v", req)

	var id int64

	builderInsert := sq.Insert(table).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "role").
		Values(
			req.GetName(),
			req.GetEmail(),
			req.GetPassword(),
			req.GetRole(),
		).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to insert to database: %v", err)
		return nil, err
	}

	res := s.db.QueryRow(ctx, query, args...)

	err = res.Scan(&id)
	if err != nil {
		log.Printf("failed to insert to database: %v", err)
		return nil, err
	}

	return &service.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *service.GetRequest) (*service.GetResponse, error) {
	log.Printf("getting user with id: %d\n", req.GetId())

	var (
		id        int64
		name      string
		email     string
		role      int32
		createdAt time.Time
		updatedAt sql.NullTime
	)

	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to select from database: %v", err)
		return nil, err
	}

	err = s.db.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to select from database: %v", err)
		return nil, err
	}

	return &service.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      service.Role(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt.Time),
	}, nil
}

func (s *server) Update(ctx context.Context, req *service.UpdateRequest) (*emptypb.Empty, error) {
	log.Println("updating user...")
	log.Printf("req: %+v", req)

	builderUpdate := sq.Update(table).PlaceholderFormat(sq.Dollar)

	if req.GetName() != nil {
		builderUpdate = builderUpdate.Set("name", req.GetName().Value)
	}

	if req.GetEmail() != nil {
		builderUpdate = builderUpdate.Set("email", req.GetEmail().Value)
	}

	builderUpdate = builderUpdate.Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to update database: %v", err)
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update database: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *service.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("deleting user with id: %d\n", req.GetId())

	builderDelete := sq.Delete(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
