package services

import (
	"beego/db"
	"fmt"
	"strconv"
	"time"

	"errors"

	"github.com/astaxie/beego"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	secret         []byte = []byte("admin@hortorgames.com")
	tokenKey              = "auth:token"
	ExpireDuration int64  = 1 * 24 * 60 * 60
)

// GenToken 生产token
func GenToken(uid int64) (token string, err error) {
	claims := make(jwt.MapClaims)
	claims["uid"] = uid

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		token = ""
		beego.Error(err)
		return
	}
	redisKey := fmt.Sprintf("%s:%d:%s", tokenKey, uid, token)
	// 写入到redis中
	err = db.Redis.HMSet(redisKey, map[string]interface{}{
		"uid":       uid,
		"lastLogin": int64(time.Now().Unix()),
	})
	db.Redis.Expire(redisKey, ExpireDuration)
	if err != nil {
		token = ""
		beego.Error(err)
		return
	}
	return
}

// ParseToken 解析token
func ParseToken(token string) (tokenInfo map[string]string, err error) {
	signedString, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		beego.Error("parse with claims failed.", err)
		return
	}
	signToken, ok := signedString.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("parse with claims failed")
		beego.Error(err)
		return
	}
	uid := int64(signToken["uid"].(float64))
	redisKey := fmt.Sprintf("%s:%d:%s", tokenKey, uid, token)
	tokenInfo, err = db.Redis.HGetAll(redisKey)
	if err != nil {
		beego.Error("parse with claims failed.", err)
		return
	}
	if tokenInfo == nil {
		err = errors.New("parse with claims failed")
		beego.Error(err)
		return
	}
	return
}

// RefreshToken 刷新token
func RefreshToken(uid int64, token string) (err error) {
	if token == "" {
		err = errors.New("refresh token failed")
		beego.Error(err)
		return
	}
	if uid == 0 {
		err = errors.New("refresh token failed")
		beego.Error(err)
		return
	}
	redisKey := fmt.Sprintf("%s:%d:%s", tokenKey, uid, token)
	err = db.Redis.HSet(redisKey, "lastLogin", strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		beego.Error("refresh token failed.", err)
		return
	}
	db.Redis.Expire(redisKey, ExpireDuration)
	return
}
