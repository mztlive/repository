package database

import (
	"sync"
	// 单元测试sqlx使用的是sqlite驱动
	_ "github.com/mattn/go-sqlite3"

	// 程序要使用mysql驱动
	_ "github.com/go-sql-driver/mysql"

	gkcfg "github.com/gookit/config/v2"
	"github.com/jmoiron/sqlx"
)

var (
	db   *sqlx.DB
	once sync.Once
)

// GetDB 获得SQLX的DB实例
func GetDB() *sqlx.DB {
	once.Do(func() {
		var err error
		dbDriver := gkcfg.String("DB.DRIVER", "sqlite3")
		dbURI := gkcfg.String("DB.URI", ":memory?parseTime=true")

		maxOpenConns := gkcfg.Int("DB.MaxOpenConns", 30)
		maxIdleConns := gkcfg.Int("DB.MaxIdleConns", 30)

		db, err = sqlx.Open(dbDriver, dbURI)

		if err != nil {
			panic(err)
		}

		db.SetMaxIdleConns(maxIdleConns)
		db.SetMaxOpenConns(maxOpenConns)
	})

	err := db.Ping()
	if err != nil {
		panic("Databaes Connection Failed.")
	}

	return db
}
