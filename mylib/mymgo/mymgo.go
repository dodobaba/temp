package mymgo

import (
	"shoppingzone/mylib/mylog"

	"gopkg.in/mgo.v2"
)

const mgoURL = "127.0.0.1"

var (
	mgoSession *mgo.Session
	db         = "shoppingzone"
)

//GetMgoSession : get mongodb conection to use
func GetMgoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mgoURL)
		if err != nil {
			mylog.Tf("[Error]", "MyMgo", "GetMgoSession", "MongoDB connection fail! %s", err.Error())
		}
		mgoSession.SetPoolLimit(20)
	}
	return mgoSession.Clone()
}

//Do : do query
func Do(c string, f func(*mgo.Collection) error) error {
	s := GetMgoSession()
	defer s.Close()
	collection := s.DB(db).C(c)
	return f(collection)
}
