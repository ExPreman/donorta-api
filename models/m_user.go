package models

import (
	"time"
	"strconv"
	"encoding/json"

	"donorta-api/lib/security"
	"donorta-api/lib/redis"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       				uint64						`json:"id"`
	Fullname 				string						`orm:"size(100)" valid:"Required" json:"fullname"`
	Handphone				string						`orm:"size(100)" valid:"Required" json:"handphone"`
	Email					string						`valid:"Email" orm:"size(100)" json:"email"`
	Password 				string						`orm:"size(100)" valid:"Required" json:"password,omitempty"`
	Salt  	 				string						`orm:"size(100)" valid:"Required" json:"salt,omitempty"`
	//TODO: Reward Point
	//Point					uint						`json:"point"`
	Birthdate 				time.Time					`orm:"type(date)" json:"birthdate"`
	BloodType				string						`orm:"size(2)" valid:"Required" json:"blood_type"`
	Gender					string						`orm:"size(10)" valid:"Required" json:"gender"`
	//TODO: Address
	//Address  				string						`orm:"size(255)" json:"address"`
	//City	 				string						`orm:"size(100)" json:"city"`
	//Province				string						`orm:"size(100)" json:"province"`
	//Zipcode				string						`orm:"size(10)" json:"zipcode"`
	//TODO: Profile Picture
	//Avatar				string						`orm:"size(65535)" json:"avatar"`
	SecurityQuestion		string						`orm:"size(100)" valid:"Required" json:"security_question"`
	SecurityAnswer			string						`orm:"size(100)" valid:"Required" json:"security_answer,omitempty"`
	SecuritySalt			string						`orm:"size(100)" valid:"Required" json:"security_salt,omitempty"`
	//TODO: OTP
	//ChallengeCode			string						`orm:"size(50)" valid:"Required" json:"challenge_code"`
	//ResponseCode			string						`orm:"size(50)" valid:"Required" json:"response_code,omitempty"`
	LoginNumber				uint32						`json:"login_number"`
	FailedAttempt			uint8 						`json:"failed_attempt"`
	TotalFailedAttempt		int 						`json:"total_failed_attempt"`
	LastLogin				time.Time					`orm:"type(datetime)" json:"last_login"`
	Longitude				string						`json:"longitude"`
	Latitude				string						`json:"latitude"`
	OneSignalId				string						`json:"one_signal_id"`
	IsLocked				uint8						`orm:"default(0)" json:"is_locked"`
	IsActive				uint8						`orm:"default(1)" json:"is_active"`
	CreatedAt				time.Time 					`orm:"auto_now_add;type(datetime)" json:"created_at"`
	CreatedBy				string	  					`json:"created_by"`
	UpdatedAt				time.Time 					`orm:"auto_now;type(datetime)" json:"updated_at"`
	UpdatedBy				string 						`orm:"null" json:"updated_by"`
	UserToken				[]*UserToken				`orm:"reverse(many)" json:"user_token,omitempty"`
	UserActivity			[]*UserActivity				`orm:"reverse(many)" json:"user_activity,omitempty"`
}

func init() {
	orm.RegisterModel(new(User))
}

func UserIsCorrectPassword (user User, pass string) (bool) {
	encodedPin := security.ShaOneEncrypt(pass + user.Salt)
	if encodedPin != user.Password {
		o := orm.NewOrm()
		beego.Debug("User Wrong Password")
		ActivityLog(user, "LOGIN", "Login gagal", user.Email, "", 0)

		user.FailedAttempt += 1
		user.TotalFailedAttempt += 1
		if user.FailedAttempt > 2 {
			user.IsLocked = 1

			// Clear all active token
			o.QueryTable("user_token").
				Filter("User", user.Id).
				Filter("is_active", 1).
				Update(orm.Params{
				"is_active": 0,
			})
		}
		o.Update(&user)
		return false
	}
	return true
}

func CleanUserData(u User) (User) {
	u.Password 			= ""
	u.Salt 				= ""
	u.SecurityAnswer	= ""
	u.SecuritySalt	 	= ""

	return u
}

func GetActiveUser (id uint64) User {
	var data User
	o := orm.NewOrm()
	key 	  := redis.UserProfile +"_"+ strconv.Itoa(int(id))
	expiry 	  := time.Hour * 24 * 7
	cData,err := redis.GetBytes("GET", key)

	if err != nil {
		o.QueryTable("user").
			Filter("id", id).
			Filter("is_active", 1).
			Filter("is_locked", 0).
			One(&data)

		plJson, _ := json.Marshal(data)
		redis.SetEx(key, string(plJson), expiry)
	} else {
		json.Unmarshal(cData, &data)
	}

	return data
}