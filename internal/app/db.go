package app

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

func InitDB(cfg configer) *sql.DB {
	mysqlConfig := mysql.Config{
		User:      cfg.DBUser(),
		Passwd:    cfg.DBPassword(),
		Net:       "tcp",
		Addr:      cfg.DBHost() + ":" + cfg.DBPort(),
		DBName:    cfg.DBName(),
		Loc:       time.Local,
		ParseTime: true,
	}
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
