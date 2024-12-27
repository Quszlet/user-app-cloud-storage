package repository

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	Port   string
	DBName string
}

func NewMySqlDB(cfg Config) (*gorm.DB, error) {
	//"username:password@tcp(127.0.0.1:3306)/dbname"
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", cfg.User, cfg.Passwd, cfg.Net, cfg.Addr, cfg.Port, cfg.DBName)
	fmt.Println(dsn)

	// Открытие соединения
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error connecting to the database:", err)
	}

	return db, nil
}
