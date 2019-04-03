package micro2

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"sync"
)

var MongoDBHost = "127.0.0.1"
var instance *mgo.Database
var once sync.Once



func GetCollection(collection string) *mgo.Collection {
	db := GetDBInstance()
	c := db.C(collection)
	return c

}


func GetDBInstance() *mgo.Database {
	once.Do(func() {
		instance = mongodbconnect()
	})
	return instance
}

func mongodbconnect() *mgo.Database {
	config := GetConfig()
	logStr := fmt.Sprintf("mongodb hostname : %s", config.DatabaseHost)
	WriteLog(logStr)
	session, err := mgo.Dial(MongoDBHost)
	if err != nil {
		WriteLog("can not connect to database")
		return nil
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(config.DatabaseName)
	return db
}
