package database

import (
	"log"
	"project/pkg/errors"

	"github.com/gocql/gocql"
)

type CassandraConnection struct {
	session *gocql.Session
}

func SetupCassandraConnection() *CassandraConnection {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "test1"
	session, _ := cluster.CreateSession()
	return &CassandraConnection{session: session}
}

func (db *CassandraConnection) ExecuteCreateUser(query string, values ...interface{}) {
		if err := db.session.Query(query).Bind(values...).Exec(); err != nil {
			log.Fatal(err)
		}
}

func (db *CassandraConnection) ExecuteGetUserByEmail(query string, values ...interface{}) (string, *errors.RestError) {
	var message string
	var password string
	if err := db.session.Query(query).Bind(values...).Scan(&password); err != nil {
		message = "Invalid username or password"
		return password, errors.NewBadRequestError(message)
	}
	return password, nil
}