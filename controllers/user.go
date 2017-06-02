package controllers

import (
	"framework/libs"
	"framework/models"
)

// UserController 用户相关控制器
type UserController struct {
	BaseController
}

// @Title Get
// @Description find user by uid
// @Param	uid		path 	int64	true		"the uid you want to get"
// @Success 200 {user} entities.user
// @Failure 403 :uid is empty
// @router /:uid [get]
func (controller *UserController) Get() {
	uid, err := controller.GetInt64(":uid")
	if uid == 0 {
		controller.Halt(libs.ErrorMissParameter)
		return
	}
	if err != nil {
		controller.Halt(libs.ErrorInternalError)
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
	controller.Json(user)
}
