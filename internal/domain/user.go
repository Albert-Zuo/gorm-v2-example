package domain

import (
	"fmt"
	"strconv"
)

type User struct {
	UserId   int64  `gorm:"primary_key" json:"user_id"`
	UserName string `json:"user_name`
	Password string `json:"password"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
}

func (u User) string() string {
	return fmt.Sprintf("user:%s[name=%s, pass=%s, email=%s, mobile=%s]",
		strconv.FormatInt(u.UserId, 10), u.UserName, u.Password, u.Email, u.Mobile)
}

func (User) TableName() string {
	return "blog_user"
}