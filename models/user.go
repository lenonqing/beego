package models

import (
	"beego/db"

	"github.com/astaxie/beego"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User 用户信息
type User struct {
	UID     int64  `bson:"uid"`
	Name    string `bson:"nickname"`
	HeadImg string `bson:"avatar"`
}

var initer = func(Mongo *mgo.Database) {
	user := &User{}
	c := Mongo.C(user.GetEntityName())
	c.EnsureIndex(mgo.Index{
		Key:        []string{`uid`},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	})
}

func init() {
	db.InsertModelIniter(&initer)
}

// GetEntityName 获取Collection名字
func (user *User) GetEntityName() string {
	return `users`
}

// FindUserByUID 通过uid来查找用户信息.
func FindUserByUID(uid int64) (user *User, err error) {
	user = &User{}
	witch := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(user)
	}
	err = db.Mongo.WitchCollection(user.GetEntityName(), witch)
	if err == mgo.ErrNotFound {
		user = nil
		err = nil
	}
	if err != nil {
		beego.Error("find user error, err:", err.Error())
		return nil, err
	}
	return
}
