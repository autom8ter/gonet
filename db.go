/*
Create a Document

For example:

type Person struct {
	bongo.DocumentBase `bson:",inline"`
	FirstName string
	LastName string
	Gender string
}
You can use child structs as well.

type Person struct {
	bongo.DocumentBase `bson:",inline"`
	FirstName string
	LastName string
	Gender string
	HomeAddress struct {
		Street string
		Suite string
		City string
		State string
		Zip string
	}
}
*/

package gonet

import (
	"github.com/autom8ter/gonet/driver"
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type MongoDB struct {
	conn *bongo.Connection
}

func NewMongoDB(connectionStr, database string) *MongoDB {
	cfg := &bongo.Config{
		ConnectionString: connectionStr,
		Database:         database,
	}
	connection, err := bongo.Connect(cfg)

	if err != nil {
		log.Fatal(err)
	}
	return &MongoDB{connection}
}

func (m *MongoDB) GetConfig() *bongo.Config {
	return m.conn.Config
}

func (m *MongoDB) GetSession() *mgo.Session {
	return m.conn.Session
}

func (m *MongoDB) Collection(name string) *bongo.Collection {
	return m.conn.Collection(name)
}

func (m *MongoDB) Connect() error {
	return m.conn.Connect()
}

func (m *MongoDB) Ping() error {
	return m.conn.Session.Ping()
}

func (m *MongoDB) Close() {
	m.conn.Session.Close()
}

func (m *MongoDB) GetDB(name string) *mgo.Database {
	return m.conn.Session.DB(name)
}

func (m *MongoDB) SaveDoc(collname string, obj bongo.Document) error {
	return m.conn.Collection(collname).Save(obj)
}

func (m *MongoDB) PreSaveDoc(collname string, obj bongo.Document) error {
	return m.conn.Collection(collname).PreSave(obj)
}

func (m *MongoDB) DeleteDoc(collname string, query bson.M) (*mgo.ChangeInfo, error) {
	return m.conn.Collection(collname).Delete(query)
}

func (m *MongoDB) FindOne(collname string, query interface{}, doc interface{}) error {
	return m.conn.Collection(collname).FindOne(query, doc)
}

func (m *MongoDB) FindById(collname string, id bson.ObjectId, doc interface{}) error {
	return m.conn.Collection(collname).FindById(id, doc)
}

func (m *MongoDB) Find(collname string, query interface{}) *bongo.ResultSet {
	return m.conn.Collection(collname).Find(query)
}

func (m *MongoDB) DeleteOne(collname string, query bson.M) error {
	return m.conn.Collection(collname).DeleteOne(query)
}

func (m *MongoDB) DeleteDocument(collname string, doc bongo.Document) error {
	return m.conn.Collection(collname).DeleteDocument(doc)
}

func (m *MongoDB) AsDB() driver.DB {
	return m
}
