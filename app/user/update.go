package user

import (
	"errors"
	"strconv"

	"donorta-api/models"
	"donorta-api/lib/security"
	"donorta-api/lib/helper"
	"donorta-api/lib/checkmail"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

func Update(current models.User, user models.User) (models.User, int, error) {
	o := orm.NewOrm()

	if user.Password != "" {
		user.Salt 		 = security.ShaOneEncrypt(helper.GetNowTime().String())
		user.Password 	 = security.ShaOneEncrypt(user.Password + user.Salt)
		current.Password = user.Password
		current.Salt 	 = user.Salt
	}

	if user.Handphone != "" && user.Handphone != current.Handphone {
		if _, err := strconv.Atoi(user.Handphone); err != nil {
			return user, 400, errors.New("Handphone harus berupa angka")
		}
		// First 2 digits should be 62
		twoFirstDigit := user.Handphone[:2]
		if twoFirstDigit != "62" {
			return user, 400, errors.New("Nomor Handphone harus nomor Indonesia")
		}
		// HP Already registered
		hpExist := o.QueryTable("user").Filter("handphone", user.Handphone).Exclude("id", current.Id).Exist()
		if hpExist {
			return user, 400, errors.New("Handphone sudah terdaftar")
		}

		current.Handphone = user.Handphone
	}

	// Email Already registered
	if user.Email != "" && user.Email != current.Email {
		valid := validation.Validation{}
		valid.Required(user.Email, "email")
		valid.Email(user.Email, "email")

		if valid.HasErrors() {
			for _, err := range valid.Errors {
				return user, 400, errors.New(err.Key +": "+ err.Message)
			}
		}

		emailExist := o.QueryTable("user").Filter("email", user.Email).Exclude("email", current.Email).Exist()
		if emailExist {
			return user, 400, errors.New("Email sudah terdaftar")
		}

		// Check false email
		result, _ := checkmail.Check(user.Email)
		if result.IsValid() == false {
			return user, 400, errors.New("Format email salah. Tolong gunakan email yang valid.")
		}

		current.Email = user.Email
	}

	current.UpdatedBy 	= user.Email

	if user.Fullname != "" && user.Fullname != current.Fullname {
		current.Fullname = user.Fullname
	}
	if user.Address != "" && user.Address != current.Address {
		current.Address	= user.Address
	}
	if user.City != "" && user.City != current.City {
		current.City	= user.City
	}
	if user.Province != "" && user.Province != current.Province {
		// Prevent city in another province
		if user.City == "" {
			current.City = ""
		}
		current.Province = user.Province
	}
	if user.Zipcode != "" && user.Zipcode != current.Zipcode {
		current.Zipcode = user.Zipcode
	}
	if beego.Substr(user.Birthdate.String(), 0, 4) != "0001" && user.Birthdate != current.Birthdate {
		current.Birthdate = user.Birthdate
	}
	//TODO: Upload Image file
	//if user.Avatar != "" {
	//	rs := cloudinary.ImageUpload(user.Avatar, "")
	//	current.Avatar = rs
	//}
	if user.OneSignalId != "" {
		current.OneSignalId = user.OneSignalId
	}
	if user.Longitude != "" {
		current.Longitude = user.Longitude
	}
	if user.Latitude != "" {
		current.Latitude = user.Latitude
	}

	_, err := o.Update(&current)
	if err != nil {
		return user, 500, errors.New(err.Error())
	}

	return current, 200, nil
}