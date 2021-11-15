package data

import (
	"fmt"
	"time"
	"usermanagement/app/internal"
	"usermanagement/app/internal/serviceerror"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Group struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name      string     `sql:"index"`
}

type UserGroup struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	UserID    uint       `sql:"index"`
	GroupID   uint       `sql:"index"`
}

type groupDataService struct {
	db *gorm.DB
}

func NewGroupService(db *gorm.DB) *groupDataService {
	db.AutoMigrate(&Group{})
	db.AutoMigrate(&UserGroup{})
	return &groupDataService{
		db: db,
	}
}

func (g *groupDataService) CreateGroup(request internal.GroupRequest) (response internal.GroupResponse, err error) {
	if request.Name == "" {
		return response, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("missing create group fields"))
	}
	var count int64
	err = g.db.Model(&Group{}).Where("name = ?", request.Name).Count(&count).Error
	if err != nil {
		return response, errors.Wrap(err, "get group with name count failed")
	}
	if count > 0 {
		return response, serviceerror.NewServiceError(serviceerror.DuplicateUser, fmt.Errorf("group with name %s is present", request.Name))
	}
	group := Group{
		Name: request.Name,
	}
	err = g.db.Create(&group).Error
	if err != nil {
		return response, errors.Wrap(err, "create group failed")
	}
	response = internal.GroupResponse{
		ID:   group.ID,
		Name: group.Name,
	}
	return response, err
}

func (g *groupDataService) UpdateGroup(request internal.UpdateGroupRequest) (err error) {
	if request.ID == 0 || request.Name == "" {
		return serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("missing update group fields"))
	}
	var count int64
	err = g.db.Model(&Group{}).Where("name = ?", request.Name).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "get group with name count failed")
	}
	if count > 0 {
		return serviceerror.NewServiceError(serviceerror.DuplicateUser, fmt.Errorf("group with name %s is present", request.Name))
	}
	group := Group{
		ID: request.ID,
	}
	err = g.db.Model(&group).Updates(User{Name: request.Name}).Error
	if err != nil {
		return errors.Wrap(err, "update group failed")
	}
	return err
}

func (g *groupDataService) DeleteGroup(id uint) (err error) {
	if id == 0 {
		return serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("group_id is 0 for delete group"))
	}
	err = g.db.Where("group_id = ?", id).Delete(&UserGroup{}).Error
	if err != nil {
		return errors.Wrap(err, "delete user group failed")
	}
	err = g.db.Delete(&Group{}, id).Error
	if err != nil {
		return errors.Wrap(err, "delete group failed")
	}
	return err
}

func (g *groupDataService) AddUser(userID uint, groupID uint) (err error) {
	if userID == 0 || groupID == 0 {
		return serviceerror.NewServiceError(serviceerror.InvalidUserGroupRequest, errors.New("group_id or user_id is 0 for adding user"))
	}
	var count int64
	err = g.db.Model(&UserGroup{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, "count usergroup failed")
	}
	if count != 0 {
		return serviceerror.NewServiceError(serviceerror.InvalidUserGroupRequest, fmt.Errorf("user_id %d is already associated with a group", userID))
	}
	userGrp := UserGroup{
		UserID:  userID,
		GroupID: groupID,
	}
	err = g.db.Create(&userGrp).Error
	if err != nil {
		return errors.Wrap(err, "create usergroup failed")
	}
	return err
}

func (g *groupDataService) RemoveUser(groupID uint, userID uint) (err error) {
	if userID == 0 || groupID == 0 {
		return serviceerror.NewServiceError(serviceerror.InvalidUserGroupRequest, errors.New("group_id or user_id is 0 for removing user"))
	}
	err = g.db.Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&UserGroup{}).Error
	if err != nil {
		return errors.Wrap(err, "remove usergroup failed")
	}
	return err
}

func (g *groupDataService) GetUsersByGroupID(groupID uint, offset uint, limit uint) (response internal.UsersResponse, err error) {
	if groupID == 0 || limit == 0 || limit > 1000 {
		return response, serviceerror.NewServiceError(serviceerror.InvalidUserGroupRequest, fmt.Errorf("group_id or limit %d is not valid for getting users", limit))
	}
	var users []User
	err = g.db.Joins("JOIN user_groups ON user_groups.group_id = ? AND user_groups.deleted_at IS NULL", groupID).
		Where("users.id = user_groups.user_id").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return response, errors.Wrap(err, "get users failed")
	}
	var count int64
	err = g.db.Model(&User{}).Joins("JOIN user_groups ON user_groups.group_id = ? AND user_groups.deleted_at IS NULL", groupID).
		Where("users.id = user_groups.user_id").Count(&count).Error
	if err != nil {
		return response, errors.Wrap(err, "get users count failed")
	}
	usersResponse := make([]internal.UserResponse, len(users))
	for i, u := range users {
		usersResponse[i] = internal.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}
	response = internal.UsersResponse{
		Users: usersResponse,
		Total: uint(count),
	}
	return response, err
}

func (g *groupDataService) GetGroups(offset uint, limit uint) (response internal.GroupsResponse, err error) {
	if limit == 0 || limit > 1000 {
		return response, serviceerror.NewServiceError(serviceerror.InvalidUserGroupRequest, fmt.Errorf("limit %d is not valid for getting groups", limit))
	}
	var groups []Group
	err = g.db.Limit(limit).Offset(offset).Find(&groups).Error
	if err != nil {
		return response, errors.Wrap(err, "get groups failed")
	}
	var count int64
	err = g.db.Model(&Group{}).Count(&count).Error
	if err != nil {
		return response, errors.Wrap(err, "get groups count failed")
	}
	groupsResponse := make([]internal.GroupResponse, len(groups))
	for i, g := range groups {
		groupsResponse[i] = internal.GroupResponse{
			ID:   g.ID,
			Name: g.Name,
		}
	}
	response = internal.GroupsResponse{
		Groups: groupsResponse,
		Total:  uint(count),
	}
	return response, err
}
