package ushare

import (
	"gopkg.in/mgo.v2"
	"ushare/core"
)

type Database struct {
	Addr     string
	Session  *mgo.Session
	Database *mgo.Database
}

func (db *Database) Connect(dbAddr, dbName string) error {
	if db.Session != nil {
		db.Session.Close()
		db.Session = nil
	}

	session, err := mgo.Dial(dbAddr)
	if err != nil {
		return core.NewStdErr(core.ERR_DBNotAvail, err.Error())
	}

	db.Addr = dbAddr
	db.Session = session
	db.Database = session.DB(dbName)
	return nil
}

func (db *Database) Insert(collectionName string, record interface{}) error {
	if db == nil || db.Database == nil {
		return core.NewStdErr(core.ERR_DatabaseInsertionFailed, "database not available for accessing.")
	}

	collection := db.Database.C(collectionName)
	// 按照 mgo 的文档，这里我们假定总是返回有效的 collection
	// if collection == nil {
	// 	return "", core.NewStdErr(core.ERR_DatabaseInsertionFailed, "collection not found. DB Collection Name: "+ushare.DBCName_SharingSession)
	// }

	if err := collection.Insert(record); err != nil {
		return core.NewStdErr(core.ERR_DatabaseInsertionFailed, "insertion failed："+err.Error())
	}

	return nil
}

func (db *Database) Close() {
	if db == nil || db.Session == nil {
		return
	}

	db.Session.Close()
}
