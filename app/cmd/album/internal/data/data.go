package data

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var DB *sql.DB
var cfg mysql.Config

func InitDB() error {

	cfg = mysql.Config{
		User:   viper.GetString("db.username"),
		Passwd: viper.GetString("db.password"),
		Net:    "tcp",
		Addr:   viper.GetString("db.addr"),
		DBName: viper.GetString("db.name"),
	}
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	DB = db

	return err
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type ResponseModel struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
