package service

import (
	"fmt"
	"usermanagement/app/internal"
	"usermanagement/app/internal/serviceerror"
)

type userService struct {
	data internal.UserData
}

func NewUserService(data internal.UserData) *userService {
	return &userService{
		data: data,
	}
}

func (u *userService) CreateUser(request internal.UserRequest) (response internal.UserResponse, err error) {
	return u.data.CreateUser(request)
}

func (u *userService) UpdateUser(request internal.UpdateUserRequest) (err error) {
	return u.data.UpdateUser(request)
}

func (u *userService) DeleteUser(id uint) (err error) {
	return u.data.DeleteUser(id)
}

func (u *userService) GetUsers(page uint, perPage uint) (response internal.UsersResponse, err error) {
	if page <= 0 || perPage == 0 {
		return response, serviceerror.NewServiceError(serviceerror.InvalidUserRequest, fmt.Errorf("page %d  or per page %d is not valid", page, perPage))
	}
	offset := perPage * (page - 1)
	response, err = u.data.GetUsers(offset, perPage)
	response.Page = page
	response.PerPage = perPage
	return response, err
}

func (u *userService) ChangePassword(userID uint, password string) (err error) {
	return u.data.ChangePassword(userID, password)
}
