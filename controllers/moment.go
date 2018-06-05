package controllers

import (
	"social-app-server/models"
	"encoding/json"

	"github.com/astaxie/beego"
	"strconv"
	"fmt"
)

// Operations about moment
type MomentController struct {
	beego.Controller
}


// @Title Create
// @Description create and post a moment
// @Param	body		body 	models.MomentContent	true		"用户发送动态包括文字和图片两部分，服务器端未检查文字和图片均为空的情况，需要在客户端进行检查"
// @Success 200 {int} models.Moment.Id
// @Failure 403 body is empty
// @router / [post]
func (controller *MomentController) Post() {
	var raw models.MomentContent
	json.Unmarshal(controller.Ctx.Input.RequestBody, &raw)

	MomentId := models.AddOne(raw)
	// id为存储的文件名
	controller.Data["json"] = map[string]int64{"MomentId": MomentId}
	controller.ServeJSON()
}

// @Title Get
// @Description find momnent by MomentId
// @Param	MomentId		path 	string	true		"输入MomentId来获取某条动态"
// @Success 200 {object} models.Moment
// @Failure 403 :MomentId is empty
// @router /:MomentId [get]
func (controller *MomentController) Get() {
	StrId := controller.Ctx.Input.Param(":MomentId")
	if StrId != "" {
		id, err := strconv.Atoi(StrId)
		if err != nil {
			fmt.Println(err)
		} else {
			ob, err := models.GetOne(id)
			if err != nil {
				controller.Data["json"] = err.Error()
			} else {
				controller.Data["json"] = ob
			}
		}
	}
	controller.ServeJSON()
}

// @Title GetAll
// @Description get all moments
// @Success 200 {object} models.Moment
// @Failure 403 :MomentId is empty
// @router / [get]
func (o *MomentController) GetAll() {
	obs := models.GetAll()
	o.Data["json"] = obs
	o.ServeJSON()
}

// @Title Delete
// @Description delete the moment
// @Param	MomentId		path 	string	true		"The MomentId of the moment you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 MomentId is empty
// @router /:objectId [delete]
func (o *MomentController) Delete() {
	StrId := o.Ctx.Input.Param(":MomentId")
	id, err := strconv.Atoi(StrId)
	if err != nil {
		fmt.Println(err)
	} else {
		models.Delete(id)
		o.Data["json"] = "delete success!"
	}
	o.ServeJSON()
}

