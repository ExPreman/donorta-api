package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(UserActivity))
}

type UserActivity struct {
	Id       		uint64		`json:"id"`
	User      		*User		`orm:"rel(fk);column(user_id)" json:"user"`
	Container      	string		`orm:"size(50)" json:"container"`
	ContainerId     uint64		`json:"container_id"`
	Name			string		`orm:"size(100)" valid:"Required" json:"name"`
	Description		string		`orm:"size(65535)" json:"description"`
	IsActive		uint8		`orm:"default(1)" json:"is_active"`
	CreatedAt		time.Time 	`orm:"auto_now_add;type(datetime)" json:"created_at"`
	CreatedBy		string	  	`json:"created_by"`
}

func ActivityLog (user User, name string, desc string, creator string, container string, cid uint64) (custAct UserActivity, code int, err error) {
	o := orm.NewOrm()
	custAct.User 		= &user
	custAct.Container 	= container
	custAct.ContainerId = cid
	custAct.Name		= name
	custAct.Description	= desc
	custAct.IsActive 	= 1
	custAct.CreatedBy	= creator

	_, err = o.Insert(&custAct)
	if err != nil {
		return custAct, 400, err
	}

	return custAct, 200, nil
}