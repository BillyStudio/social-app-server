package controllers

import (
	"social-app-server/models"
	"encoding/json"

	"github.com/astaxie/beego"
	"fmt"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.UserBasic	true		"body for user content"
// @Success 200 {string} models.UserBasic.id
// @Failure 400 no enough input
// @Failure 403 body is empty
// @Failure 500 get products common error
// @router /CreateUser [post]
func (u *UserController) Post() {
	var user models.UserBasic
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	userId, err := models.AddUser(user)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = userId
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} []models.UserBasic
// @Failure 500 server internal error
// @router /GetAllUser [get]
func (u *UserController) GetAllUsers() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	for _, u := range users {
		fmt.Printf("user: %#v\n", u);
	}
	u.ServeJSON()
}

// @Title Get
// @Description get user by phone id
// @Param	userId	path 	string	true "电话号码作为主键,但是并没有对字符串进行是否为电话号码的检查"
// @Success 200 {object} models.UserBasic
// @Failure 403 phone id is empty
// @router /:UserId [get]
func (u *UserController) Get() {
	userId := u.GetString(":userId")
	fmt.Printf("Debug phone:%v\n", userId)
	if userId != "" {
		user, err := models.GetUser(userId)
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
// @Param	uid		path 	string	true		"需要更新信息的用户手机号"
// @Param	body		body 	models.UserBasic	true		"body for user content"
// @Success 200 {object} models.UserBasic
// @Failure 403 :uid is not int
// @router /:UserId [put]
func (u *UserController) Put() {
	userId := u.GetString("uid")
	if userId != "" {
		var user models.UserBasic
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(userId, &user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	phone		query 	string	true		"手机号登陆"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} token
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	userId := u.GetString("phone")
	password := u.GetString("password")
	token, err := models.Login(userId, password)
	if err != nil {
		u.Data["json"] = "user not exist"
	} else {
		u.Data["json"] = token
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Param	token	query	string	true	"使用token（登陆令牌）注销用户"
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	Token := u.GetString("token")
	err := models.Logout(Token)
	if err != nil {
		fmt.Println(err.Error())
		u.Data["json"] = err
	} else {
		u.Data["json"] = "logout success"
	}
	u.ServeJSON()
}

