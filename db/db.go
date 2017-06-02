package db

import "github.com/astaxie/beego"

// Redis Redis实例
var Redis *RedisDB

// Mongo Mongo实例
var Mongo *MongoDB

// Init 初始化数据库
func Init() (err error) {
	Redis, err = CreateRedisDB(beego.AppConfig.String(`redisHost`), beego.AppConfig.String(`redisAuth`))
	if err != nil {
		return err
	}
	Mongo, err = CreateMongoDB(beego.AppConfig.String(`mongoHost`), beego.AppConfig.String(`mongoDBName`))

	return nil
}
