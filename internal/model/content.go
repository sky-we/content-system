package model

import (
	"time"
)

type ContentDetail struct {
	ID             int           `gorm:"column:id;primary_key"`  // 内容ID
	Title          string        `gorm:"column:title"`           // 内容标题
	Description    string        `gorm:"column:description"`     // 内容描述
	Author         string        `gorm:"column:author"`          // 作者
	VideoURL       string        `gorm:"column:video_url"`       // 视频播放URL
	Thumbnail      string        `gorm:"column:thumbnail"`       // 封面图URL
	Category       string        `gorm:"column:category"`        // 内容分类
	Duration       time.Duration `gorm:"column:duration"`        // 内容时长
	Resolution     string        `gorm:"column:resolution"`      // 分辨率 如720p、1080p
	FileSize       int64         `gorm:"column:file_size"`       // 文件大小
	Format         string        `gorm:"column:format"`          // 文件格式 如MP4、AVI
	Quality        int32         `gorm:"column:quality"`         // 视频质量 1-高清 2-标清
	ApprovalStatus int32         `gorm:"column:approval_status"` // 审核状态 1-审核中 2-审核通过 3-审核不通过
	UpdatedAt      time.Time     `gorm:"column:updated_at"`      // 内容更新时间
	CreatedAt      time.Time     `gorm:"column:created_at"`      // 内容创建时间
}

func (c *ContentDetail) TableName() string {
	table := "cms_content.content_details"
	return table

}
