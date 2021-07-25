package vars

import (
	"gitee.com/cristiane/micro-mall-shop/pkg/util/goroutine"
	"gitee.com/kelvins-io/kelvins"
)

var (
	App                                *kelvins.GRPCApplication
	EmailConfigSetting                 *EmailConfigSettingS
	EmailNoticeSetting                 *EmailNoticeSettingS
	GPool                              *goroutine.Pool
)
