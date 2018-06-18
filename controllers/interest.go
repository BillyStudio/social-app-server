package controllers

import (
	"github.com/astaxie/beego"
	"social-app-server/models"
	"encoding/json"
)

// Operations about interest
type InterestController struct {
	beego.Controller
}

// @Title GetAll
// @Description get all interest areas
// @Success 200 {object} models.Area
// @Failure 403 :InterestId is empty
// @router / [get]
func (o *InterestController) GetAll() {
	Arrays := models.GetAllInterests()
	o.Data["json"] = Arrays
	o.ServeJSON()
}

// @Title Upload
// @Description upload interests of users
// @Param	body	body 	models.UserInterest true "用户返回所有可用标签中感兴趣的标签，每个标签以＃开头"
// @Success 200 {int64} models.Moment.Id
// @Failure 403 body is empty
// @router /:UploadInterest [post]
func (controller *InterestController) Post() {
	var userInterest models.UserInterest
	json.Unmarshal(controller.Ctx.Input.RequestBody, &userInterest)
	if models.UploadInterest(userInterest) {
		// id为存储的文件名
		controller.Data["json"] = "upload success!"
	} else {
		controller.Data["json"] = "upload failure"
	}
	controller.ServeJSON()
}
