package middlewares

import (
	"beego/libs"
	"beego/services"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var whiteList = map[string][]string{
	"/v1/auth": []string{"GET"},
}

var (
	secret []byte = []byte("admin@xxx.com")
)

func isMatchRoute(route, path string) bool {
	re, _ := regexp.Compile(fmt.Sprintf("^%s$", route))
	return re.MatchString(path)
}

func isMatchMethod(methods []string, method string) bool {
	for i := 0; i < len(methods); i++ {
		if strings.ToUpper(methods[i]) == strings.ToUpper(method) {
			return true
		}
	}
	return false
}

var tokenValidateFilter = func(ctx *context.Context) {
	re, _ := regexp.Compile("/$")
	path := re.ReplaceAllString(ctx.Input.URL(), "")
	method := ctx.Input.Method()
	inWhiteList := false

	for route, methods := range whiteList {
		if !isMatchRoute(route, path) {
			continue
		}
		inWhiteList = isMatchMethod(methods, method)
		break
	}

	beego.Debug(`Router: [`, method, `::`, path, `]`, `Need token:`, !inWhiteList)
	if inWhiteList {
		// 不需要校验token
		return
	}
	token := ctx.Input.Header("X-Authorization")
	if token == "" {
		token = ctx.Input.Query("token")
	}
	//token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjF9.4Cjddbw39RW7DjCZ3_rmbk9gbRbBkBWtziouMS1idhg"
	tokenInfo, err := services.ParseToken(token)
	if err != nil || tokenInfo == nil {
		ctx.Output.SetStatus(libs.ErrorInvalidToken.StatusCode)
		ctx.Output.Body([]byte(libs.ErrorInvalidToken.ErrorMsg))
		return
	}

	lastLogin, err := strconv.ParseInt(tokenInfo["lastLogin"], 10, 64)
	uid, err := strconv.ParseInt(tokenInfo["uid"], 10, 64)
	if err != nil {
		ctx.Output.SetStatus(libs.ErrorInvalidToken.StatusCode)
		ctx.Output.Body([]byte(libs.ErrorInvalidToken.ErrorMsg))
		return
	}
	if lastLogin+services.ExpireDuration < int64(time.Now().Unix()) {
		// 过期
		ctx.Output.SetStatus(libs.ErrorInvalidToken.StatusCode)
		ctx.Output.Body([]byte(libs.ErrorInvalidToken.ErrorMsg))
		return
	}
	services.RefreshToken(uid, token)
	ctx.Input.SetData("uid", uid)
}

func init() {
	beego.InsertFilter(`*`, beego.BeforeRouter, tokenValidateFilter)
}
