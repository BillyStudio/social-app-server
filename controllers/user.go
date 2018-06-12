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
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {string} models.User.PhoneId
// @Failure 400 no enough input
// @Failure 403 body is empty
// @Failure 500 get products common error
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	PhoneId, err := models.AddUser(user)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = PhoneId
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} []models.User
// @Failure 500 server internal error
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by phone id
// @Param	PhoneId	path 	string	true "电话号码作为主键,但是并没有对字符串进行是否为电话号码的检查"
// @Success 200 {object} models.User
// @Failure 403 phone id is empty
// @router /:PhoneId [get]
func (u *UserController) Get() {
	PhoneId := u.GetString(":PhoneId")
	fmt.Printf("Debug phone:%v\n", PhoneId)
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
// @Param	phone		path 	string	true		"需要更新信息的用户手机号"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:PhoneId [put]
func (u *UserController) Put() {
	PhoneId := u.GetString("phone")
	if PhoneId != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(PhoneId, &user)
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
	PhoneId := u.GetString("phone")
	password := u.GetString("password")
	token, err := models.Login(PhoneId, password)
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

