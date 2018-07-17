package controllers

import (
	"strings"
	"encoding/json"

	"donorta-api/models"
	"donorta-api/lib/helper"
	"donorta-api/lib/security"
	app "donorta-api/app/user"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// User Modules
type UserController struct {
	BaseController
}

// @Title Create User
// @Description register
// @router / [post]
func (c *UserController) Register() {
	c.isValidSignature()

	var user models.User
	json.Unmarshal(c.Ctx.Input.RequestBody, &user)

	// Clean Input
	user.Handphone 	  = helper.CleanHPNo(user.Handphone)
	user.Email 	   	  = strings.ToLower(user.Email)

	result, code, err := app.Register(user)

	if err != nil {
		beego.Error(err.Error())
		beego.Error(result.Email)
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = c.getErrorResponse(code, err.Error(), EmptyStruct{})
		c.ServeJSON()
		return
	}

	// Generate token for session
	token, code, err := models.TokenGenerate(result)
	if err != nil {
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = c.getErrorResponse(code, err.Error(), EmptyStruct{})
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = c.getSuccessResponse(201, "object", token)
	c.ServeJSON()
}

// @Title Activate user using Response Code
// @Description Activate user
// TODO: @router /activate [post]
//func (c *UserController) Activate() {
//	var user models.User
//	json.Unmarshal(c.Ctx.Input.RequestBody, &user)
//	c.isValidSignature()
//
//	// Clean Input
//	user.Handphone 	  = helper.CleanHPNo(user.Handphone)
//	user.Email 	   	  = strings.ToLower(user.Email)
//
//	result, code, err := user.Activate(user)
//
//	if err != nil {
//		beego.Error(err.Error())
//		beego.Error(result.Email)
//		c.Ctx.Output.SetStatus(code)
//		c.Data["json"] = c.getErrorResponse(code, err.Error(), EmptyStruct{})
//		c.ServeJSON()
//		return
//	}
//
//	token, code, err := models.TokenGenerate(result)
//	if err != nil {
//		c.Ctx.Output.SetStatus(code)
//		c.Data["json"] = c.getErrorResponse(code, err.Error(), EmptyStruct{})
//		c.ServeJSON()
//		return
//	}
//
//	c.Ctx.Output.SetStatus(201)
//	c.Data["json"] = c.getSuccessResponse(201, "object", token)
//	c.ServeJSON()
//}

// @Title Login
// @Description customer login
// @router /login [post]
func (c *UserController) Login() {
	var user models.User
	json.Unmarshal(c.Ctx.Input.RequestBody, &user)

	user.Handphone = helper.CleanHPNo(user.Handphone)
	result, code, err := app.Login(user)

	if err != nil {
		beego.Error(err.Error())
		beego.Error(result.Email)
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = c.getErrorResponse(code, err.Error(), EmptyStruct{})
		c.ServeJSON()
		return
	}

	token, code, err := models.TokenGenerate(result)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = c.getErrorResponse(code, err.Error(), EmptyStruct{})
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = c.getSuccessResponse(200, "object", token)
	c.ServeJSON()
}

// @Title Logout
// @Description customer logout
// @router /logout [post]
func (c *UserController) Logout() {
	_, tok, _ := c.isTokenValid()

	tok.IsActive = 0
	orm.NewOrm().Update(&tok, "IsActive")

	c.Ctx.Output.SetStatus(204)
	c.Data["json"] = ""
	c.ServeJSON()
}

// @Title Forgot password get security question
// @Description Activate customer
// @router /forgot/:hp_or_email [get]
func (c *UserController) ForgotPassGet() {
	c.isValidSignature()

	id := c.Ctx.Input.Param(":hp_or_email")
	type Resp struct {
		SecurityQuestion 	string 		`json:"security_question"`
		ChallengeCode 		string 		`json:"challenge_code"`
		Handphone	 		string 		`json:"handphone"`
		Email	 			string 		`json:"email"`
	}
	var user models.User
	var resp Resp
	o := orm.NewOrm()

	id = helper.CleanHPNo(id)

	err := o.QueryTable("user").Filter("handphone", id).One(&user)
	if err == orm.ErrNoRows {
		err := o.QueryTable("user").Filter("email", id).One(&user)
		if err == orm.ErrNoRows {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = c.getErrorResponse(404, "Data anda tidak ditemukan", EmptyStruct{})
			c.ServeJSON()
			return
		}
	}

	resp.SecurityQuestion 	= user.SecurityQuestion
	resp.Handphone 			= user.Handphone
	resp.Email 				= user.Email

	c.Data["json"] = c.getSuccessResponse(200, "object", resp)
	c.ServeJSON()
}

// @Title Change password after forgot
// @Description Get Security Answer
// @router /forgot [post]
func (c *UserController) ChangePassAfterForgot() {
	c.isValidSignature()

	type Req struct {
		Username 			string 		`json:"username"`
		SecurityQuestion 	string 		`json:"security_question"`
		SecurityAnswer		string 		`json:"security_answer"`
		Password			string 		`json:"password"`
	}
	var req Req
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	var user models.User
	o := orm.NewOrm()

	// Clean HP number
	if beego.Substr(req.Username, 0, 1) == "0" {
		req.Username = "62"+ beego.Substr(req.Username, 1, 15)
	}

	// Validate User
	err := o.QueryTable("user").Filter("handphone", req.Username).One(&user)
	if err == orm.ErrNoRows {
		err := o.QueryTable("user").Filter("email", req.Username).One(&user)
		if err == orm.ErrNoRows {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = c.getErrorResponse(404, "Data anda tidak ditemukan", EmptyStruct{})
			c.ServeJSON()
			return
		}
	}

	ans := security.ShaOneEncrypt(strings.ToLower(req.SecurityAnswer) + user.SecuritySalt)
	if user.SecurityQuestion != req.SecurityQuestion || user.SecurityAnswer != ans {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = c.getErrorResponse(400, "Jawaban anda salah", EmptyStruct{})
		c.ServeJSON()
		return
	}

	user.Salt 		= security.ShaOneEncrypt(helper.GetNowTime().String())
	user.Password 	= security.ShaOneEncrypt(req.Password + user.Salt)
	user.UpdatedBy 	= user.Email
	o.Update(&user,"Salt","Password","UpdatedBy","UpdatedAt")

	c.Data["json"] = c.getSuccessResponse(200, "object", models.CleanUserData(user))
	c.ServeJSON()
}

// @Title Update User
// @Description Update my profile
// @router / [put]
func (c *UserController) Put() {
	current, _, _	:= c.isTokenValid()

	var UpdateValidate struct{
		OldPassword 	string `json:"old_password"`
	}
	var user models.User
	json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	json.Unmarshal(c.Ctx.Input.RequestBody, &UpdateValidate)

	// If change password then validate here
	if user.Password != "" && security.ShaOneEncrypt(UpdateValidate.OldPassword + current.Salt) != current.Password {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = c.getErrorResponse(400, "Password lama salah", EmptyStruct{})
		c.ServeJSON()
		return
	}

	result, code, err := app.Update(current, user)
	if err != nil {
		beego.Error(err.Error())
		beego.Error(result.Email)
		c.Ctx.Output.SetStatus(code)
		c.Data["json"] = c.getErrorResponse(code, err.Error(), EmptyStruct{})
		c.ServeJSON()
		return
	}

	c.Data["json"] = c.getSuccessResponse(200, "object", models.CleanUserData(result))
	c.ServeJSON()
}

// @Title Get User Activity Data
// @Description Get User Activity Data
// @router /activity [get]
func (c *UserController) Get() {
	user, _, _	:= c.isTokenValid()

	var userAct []models.UserActivity
	limit, _	:= c.GetInt("limit", 20)
	page, _		:= c.GetInt("page", 1)
	offset 		:= (page - 1) * limit

	o := orm.NewOrm()
	o.QueryTable("user_activity").
		Filter("user_id", user.Id).
		OrderBy("-id").
		Limit(limit).
		Offset(offset).
		All(&userAct)

	c.Data["json"] = c.getSuccessResponse(200, "array", userAct)
	c.ServeJSON()
}