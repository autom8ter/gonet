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

package db

import (
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type MongoDB struct {
	collection string
	conn       *bongo.Connection
}

func NewMongoDB(collName, connectionStr, database string) *MongoDB {
	cfg := &bongo.Config{
		ConnectionString: connectionStr,
		Database:         database,
	}
	connection, err := bongo.Connect(cfg)

	if err != nil {
		log.Fatal(err)
	}
	return &MongoDB{
		collection: collName,
		conn:       connection,
	}
}

func (m *MongoDB) SwitchCollection(name string) {
	m.collection = name
}

func (m *MongoDB) SwitchConnection(connection *bongo.Connection) {
	m.conn = connection
}

func (m *MongoDB) Config() *bongo.Config {
	return m.conn.Config
}

func (m *MongoDB) Session() *mgo.Session {
	return m.conn.Session
}

func (m *MongoDB) Collection() *bongo.Collection {
	return m.conn.Collection(m.collection)
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

func (m *MongoDB) SaveDoc(obj bongo.Document) error {
	return m.conn.Collection(m.collection).Save(obj)
}

func (m *MongoDB) PreSaveDoc(obj bongo.Document) error {
	return m.conn.Collection(m.collection).PreSave(obj)
}

func (m *MongoDB) DeleteDoc(query bson.M) (*mgo.ChangeInfo, error) {
	return m.conn.Collection(m.collection).Delete(query)
}

func (m *MongoDB) FindDoc(query interface{}, doc interface{}) error {
	return m.conn.Collection(m.collection).FindOne(query, doc)
}

func (m *MongoDB) FindByDocId(id bson.ObjectId, doc interface{}) error {
	return m.conn.Collection(m.collection).FindById(id, doc)
}

func (m *MongoDB) Query(query interface{}) *bongo.ResultSet {
	return m.conn.Collection(m.collection).Find(query)
}

func (m *MongoDB) DeleteDocument(doc bongo.Document) error {
	return m.conn.Collection(m.collection).DeleteDocument(doc)
}
