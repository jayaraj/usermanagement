package service

import (
	"fmt"
	"usermanagement/app/internal"
	"usermanagement/app/internal/serviceerror"
)

type groupService struct {
	data internal.GroupData
}

func NewGroupService(data internal.GroupData) *groupService {
	return &groupService{
		data: data,
	}
}

func (g *groupService) CreateGroup(request internal.GroupRequest) (response internal.GroupResponse, err error) {
	return g.data.CreateGroup(request)
}

func (g *groupService) UpdateGroup(request internal.UpdateGroupRequest) (err error) {
	return g.data.UpdateGroup(request)
}

func (g *groupService) DeleteGroup(id uint) (err error) {
	return g.data.DeleteGroup(id)
}

func (g *groupService) GetUsersByGroupID(groupID uint, page uint, perPage uint) (response internal.UsersResponse, err error) {
	if page <= 0 || perPage == 0 {
		return response, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, fmt.Errorf("page %d  or per page %d is not valid", page, perPage))
	}
	offset := perPage * (page - 1)
	response, err = g.data.GetUsersByGroupID(groupID, offset, perPage)
	response.Page = page
	response.PerPage = perPage
	return response, err
}

func (g *groupService) GetGroups(page uint, perPage uint) (response internal.GroupsResponse, err error) {
	if page <= 0 || perPage == 0 {
		return response, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, fmt.Errorf("page %d  or per page %d is not valid", page, perPage))
	}
	offset := perPage * (page - 1)
	response, err = g.data.GetGroups(offset, perPage)
	response.Page = page
	response.PerPage = perPage
	return response, err
}

func (g *groupService) AddUser(userID uint, groupID uint) (err error) {
	return g.data.AddUser(userID, groupID)
}

func (g *groupService) RemoveUser(groupID uint, userID uint) (err error) {
	return g.data.RemoveUser(groupID, userID)
}
