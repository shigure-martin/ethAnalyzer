package mysqlconn

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Tokens struct {
	Id         int
	Token_name string
	Token_addr string
}

type Pools struct {
	Id           int
	Pool_addr    string
	Factory_addr string
	Protocol     string
}

func panic_err(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func Update() {

}

func Select_token(db *sql.DB, token Tokens) Tokens {
	var to Tokens
	to.Id = 0

	sql := "select * from tokens where token_addr = ?"

	rows, err := db.Query(sql, token.Token_addr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&to.Id, &to.Token_name, &to.Token_addr); err != nil {
			log.Fatal(err)
		}
	}

	return to
}

func Select_pool(db *sql.DB, pool Pools) Pools {
	var po Pools
	po.Id = 0

	sql := "select * from pools where pool_addr = ?"

	rows, err := db.Query(sql, pool.Pool_addr)
	panic_err(err)

	for rows.Next() {
		if err := rows.Scan(&po.Id, &po.Pool_addr, &po.Protocol, &po.Factory_addr); err != nil {
			log.Fatal(err)
		}
	}

	return po
}

func Select_token_by_id(db *sql.DB, id int) Tokens {
	var token Tokens

	sql := "select * from tokens where id = ?"

	rows, err := db.Query(sql, id)
	panic_err(err)
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&token.Id, &token.Token_name, &token.Token_addr); err != nil {
			log.Fatal(err)
		}
	}

	return token
}

func Select_pool_by_id(db *sql.DB, id int) Pools {
	var pool Pools

	sql := "select * from pools where id = ?"

	rows, err := db.Query(sql, id)
	panic_err(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&pool.Id, &pool.Pool_addr, &pool.Protocol, &pool.Factory_addr)
		panic_err(err)

	}

	return pool
}

func Delete() {

}

func Insert_token(db *sql.DB, token Tokens) int {

	sql := "insert into tokens(token_name, token_addr) values(?,?)"
	stmt, err := db.Prepare(sql)
	panic_err(err)
	defer stmt.Close()

	result, err := stmt.Exec(token.Token_name, token.Token_addr)
	panic_err(err)

	row_number, err := result.LastInsertId()
	panic_err(err)

	log.Println("Token insert success ", row_number)
	return int(row_number)
}

func Insert_pool(db *sql.DB, pool Pools) {
	sql := "insert into pools(pool_addr, protocol, factory_addr) values(?,?,?)"
	stmt, err := db.Prepare(sql)
	panic_err(err)
	defer stmt.Close()

	result, err := stmt.Exec(pool.Pool_addr, pool.Protocol, pool.Factory_addr)
	panic_err(err)

	row_number, err := result.LastInsertId()
	panic_err(err)

	log.Println("Pool insert success ", row_number)
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

	var pool Pools
	pool.Id = 1

	pool = Select_pool_by_id(db, pool.Id)

	fmt.Println("Pool id: ", pool.Id)
	fmt.Println("Pool addr: ", pool.Pool_addr)
	fmt.Println("Pool prot: ", pool.Protocol)
	fmt.Println("Pool fact: ", pool.Factory_addr)

	db.Close()
}
