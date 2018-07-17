package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["donorta-api/controllers:UserController"] = append(beego.GlobalControllerRouter["donorta-api/controllers:UserController"],
		beego.ControllerComments{
			Method: "Register",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["donorta-api/controllers:UserController"] = append(beego.GlobalControllerRouter["donorta-api/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["donorta-api/controllers:UserController"] = append(beego.GlobalControllerRouter["donorta-api/controllers:UserController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/activity`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["donorta-api/controllers:UserController"] = append(beego.GlobalControllerRouter["donorta-api/controllers:UserController"],
		beego.ControllerComments{
			Method: "ChangePassAfterForgot",
			Router: `/forgot`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["donorta-api/controllers:UserController"] = append(beego.GlobalControllerRouter["donorta-api/controllers:UserController"],
		beego.ControllerComments{
			Method: "ForgotPassGet",
			Router: `/forgot/:hp_or_email`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["donorta-api/controllers:UserController"] = append(beego.GlobalControllerRouter["donorta-api/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["donorta-api/controllers:UserController"] = append(beego.GlobalControllerRouter["donorta-api/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
