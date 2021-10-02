package args

const (
	RpcServiceMicroMallUsers  = "micro-mall-users"
	RpcServiceMicroMallSearch = "micro-mall-search"
	RpcServiceMicroMallPay    = "micro-mall-pay"
)

const (
	UserApplyShopTemplate  = "尊敬的商户【%d】你好，恭喜你于：%v开通微店铺【%v】成功，并初始交易金额【%v】"
	UserModifyShopTemplate = "尊敬的商户【%d】你好，恭喜你于：%v变更微店铺资料【%v】成功"
	UserCloseShopTemplate  = "尊敬的商户【%d】你好，恭喜你于：%v关闭微店铺资料【%v】成功，期待你再来"
)

const (
	ShopInfoSearchNoticeTag    = "shop_info_search_notice"
	ShopInfoSearchNoticeTagErr = "shop_info_search_notice_err"
)

const (
	ShopInfoSearchNoticeType = 10001
)

type SearchStoreShop struct {
	ShopId       int64  `json:"shop_id,omitempty"`
	NickName     string `json:"nick_name,omitempty"`
	FullName     string `json:"full_name,omitempty"`
	ShopCode     string `json:"shop_code,omitempty"`
	RegisterAddr string `json:"register_addr,omitempty"`
	BusinessAddr string `json:"business_addr,omitempty"`
	BusinessDesc string `json:"business_desc,omitempty"`
}

type CommonBusinessMsg struct {
	Type    int    `json:"type"`
	Tag     string `json:"tag"`
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}
