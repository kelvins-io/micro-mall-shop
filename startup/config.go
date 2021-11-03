package startup

import (
	"gitee.com/cristiane/micro-mall-shop/vars"
	"gitee.com/kelvins-io/kelvins/config"
	"gitee.com/kelvins-io/kelvins/config/setting"
)

const (
	SectionEmailConfig          = "email-config"
	EmailNotice                 = "email-notice"
	SectionShopInfoSearchNotice = "shop-info-search-notice"
	SectionG2Cache              = "micro-mall-g2cache"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	// 邮件通知
	vars.EmailNoticeSetting = new(vars.EmailNoticeSettingS)
	config.MapConfig(EmailNotice, vars.EmailNoticeSetting)
	// 店铺搜索通知
	vars.ShopInfoSearchNoticeSetting = new(setting.QueueAMQPSettingS)
	config.MapConfig(SectionShopInfoSearchNotice, vars.ShopInfoSearchNoticeSetting)
	//加载G2Cache二级缓存配置
	vars.G2CacheSetting = new(vars.G2CacheSettingS)
	config.MapConfig(SectionG2Cache, vars.G2CacheSetting)
	return nil
}
