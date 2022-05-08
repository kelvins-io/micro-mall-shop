package service

import (
	"context"
	"fmt"
	"time"

	"gitee.com/cristiane/micro-mall-shop/model/args"
	"gitee.com/cristiane/micro-mall-shop/pkg/util"
	"gitee.com/cristiane/micro-mall-shop/pkg/util/email"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_shop_proto/shop_business"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-shop/repository"
	"gitee.com/cristiane/micro-mall-shop/vars"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
)

func shopInfoEventNotice(info *args.ShopEventNotice) {
	kelvins.GPool.SendJob(func() {
		var ctx = context.TODO()
		shopIdInfo, err := repository.GetShopBusinessInfoByShopId(info.ShopId)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetShopBusinessInfoByShopId err: %v, info.ShopId: %v", err, info.ShopId)
			return
		}
		info.MerchantId = shopIdInfo.LegalPerson
		info.NickName = shopIdInfo.NickName
		info.FullName = shopIdInfo.FullName
		info.ShopCode = shopIdInfo.ShopCode
		info.BusinessAddr = shopIdInfo.BusinessAddr
		info.RegisterAddr = shopIdInfo.RegisterAddr
		info.BusinessDesc = shopIdInfo.BusinessDesc
		// 1 搜索事件
		if info.OperationType != (shop_business.OperationType_DELETE) {
			var msg = &args.CommonBusinessMsg{
				Type:    args.ShopInfoSearchNoticeType,
				Tag:     "店铺搜索通知",
				UUID:    uuid.New().String(),
				Content: json.MarshalToStringNoError(info),
			}
			vars.ShopInfoSearchNoticePusher.PushMessage(ctx, msg)
		}

		// 2 发邮件
		serverName := args.RpcServiceMicroMallUsers
		conn, err := util.GetGrpcClient(ctx, serverName)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
			return
		}
		client := users.NewMerchantsServiceClient(conn)
		material, err := client.GetMerchantsMaterial(ctx, &users.GetMerchantsMaterialRequest{MaterialId: info.MerchantId})
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterial %v,err: %v", serverName, err)
			return
		}
		if material.GetCommon().GetCode() != users.RetCode_SUCCESS {
			kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterial req %v,rsp: %v", json.MarshalToStringNoError(info.MerchantId), json.MarshalToStringNoError(material))
			return
		}
		now := util.ParseTimeOfStr(time.Now().Unix())
		if material.GetInfo() != nil && material.GetInfo().GetMerchantEmail() != "" {
			var emailNotice string
			switch info.OperationType {
			case shop_business.OperationType_CREATE:
				emailNotice = fmt.Sprintf(args.UserApplyShopTemplate, material.GetInfo().GetMerchantName(), now, info.NickName, info.ChargeBalance)
			case shop_business.OperationType_UPDATE:
				emailNotice = fmt.Sprintf(args.UserModifyShopTemplate, material.GetInfo().GetMerchantName(), now, info.NickName)
			case shop_business.OperationType_DELETE:
				emailNotice = fmt.Sprintf(args.UserCloseShopTemplate, material.GetInfo().GetMerchantName(), now, info.NickName)
			default:
			}
			err = email.SendEmailNotice(ctx, material.GetInfo().MerchantEmail, kelvins.AppName, emailNotice)
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
				return
			}
		}
	})
	return
}
