package mysql

import (
	"database/sql"
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"crawl-worker/pkg/l"
)

type Config struct {
	Username string `json:"username,omitempty" mapstructure:"username"`
	Password string `json:"password,omitempty" mapstructure:"password"`
	Host     string `json:"host,omitempty" mapstructure:"host"`
	Port     string `json:"port,omitempty" mapstructure:"port"`
	Database string `json:"database,omitempty" mapstructure:"database"`
}

func New(cfg Config, ll l.Logger) *gorm.DB {
	createDatabase(cfg, ll)

	db := getConnection(cfg, ll)

	return db
}

func AutoMigration(db *gorm.DB, modelList []any, ll l.Logger) {
	for i, _ := range modelList {
		err := db.AutoMigrate(modelList[i])
		if err != nil {
			ll.Fatal("error when auto migrate db", l.Error(err))
		}
	}
}

func getConnection(cfg Config, ll l.Logger) *gorm.DB {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	dsnConfig := mysqldriver.Config{
		User:                 cfg.Username,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DBName:               cfg.Database,
		Params:               map[string]string{"charset": "utf8"},
		ParseTime:            true,
		Loc:                  loc,
		AllowNativePasswords: true,
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsnConfig.FormatDSN(),
	}), &gorm.Config{})
	if err != nil {
		ll.Fatal("error when init mysql db connection", l.Error(err))
	}

	sqlDb, err := db.DB()
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)
	return db
}

func createDatabase(cfg Config, ll l.Logger) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		ll.Fatal("cannot open connection to init database", l.Error(err))
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.Database)
	if err != nil {
		ll.Fatal("cannot create database when init", l.Error(err))
	}
}
