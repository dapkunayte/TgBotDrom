package mysql_db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func mysqlConn(dbName string) *sql.DB {
	db, err := sql.Open("mysql", "root:447595Sql!@/users_state")
	if err != nil {
		panic(err)
	}
	return db
}

func InsertUserInfo(chatID int) {
	database := mysqlConn("users_state")
	result, err := database.Prepare("INSERT INTO users (id, state, town) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}
	_, err = result.Exec(chatID, "start", "tomsk")
	if err != nil {
		panic(err)
	}
	defer database.Close()
}
