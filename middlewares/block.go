package middlewares

import (
	"framework/db"
	"framework/libs"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

const (
	blockUID    = "block:uid:list"
	blockOpenID = "block:openid:list"
	blockIP     = "block:ip:list"
)

const (
	blockTypeUID = iota
	blockTypeOpenID
	blockTypeIP
)

func isBlock(blockType int, val string) bool {
	var isMember bool
	var err error
	if blockType == blockTypeUID {
		isMember, err = db.Redis.SIsMember(blockUID, val)
	} else if blockType == blockTypeOpenID {
		isMember, err = db.Redis.SIsMember(blockOpenID, val)
	} else if blockType == blockTypeIP {
		isMember, err = db.Redis.SIsMember(blockIP, val)
	}
	return err != nil || isMember
}

var blockFilter = func(ctx *context.Context) {
	uid, _ := ctx.Input.GetData("uid").(int64)
	if uid != 0 && isBlock(blockTypeUID, strconv.FormatInt(uid, 10)) {
		ctx.Output.SetStatus(libs.ErrorForbiddenAccess.StatusCode)
		ctx.Output.Body([]byte(libs.ErrorForbiddenAccess.ErrorMsg))
		return
	}

	openid, _ := ctx.Input.GetData("openid").(string)
	if openid != "" && isBlock(blockTypeOpenID, openid) {
		ctx.Output.SetStatus(libs.ErrorForbiddenAccess.StatusCode)
		ctx.Output.Body([]byte(libs.ErrorForbiddenAccess.ErrorMsg))
		return
	}

	ip := ctx.Input.IP()
	if ip != "" && isBlock(blockTypeIP, ip) {
		ctx.Output.SetStatus(libs.ErrorForbiddenAccess.StatusCode)
		ctx.Output.Body([]byte(libs.ErrorForbiddenAccess.ErrorMsg))
		return
	}
}

func init() {
	beego.InsertFilter(`*`, beego.BeforeExec, blockFilter)
}
