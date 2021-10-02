package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-shop/model/args"
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"gitee.com/cristiane/micro-mall-shop/pkg/code"
	"gitee.com/cristiane/micro-mall-shop/pkg/util"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_search_proto/search_business"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_shop_proto/shop_business"
	"gitee.com/cristiane/micro-mall-shop/repository"
	"gitee.com/cristiane/micro-mall-shop/vars"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"strconv"
)

func SearchShopSync(ctx context.Context, shopId int64, pageSize, pageNum int) ([]*shop_business.SearchSyncShopEntry, int) {
	result := make([]*shop_business.SearchSyncShopEntry, 0)
	var shopIds []int64
	if shopId > 0 {
		shopIds = append(shopIds, shopId)
	}
	shopInfoList, err := repository.GetShopInfoList("*", shopIds, pageSize, pageNum)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "SearchShopSync err: %v, shopId: %v", err, shopId)
		return result, code.ErrorServer
	}
	result = make([]*shop_business.SearchSyncShopEntry, len(shopInfoList))
	for i := 0; i < len(shopInfoList); i++ {
		entry := &shop_business.SearchSyncShopEntry{
			ShopId:       shopInfoList[i].ShopId,
			NickName:     shopInfoList[i].NickName,
			FullName:     shopInfoList[i].FullName,
			ShopCode:     shopInfoList[i].ShopCode,
			RegisterAddr: shopInfoList[i].RegisterAddr,
			BusinessAddr: shopInfoList[i].BusinessAddr,
			BusinessDesc: shopInfoList[i].BusinessDesc,
		}
		result[i] = entry
	}
	return result, code.Success
}

func GetShopInfoList(ctx context.Context, shopIds []int64) ([]mysql.ShopBusinessInfo, int) {
	shopInfoList, err := repository.GetShopInfoList("*", shopIds, 0, 0)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetShopInfo err: %v, shopIds: %v", err, json.MarshalToStringNoError(shopIds))
		return shopInfoList, code.ErrorServer
	}
	return shopInfoList, code.Success
}

func SearchShop(ctx context.Context, req *shop_business.SearchShopRequest) (result []*shop_business.SearchShopInfo, retCode int) {
	result = make([]*shop_business.SearchShopInfo, 0)
	retCode = code.Success
	serverName := args.RpcServiceMicroMallSearch
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		retCode = code.ErrorServer
		return
	}
	//defer conn.Close()
	client := search_business.NewSearchBusinessServiceClient(conn)
	searchReq := search_business.ShopSearchRequest{
		ShopKey: req.Keyword,
	}
	searchRsp, err := client.ShopSearch(ctx, &searchReq)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "ShopSearch %v,err: %v", serverName, err)
		retCode = code.ErrorServer
		return
	}
	if searchRsp.Common.Code != search_business.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "ShopSearch  req: %+v, rsp: %+v", json.MarshalToStringNoError(searchReq), json.MarshalToStringNoError(searchRsp))
		retCode = code.ErrorServer
		return
	}
	if len(searchRsp.List) == 0 {
		return
	}
	shopIds := make([]int64, 0)
	for i := range searchRsp.List {
		if searchRsp.List[i].ShopId == "" {
			continue
		}
		shopId, err := strconv.ParseInt(searchRsp.List[i].ShopId, 10, 64)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "ShopSearch  ParseInt err: %v, shopId: %s", err, searchRsp.List[i].ShopId)
			retCode = code.ErrorServer
		}
		if shopId > 0 {
			shopIds = append(shopIds, shopId)
		}
	}
	shopList, retCode := GetShopInfoList(ctx, shopIds)
	if retCode != code.Success {
		return
	}
	if len(shopList) == 0 {
		return
	}
	shopIdToInfo := map[int64]mysql.ShopBusinessInfo{}
	for i := 0; i < len(shopList); i++ {
		shopIdToInfo[shopList[i].ShopId] = shopList[i]
	}
	result = make([]*shop_business.SearchShopInfo, 0, len(searchRsp.List))
	for i := range searchRsp.List {
		if searchRsp.List[i].ShopId == "" {
			continue
		}
		shopId, err := strconv.ParseInt(searchRsp.List[i].ShopId, 10, 64)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "SearchShop ParseInt err: %v, shopId: %v", err, searchRsp.List[i].ShopId)
			retCode = code.ErrorServer
			return
		}
		if _, ok := shopIdToInfo[shopId]; !ok {
			continue
		}
		entry := &shop_business.SearchShopInfo{
			Info: &shop_business.ShopMaterial{
				ShopId:           shopIdToInfo[shopId].ShopId,
				MerchantId:       shopIdToInfo[shopId].LegalPerson,
				NickName:         shopIdToInfo[shopId].NickName,
				FullName:         shopIdToInfo[shopId].FullName,
				RegisterAddr:     shopIdToInfo[shopId].RegisterAddr,
				BusinessAddr:     shopIdToInfo[shopId].BusinessAddr,
				BusinessLicense:  shopIdToInfo[shopId].BusinessLicense,
				TaxCardNo:        shopIdToInfo[shopId].TaxCardNo,
				BusinessDesc:     shopIdToInfo[shopId].BusinessDesc,
				SocialCreditCode: shopIdToInfo[shopId].SocialCreditCode,
				OrganizationCode: shopIdToInfo[shopId].OrganizationCode,
				ShopCode:         shopIdToInfo[shopId].ShopCode,
			},
			Score: searchRsp.List[i].Score,
		}
		result = append(result, entry)
	}

	return result, retCode
}

func shopInfoNoticeSearchNotice(info *args.SearchStoreShop) {
	var msg = &args.CommonBusinessMsg{
		Type:    args.ShopInfoSearchNoticeType,
		Tag:     "店铺搜索通知",
		UUID:    uuid.New().String(),
		Content: json.MarshalToStringNoError(info),
	}
	kelvins.GPool.SendJob(func() {
		var ctx = context.TODO()
		vars.ShopInfoSearchNoticePusher.PushMessage(ctx, msg)
	})
	return
}
