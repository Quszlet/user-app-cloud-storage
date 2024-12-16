package repository

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

type Config struct {
	User    string
	Passwd  string
	Net     string
	Addr    string
	Port 	string
	DBName  string
}

func NewMySqlDB(cfg Config) (*sqlx.DB, error) {
	//"username:password@tcp(127.0.0.1:3306)/dbname"
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", cfg.User, cfg.Passwd, cfg.Net, cfg.Addr, cfg.Port, cfg.DBName)

	// Открытие соединения
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln("Error connecting to the database:", err)
	}
	defer db.Close() // Обратить внимание при работе

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
