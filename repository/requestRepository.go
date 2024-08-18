package repository

import (
	"fmt"
	"splitz/domain"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestInterface interface {
	CreateRequest(requestType string, note string, amount string, group uuid.UUID, owner uuid.UUID, users []map[string]string) (errorMessage string)
	PayShare(request uuid.UUID, group uuid.UUID, user uuid.UUID, amount float64) (errorMessage error)
	ListRequest(group uuid.UUID) (errorMessage error, userReqs []domain.RequestList)
}

type requestDbStruct struct {
	DB *gorm.DB
}

func InitRequestRepo(db *gorm.DB) RequestInterface {
	return &requestDbStruct{DB: db}
}

func (d *requestDbStruct) CreateRequest(requestType string, note string, amount string, group uuid.UUID, owner uuid.UUID, users []map[string]string) (errorMessage string) {
	tx := d.DB.Begin()

	defer func() {
		if r := recover(); r != nil {

			tx.Rollback()
			errorMessage = "Internal Server Error"

		}
	}()

	totalAmount, strConvErr := strconv.ParseFloat(amount, 64)

	if strConvErr != nil {
		tx.Rollback()
		errorMessage = "Unable to convert amount"
		return
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
		tx.Rollback()
		errorMessage = "Unable to Create Request"
		return
	}

	sharesum := 0.0

	for _, user := range users {
		userId, userErr := uuid.Parse(user["id"])

		if userErr != nil {
			tx.Rollback()
			return "Unable to parse user"
		}

		shareAmount, strConvErr := strconv.ParseFloat(user["amount"], 64)

		if strConvErr != nil {
			tx.Rollback()
			return "Unable to convert share"
		}

		sharesum += shareAmount

		if sharesum > request.Amount {
			tx.Rollback()
			return "User shares are greater than the Request amount"
		}

		request := domain.UserRequest{
			RequestID: request.ID,
			Share:     shareAmount,
			UserID:    userId,
		}

		dbErr := d.DB.Create(&request).Error

		if dbErr != nil {
			tx.Rollback()
			errorMessage = "Unable to Create Request User"
			return
		}
	}

	tx.Commit()

	return ""
}

func (d *requestDbStruct) PayShare(request uuid.UUID, group uuid.UUID, user uuid.UUID, amount float64) (errorMessage error) {

	var userReq domain.UserRequest

	getErr := d.DB.Where("id = ?", request).Find(&userReq).Error

	if getErr != nil {
		return getErr
	}

	if userReq.Share == userReq.Paid {
		return fmt.Errorf("Share is already paid")

	}

	userReq.Paid += amount

	if userReq.Share < userReq.Paid {
		return fmt.Errorf("Amount is greater than share")
	}

	err := d.DB.Model(userReq).Updates(userReq).Error

	if err != nil {
		return err
	}

	return nil
}

func (d *requestDbStruct) ListRequest(group uuid.UUID) (errorMessage error, userRequests []domain.RequestList) {

	var requests []domain.Request

	var groupRequests []domain.RequestList

	getErr := d.DB.Where("group_id = ?", group).Find(&requests).Error

	if getErr != nil {
		return getErr, groupRequests
	}

	for _, request := range requests {

		var req domain.RequestList

		req.Amount = request.Amount
		req.ID = request.ID
		req.Note = request.Note
		req.Type = request.Type

		var userReqs []domain.UserRequest

		getErr := d.DB.Preload("User").Where("request_id = ?", request.ID).Find(&userReqs).Error

		if getErr != nil {
			return getErr, groupRequests
		}

		var userReqsList []domain.UserRequestList

		paidAmount := 0.0

		for _, userReq := range userReqs {
			paidAmount += userReq.Paid
			userReqsList = append(userReqsList, domain.UserRequestList{
				ID:    userReq.ID,
				User:  userReq.User.ToUserListResponse(),
				Share: userReq.Share,
				Paid:  userReq.Paid,
			})
		}

		req.Splits = userReqsList
		req.Paid = paidAmount

		groupRequests = append(groupRequests, req)

	}

	return nil, groupRequests
}
