package service

import (
	"context"
	"strings"
	"time"

	"gitee.com/cristiane/micro-mall-shop/model/args"
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"gitee.com/cristiane/micro-mall-shop/pkg/code"
	"gitee.com/cristiane/micro-mall-shop/pkg/util"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_shop_proto/shop_business"
	"gitee.com/cristiane/micro-mall-shop/repository"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
)

func CreateShopBusiness(ctx context.Context, req *shop_business.ShopApplyRequest) (shopId int64, retCode int) {
	retCode = code.Success
	if req.OperationType == shop_business.OperationType_CREATE {
		// 创建账户前检查是为了防止数据表没有唯一性检查
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
		shopCode := uuid.New().String()
		shopInfo := mysql.ShopBusinessInfo{
			NickName:         req.NickName,
			FullName:         req.FullName,
			ShopCode:         shopCode,
			RegisterAddr:     req.RegisterAddr,
			BusinessAddr:     req.BusinessAddr,
			LegalPerson:      req.MerchantId,
			BusinessLicense:  req.BusinessLicense,
			TaxCardNo:        req.TaxCardNo,
			BusinessDesc:     req.BusinessDesc,
			SocialCreditCode: req.SocialCreditCode,
			OrganizationCode: req.OrganizationCode,
			State:            2, // 线上流程禁止
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
		}
		tx := kelvins.XORM_DBEngine.NewSession()
		err = tx.Begin()
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo Begin err: %v", err)
			retCode = code.ErrorServer
			return
		}
		defer func() {
			if retCode != code.Success {
				err := tx.Rollback()
				if err != nil {
					kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo Rollback err: %v", err)
					return
				}
			}
		}()
		// 创建店铺账户
		err = repository.CreateShopBusinessInfo(tx, &shopInfo)
		if err != nil {
			if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
				retCode = code.ShopBusinessExist
				return
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo err: %v, shopInfo: %v", err, json.MarshalToStringNoError(shopInfo))
			retCode = code.ErrorServer
			return
		}
		serverName := args.RpcServiceMicroMallPay
		conn, err := util.GetGrpcClient(ctx, serverName)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
			retCode = code.ErrorServer
			return
		}
		//defer conn.Close()
		client := pay_business.NewPayBusinessServiceClient(conn)
		balance := "1.9999"
		accountReq := pay_business.CreateAccountRequest{
			Owner:       shopCode,
			AccountType: pay_business.AccountType_Company,
			CoinType:    pay_business.CoinType_CNY,
			Balance:     balance,
		}
		rsp, err := client.CreateAccount(ctx, &accountReq)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount %v,err: %v", serverName, err)
			retCode = code.ErrorServer
			return
		}
		if rsp.Common.Code != pay_business.RetCode_SUCCESS {
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount req %v,rsp: %v", json.MarshalToStringNoError(req), json.MarshalToStringNoError(rsp))
			retCode = code.ErrorServer
			return
		}
		err = tx.Commit()
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "Commit err: %v", err)
			retCode = code.TransactionFailed
			return
		}
		shopIdInfo, err := repository.GetShopBusinessInfo("shop_id", shopCode)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetShopBusinessInfo err: %v, shopCode: %v", err, shopCode)
			retCode = code.ErrorServer
			return
		}
		shopId = shopIdInfo.ShopId
		shopInfoEventNotice(&args.ShopEventNotice{
			ShopId:        shopIdInfo.ShopId,
			MerchantId:    req.MerchantId,
			NickName:      req.GetNickName(),
			FullName:      req.GetFullName(),
			ShopCode:      shopCode,
			ChargeBalance: balance,
			OperationType: req.OperationType,
			RegisterAddr:  req.GetRegisterAddr(),
			BusinessAddr:  req.GetBusinessAddr(),
			BusinessDesc:  req.GetBusinessDesc(),
		})
	} else if req.OperationType == shop_business.OperationType_UPDATE {
		shopId = req.ShopId
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
			kelvins.ErrLogger.Errorf(ctx, "UpdateShopBusinessInfo err: %v, query: %v, maps: %v", err, json.MarshalToStringNoError(query), json.MarshalToStringNoError(maps))
			retCode = code.ErrorServer
			return
		}
		shopInfoEventNotice(&args.ShopEventNotice{
			OperationType: req.OperationType,
			MerchantId:    req.MerchantId,
			ShopId:        req.GetShopId(),
			NickName:      req.GetNickName(),
			FullName:      req.GetFullName(),
			RegisterAddr:  req.GetRegisterAddr(),
			BusinessAddr:  req.GetBusinessAddr(),
			BusinessDesc:  req.GetBusinessDesc(),
		})
	} else if req.OperationType == shop_business.OperationType_DELETE {
		shopId = req.ShopId
		where := map[string]interface{}{
			"shop_id": req.ShopId,
		}
		if req.MerchantId > 0 {
			where["legal_person"] = req.MerchantId
		}
		err := repository.DeleteShopBusinessInfo(where)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "DeleteShopBusinessInfo err: %v, where: %v", err, json.MarshalToStringNoError(where))
			retCode = code.ErrorServer
			return
		}

		shopInfoEventNotice(&args.ShopEventNotice{
			MerchantId:    req.MerchantId,
			OperationType: req.OperationType,
			ShopId:        req.GetShopId(),
		})
		return
	}
	return
}

func GetShopMaterial(ctx context.Context, shopId int64) (*mysql.ShopBusinessInfo, int) {
	shopInfo, err := repository.GetShopBusinessInfoByShopId(shopId)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetShopBusinessInfoByShopId err: %v, shopId: %v", err, shopId)
		return shopInfo, code.ErrorServer
	}
	return shopInfo, code.Success
}

const sqlSelectShopMajorInfo = "shop_id,nick_name,shop_code,state"

func GetShopMajorInfo(ctx context.Context, req *shop_business.GetShopMajorInfoRequest) (result []*shop_business.ShopMajorInfo, retCode int) {
	retCode = code.Success
	shopInfoList, err := repository.GetShopInfoList(sqlSelectShopMajorInfo, req.ShopIds, 0, 0)
	result = make([]*shop_business.ShopMajorInfo, len(shopInfoList))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetShopInfoList err: %v, shopId: %s", err, req.GetShopIds())
		retCode = code.ErrorServer
		return
	}
	if len(shopInfoList) == 0 {
		retCode = code.ShopBusinessNotExist
		return
	}
	if len(shopInfoList) != len(req.GetShopIds()) {
		retCode = code.ShopBusinessNotExist
		return
	}
	for i := 0; i < len(shopInfoList); i++ {
		if shopInfoList[i].State != 2 {
			retCode = code.ShopBusinessStateNotVerify
			return
		}
		majorInfo := &shop_business.ShopMajorInfo{
			ShopId:   shopInfoList[i].ShopId,
			ShopCode: shopInfoList[i].ShopCode,
			ShopName: shopInfoList[i].NickName,
		}
		result[i] = majorInfo
	}
	return
}
