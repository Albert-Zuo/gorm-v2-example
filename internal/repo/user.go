package repo

import (

	"fmt"
	"github.com/gorm-v2-example/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var key = "xiaozuo1221@gmail.com"

// userRepo repo
type userRepo struct {
	db *gorm.DB
}

// NewUserRepo NewUserRepo
func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db}
}

// ListUser ListUser
func (ur userRepo) ListUser() ([]*domain.User, error) {
	users := make([]*domain.User, 0)

	rows, err := ur.db.Raw(`SELECT 
			user_id,
			CAST(AES_DECRYPT(UNHEX(user_name), ?) AS CHAR) AS user_name,
			CAST(AES_DECRYPT(UNHEX(password), ?) AS CHAR) AS password,
			CAST(AES_DECRYPT(UNHEX(email), ?) AS CHAR) AS email,
			CAST(AES_DECRYPT(UNHEX(mobile), ?) AS CHAR) AS mobile
		FROM blog_user`, key, key, key, key).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user domain.User
		ur.db.ScanRows(rows, &user)
		users = append(users, &user)
	}
	return users, nil
}

// CreateUser CreateUser
func (ur userRepo) CreateUser(user *domain.User) (bool, error) {
	ur.db.Model(domain.User{}).Create(map[string]interface{}{
		"UserName": clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.UserName, key}},
		"Password": clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.Password, key}},
		"Email":    clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.Email, key}},
		"Mobile":   clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.Mobile, key}},
	})

	return true, nil
}

// FindUser FindUser
func (ur userRepo) FindUser(s string) ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	rows, err := ur.db.Raw(`SELECT 
			user_id,
			CAST(AES_DECRYPT(UNHEX(user_name), ?) AS CHAR) AS user_name,
			CAST(AES_DECRYPT(UNHEX(password), ?) AS CHAR) AS password,
			CAST(AES_DECRYPT(UNHEX(email), ?) AS CHAR) AS email,
			CAST(AES_DECRYPT(UNHEX(mobile), ?) AS CHAR) AS mobile
		FROM
			blog_user
		WHERE
			CAST(AES_DECRYPT(UNHEX(user_name), ?) AS CHAR) LIKE ?`, key, key, key, key, key, fmt.Sprint("%", s, "%")).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user domain.User
		ur.db.ScanRows(rows, &user)
		users = append(users, &user)
	}
	return users, nil
}

// UpdateUser UpdateUser
func (ur userRepo) UpdateUser(id int64, user *domain.User) error {

	ur.db.Exec("UPDATE blog_user SET user_name=HEX(AES_ENCRYPT(?, ?)),"+
		" password=HEX(AES_ENCRYPT(?, ?)), "+
		"email=HEX(AES_ENCRYPT(?, ?)), "+
		"mobile=HEX(AES_ENCRYPT(?, ?)) WHERE user_id = ?",
		user.UserName, key, user.Password, key,
		user.Email, key, user.Mobile, key, id,
	)
	
	return nil
}
