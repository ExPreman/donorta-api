package tasks

import (
	"time"

	"donorta-api/models"
	"donorta-api/lib/helper"
	"donorta-api/lib/constanta"

	"github.com/astaxie/beego/toolbox"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// Run every minutes
func init () {
	unlockUser := toolbox.NewTask("unlock_user", "2 * * * * *", func() error {
		beego.Debug("RUN TASK unlock_user")

		o := orm.NewOrm()
		var users []models.User

		// 5 Minutes ago
		ago := helper.GetNowTime().Add(-time.Minute * time.Duration(5)).Format(constanta.MYSQL_DATETIME_FORMAT)
		o.QueryTable("user").Filter("is_locked", 1).Filter("updated_at__lt", ago).All(&users)

		for _, row := range users {
			row.IsLocked 	  = 0
			row.FailedAttempt = 0
			o.Update(&row,"IsLocked","FailedAttempt","UpdatedAt")
			beego.Debug(row.Email + " is unlocked")
		}

		return nil
	})

	beego.Debug("INIT TASK unlock_user")
	toolbox.AddTask("unlock_user", unlockUser)

	defer toolbox.StopTask()
}