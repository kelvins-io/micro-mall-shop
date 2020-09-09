package server

import (
	"context"
	"gitee.com/cristiane/micro-mall-shop/pkg/code"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_shop_proto/shop_business"
	"gitee.com/cristiane/micro-mall-shop/service"
	"gitee.com/kelvins-io/common/errcode"
)

type ShopBusinessServer struct {
}

func NewShopBusinessServer() shop_business.ShopBusinessServiceServer {
	return new(ShopBusinessServer)
}

func (s *ShopBusinessServer) ShopApply(ctx context.Context, req *shop_business.ShopApplyRequest) (*shop_business.ShopApplyResponse, error) {
	var result = shop_business.ShopApplyResponse{
		Common: &shop_business.CommonResponse{
			Code: 0,
			Msg:  "",
		},
		ShopId: 0,
	}
	var shopId int64
	var retCode int
	shopId, retCode = service.CreateShopBusiness(ctx, req)
	result.ShopId = shopId
	if retCode == code.ShopBusinessExist {
		result.Common.Code = shop_business.RetCode_SHOP_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.ShopBusinessExist)
	} else if retCode == code.ShopBusinessNotExist {
		result.Common.Code = shop_business.RetCode_SHOP_NOT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.ShopBusinessNotExist)
	} else if retCode == code.MerchantNotExist {
		result.Common.Code = shop_business.RetCode_MERCHANT_NOT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.MerchantNotExist)
	} else if retCode == code.MerchantExist {
		result.Common.Code = shop_business.RetCode_MERCHANT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.MerchantExist)
	} else {
		result.Common.Code = shop_business.RetCode_SUCCESS
		result.Common.Msg = errcode.GetErrMsg(code.Success)
	}

	return &result, nil
}

func (s *ShopBusinessServer) ShopPledge(ctx context.Context, req *shop_business.ShopPledgeRequest) (*shop_business.ShopPledgeResponse, error) {
	return &shop_business.ShopPledgeResponse{}, nil
}

func (s *ShopBusinessServer) GetShopMaterial(ctx context.Context, req *shop_business.GetShopMaterialRequest) (*shop_business.GetShopMaterialResponse, error) {
	var result shop_business.GetShopMaterialResponse
	result.Material = &shop_business.ShopMaterial{}
	shopInfo, retCode := service.GetShopMaterial(ctx, req.ShopId)
	if retCode != code.Success {
		return &result, errcode.TogRPCError(code.ErrorServer)
	}

	result.Material = &shop_business.ShopMaterial{
		ShopId:           shopInfo.ShopId,
		MerchantId:       shopInfo.LegalPerson,
		NickName:         shopInfo.NickName,
		FullName:         shopInfo.FullName,
		RegisterAddr:     shopInfo.RegisterAddr,
		BusinessAddr:     shopInfo.BusinessAddr,
		BusinessLicense:  shopInfo.BusinessLicense,
		TaxCardNo:        shopInfo.TaxCardNo,
		BusinessDesc:     shopInfo.BusinessDesc,
		SocialCreditCode: shopInfo.SocialCreditCode,
		OrganizationCode: shopInfo.OrganizationCode,
	}

	return &result, nil
}
