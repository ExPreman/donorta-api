package models

import (
	"time"
	"encoding/json"

	"donorta-api/lib/redis"
	"donorta-api/lib/helper"
	"donorta-api/lib/security"

	"github.com/astaxie/beego/orm"
)

type UserToken struct {
	Id       		uint64			`json:"id"`
	Token			string			`orm:"size(100)" valid:"Required" json:"token"`
	User      		*User			`orm:"rel(fk);column(user_id)" json:"user"`
	Expires			time.Time		`orm:"type(datetime)" json:"expires"`
	IsActive		uint8			`orm:"default(1)" json:"is_active"`
	CreatedAt		time.Time 		`orm:"auto_now_add;type(datetime)" json:"created_at"`
	CreatedBy		string	  		`json:"created_by"`
}

func init() {
	orm.RegisterModel(new(UserToken))
}

func TokenGenerate (user User) (token UserToken, code int, err error) {
	o := orm.NewOrm()
	token.CreatedBy = "SYSTEM"
	token.Token 	= security.ShaOneEncrypt(helper.GetNowTime().String() + helper.StringRandom(20))
	token.IsActive 	= 1
	token.User		= &user
	token.Expires	= helper.GetNowTime().Add(time.Minute * time.Duration(180000))

	_, err = o.Insert(&token)
	if err != nil {
		return token, 500, err
	}

	return token, 200, nil
}

func GetUserToken(token string) UserToken {
	var data UserToken
	o := orm.NewOrm()
	key 	  := redis.UserToken +"_"+ token
	expiry 	  := time.Hour * 24 * 7
	cData,err := redis.GetBytes("GET", key)

	if err != nil {
		o.QueryTable("user_token").
			Filter("token", token).
			Filter("is_active", 1).
			One(&data)

		// Check token expires
		if helper.GetNowTime().After(data.Expires) {
			o.Delete(&data)
			return data
		}
		plJson, _ := json.Marshal(data)
		redis.SetEx(key, string(plJson), expiry)
	} else {
		json.Unmarshal(cData, &data)
	}

	return data
}