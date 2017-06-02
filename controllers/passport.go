package controllers

import (
	"beego/libs"
	"beego/models"
	"beego/services"
)

// PassportController 通行证相关控制器
type PassportController struct {
	BaseController
}

// @Title Get
// @Description auth user by uid
// @Param	uid	    query 	int64	true		"the uid you want to auth"
// @Success 200 {token} string
// @Failure 403 :uid is empty
// @router / [get]
func (controller *PassportController) Get() {
	uid, _ := controller.GetInt64("uid")
	if uid == 0 {
		controller.Halt(libs.ErrorMissParameter)
		return
	}
	user, err := models.FindUserByUID(uid)
	if err != nil {
		controller.Halt(libs.ErrorInternalError)
		return
	}
	if user == nil {
		controller.Halt(libs.ErrorNotFound)
		return
	}
	token, err := services.GenToken(uid)
	if err != nil {
		controller.Halt(libs.ErrorInternalError)
		return
	}
	if token == "" {
		controller.Halt(libs.ErrorNotFound)
		return
	}
	controller.Json(token)
}
