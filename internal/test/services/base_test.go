package services

import (
	"context"
	"fmt"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	LoadFakeDBConfig()
	code := m.Run()
	os.Exit(code)
}

func GetMysqlFakeDBServer(dbName string, tableName string) (s *server.Server, provider *memory.DbProvider, table *memory.Table) {
	pro, newTable := createTestDatabase(dbName, tableName)
	engine := sqle.NewDefault(pro)

	session := memory.NewSession(sql.NewBaseSession(), pro)
	ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	ctx.SetCurrentDatabase("test")
	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", "localhost", 3307),
	}
	s, err := server.NewServer(config, engine, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		panic(err)
	}
	return s, pro, newTable

}

func createTestDatabase(dbName string, tableName string) (provider *memory.DbProvider, table *memory.Table) {
	db := memory.NewDatabase(dbName)
	db.BaseDatabase.EnablePrimaryKeyIndexes()
	pro := memory.NewDBProvider(db)
	newTable := memory.NewTable(db, tableName, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "id", Type: types.Int32, Nullable: false, Source: tableName, PrimaryKey: true, Comment: "主键ID", AutoIncrement: true},
		{Name: "user_id", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "用户id"},
		{Name: "pass_word", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "密码"},
		{Name: "nick_name", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "昵称"},
		{Name: "created_at", Type: types.MustCreateDatetimeType(query.Type_DATETIME, 6), Nullable: false, Source: tableName},
		{Name: "updated_at", Type: types.MustCreateDatetimeType(query.Type_DATETIME, 6), Nullable: false, Source: tableName},
	}), db.GetForeignKeyCollection())
	db.AddTable(tableName, newTable)
	return pro, newTable
}

func InsertData(dbName string, pro *memory.DbProvider, table *memory.Table, rowData sql.Row) {
	db := memory.NewDatabase(dbName)
	db.BaseDatabase.EnablePrimaryKeyIndexes()
	session := memory.NewSession(sql.NewBaseSession(), pro)
	ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	if err := table.Insert(ctx, rowData); err != nil {
		fmt.Println("insert data err", err)
		panic(err)
	}

}
