package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"gitee.com/cristiane/micro-mall-shop/pkg/code"
	"gitee.com/cristiane/micro-mall-shop/repository"
	"gitee.com/kelvins-io/kelvins"
)

func GetUserInfo(ctx context.Context, uid int) (*mysql.UserInfo, int) {
	user, err := repository.GetUserByUid(uid)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, uid)
		return user, code.ErrorServer
	}
	return user, code.Success
}
