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

	userApi "github.com/evg555/auth/internal/api"
	"github.com/evg555/auth/internal/model"
	"github.com/evg555/auth/internal/service"
	serviceMock "github.com/evg555/auth/internal/service/mocks"
	proto "github.com/evg555/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *proto.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, false, false, 10)
		role     = gofakeit.IntRange(0, 1)

		serviceErr = errors.New("service error")

		req = &proto.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            proto.Role(role),
		}

		res = &proto.CreateResponse{Id: id}

		user = &model.User{
			ID: 0,
			Name: sql.NullString{
				String: name,
				Valid:  false,
			},
			Password: password,
			Email: sql.NullString{
				String: email,
				Valid:  false,
			},
			Role:      int32(role),
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *proto.CreateResponse
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
				mock.CreateMock.Expect(ctx, user).Return(id, nil)
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
				mock.CreateMock.Expect(ctx, user).Return(0, serviceErr)
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

			resp, err := api.Create(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
