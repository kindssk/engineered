// +build wireinject

package main

import (
	"engineered/internel/service/biz"
	"engineered/internel/service/data"
	"engineered/internel/service/service"
)

func InitUserService() *service.UserService {
	wire.Build(service.NewUserService, biz.NewUserBiz, data.NewDateUser, data.NewMysqlDB)
	return &service.UserService{}
}
