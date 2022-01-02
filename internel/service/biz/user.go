package biz

import "engineered/internel/service/data"

type UserBiz struct {
	repo *data.UserRepo
}

func NewUserBiz(repo *data.UserRepo) *UserBiz {
	return &UserBiz{repo: repo}
}

func (u *UserBiz)InsertUser(name string,age int32){
	u.repo.InsertUser(name,age)
}
