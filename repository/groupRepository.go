package repository

import (
	"splitz/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepoInterface interface {
	CreateGroup(name string, image string, admin string) (id uuid.UUID, err error)
	CreateUserGroup(name string, image string, user string, groupId uuid.UUID) (err error)
	DeleteGroup(id string) (data domain.GroupListResponse)
	GetSingleUsers(id string) (status bool, erroMsg string)
	GroupDetails(id string) (data domain.GroupListResponse)
	ListGroup(userId string) (data []domain.GroupListResponse, err error)
	UpdateGroup(id string) (data domain.GroupListResponse)
}

type groupDbStruct struct {
	DB *gorm.DB
}

func InitGroupRepo(db *gorm.DB) GroupRepoInterface {
	return &groupDbStruct{DB: db}
}

func (d *groupDbStruct) CreateGroup(name string, image string, admin string) (id uuid.UUID, errorMsg error) {
	adminId, adminErr := uuid.Parse(admin)

	if adminErr != nil {
		return uuid.Max, adminErr
	}

	group := domain.Group{
		Name:    name,
		Image:   image,
		AdminID: adminId,
	}

	dbErr := d.DB.Create(&group).Error
	if dbErr != nil {
		return uuid.Max, dbErr
	} else {
		return group.ID, nil
	}

}

func (d *groupDbStruct) CreateUserGroup(name string, image string, user string, groupId uuid.UUID) (errorMsg error) {
	userId, userErr := uuid.Parse(user)

	if userErr != nil {
		return userErr
	}

	group := domain.UserGroup{
		UserID:  userId,
		GroupID: groupId,
	}

	dbErr := d.DB.Create(&group).Error
	if dbErr != nil {
		return dbErr
	} else {
		return nil
	}

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

func (d *groupDbStruct) GroupDetails(id string) (data domain.GroupListResponse) {
	println(id)
	panic("unimplemented")
}

func (d *groupDbStruct) ListGroup(userId string) (data []domain.GroupListResponse, groupError error) {

	var groups []domain.GroupListResponse

	var userGroups []domain.UserGroup

	err := d.DB.Preload("Group").Find(&userGroups).Error

	for _, userGroup := range userGroups {
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
