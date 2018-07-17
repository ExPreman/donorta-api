package main

import (
	"time"

	"donorta-api/lib/helper"
	_ "donorta-api/routers"
	_ "donorta-api/models"
	_ "donorta-api/tasks"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql" // import your required driver
)

func init() {
	// set default database
	configDb := beego.AppConfig.String("mysql_user") +":"+ beego.AppConfig.String("mysql_pass") +"@tcp("+
		beego.AppConfig.String("mysql_url") +")/"+ beego.AppConfig.String("mysql_db") +"?charset=utf8&parseTime=true&loc=Asia%2fJakarta"

	orm.RegisterDriver("mysql", orm.DRMySQL)
	// GMT +7
	loc, _ := helper.GetLoc()
	orm.DefaultTimeLoc = loc
	// MaxIdleConns = 0, SetMaxOpenConns = 2400
	orm.RegisterDataBase("default", "mysql", configDb, 0, 2400)
	orm.SetDataBaseTZ("default", loc)

	// Connection Setting
	db, err := orm.GetDB("default")
	if err == nil {
		db.SetConnMaxLifetime(time.Second * 10)
	} else {
		beego.Error(err)
	}

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	beego.BConfig.WebConfig.AutoRender = false
	toolbox.StartTask()
}


func main() {
	beego.BConfig.EnableGzip = true
	// CORS for https://* origins, allowing:
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	beego.Run()
}

