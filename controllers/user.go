package controllers

import (
	"social-app-server/models"
	"encoding/json"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.PhoneId
// @Failure 400 no enough input
// @Failure 403 body is empty
// @Failure 500 get products common error
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	PhoneId := models.AddUser(user)
	u.Data["json"] = map[string]string{"phone": PhoneId}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @Failure 500 server internal error
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by user PhoneId
// @Param	phone		path 	string	true		"电话号码作为主键"
// @Success 200 {object} models.User
// @Failure 403 :phone is empty
// @router /:PhoneId [get]
func (u *UserController) Get() {
	PhoneId := u.GetString(":phone")
	if PhoneId != "" {
		user, err := models.GetUser(PhoneId)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	PhoneId		path 	string	true		"The user you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:PhoneId [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	PhoneId		path 	string	true		"The user you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:PhoneId [delete]
func (u *UserController) Delete() {
	PhoneId := u.GetString(":phone")
	models.DeleteUser(PhoneId)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	phone		query 	string	true		"手机号登陆"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	PhoneId := u.GetString("phone")
	password := u.GetString("password")
	if models.Login(PhoneId, password) {
		u.Data["json"] = "login success"
	} else {
		u.Data["json"] = "user not exist"
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}
