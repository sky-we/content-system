package dao

import (
	"content-system/internal/middleware"
	"content-system/internal/model"
	"errors"
	"gorm.io/gorm"
)

var Logger = middleware.GetLogger()

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
		Logger.Error(err)
		return false, err
	}
	return true, nil

}
func (a *AccountDao) Create(account model.Account) error {
	if err := a.db.Create(&account).Error; err != nil {
		Logger.Error("AccountDao create error =[%v]", err)
		return err
	}
	return nil
}

func (a *AccountDao) FindByUserId(userId string) (*model.Account, error) {
	var account model.Account

	if err := a.db.Where("user_id = ?", userId).First(&account).Error; err != nil {
		Logger.Error("account dao FindByUserId error [%v]\n", err)
		return nil, err
	}
	return &account, nil

}
