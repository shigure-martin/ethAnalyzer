package mysqlconn

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Update() {

}

func Select() {

}

func Delete() {

}

func Insert(db *sql.DB) bool {

	return true
}

func connectDB(uri string) *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", uri)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Main() {
	db := connectDB("root:123456@tcp(172.19.0.1:3306)/mev_bot_db")

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Fatal("connect success\n")
	}

	db.Close()
}
