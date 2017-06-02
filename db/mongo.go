package db

import (
	"gopkg.in/mgo.v2"
)

// MongoDB 数据库实例结构.
type MongoDB struct {
	session *mgo.Session
	host    string
	dbName  string
	mongo   *mgo.Database
}

// modelIniters 实体初始化函数集合
var modelIniters = []*func(*mgo.Database){}

// InsertModelIniter 插入初始化函数
func InsertModelIniter(initer *func(*mgo.Database)) {
	modelIniters = append(modelIniters, initer)
}

// CreateMongoDB 创建mongo数据库连接实例.
func CreateMongoDB(url, dbname string) (db *MongoDB, err error) {
	var session *mgo.Session
	session, err = mgo.Dial(url)
	if err != nil {
		return
	}

	db = &MongoDB{
		session: session,
		host:    url,
		dbName:  dbname,
		mongo:   session.DB(dbname),
	}
	session.SetMode(mgo.Strong, true)

	for i := 0; i < len(modelIniters); i++ {
		(*modelIniters[i])(db.mongo)
	}
	return
}

// WitchCollection collection执行语句
func (Mongo *MongoDB) WitchCollection(collection string, witch func(*mgo.Collection) error) error {
	session := Mongo.session.Clone()
	defer session.Close()
	c := session.DB(Mongo.dbName).C(collection)
	return witch(c)
}
