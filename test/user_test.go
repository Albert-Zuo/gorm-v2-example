package test

import (
	"fmt"
	"github.com/gorm-v2-example/internal/domain"
	"github.com/gorm-v2-example/internal/repo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var (
	db, _ = gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	userRepo = repo.NewUserRepo(db)
)


// TestCreateUser init data
func TestCreateUser(t *testing.T) {

	users := make([]*domain.User, 0)
	for i := 0; i < 100; i++ {
		users = append(users, &domain.User{
			UserName: fmt.Sprint("test", i),
			Password: "00000",
			Email:    "xiaozuo1221@gmail.com",
			Mobile:   "12345678900",
		})
	}

	users = append(users, &domain.User{
		UserName: "Albert",
		Password: "albert-zuo",
		Email:    "xiaozuo1221@gmail.com",
		Mobile:   "12345678900",
	})

	for _, user := range users {
		if ok, err := userRepo.CreateUser(user); !ok || err != nil {
			t.Error("TestCreateUser error: ", err)
		}
	}
}


func TestUserList(t *testing.T)  {

	users, err := userRepo.ListUser()
	if err != nil {
		t.Error(err.Error())
	}

	for _, user := range users {
		t.Log(user)
	}

}

func TestFindUser(t *testing.T) {

	ts := []string{"ber", "es"}

	for _, s := range ts {
		users, err := userRepo.FindUser(s)
		if err != nil {
			t.Error(err.Error())
		}
		for _, user := range users {
			t.Log(user)
		}
		t.Log("模糊查找", s, "完成=====================================")
	}
}
