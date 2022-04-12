package database

import (
	"database/sql"
	"fmt"
	"project/pkg/messages"

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
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return r, nil
}

func (s *MySQLConnection) ExecuteGetUserPasswordFromEmail(query string, values ...interface{}) (string, error) {
	var password string

	err := s.db.QueryRow(query, values...).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (s *MySQLConnection) ExecuteGetProfileImageFromUserId(query string, userId int) (string, error) {
	var imageUrl sql.NullString

	err := s.db.QueryRow(query, userId).Scan(&imageUrl)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return imageUrl.String, nil
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

func (s *MySQLConnection) ExecuteUpdateUserThreadQuery(query string, values []int) (sql.Result, error) {
	conv := make([]interface{}, len(values))
	for i, v := range values {
		conv[i] = v
	}

	stmt, err := s.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	r, err := stmt.Exec(conv...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return r, nil
}

func (s *MySQLConnection) ExecuteGetUsersInThread(query string, threadId int) ([]int, error) {
	userIds := []int{}

	rows, err := s.db.Query(query, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		err := rows.Scan(&i)
		if err != nil {
			fmt.Println(err)
			return []int{}, nil
		}

		userIds = append(userIds, i)
	}

	return userIds, nil
}

func (s *MySQLConnection) ExecuteGetMessagesFromThread(query string, threadid int) ([]messages.MessageGet, error) {
	var messagesArray []messages.MessageGet
	rows, err := s.db.Query(query, threadid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message messages.MessageGet
		err := rows.Scan(&message.Id, &message.ThreadId, &message.FromId, &message.Date, &message.MessageText)
		if err != nil {
			fmt.Println(err)
		}
		messagesArray = append(messagesArray, message)
	}
	return messagesArray, nil
}

func (s *MySQLConnection) ExecuteGetLatestThreads(query string, userid int) ([]messages.ThreadGet, error) {
	var threadsArray []messages.ThreadGet

	rows, err := s.db.Query(query, userid, userid, userid, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var thread messages.ThreadGet
		err := rows.Scan(&thread.Id, &thread.Name, &thread.Seen, &thread.LastSender, &thread.LastMessage, &thread.LastUpdated, &thread.Type, &thread.Users)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		threadsArray = append(threadsArray, thread)
	}

	return threadsArray, nil
}
