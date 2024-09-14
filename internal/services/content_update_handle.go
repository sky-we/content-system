package services

import (
	"content-system/internal/dao"
	"content-system/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ContentUpdateReq struct {
	ID             int           `json:"id" binding:"required"`        // 内容ID
	Title          string        `json:"title" binding:"required"`     // 内容标题
	VideoURL       string        `json:"video_url" binding:"required"` // 视频播放URL
	Author         string        `json:"author" binding:"required"`    // 作者
	Description    string        `json:"description"`                  // 内容描述
	Thumbnail      string        `json:"thumbnail"`                    // 封面图URL
	Category       string        `json:"category"`                     // 内容分类
	Duration       time.Duration `json:"duration"`                     // 内容时长
	Resolution     string        `json:"resolution"`                   // 分辨率 如720p、1080p
	FileSize       int64         `json:"file_size"`                    // 文件大小
	Format         string        `json:"format"`                       // 文件格式 如MP4、AVI
	Quality        int32         `json:"quality"`                      // 视频质量 1-高清 2-标清
	ApprovalStatus int32         `json:"approval_status"`              // 审核状态 1-审核中 2-审核通过 3-审核不通过
	UpdatedAt      time.Time     `json:"updated_at"`                   // 内容更新时间
	CreatedAt      time.Time     `json:"created_at"`                   // 内容创建时间
}

type ContentUpdateRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

func (app *CmsApp) ContentUpdate(ctx *gin.Context) {
	// 参数校验
	var contentUpdateReq ContentUpdateReq
	if err := ctx.ShouldBindJSON(&contentUpdateReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "参数错误", "err": err.Error()})
		return
	}

	// 内容是否存在
	contentDetailDao := dao.NewContentDetailDao(app.db)
	exists, err := contentDetailDao.IsExist(contentUpdateReq.ID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": fmt.Sprintf("[ID=%d]内容不存在", contentUpdateReq.ID)})
		return
	}
	contentDetail := &model.ContentDetail{
		Title:          contentUpdateReq.Title,
		Description:    contentUpdateReq.Description,
		Author:         contentUpdateReq.Author,
		VideoURL:       contentUpdateReq.VideoURL,
		Thumbnail:      contentUpdateReq.Thumbnail,
		Category:       contentUpdateReq.Category,
		Duration:       contentUpdateReq.Duration,
		Resolution:     contentUpdateReq.Resolution,
		FileSize:       contentUpdateReq.FileSize,
		Format:         contentUpdateReq.Format,
		Quality:        contentUpdateReq.Quality,
		ApprovalStatus: contentUpdateReq.ApprovalStatus,
		UpdatedAt:      contentUpdateReq.UpdatedAt,
		CreatedAt:      contentUpdateReq.CreatedAt,
	}

	// 更新
	if err := contentDetailDao.Update(contentUpdateReq.ID, contentDetail); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &ContentCreateRsp{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"ID": contentUpdateReq.ID,
		},
	})

}
