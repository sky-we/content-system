package dao

import (
	"content-system/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ContentDetailDao struct {
	db *gorm.DB
}

func NewContentDetailDao(db *gorm.DB) *ContentDetailDao {
	return &ContentDetailDao{db: db}
}

func (c *ContentDetailDao) IsVideoRepeat(videoUrl string) (bool, error) {
	var contentDetail model.ContentDetail
	err := c.db.Where("video_url = ?", videoUrl).First(&contentDetail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *ContentDetailDao) IsExist(id int) (bool, error) {
	var contentDetail model.ContentDetail
	err := c.db.Where("id = ?", id).First(&contentDetail).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *ContentDetailDao) Create(detail *model.ContentDetail) (int, error) {
	result := c.db.Create(&detail)
	if result.Error != nil {
		fmt.Printf("ContentDetailDao create error,%s", result.Error)
		return 0, result.Error
	}
	return detail.ID, nil
}

func (c *ContentDetailDao) Update(id int, detail *model.ContentDetail) error {

	if err := c.db.Where("id = ?", id).Updates(&detail).Error; err != nil {
		fmt.Printf("ContentDetailDao update error, %s", err)
		return err
	}
	return nil
}

func (c *ContentDetailDao) Delete(id int, detail *model.ContentDetail) error {
	if err := c.db.Delete(&detail, id).Error; err != nil {
		fmt.Printf("ContentDetailDao delete error, %s", err)
		return err
	}
	return nil

}
