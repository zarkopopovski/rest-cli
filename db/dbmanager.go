package db;

import (
	"database/sql"
	//"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DBManager struct {
	DB *sql.DB
}

func (dbConnection *DBManager) OpenConnection() (err error) {
	db, err := sql.Open("sqlite3", "./rcli.db")
	if err != nil {
		panic(err)
	}

	dbConnection.DB = db

	dbConnection.setupInitialDatabase()

	return
}

func (dbConnection *DBManager) setupInitialDatabase() (err error) {
	statement, _ := dbConnection.DB.Prepare("CREATE TABLE IF NOT EXISTS collection (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR, date_created VARCHAR, date_modified VARCHAR)")
	statement.Exec()
	
	statement, _ = dbConnection.DB.Prepare("CREATE TABLE IF NOT EXISTS url_request (id INTEGER PRIMARY KEY AUTOINCREMENT, collection_id INTEGER, name VARCHAR, url VARCHAR, method VARCHAR, params_data TEXT, header_data TEXT, cookie_data TEXT, body_data TEXT, date_created VARCHAR, date_modified VARCHAR)")
	statement.Exec()
	
	return
}

