package user

import (
	"errors"
	"strings"
	"strconv"

	"donorta-api/models"
	"donorta-api/lib/helper"
	"donorta-api/lib/security"
	"donorta-api/lib/checkmail"
	"donorta-api/lib/constanta"

	"github.com/astaxie/beego/validation"
	"github.com/astaxie/beego/orm"
)

func Register(user models.User) (models.User, int, error) {
	o := orm.NewOrm()
	valid := validation.Validation{}
	valid.Required(user.Password, "password")
	valid.Required(user.Fullname, "fullname")
	valid.Required(user.Birthdate, "birthdate")
	valid.Required(user.BloodType, "blood_type")
	valid.Required(user.Gender, "gender")
	valid.Required(user.Handphone, "handphone")
	valid.Required(user.Email, "email")
	valid.Email(user.Email, "email")
	valid.Required(user.SecurityQuestion, "security_question")
	valid.Required(user.SecurityAnswer, "security_answer")
	valid.Required(user.OneSignalId, "one_signal_id")
	valid.Required(user.Longitude, "longitude")
	valid.Required(user.Latitude, "latitude")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return user, 400, errors.New(err.Key +": "+ err.Message)
		}
	}

	// Not a number
	if _, err := strconv.Atoi(user.Handphone); err != nil {
		return user, 400, errors.New("Nomor telepon harus berupa angka")
	}

	// String length not 10
	if len(user.Handphone) < 10 {
		return user, 400, errors.New("Format nomor Handphone salah")
	}

	// First 2 digits should be 62
	twoFirstDigit := user.Handphone[:2]
	if twoFirstDigit != "62" {
		return user, 400, errors.New("Nomor Handphone harus nomor Indonesia")
	}

	// HP Already registered
	var exist models.User
	o.QueryTable("user").Filter("handphone", user.Handphone).One(&exist)
	if exist.Id != 0 {
		if exist.IsActive == 0 {
			return exist, 400, errors.New("Account anda belum aktif")
		} else {
			return exist, 400, errors.New("Nomor telepon/ email Anda telah terdaftar. Silakan login.")
		}
	}

	// Check Email Already registered and false account
	if user.Email != "" {
		emailExist := o.QueryTable("user").Filter("email", user.Email).Exist()
		if emailExist {
			return user, 400, errors.New("Email sudah terdaftar. Silakan login.")
		}

		// Check false email
		result, _ := checkmail.Check(user.Email)
		if result.IsValid() == false {
			return user, 400, errors.New("Format email salah. Tolong gunakan email yang valid.")
		}
	}


	if !helper.InArray(user.BloodType, constanta.BloodList) {
		return user, 400, errors.New("Golongan darah tidak valid.")
	}
	if user.Gender != constanta.USER_FEMALE && user.Gender != constanta.USER_MALE {
		return user, 400, errors.New("Jenis Kelamin tidak valid.")
	}

	user.SecuritySalt 	= security.ShaOneEncrypt(helper.GetNowTime().String() + helper.StringRandomWithCharset(30, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	user.SecurityAnswer = security.ShaOneEncrypt(strings.ToLower(user.SecurityAnswer) + user.SecuritySalt)
	user.Salt 			= security.ShaOneEncrypt(helper.GetNowTime().String() + helper.StringRandomWithCharset(32, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	user.Password 		= security.ShaOneEncrypt(user.Password + user.Salt)
	user.CreatedBy 		= helper.StringRandomWithCharset(32, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	user.LastLogin		= helper.GetNowTime()
	user.IsActive 		= 1
	user.LoginNumber++
	//TODO: Need to activate using SMS OTP
	//user.ChallengeCode	= helper.StringRandomWithCharset(8, helper.STRING_NUMBER)
	//user.ResponseCode	= helper.StringRandomWithCharset(6, helper.STRING_NUMBER)

	_, err := o.Insert(&user)
	if err != nil {
		return user, 500, errors.New(err.Error())
	}

	return models.CleanUserData(user), 0, nil
}