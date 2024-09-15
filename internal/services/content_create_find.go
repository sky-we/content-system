package services

import (
	"content-system/internal/dao"
	"content-system/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Content struct {
	ID             int           // 内容ID
	Title          string        // 内容标题
	Description    string        // 内容描述
	Author         string        // 作者
	VideoURL       string        // 视频播放URL
	Thumbnail      string        // 封面图URL
	Category       string        // 内容分类
	Duration       time.Duration // 内容时长
	Resolution     string        // 分辨率 如720p、1080p
	FileSize       int64         // 文件大小
	Format         string        // 文件格式 如MP4、AVI
	Quality        int32         // 视频质量 1-高清 2-标清
	ApprovalStatus int32         // 审核状态 1-审核中 2-审核通过 3-审核不通过
}
type ContentFindReq struct {
	ID       int    `json:"id"`       // 内容ID
	Author   string `json:"author"`   // 内容ID
	Title    string `json:"title"`    // 内容ID
	Page     int    `json:"page"`     // 页数
	PageSize int    `json:"pageSize"` // 页大小
}

type ContentFindRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    gin.H  `json:"data"`
}

func (app *CmsApp) ContentFind(ctx *gin.Context) {
	var contentFindReq ContentFindReq

	// 入参校验
	if err := ctx.ShouldBindJSON(&contentFindReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "参数错误", "error": err.Error()})
		return
	}

	contentDetailDao := dao.NewContentDetailDao(app.db)

	findParams := &dao.FindParams{
		ID:       contentFindReq.ID,
		Author:   contentFindReq.Author,
		Title:    contentFindReq.Title,
		Page:     contentFindReq.Page,
		PageSize: contentFindReq.PageSize,
	}
	contentList, total, err := contentDetailDao.Find(findParams, &model.ContentDetail{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}

	contents := make([]Content, 0, len(*contentList))

	for _, content := range *contentList {
		contents = append(contents, Content{
			ID:             content.ID,
			Title:          content.Title,
			Description:    content.Description,
			Author:         content.Author,
			VideoURL:       content.VideoURL,
			Thumbnail:      content.Thumbnail,
			Category:       content.Category,
			Duration:       content.Duration,
			Resolution:     content.Resolution,
			FileSize:       content.FileSize,
			Format:         content.Format,
			Quality:        content.Quality,
			ApprovalStatus: content.ApprovalStatus,
		})
	}

	ctx.JSON(http.StatusOK, &ContentFindRsp{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"content": contents,
			"total":   total,
		},
	})

}
