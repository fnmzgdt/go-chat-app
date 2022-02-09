package database

import (
	"log"

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

func (db *CassandraConnection) ExecuteQuery(query string, values ...interface{}) {
	if err := db.session.Query(query).Bind(values...).Exec(); err != nil {
		log.Fatal(err)
	}
}