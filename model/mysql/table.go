package mysql

import (
	"time"
)

const (
	TableShopBusinessInfo = "shop_business"
)

type ShopBusinessInfo struct {
	ShopId           int64     `xorm:"'shop_id' not null pk autoincr comment('店铺ID') BIGINT"`
	NickName         string    `xorm:"'nick_name' not null comment('简称') VARCHAR(36)"`
	ShopCode         string    `xorm:"'shop_code' not null comment('店铺唯一code') CHAR(36)"`
	FullName         string    `xorm:"'full_name' not null comment('店铺全称') TEXT"`
	RegisterAddr     string    `xorm:"'register_addr' not null comment('注册地址') TEXT"`
	BusinessAddr     string    `xorm:"'business_addr' not null comment('实际经营地址') TEXT"`
	LegalPerson      int64     `xorm:"'legal_person' not null comment('店铺法人') index BIGINT"`
	BusinessLicense  string    `xorm:"'business_license' not null comment('经营许可证') CHAR(36)"`
	TaxCardNo        string    `xorm:"'tax_card_no' not null comment('纳税号') CHAR(36)"`
	BusinessDesc     string    `xorm:"'business_desc' not null comment('经营描述') TEXT"`
	SocialCreditCode string    `xorm:"'social_credit_code' not null comment('统一社会信用代码') CHAR(36)"`
	OrganizationCode string    `xorm:"'organization_code' not null comment('组织机构代码') CHAR(36)"`
	State            int       `xorm:"'state' not null default 0 comment('状态，0-未审核，1-审核不通过，2-审核通过') TINYINT"`
	CreateTime       time.Time `xorm:"'create_time' not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime       time.Time `xorm:"'update_time' not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}
