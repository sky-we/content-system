package services

import (
	"content-system/internal/dao"
	"content-system/internal/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	goflow "github.com/s8sg/goflow/v1"
	"net/http"
	"time"
)

type ContentCreateReq struct {
	ID             int           `json:"id"`                           // 内容ID
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

type ContentCreateRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

func (app *CmsApp) ContentCreate(ctx *gin.Context) {
	var contentCreateReq ContentCreateReq

	// 入参校验
	if err := ctx.ShouldBindJSON(&contentCreateReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "参数错误", "error": err.Error()})
		return
	}

	contentDetailDao := dao.NewContentDetailDao(app.db)

	// 内容是否重复上传
	exists, err := contentDetailDao.IsVideoRepeat(contentCreateReq.VideoURL)
	if exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": fmt.Sprintf("[video_url=%s]内容已存在", contentCreateReq.VideoURL)})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}
	contentDetail := &model.ContentDetail{
		Title:          contentCreateReq.Title,
		Description:    contentCreateReq.Description,
		Author:         contentCreateReq.Author,
		VideoURL:       contentCreateReq.VideoURL,
		Thumbnail:      contentCreateReq.Thumbnail,
		Category:       contentCreateReq.Category,
		Duration:       contentCreateReq.Duration,
		Resolution:     contentCreateReq.Resolution,
		FileSize:       contentCreateReq.FileSize,
		Format:         contentCreateReq.Format,
		Quality:        contentCreateReq.Quality,
		ApprovalStatus: contentCreateReq.ApprovalStatus,
		UpdatedAt:      contentCreateReq.UpdatedAt,
		CreatedAt:      contentCreateReq.CreatedAt,
	}

	contentId, err := contentDetailDao.Create(contentDetail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}

	body, err := json.Marshal(map[string]int{"input": contentId})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}

	// 数据加工开始
	if execErr := app.flowService.Execute("content-flow", &goflow.Request{
		Body: body,
	}); execErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "call flow service failed", "err": execErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &ContentCreateRsp{
		Code:    0,
		Message: "success",
		Data:    gin.H{"ID": contentId},
	})

}
