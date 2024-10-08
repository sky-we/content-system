package process

import (
	"content-system/internal/dao"
	"content-system/internal/middleware"
	"encoding/json"
	flow "github.com/s8sg/goflow/flow/v1"
	"gorm.io/gorm"
)

type ContentFlow struct {
	ContentDao *dao.ContentDetailDao
}

var Logger = middleware.GetLogger()

func NewContentFlow(db *gorm.DB) *ContentFlow {
	detailDao := dao.NewContentDetailDao(db)
	return &ContentFlow{ContentDao: detailDao}
}

func (c *ContentFlow) ContentFlowHandle(workflow *flow.Workflow, context *flow.Context) error {
	dag := workflow.Dag()
	dag.Node("input", c.input)
	dag.Node("verify", c.verify)
	dag.Node("finish", c.finish)
	// 定义所有的分支类型
	branches := dag.ConditionalBranch("branches", []string{"category", "thumbnail", "pass", "format", "fail"},
		// 根据审核状态，返回分支类型
		func(bytes []byte) []string {
			var data map[string]interface{}
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil
			}
			if data["approval_status"].(float64) == 2 {
				return []string{"category", "thumbnail", "pass", "format"}
			}
			return []string{"fail"}
			// 分支结果聚合
		}, flow.Aggregator(func(m map[string][]byte) ([]byte, error) {
			return []byte("ok"), nil
		}),
	)
	branches["category"].Node("category", c.category)
	branches["thumbnail"].Node("thumbnail", c.thumbnail)
	branches["pass"].Node("category", c.pass)
	branches["format"].Node("format", c.format)
	branches["fail"].Node("fail", c.fail)

	dag.Edge("input", "verify")
	dag.Edge("verify", "branches")
	dag.Edge("branches", "finish")

	return nil

}

func (c *ContentFlow) input(data []byte, options map[string][]string) ([]byte, error) {
	Logger.Info("exec input node...")

	var d map[string]int
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	id := d["input"]
	detail, err := c.ContentDao.First(id)
	if err != nil {
		return nil, err
	}
	result, err := json.Marshal(map[string]interface{}{
		"title":     detail.Title,
		"video_url": detail.VideoURL,
		"id":        detail.ID,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *ContentFlow) verify(data []byte, options map[string][]string) ([]byte, error) {
	Logger.Info("exec verify node...")
	var detail map[string]interface{}

	if err := json.Unmarshal(data, &detail); err != nil {
		return nil, err
	}
	var (
		title    = detail["title"]
		videoUrl = detail["video_url"]
		id       = detail["id"]
	)
	if int(id.(float64))%2 == 0 {
		detail["approval_status"] = 3
	} else {
		detail["approval_status"] = 2
	}
	Logger.Info(id, title, videoUrl)
	return json.Marshal(detail)
}

func (c *ContentFlow) category(data []byte, options map[string][]string) ([]byte, error) {
	Logger.Info("exec category node...")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentId := int(input["id"].(float64))
	err := c.ContentDao.UpdateColById(contentId, "category", "category-workflow")
	if err != nil {
		return nil, err
	}
	return []byte("category"), nil
}
func (c *ContentFlow) thumbnail(data []byte, options map[string][]string) ([]byte, error) {
	Logger.Info("exec thumbnail node...")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentId := int(input["id"].(float64))
	err := c.ContentDao.UpdateColById(contentId, "thumbnail", "thumbnail-workflow")
	if err != nil {
		return nil, err
	}
	return []byte("thumbnail"), nil
}
func (c *ContentFlow) format(data []byte, options map[string][]string) ([]byte, error) {
	Logger.Info("exec format node...")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentId := int(input["id"].(float64))
	err := c.ContentDao.UpdateColById(contentId, "format", "format-workflow")
	if err != nil {
		return nil, err
	}
	return []byte("format"), nil
}
func (c *ContentFlow) pass(data []byte, option map[string][]string) ([]byte, error) {
	Logger.Info("exec pass node...")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentID := int(input["id"].(float64))
	// 审核成功
	if err := c.ContentDao.UpdateColById(contentID, "approval_status", 2); err != nil {
		return nil, err
	}
	return data, nil
}
func (c *ContentFlow) fail(data []byte, options map[string][]string) ([]byte, error) {
	Logger.Info("exec fail node...")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentId := int(input["id"].(float64))
	// 审核失败
	if err := c.ContentDao.UpdateColById(contentId, "approval_status", 3); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *ContentFlow) finish(data []byte, options map[string][]string) ([]byte, error) {
	Logger.Info("exec finish node...")
	Logger.Info(string(data))
	return data, nil
}
