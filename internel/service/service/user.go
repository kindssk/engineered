package service

import (
	"context"
	s "engineered/api"
	"engineered/internel/service/biz"
)

type UserService struct {
	s.UnimplementedServiceServer
	biz *biz.UserBiz
}

func (u *UserService) InsertUser(ctx context.Context, user *s.User) (*s.Res, error) {
	u.biz.InsertUser(user.Name,user.Age)
	return nil,nil
}

func (u *UserService) UpdateUser(ctx context.Context, user *s.User) (*s.Res, error) {
	panic("implement me")
}

func (u *UserService) ShowUser(ctx context.Context, id *s.Id) (*s.Res, error) {
	panic("implement me")
}

//data.NewDateUser(user.GetName(), user.GetAge()).InsertUser(ctx)

func NewUserService(biz *biz.UserBiz) *UserService {
	return &UserService{biz: biz}
}

