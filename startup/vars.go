package startup

import (
	"gitee.com/cristiane/micro-mall-shop/model/args"
	"gitee.com/cristiane/micro-mall-shop/vars"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/setup"
	"gitee.com/kelvins-io/kelvins/util/queue_helper"
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	// 1  shop info search
	err = setupQueueShopInfoSearchNotice()

	return err
}

func setupQueueShopInfoSearchNotice() error {
	var err error
	if vars.ShopInfoSearchNoticeSetting != nil {
		vars.ShopInfoSearchNoticeServer, err = setup.NewAMQPQueue(vars.ShopInfoSearchNoticeSetting, nil)
		if err != nil {
			return err
		}
		vars.ShopInfoSearchNoticePusher, err = queue_helper.NewPublishService(
			vars.ShopInfoSearchNoticeServer, &queue_helper.PushMsgTag{
				DeliveryTag:    args.ShopInfoSearchNoticeTag,
				DeliveryErrTag: args.ShopInfoSearchNoticeTagErr,
				RetryCount:     vars.ShopInfoSearchNoticeSetting.TaskRetryCount,
				RetryTimeout:   vars.ShopInfoSearchNoticeSetting.TaskRetryTimeout,
			}, kelvins.BusinessLogger)
		if err != nil {
			return err
		}
	}
	return err
}

func StopFunc() error {

	return nil
}
