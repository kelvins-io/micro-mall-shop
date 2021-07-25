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
			Code: shop_business.RetCode_SUCCESS,
		},
	}
	var shopId int64
	var retCode int
	shopId, retCode = service.CreateShopBusiness(ctx, req)
	result.ShopId = shopId
	if retCode != code.Success {
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
			result.Common.Code = shop_business.RetCode_ERROR
			result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		}
	}
	return &result, nil
}

func (s *ShopBusinessServer) ShopPledge(ctx context.Context, req *shop_business.ShopPledgeRequest) (*shop_business.ShopPledgeResponse, error) {
	return &shop_business.ShopPledgeResponse{}, nil
}

func (s *ShopBusinessServer) GetShopMaterial(ctx context.Context, req *shop_business.GetShopMaterialRequest) (*shop_business.GetShopMaterialResponse, error) {
	var result shop_business.GetShopMaterialResponse
	result.Common = &shop_business.CommonResponse{
		Code: shop_business.RetCode_SUCCESS,
		Msg:  "",
	}
	result.Material = &shop_business.ShopMaterial{}
	shopInfo, retCode := service.GetShopMaterial(ctx, req.ShopId)
	if retCode != code.Success {
		result.Common.Code = shop_business.RetCode_ERROR
		result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		return &result, nil
	}

	result.Material = &shop_business.ShopMaterial{
		ShopCode:         shopInfo.ShopCode,
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

func (s *ShopBusinessServer) GetShopInfo(ctx context.Context, req *shop_business.GetShopInfoRequest) (*shop_business.GetShopInfoResponse, error) {
	var result shop_business.GetShopInfoResponse
	result.Common = &shop_business.CommonResponse{
		Code: shop_business.RetCode_SUCCESS,
		Msg:  "",
	}
	shopInfoList, retCode := service.GetShopInfoList(ctx, req.ShopIds)
	if retCode != code.Success {
		result.Common.Code = shop_business.RetCode_ERROR
		result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		return &result, nil
	}
	result.InfoList = make([]*shop_business.ShopInfo, len(shopInfoList))
	for i := 0; i < len(shopInfoList); i++ {
		shopInfo := &shop_business.ShopInfo{
			ShopId:     shopInfoList[i].ShopId,
			MerchantId: shopInfoList[i].LegalPerson,
			FullName:   shopInfoList[i].FullName,
			ShopCode:   shopInfoList[i].ShopCode,
		}
		result.InfoList[i] = shopInfo
	}

	return &result, nil
}

func (s *ShopBusinessServer) SearchSyncShop(ctx context.Context, req *shop_business.SearchSyncShopRequest) (*shop_business.SearchSyncShopResponse, error) {
	result := &shop_business.SearchSyncShopResponse{
		Common: &shop_business.CommonResponse{
			Code: shop_business.RetCode_SUCCESS,
		},
		List: nil,
	}
	list, retCode := service.SearchShopSync(ctx, req.ShopId, int(req.PageSize), int(req.PageNum))
	if retCode != code.Success {
		result.Common.Code = shop_business.RetCode_ERROR
		return result, nil
	}
	result.List = list
	return result, nil
}

func (s *ShopBusinessServer) SearchShop(ctx context.Context, req *shop_business.SearchShopRequest) (*shop_business.SearchShopResponse, error) {
	result := &shop_business.SearchShopResponse{
		Common: &shop_business.CommonResponse{
			Code: shop_business.RetCode_SUCCESS,
		},
	}
	list, retCode := service.SearchShop(ctx, req)
	if retCode != code.Success {
		result.Common.Code = shop_business.RetCode_ERROR
		return result, nil
	}
	result.List = list
	return result, nil
}

func (s *ShopBusinessServer) GetShopMajorInfo(ctx context.Context, req *shop_business.GetShopMajorInfoRequest) (*shop_business.GetShopMajorInfoResponse, error) {
	result := &shop_business.GetShopMajorInfoResponse{
		Common: &shop_business.CommonResponse{
			Code: shop_business.RetCode_SUCCESS,
		},
		InfoList: nil,
	}
	list, retCode := service.GetShopMajorInfo(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.ShopBusinessNotExist:
			result.Common.Code = shop_business.RetCode_SHOP_NOT_EXIST
		case code.ShopBusinessStateNotVerify:
			result.Common.Code = shop_business.RetCode_SHOP_STATE_NOT_VERIFY
		default:
			result.Common.Code = shop_business.RetCode_ERROR
		}
		return result, nil
	}
	result.InfoList = list
	return result, nil
}
