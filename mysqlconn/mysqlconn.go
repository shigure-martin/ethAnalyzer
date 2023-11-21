package mysqlconn

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Tokens struct {
	Token_name string
	Token_addr string
}

func panic_err(err error) {
	if err != nil {
		panic(err)
	}
}

func Update() {

}

func Select_token(db *sql.DB, token Tokens) []Tokens {
	var tokens []Tokens

	sql := "select * from tokens where token_addr = ?"

	rows, err := db.Query(sql, token.Token_addr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var to Tokens
		var id int
		if err := rows.Scan(&id, &to.Token_name, &to.Token_addr); err != nil {
			log.Fatal(err)
		}
		tokens = append(tokens, to)
	}

	return tokens
}

func Delete() {

}

func Insert_token(db *sql.DB, token Tokens) {

	sql := "insert into tokens(token_name, token_addr) values(?,?)"
	stmt, err := db.Prepare(sql)
	panic_err(err)
	defer stmt.Close()

	result, err := stmt.Exec(token.Token_name, token.Token_addr)
	panic_err(err)

	row_number, err := result.LastInsertId()
	panic_err(err)

	log.Println("insert success ", row_number)
}

func ConnectDB(uri string) *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", uri)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Main() {
	db := ConnectDB("root:123456@tcp(172.19.0.1:3306)/mev_bot_db")

	var token Tokens
	token.Token_addr = "test_addr"

	result := Select_token(db, token)
	fmt.Println(len(result))
	for _, value := range result {
		fmt.Println(value.Token_name, value.Token_addr)
	}

	db.Close()
}
