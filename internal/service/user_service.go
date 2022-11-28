package service

import (
	"context"
	userpb "github.com/leondevpt/resume-backend/apigen/go/user/v1"
	"github.com/leondevpt/resume-backend/internal/biz"
)

var _ userpb.UserServiceServer = &UserServiceImpl{}

type UserServiceImpl struct {
	uc *biz.UserUseCase
	userpb.UnimplementedUserServiceServer
}

func NewUserService(uc *biz.UserUseCase) userpb.UserServiceServer {
	return &UserServiceImpl{uc: uc}
}

func (u UserServiceImpl) SignUpUser(ctx context.Context, input *userpb.SignUpUserInput) (*userpb.GenericResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) SignInUser(ctx context.Context, input *userpb.SignInUserInput) (*userpb.SignInUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) GetMe(ctx context.Context, request *userpb.GetMeRequest) (*userpb.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}
