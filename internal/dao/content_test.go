package dao

import (
	"content-system/internal/config"
	"content-system/internal/model"
	"gorm.io/gorm"
	"testing"
)

// Goland右键自动生成测试代码
func TestContentDetailDao_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		detail *model.ContentDetail
	}
	config.LoadDBConfig()
	mysqldb := config.NewMySqlDB(config.DBConfig.MySQL)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "insert with id", // 测试名称
			fields: fields{mysqldb},  // 测试依赖项，方法绑定的结构体
			args: args{ // 函数传入的参数
				detail: &model.ContentDetail{
					ID: 0,
				},
			},
			wantErr: false, // 是否期望产生错误
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ContentDetailDao{
				db: tt.fields.db,
			}
			if err := c.Create(tt.args.detail); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
