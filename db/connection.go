package db

import (
	"fmt"
	"time"
	
	"github.com/gocraft/dbr"
)

var defaultEventReceiver = EventReceiver{}

type MysqlConfig struct {
	Host     string 
	Port     string 
	User     string 
	Password string 
	Database string 
}

func (m *MysqlConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

func OpenDatabase(cfg MysqlConfig) (*Database, error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	conn, err := dbr.Open("mysql", cfg.GetUrl()+"?parseTime=1&multiStatements=1&charset=utf8mb4&collation=utf8mb4_unicode_ci", &defaultEventReceiver)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(100)
	conn.SetMaxOpenConns(100)
	conn.SetConnMaxLifetime(10 * time.Second)
	return &Database{conn.NewSession(nil)}, nil
}