package repository

import (
	"fmt"
	"splitz/domain"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepoInterface interface {
	CreateGroup(name string, image string, admin string, users []string) (errorMsg string)
	DeleteGroup(id string) (data domain.GroupListResponse)
	GetSingleUsers(id string) (status bool, erroMsg string)
	GroupDetails(id string, userId uuid.UUID) (data domain.GroupDetailsResponse, errorMessage error)
	ListGroup(userId string) (data []domain.GroupListResponse, err error)
	UpdateGroup(id string) (data domain.GroupListResponse)
	GetUserList(contacts []string) (response gin.H, listError error)
}

type groupDbStruct struct {
	DB *gorm.DB
}

func InitGroupRepo(db *gorm.DB) GroupRepoInterface {
	return &groupDbStruct{DB: db}
}

func (d *groupDbStruct) CreateGroup(name string, image string, admin string, users []string) (errorMsg string) {

	tx := d.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errorMsg = "Internal Server error"
		}
	}()

	adminId, adminErr := uuid.Parse(admin)

	if adminErr != nil {
		tx.Rollback()
		return adminErr.Error()
	}

	group := domain.Group{
		Name:    name,
		Image:   image,
		AdminID: adminId,
	}

	dbErr := tx.Create(&group).Error

	if dbErr != nil {
		tx.Rollback()
		return dbErr.Error()
	}

	tx.Commit()

	for _, user := range users {

		userId, userErr := uuid.Parse(user)

		if userErr != nil {
			tx.Rollback()
			tx.Delete(&group)
			return userErr.Error()
		}

		userGroup := domain.UserGroup{
			UserID:  userId,
			GroupID: group.ID,
		}

		dbErr := d.DB.Create(&userGroup).Error

		if dbErr != nil {
			fmt.Printf("dbErr.Error(): %v\n", dbErr.Error())
			tx.Rollback()
			tx.Delete(&group)
			return dbErr.Error()
		}
	}

	tx.Commit()

	return ""
}

func (d *groupDbStruct) DeleteGroup(id string) (data domain.GroupListResponse) {
	panic("unimplemented")
}

func (d *groupDbStruct) GetSingleUsers(id string) (status bool, errorMsg string) {
	defer func() {
		if r := recover(); r != nil {
			status = false
			errorMsg = "Internal Server error"
		}
	}()

	var user domain.User
	err := d.DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		return false, "No User Found"
	} else {
		return true, ""
	}
}

func (d *groupDbStruct) GroupDetails(id string, userId uuid.UUID) (data domain.GroupDetailsResponse, errorMessage error) {

	defer func() {
		if r := recover(); r != nil {
			errorMessage = fmt.Errorf("group does not exist")
		}
	}()

	var users []domain.UserListResponse
	var userGroups []domain.UserGroup

	err := d.DB.Preload("Group").Preload("User").Where("group_id = ?", id).Find(&userGroups).Error

	if err != nil {
		return domain.GroupDetailsResponse{}, err
	}

	for _, userGroup := range userGroups {
		users = append(users, userGroup.User.ToGroupDetailsUser(userGroup.Group.AdminID == userGroup.UserID))
	}

	return domain.GroupDetailsResponse{
		Name:  userGroups[len(userGroups)-1].Group.Name,
		Image: userGroups[len(userGroups)-1].Group.Image,
		Users: users,
	}, nil
}

func (d *groupDbStruct) ListGroup(userId string) (data []domain.GroupListResponse, groupError error) {

	var groups []domain.GroupListResponse

	var userGroups []domain.UserGroup

	userUUID, uuidError := uuid.Parse(userId)

	if uuidError != nil {
		return groups, uuidError
	}

	err := d.DB.Preload("Group").Where("user_id = ?", userUUID).Find(&userGroups).Error

	for _, userGroup := range userGroups {
		fmt.Println(userGroup.UserID)
		fmt.Println(userGroup.ID)
		groups = append(groups, userGroup.Group.ToGroupListResponse())
	}

	if err != nil {
		return groups, err
	}

	return groups, nil
}

func (d *groupDbStruct) UpdateGroup(id string) (data domain.GroupListResponse) {

	panic("unimplemented")
}

func (d *groupDbStruct) GetUserList(contacts []string) (response gin.H, listError error) {

	const chunkSize = 100
	var wg sync.WaitGroup
	var mu sync.Mutex

	var users []domain.User
	var responseUsers []domain.UserListResponse
	var contactsNotInDB []string

	errChan := make(chan error, len(contacts)/chunkSize+1)

	for i := 0; i < len(contacts); i += chunkSize {
		end := i + chunkSize
		if end > len(contacts) {
			end = len(contacts)
		}
		chunk := contacts[i:end]

		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			var chunkUsers []domain.User

			err := d.DB.Where("mobile IN ?", chunk).Find(&chunkUsers).Error
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			users = append(users, chunkUsers...)
			mu.Unlock()
		}(chunk)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	foundMobiles := make(map[string]bool)
	for _, user := range users {
		responseUsers = append(responseUsers, user.ToUserListResponse())
		foundMobiles[user.Mobile] = true
	}

	for _, contact := range contacts {
		if !foundMobiles[contact] {
			contactsNotInDB = append(contactsNotInDB, contact)
		}
	}

	return gin.H{
		"existing_users":     responseUsers,
		"non_existing_users": contactsNotInDB,
	}, nil

}
