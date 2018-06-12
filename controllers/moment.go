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
// @Param	body	body 	models.MomentContent true "用户发送动态包括文字,图片,标签，标签以空格分开，服务器端未检查三部分均为空的情况，需要在客户端进行检查"
// @Success 200 {int64} models.Moment.Id
// @Failure 403 body is empty
// @router / [post]
func (controller *MomentController) Post() {
	var raw models.MomentContent
	json.Unmarshal(controller.Ctx.Input.RequestBody, &raw)
	MomentId := models.AddOne(raw)
	// id为存储的文件名
	controller.Data["json"] = map[string]string{"MomentId": strconv.FormatInt(MomentId, 10)}
	controller.ServeJSON()
}

// @Title Get
// @Description find moment by MomentId
// @Param	Token		query	string 	true 	"通过Token确认用户已登陆"
// @Param	MomentId		path 	string	true		"输入MomentId来获取某条动态"
// @Success 200 {object} models.MomentContent
// @Failure 403 :MomentId is empty
// @router /:MomentId [get]
func (controller *MomentController) Get() {
	Token := controller.GetString("Token")
	fmt.Println("Get Token: ", Token)
	StrId := controller.Ctx.Input.Param(":MomentId")
	if StrId != "" && Token != "" {
		MomentId, err := strconv.ParseInt(StrId, 10, 64)
		if err != nil {
			fmt.Println(err)
		} else {
			mo, err := models.GetOne(Token, MomentId)
			if err != nil {
				controller.Data["json"] = err.Error()
			} else {
				controller.Data["json"] = mo
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
// @router /:MomentId [delete]
func (o *MomentController) Delete() {
	StrId := o.Ctx.Input.Param(":MomentId")
	id, err := strconv.ParseInt(StrId, 10, 64)
	if err != nil {
		panic(err.Error())
	} else if models.Delete(id) == true {
		o.Data["json"] = "delete success!"
	} else {
		o.Data["json"] = "delete failure!"
	}
	o.ServeJSON()
}

