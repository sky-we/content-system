package services

import (
	"content-system/internal/api"
	"content-system/internal/services"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ContentTestSuite struct {
	suite.Suite
	DbName    string
	Provider  *memory.DbProvider
	Table     *memory.Table
	Root      gin.IRoutes
	GinEngine *gin.Engine
	Db        *gorm.DB
	Rdb       *redis.Client
	App       *services.CmsApp
}

func (suite *ContentTestSuite) SetupTest() {
	suite.T().Log("Load go-mysql-server miniredis config")
	LoadFakeDBConfig()
	suite.T().Log("create cms_content.content_details table in go-mysql-server")
	dbName := "cms_content"
	tableName := "content_details"
	schema := sql.Schema{
		{Name: "id", Type: types.Int32, Nullable: false, Source: tableName, PrimaryKey: true, Comment: "主键ID", AutoIncrement: true},
		{Name: "title", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "内容标题"},
		{Name: "description", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "内容描述"},
		{Name: "author", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "作者"},
		{Name: "video_url", Type: types.Text, Nullable: false, Source: tableName, PrimaryKey: false, Comment: "视频播放URL"},
		{Name: "thumbnail", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "封面图URL"},
		{Name: "category", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "内容分类"},
		{Name: "duration", Type: types.Int64, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "内容时长"},
		{Name: "resolution", Type: types.Text, Nullable: false, Source: tableName, PrimaryKey: false, Comment: "分辨率 如720p、1080p"},
		{Name: "file_size", Type: types.Int64, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "文件大小"},
		{Name: "format", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "文件格式 如MP4、AVI"},
		{Name: "quality", Type: types.Int8, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "视频质量 1-高清 2-标清"},
		{Name: "approval_status", Type: types.Int8, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "审核状态 1-审核中 2-审核通过 3-审核不通过"},
		{Name: "created_at", Type: types.MustCreateDatetimeType(query.Type_DATETIME, 6), Nullable: false, Source: tableName},
		{Name: "updated_at", Type: types.MustCreateDatetimeType(query.Type_DATETIME, 6), Nullable: false, Source: tableName},
	}
	pro, table := CreateTestDatabase(dbName, tableName, schema)

	server, _ := FakeMysqlServer(pro)

	go func() {
		suite.T().Log("start go-mysql-server")
		if err := server.Start(); err != nil {
			fmt.Println("mysql fake server start error")
			panic(err)
		}
	}()
	defer func() {
		if err := server.Close(); err != nil {
			panic(err)
		}
	}()

	redisServer := miniredis.NewMiniRedis()
	suite.T().Log("start MiniRedis")

	if err := redisServer.StartAddr("localhost:6380"); err != nil {
		fmt.Println("redis fake server start error")
		panic(err)
	}
	suite.T().Cleanup(redisServer.Close)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	db := NewFakeMySqlDB(FakeDbCfg.MySQL)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	app := services.NewCmsApp(db, rdb)
	sessionMiddleware := &api.SessionAuth{Rdb: rdb}
	root := r.Group(RootPath).Use(sessionMiddleware.Auth)
	suite.DbName = dbName
	suite.Provider = pro
	suite.Table = table
	suite.GinEngine = r
	suite.Root = root
	suite.Db = db
	suite.Rdb = rdb
	suite.App = app

}
