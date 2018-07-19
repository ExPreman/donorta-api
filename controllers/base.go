package controllers

import (
	"time"
	"math"
	"bytes"
	"strings"
	"strconv"

	"donorta-api/models"
	"donorta-api/lib/helper"
	"donorta-api/lib/security"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)

type BaseController struct {
	beego.Controller
}

type AdminSuccessResponse struct {
	Code 		int			`json:"code"`		// http code for easier read, must be the same with HEADER
	Type 		string		`json:"type"`		// return type, array or object
	Data 		interface{}	`json:"data"`		// can be an array of object or single object
	Total 		int			`json:"total"`		// total data
	FirstPage 	string 		`json:"first_page"`	// first page if type = array
	LastPage 	string 		`json:"last_page"`	// last page if type = array
	NextPage 	string 		`json:"next_page"`	// next page if type = array
	Timestamp 	time.Time	`json:"timestamp"`	// current timestamp
}

type SuccessResponse struct {
	Code 		int			`json:"code"`		// http code for easier read, must be the same with HEADER
	Type 		string		`json:"type"`		// return type, array or object
	Data 		interface{}	`json:"data"`		// can be an array of object or single object
	Timestamp 	time.Time	`json:"timestamp"`	// current timestamp
}

type ErrorResponse struct {
	ErrorCode 		int			`json:"error_code"`
	ErrorMessage 	string		`json:"error_message"`
	ErrorData 		interface{}	`json:"error_data"`
	Timestamp 		time.Time	`json:"timestamp"`
}

type EmptyStruct struct {

}

// init function, preRequest
func (c BaseController) Prepare (){
	beego.SetLogger(logs.AdapterFile, `{"filename":"logs/main.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":30}`)
	beego.SetLogFuncCall(true)

	//uri := c.Ctx.Input.URI()

	if beego.BConfig.RunMode == "dev" {
		buf := new(bytes.Buffer)
		beego.Debug(c.Ctx.Request.Header)

		// debug body if exist
		if c.Ctx.Request.Body != nil {
			buf.ReadFrom(c.Ctx.Request.Body)
			newStr := buf.String()
			beego.Debug(newStr)
		}
	}

	// TODO: Admin security
	//if strings.Contains(uri, "/admin") == true {
	//	if strings.Contains(uri, "login") == true {
	//		c.isValidSignature()
	//	} else if strings.Contains(uri, "logout") == true {
	//		c.isValidSignature()
	//		//c.isAdminValid()
	//	} else {
	//		c.isValidSignature()
	//	}
	//}
}

// Get app success response
func (c *BaseController) getSuccessResponse(code int, tipe string, data interface{}) SuccessResponse{
	s := SuccessResponse{}
	s.Code 		= code
	s.Type 		= tipe
	s.Data 		= data
	s.Timestamp = helper.GetNowTime()

	return s
}

// Get admin success response with pagination
func (c *BaseController) getAdminSuccessResponse(code int, tipe string, data interface{}, total int, params []string) AdminSuccessResponse{
	s := AdminSuccessResponse{}
	s.Code 		= code
	s.Type 		= tipe
	s.Data 		= data
	s.Total 	= total
	s.NextPage 	= ""
	s.FirstPage = ""
	s.LastPage	= ""

	// URL for Pagination
	if len(params) > 0 {
		baseUrl := c.Ctx.Input.Site() + c.Ctx.Input.URL() + "?"
		page, _  := strconv.Atoi(params[0])
		limit, _ := strconv.Atoi(params[1])
		lp := math.Ceil(float64(total)/float64(limit))
		np := page + 1

		// Delete page and limit params
		params = append(params[:0], params[2:]...)
		if np > page && np <= int(lp) {
			s.NextPage = baseUrl +"limit="+ strconv.Itoa(limit) +"&page="+ strconv.Itoa(np) + strings.Join(params,"")
		}

		s.FirstPage = baseUrl +"limit="+ strconv.Itoa(limit) +"&page="+ strconv.Itoa(1) + strings.Join(params,"")
		if lp != 0 {
			s.LastPage  = baseUrl +"limit="+ strconv.Itoa(limit) +"&page="+ strconv.FormatFloat(lp, 'f', -1, 64) + strings.Join(params,"")
		}
	}
	s.Timestamp = helper.GetNowTime()

	return s
}

// Get default error response
func (c *BaseController) getErrorResponse(code int, message string, data interface{}) ErrorResponse {
	err := ErrorResponse{}
	err.Timestamp 	 = helper.GetNowTime()
	err.ErrorCode 	 = code
	err.ErrorMessage = message
	err.ErrorData	 = data

	return err
}

// Check whether signature is valid
func (c *BaseController) isValidSignature() bool {
	header 		:= c.Ctx.Request.Header
	timestamp	:= header.Get("Timestamp")
	signature 	:= header.Get("Signature")

	if timestamp == "" || signature == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = c.getErrorResponse(401, "Unauthorized Signature", EmptyStruct{})
		c.ServeJSON()
		c.StopRun()

		return false
	}

	if signature != security.ShaOneEncrypt(timestamp + beego.AppConfig.String("SIGNATURE")) {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = c.getErrorResponse(401, "Unauthorized Signature", EmptyStruct{})
		c.ServeJSON()
		c.StopRun()
		return false
	}

	return true
}

// Check whether customer token is valid and return it
func (c *BaseController) isTokenValid() (models.User, models.UserToken, bool) {
	var user models.User
	var tok models.UserToken
	header 	:= c.Ctx.Request.Header
	token	:= header.Get("User-Token")

	if token == "" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = c.getErrorResponse(401, "Unauthorized Token 1", EmptyStruct{})
		c.ServeJSON()
		c.StopRun()
	}

	o := orm.NewOrm()
	// Get token user
	tok = models.GetUserToken(token)
	if tok.Id == 0 {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = c.getErrorResponse(401, "Unauthorized Token 2", EmptyStruct{})
		c.ServeJSON()
		c.StopRun()
	}

	// Get user base on token
	user = models.GetActiveUser(tok.User.Id)
	if user.Id == 0 {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = c.getErrorResponse(401, "Unauthorized Token 3", EmptyStruct{})

		// For Single User
		o.Delete(&tok)

		// For Multi User
		//tok.IsActive = 0
		//o.Update(&tok,"IsActive")

		c.ServeJSON()
		c.StopRun()
	}

	return user, tok, true
}

// Check whether admin token is valid
//func (c *BaseController) isAdminValid() (models.Admin, bool) {
//	var admin models.Admin
//	o := orm.NewOrm()
//	header 	:= c.Ctx.Request.Header
//	token	:= header.Get("Admin-Token")
//
//	if token == "" {
//		c.Data["json"] = c.getErrorResponse(401, "Unauthorized Admin Token", EmptyStruct{})
//		c.ServeJSON()
//		c.StopRun()
//		return admin, false
//	}
//
//	err := o.QueryTable("admin").Filter("token", token).Filter("is_active", 1).One(&admin)
//	if err == orm.ErrNoRows {
//		c.Data["json"] = c.getErrorResponse(401, "Unauthorized Admin Token.", EmptyStruct{})
//		c.ServeJSON()
//		c.StopRun()
//		return admin, false
//	}
//
//	admin.UpdatedBy = admin.Username
//	admin.UpdatedAt = helper.GetNowTime()
//	o.Update(&admin)
//
//	return admin, true
//}