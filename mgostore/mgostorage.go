package mgostore

import (
	"github.com/Scalingo/osin"

	"gopkg.in/mgo.v2"
)

// collection names for the entities
const (
	CLIENT_COL    = "clients"
	AUTHORIZE_COL = "authorizations"
	ACCESS_COL    = "accesses"
)

type MongoStorage struct {
	dbName  string
	session *mgo.Session
}

func New(session *mgo.Session, dbName string) *MongoStorage {
	storage := &MongoStorage{
		dbName:  dbName,
		session: session.Copy(),
	}
	return storage
}

func (store *MongoStorage) Clone() osin.Storage {
	return &MongoStorage{
		dbName:  store.dbName,
		session: store.session.Clone(),
	}
}

func (store *MongoStorage) Close() {
	store.session.Close()
}
