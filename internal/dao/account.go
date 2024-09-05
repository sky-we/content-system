package dao

import (
	"content-system/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AccountDao struct {
	db *gorm.DB
}

func NewAccountDao(db *gorm.DB) *AccountDao {
	return &AccountDao{db: db}

}

func (a *AccountDao) IsExist(userId string) (bool, error) {
	var account model.Account

	err := a.db.Where("user_id = ?", userId).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil

}
func (a *AccountDao) Create(account model.Account) error {
	if err := a.db.Create(&account).Error; err != nil {
		fmt.Printf("AccountDao create error =[%v]", err)
		return err
	}
	return nil
}
