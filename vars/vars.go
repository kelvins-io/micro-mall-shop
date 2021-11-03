package vars

import (
	"gitee.com/kelvins-io/common/queue"
	"gitee.com/kelvins-io/g2cache"
	"gitee.com/kelvins-io/kelvins/config/setting"
	"gitee.com/kelvins-io/kelvins/util/queue_helper"
)

var (
	EmailConfigSetting          *EmailConfigSettingS
	EmailNoticeSetting          *EmailNoticeSettingS
	ShopInfoSearchNoticeSetting *setting.QueueAMQPSettingS
	ShopInfoSearchNoticeServer  *queue.MachineryQueue
	ShopInfoSearchNoticePusher  *queue_helper.PublishService
	G2CacheSetting              *G2CacheSettingS
	G2CacheEngine               *g2cache.G2Cache
)
