package startup

import (
	"gitee.com/cristiane/micro-mall-shop/vars"
	"gitee.com/kelvins-io/kelvins/config"
	"gitee.com/kelvins-io/kelvins/config/setting"
)

const (
	SectionEmailConfig          = "email-config"
	SectionShopInfoSearchNotice = "shop-info-search-notice"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	// 店铺搜索通知
	vars.ShopInfoSearchNoticeSetting = new(setting.QueueAMQPSettingS)
	config.MapConfig(SectionShopInfoSearchNotice, vars.ShopInfoSearchNoticeSetting)
	return nil
}
