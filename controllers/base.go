package controllers

import (
	"framework/libs"

	"github.com/astaxie/beego"
)

// BaseController 用户相关控制器
type BaseController struct {
	beego.Controller
}

// Json 发送JSON数据
func (controller *BaseController) Json(data interface{}) {
	controller.Data["json"] = map[string]interface{}{
		"errcode": 0,
		"errmsg":  "",
		"data":    data,
	}
	controller.ServeJSON()
}

// Halt 错误终止
func (controller *BaseController) Halt(e libs.ErrorType) {
	controller.Ctx.ResponseWriter.WriteHeader(e.StatusCode)
	controller.Data["json"] = map[string]interface{}{
		"errcode": e.ErrorCode,
		"errmsg":  e.ErrorMsg,
	}
	controller.ServeJSON()
}

// HaltError 错误终止
func (controller *BaseController) HaltError(err error) {
	controller.Halt(libs.MakeError(err))
}
