// @APIVersion 0.1.3
// @Title social app api server
// @Description 使用Go语言编写的beego作为后端框架，连接MySQL数据库，Swagger UI作为API调试界面，Android为前端
// @Contact billy.ustb2016@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"social-app-server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/moment",
			beego.NSInclude(
				&controllers.MomentController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
