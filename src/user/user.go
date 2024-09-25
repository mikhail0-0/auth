package user

import (
	"auth/src/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email    string    `gorm:"type:varchar"`
	Password string    `gorm:"type:varchar"`
}

var db *gorm.DB

func Init(refDB *gorm.DB) error {
	db = refDB
	if !db.Migrator().HasTable(&User{}) {
		db.Migrator().CreateTable(&User{})

		hashString, err := common.GetHash("test")
		if err != nil {
			return err
		}

		db.Create(&User{
			Email:    "test@mail.com",
			Password: hashString,
		})
	}
	return nil
}

func FindById(id string) (*User, error) {
	var user User
	result := db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, common.ErrNotFound
	}

	return &user, nil
}

func Verify(id, password string) (*User, error) {
	user, err := FindById(id)
	if err != nil {
		return nil, err
	}

	err = common.CompareHashAndString(user.Password, password)
	if err != nil {
		return nil, common.ErrWrongPassword
	}

	return user, nil
}
