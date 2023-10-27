package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	userApi "github.com/evg555/auth/internal/api"
	"github.com/evg555/auth/internal/model"
	"github.com/evg555/auth/internal/service"
	serviceMock "github.com/evg555/auth/internal/service/mocks"
	proto "github.com/evg555/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *proto.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id         = gofakeit.Int64()
		name       = gofakeit.Name()
		email      = gofakeit.Email()
		role       = gofakeit.IntRange(0, 1)
		serviceErr = errors.New("service error")
		timeNow    = time.Now()

		req = &proto.GetRequest{Id: id}

		res = &proto.GetResponse{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      proto.Role(role),
			CreatedAt: timestamppb.New(timeNow),
			UpdatedAt: timestamppb.New(timeNow),
		}

		user = &model.User{
			ID: id,
			Name: sql.NullString{
				String: name,
				Valid:  true,
			},
			Email: sql.NullString{
				String: email,
				Valid:  true,
			},
			Role:      int32(role),
			CreatedAt: timeNow,
			UpdatedAt: sql.NullTime{
				Time:  timeNow,
				Valid: true,
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *proto.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMock.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMock.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := tt.userServiceMock(mc)
			api := userApi.NewServer(userServiceMock)

			resp, err := api.Get(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
