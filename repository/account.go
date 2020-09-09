package repository

import (
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"xorm.io/xorm"
)

func CreateAccount(tx *xorm.Session, model *mysql.Account) (err error) {
	_, err = tx.Table(mysql.TableAccount).Insert(model)
	return
}
