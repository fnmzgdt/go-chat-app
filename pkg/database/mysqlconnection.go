package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConnection struct {
	db *sql.DB
}

func SetupMySQLConnection() (*MySQLConnection, error) {
	db, err := sql.Open("mysql", "root:+Zrtp2B&Eur27@/go_chat_app")

	if err != nil {
		return nil, err
	}
	return &MySQLConnection{db: db}, nil
}

func (s *MySQLConnection) ExecuteQuery(query string, values ...interface{}) (sql.Result, error) {
	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	r, err := stmt.Exec(values...)
	return r, nil
}

func (s *MySQLConnection) ExecuteGetUserByEmail(query string, values ...interface{}) (string, error) {
	var password string

	err := s.db.QueryRow(query, values...).Scan(&password)
	if err != nil {
		fmt.Println(err)
        return "", err
    }
	return password, nil
}
////////////////////////////

func (s *MySQLConnection) ExecuteGetUserThread(query string, values ...interface{}) (int64, error) {
	var threadId int64

	err := s.db.QueryRow(query, values...).Scan(&threadId)
	if err != nil {
        return 0, err
    }
	return threadId, nil
}