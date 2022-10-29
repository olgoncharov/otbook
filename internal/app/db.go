package app

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/olgoncharov/otbook/config"
)

func initDB(cfg config.DBInstanceConfig) (*sql.DB, error) {
	mysqlConfig := mysql.Config{
		User:      cfg.User,
		Passwd:    cfg.Password,
		Net:       "tcp",
		Addr:      cfg.Host + ":" + cfg.Port,
		DBName:    cfg.DBName,
		Loc:       time.Local,
		ParseTime: true,
	}
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
