package controllers

import (
	"github.com/astaxie/beego"
	"social-app-server/models"
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
