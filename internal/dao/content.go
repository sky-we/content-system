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

func (c *ContentDetailDao) IsExist(videoUrl string) (bool, error) {
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

func (c *ContentDetailDao) Create(detail *model.ContentDetail) error {
	if err := c.db.Create(&detail).Error; err != nil {
		fmt.Printf("ContentDetailDao create error,%s", err)
		return err
	}
	return nil
}

func (c *ContentDetailDao) Update(videoUrl string, detail *model.ContentDetail) error {

	if err := c.db.Where("video_url = ?", videoUrl).Updates(&detail).Error; err != nil {
		fmt.Printf("ContentDetailDao update error, %s", err)
		return err
	}
	return nil
}
