package repository

import (
	"splitz/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepoInterface interface {
	Login(mobile string, password string) (user *domain.User, errorMsg string)
	Register(user domain.User) (id uuid.UUID, errorMsg string)
}

type authDbStruct struct {
	DB *gorm.DB
}

func InitAuthRepo(db *gorm.DB) AuthRepoInterface {
	return &authDbStruct{DB: db}
}

func (a *authDbStruct) Login(mobile string, password string) (u *domain.User, errorMsg string) {

	defer func() {
		if r := recover(); r != nil {
			u = nil
			errorMsg = "Internal Server error"
		}
	}()

	var user domain.User
	err := a.DB.Where("mobile = ?", mobile).First(&user).Error

	if err != nil {
		return nil, "No user with the given number exists"
	} else {
		if user.Password != password {
			return nil, "Password did not match"
		} else {
			return &user, ""
		}

	}

}

func (a *authDbStruct) Register(user domain.User) (id uuid.UUID, errorMsg string) {
	defer func() {
		if r := recover(); r != nil {
			errorMsg = "Internal Server error"
		}
	}()
	err := a.DB.Where("mobile = ?", user.Mobile).First(&user).Error

	if err != nil && err.Error() == "record not found" {
		user = domain.User{
			Mobile:    user.Mobile,
			Password:  user.Password,
			Name:      user.Name,
			Image:     user.Image,
			FcmTokens: user.FcmTokens,
		}
		err := a.DB.Create(&user).Error
		if err != nil {
			return uuid.Max, err.Error()
		} else {
			return user.ID, ""
		}
	} else {
		return uuid.Max, "User already exists with the given number"
	}
}
