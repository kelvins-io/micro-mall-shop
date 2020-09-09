package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-shop/model/args"
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"gitee.com/cristiane/micro-mall-shop/pkg/code"
	"gitee.com/cristiane/micro-mall-shop/pkg/util"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_shop_proto/shop_business"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-shop/repository"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/random"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"strings"
	"time"
)

func CreateShopBusiness(ctx context.Context, req *shop_business.ShopApplyRequest) (shopId int64, retCode int) {
	retCode = code.Success
	if req.MerchantId > 0 {
		serverName := args.RpcServiceMicroMallUsers
		conn, err := util.GetGrpcClient(serverName)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
			retCode = code.ErrorServer
			return
		}
		defer conn.Close()

		client := users.NewMerchantsServiceClient(conn)
		r := users.GetMerchantsMaterialRequest{
			MaterialId: req.MerchantId,
		}
		rsp, err := client.GetMerchantsMaterial(ctx, &r)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterial %v,err: %v", serverName, err)
			retCode = code.ErrorServer
			return
		}
		if rsp == nil || rsp.Common.Code != users.RetCode_SUCCESS {
			retCode = code.ErrorServer
			return
		}
		if rsp.Info == nil || rsp.Info.MaterialId <= 0 {
			retCode = code.MerchantNotExist
			return
		}
	}

	if req.OperationType == shop_business.OperationType_CREATE {
		exist, err := repository.CheckShopBusinessInfoExist(int(req.MerchantId), req.NickName)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CheckShopBusinessInfoExistByMerchantId err: %v, MerchantId: %v", err, req.MerchantId)
			retCode = code.ErrorServer
			return
		}
		if exist {
			retCode = code.ShopBusinessExist
			return
		}
		model := mysql.ShopBusinessInfo{
			NickName:         req.NickName,
			FullName:         req.FullName,
			ShopCode:         uuid.New().String(),
			RegisterAddr:     req.RegisterAddr,
			BusinessAddr:     req.BusinessAddr,
			LegalPerson:      req.MerchantId,
			BusinessLicense:  req.BusinessLicense,
			TaxCardNo:        req.TaxCardNo,
			BusinessDesc:     req.BusinessDesc,
			SocialCreditCode: req.SocialCreditCode,
			OrganizationCode: req.OrganizationCode,
			State:            0,
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
		}
		tx := kelvins.XORM_DBEngine.NewSession()
		// 创建店铺账户
		err = repository.CreateShopBusinessInfo(tx, &model)
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
				retCode = code.ShopBusinessExist
				return
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo err: %v, model: %+v", err, model)
			retCode = code.ErrorServer
			return
		}
		// 获取店铺信息
		shopBusinessInfo, err := repository.GetShopBusinessInfo(tx, int(req.MerchantId), req.NickName)
		if err != nil {
			tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "GetShopBusinessInfoByMerchantId err: %v, MerchantId: %+v", err, req.MerchantId)
			retCode = code.ErrorServer
			return
		}
		shopId = shopBusinessInfo.ShopId

		account := mysql.Account{
			AccountCode: uuid.New().String(),
			Owner:       shopBusinessInfo.ShopCode,
			Balance:     random.KrandNum(5),
			CoinType:    1,
			CoinDesc:    "CNY",
			State:       0,
			AccountType: 2,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		}
		err = repository.CreateAccount(tx, &account)
		if err != nil {
			tx.Rollback()
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount err: %v, account: %+v", err, account)
			retCode = code.ErrorServer
			return
		}
		tx.Commit()
	} else if req.OperationType == shop_business.OperationType_UPDATE {
		query := map[string]interface{}{}
		if req.ShopId > 0 {
			query["shop_id"] = req.ShopId
		}
		if req.MerchantId > 0 {
			query["legal_person"] = req.MerchantId
		}
		maps := map[string]interface{}{
			"nick_name":     req.NickName,
			"full_name":     req.FullName,
			"register_addr": req.RegisterAddr,
			"business_addr": req.BusinessAddr,
			"business_desc": req.BusinessDesc,
		}
		err := repository.UpdateShopBusinessInfo(query, maps)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "UpdateShopBusinessInfo err: %v, query: %+v, maps: %+v", err, query, maps)
			retCode = code.ErrorServer
			return
		}
	} else if req.OperationType == shop_business.OperationType_DELETE {

	}
	return
}

func GetShopMaterial(ctx context.Context, shopId int64) (*mysql.ShopBusinessInfo, int) {
	shopInfo, err := repository.GetShopBusinessInfoByShopId(shopId)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetShopBusinessInfoByShopId err: %v, shopId: %+v", err, shopId)
		return shopInfo, code.ErrorServer
	}
	return shopInfo, code.Success
}
