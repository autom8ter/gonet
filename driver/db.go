package driver

import (
	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB interface {
	GetConfig() *bongo.Config
	GetSession() *mgo.Session
	Connect() error
	Ping() error
	Close()
	GetDB(name string) *mgo.Database
	Collection(collname string) *bongo.Collection
	SaveDoc(collname string, obj bongo.Document) error
	PreSaveDoc(collname string, obj bongo.Document) error
	DeleteDoc(collname string, query bson.M) (*mgo.ChangeInfo, error)
	FindOne(collname string, query interface{}, doc interface{}) error
	FindById(collname string, id bson.ObjectId, doc interface{}) error
	Find(collname string, query interface{}) *bongo.ResultSet
	DeleteOne(collname string, query bson.M) error
	DeleteDocument(collname string, doc bongo.Document) error
	AsDB() DB
}

type SaveCollectionHook interface {
	BeforeSave(*bongo.Collection) error
	AfterSave(*bongo.Collection) error
}

type DeleteCollectionHook interface {
	BeforeDelete(*bongo.Collection) error
	AfterDelete(*bongo.Collection) error
}

type FindCollectionHook interface {
	AfterFind(*bongo.Collection) error
}

type ValidateCollectionHook interface {
	Validate(*bongo.Collection) []error
}
