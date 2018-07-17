package user

import (
	"errors"

	"donorta-api/models"
	"donorta-api/lib/helper"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

func Login (user models.User) (models.User, int, error) {
	o := orm.NewOrm()
	valid := validation.Validation{}
	valid.Required(user.Handphone, "handphone")
	valid.Required(user.Password, "password")
	valid.Required(user.OneSignalId, "one_signal_id")
	valid.Required(user.Longitude, "longitude")
	valid.Required(user.Latitude, "latitude")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return user, 400, errors.New(err.Key +": "+ err.Message)
		}
	}

	// Find record
	var tempUser models.User
	err := o.QueryTable("user").Filter("handphone", user.Handphone).One(&tempUser)
	if err == orm.ErrNoRows {
		err = o.QueryTable("user").Filter("email", user.Handphone).One(&tempUser)
		if err == orm.ErrNoRows {
			return user, 400, errors.New("No HP atau email salah!")
		}
	}

	// Not Active
	if tempUser.IsActive != 1 {
		return tempUser, 400, errors.New("Anda harus aktivasi akun anda terlebih dahulu")
	}

	// Account Locked
	if tempUser.IsLocked == 1 {
		return user, 400, errors.New("Akun Anda terkunci. Coba lagi dalam 5 menit.")
	}

	// Account Banned
	if tempUser.IsLocked == 2 {
		return user, 400, errors.New("Akun Anda di blokir.")
	}

	// Wrong Password
	if models.UserIsCorrectPassword(tempUser, user.Password) != true {
		return user, 400, errors.New("Password yang Anda masukkan salah")
	}

	tempUser.LoginNumber   += 1
	tempUser.FailedAttempt 	= 0
	tempUser.LastLogin 		= helper.GetNowTime()
	tempUser.OneSignalId	= user.OneSignalId
	tempUser.Longitude		= user.Longitude
	tempUser.Latitude		= user.Latitude
	_, err = o.Update(&tempUser, "LoginNumber","FailedAttempt","LastLogin","OneSignalId","Longitude","Latitude","UpdatedAt")
	if err != nil {
		return user, 500, errors.New(err.Error())
	}

	return tempUser, 0, nil
}