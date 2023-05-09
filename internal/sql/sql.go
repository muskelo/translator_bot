package sql

import (
	"fmt"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var (
	ErrUserNotFound error = fmt.Errorf("User not found")
)

func New(DatabaseConnStr string) (*xorm.Engine, error) {
	return xorm.NewEngine("postgres", DatabaseConnStr)
}

type User struct {
	TelegramId        int64 `xorm:"pk unique"`
	PrimaryLanguage   string
	SecondaryLanguage string
}

func GetUserByID(engine *xorm.Engine, id int64) (User, error) {
	user := User{TelegramId: id}
	ok, err := engine.Get(&user)
	if err == nil && !ok {
		err = ErrUserNotFound
	}
	return user, err
}

func UpdateUser(engine *xorm.Engine, user User) error {
	affected, err := engine.Update(user)
	if err == nil && affected == 0 {
		err = ErrUserNotFound
	}
	return err
}

func CreateUser(engine *xorm.Engine, user User) error {
	_, err := engine.Insert(user)
	return err
}
