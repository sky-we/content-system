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

type FindParams struct {
	ID       int
	Author   string
	Title    string
	Page     int
	PageSize int
}

func (c *ContentDetailDao) Find(params *FindParams, detail *model.ContentDetail) (details *[]model.ContentDetail, total int64, err error) {
	query := c.db.Model(&detail)
	if params.ID != 0 {
		query = query.Where("id = ?", params.ID)
	}
	if params.Author != "" {
		query = query.Where("author = ?", params.Author)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}

	var count int64
	query.Count(&count)

	page := 1
	pageSize := 10

	if params.Page > 0 {
		page = params.Page
	}
	if params.PageSize > 0 {
		pageSize = params.PageSize
	}

	offset := (page - 1) * pageSize

	var results []model.ContentDetail
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, 0, err
	}
	return &results, count, nil

}
