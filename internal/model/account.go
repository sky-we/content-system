package model

import (
	"time"
)

type Account struct {
	Id       int64     `gorm:"column:id;primary_key"`
	UserId   string    `gorm:"column:user_id"`
	Password string    `gorm:"column:pass_word"`
	NickName string    `gorm:"column:nick_name"`
	Ct       time.Time `gorm:"column:created_at"`
	Ut       time.Time `gorm:"column:updated_at"`
}

func (a *Account) TableName() string {
	table := "cms_account.account"
	return table

}
