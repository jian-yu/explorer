package db

import "gopkg.in/mgo.v2"

type MgoOperator interface {
	GetSession() *mgo.Session
	GetDBConn() *mgo.Database
}
