package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userApi "github.com/evg555/auth/internal/api"
	"github.com/evg555/auth/internal/model"
	"github.com/evg555/auth/internal/service"
	serviceMock "github.com/evg555/auth/internal/service/mocks"
	proto "github.com/evg555/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *proto.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name       = gofakeit.Name()
		email      = gofakeit.Email()
		serviceErr = errors.New("service error")

		req = &proto.UpdateRequest{
			Id:    0,
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(email),
		}

		res = &emptypb.Empty{}

		user = &model.User{
			ID: 0,
			Name: sql.NullString{
				String: name,
				Valid:  true,
			},
			Email: sql.NullString{
				String: email,
				Valid:  true,
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
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
				mock.UpdateMock.Expect(ctx, user).Return(nil)
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
				mock.UpdateMock.Expect(ctx, user).Return(serviceErr)
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

			resp, err := api.Update(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}

}
