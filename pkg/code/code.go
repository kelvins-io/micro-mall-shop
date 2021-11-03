package code

import "gitee.com/kelvins-io/common/errcode"

const (
	Success                    = 29000000
	ErrorServer                = 29000001
	TransactionFailed          = 29000002
	UserNotExist               = 29000005
	UserExist                  = 29000006
	DBDuplicateEntry           = 29000007
	MerchantExist              = 29000008
	MerchantNotExist           = 29000009
	ShopBusinessExist          = 29000010
	ShopBusinessNotExist       = 29000011
	ShopBusinessStateNotVerify = 29000012
)

var ErrMap = make(map[int]string)

func init() {
	dict := map[int]string{
		Success:                    "OK",
		ErrorServer:                "服务器错误",
		TransactionFailed:          "事务执行失败",
		UserNotExist:               "用户不存在",
		DBDuplicateEntry:           "Duplicate entry",
		UserExist:                  "已存在用户记录，请勿重复创建",
		MerchantExist:              "商户认证材料已存在",
		MerchantNotExist:           "商户未提交材料",
		ShopBusinessExist:          "店铺申请材料已存在",
		ShopBusinessNotExist:       "商户未提交店铺材料",
		ShopBusinessStateNotVerify: "店铺状态未审核或冻结",
	}
	errcode.RegisterErrMsgDict(dict)
	for key, _ := range dict {
		ErrMap[key] = dict[key]
	}
}
