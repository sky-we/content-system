package services

import (
	"context"
	"fmt"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

const (
	RootPath    = "/api"
	OutRootPath = "/out/api"
)

func FakeMysqlServer(pro *memory.DbProvider) (s *server.Server, sqlEngine *sqle.Engine) {
	engine := sqle.NewDefault(pro)
	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", "localhost", 3307),
	}
	s, err := server.NewServer(config, engine, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		panic(err)
	}
	return s, engine

}

func CreateTestDatabase(dbName string, tableName string, schema sql.Schema) (provider *memory.DbProvider, table *memory.Table) {
	db := memory.NewDatabase(dbName)
	db.BaseDatabase.EnablePrimaryKeyIndexes()
	pro := memory.NewDBProvider(db)
	newTable := memory.NewTable(db, tableName, sql.NewPrimaryKeySchema(schema), db.GetForeignKeyCollection())
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
