package repository

import (
	"fmt"
	"splitz/domain"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestInterface interface {
	CreateRequest(requestType string, note string, amount string, group uuid.UUID, owner uuid.UUID) (id uuid.UUID, err error)
	CreateUserRequest(requestId uuid.UUID, share string, user string) (err error)
}

type requestDbStruct struct {
	DB *gorm.DB
}

func InitRequestRepo(db *gorm.DB) RequestInterface {
	return &requestDbStruct{DB: db}
}

func (d *requestDbStruct) CreateRequest(requestType string, note string, amount string, group uuid.UUID, owner uuid.UUID) (id uuid.UUID, errorMsg error) {

	totalAmount, strConvErr := strconv.ParseFloat(amount, 64)

	if strConvErr != nil {
		return uuid.Max, strConvErr
	}

	request := domain.Request{
		Note:    note,
		Amount:  totalAmount,
		OwnerID: owner,
		GroupID: group,
		Type:    requestType,
	}

	dbErr := d.DB.Create(&request).Error

	if dbErr != nil {
		return uuid.Max, dbErr
	}

	return request.ID, nil
}

func (d *requestDbStruct) CreateUserRequest(requestId uuid.UUID, share string, user string) (err error) {
	userId, userErr := uuid.Parse(user)

	if userErr != nil {
		return userErr
	}

	shareAmount, strConvErr := strconv.ParseFloat(share, 64)

	if strConvErr != nil {
		return strConvErr
	}

	request := domain.UserRequest{
		RequestID: requestId,
		Share:     shareAmount,
		UserID:    userId,
	}

	fmt.Printf("request: %v\n", request.UserID)

	dbErr := d.DB.Create(&request).Error

	if dbErr != nil {
		return dbErr
	}

	return nil
}
