package repository

import (
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

func CreateShopBusinessInfo(tx *xorm.Session, model *mysql.ShopBusinessInfo) (err error) {
	_, err = tx.Table(mysql.TableShopBusinessInfo).Insert(model)
	return
}

func UpdateShopBusinessInfo(query, maps map[string]interface{}) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableShopBusinessInfo).Where(query).Update(maps)
	return
}

func DeleteShopBusinessInfo(query interface{}) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableShopBusinessInfo).Delete(query)
	return
}

func GetShopBusinessInfoByShopId(shopId int64) (*mysql.ShopBusinessInfo, error) {
	var model mysql.ShopBusinessInfo
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableShopBusinessInfo).
		Where("shop_id = ? ", shopId).
		Get(&model)
	return &model, err
}

func GetShopInfoList(shopIds []int64) ([]mysql.ShopBusinessInfo, error) {
	var result = make([]mysql.ShopBusinessInfo, 0)
	var err error
	err = kelvins.XORM_DBEngine.Table(mysql.TableShopBusinessInfo).In("shop_id", shopIds).Find(&result)
	return result, err
}

func GetShopBusinessInfo(shopCode string) (*mysql.ShopBusinessInfo, error) {
	var model mysql.ShopBusinessInfo
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableShopBusinessInfo).
		Select("shop_id,shop_code").
		Where("shop_code = ?", shopCode).Get(&model)
	return &model, err
}

func CheckShopBusinessInfoExist(merchantId int, nickName string) (exist bool, err error) {
	var model mysql.ShopBusinessInfo
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableShopBusinessInfo).
		Select("shop_id").
		Where("legal_person = ? and nick_name = ?", merchantId, nickName).Get(&model)
	if err != nil {
		return false, err
	}
	if model.ShopId > 0 {
		return true, nil
	}
	return false, nil
}
