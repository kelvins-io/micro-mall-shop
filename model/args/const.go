package args

import "gitee.com/cristiane/micro-mall-shop/proto/micro_mall_shop_proto/shop_business"

const (
	RpcServiceMicroMallUsers  = "micro-mall-users"
	RpcServiceMicroMallSearch = "micro-mall-search"
	RpcServiceMicroMallPay    = "micro-mall-pay"
)

const (
	UserApplyShopTemplate  = "尊敬的商户【%v】你好，恭喜你于：%v 开通微店铺【%v】成功，并初始公司交易账户成功，初始金额为【%v】"
	UserModifyShopTemplate = "尊敬的商户【%v】你好，恭喜你于：%v 变更微店铺【%v】资料成功"
	UserCloseShopTemplate  = "尊敬的商户【%v】你好，恭喜你于：%v 关闭微店铺【%v】成功，期待你再来"
)

const (
	ShopInfoSearchNoticeTag    = "shop_info_search_notice"
	ShopInfoSearchNoticeTagErr = "shop_info_search_notice_err"
)

const (
	ShopInfoSearchNoticeType = 10001
)

type ShopEventNotice struct {
	ShopId        int64                       `json:"shop_id,omitempty"`
	MerchantId    int64                       `json:"-"`
	OperationType shop_business.OperationType `json:"-"`
	ChargeBalance string                      `json:"-"`
	NickName      string                      `json:"nick_name,omitempty"`
	FullName      string                      `json:"full_name,omitempty"`
	ShopCode      string                      `json:"shop_code,omitempty"`
	RegisterAddr  string                      `json:"register_addr,omitempty"`
	BusinessAddr  string                      `json:"business_addr,omitempty"`
	BusinessDesc  string                      `json:"business_desc,omitempty"`
}

type CommonBusinessMsg struct {
	Type    int    `json:"type"`
	Tag     string `json:"tag"`
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}
