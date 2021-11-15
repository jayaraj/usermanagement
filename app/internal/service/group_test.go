package service_test

import (
	"fmt"
	"testing"
	"usermanagement/app/internal"
	"usermanagement/app/internal/mock"
	"usermanagement/app/internal/service"
	"usermanagement/app/internal/serviceerror"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	data := mock.NewMockGroupData(mockCtrl)
	handler := service.NewGroupService(data)

	t.Run("get groups successfully", func(t *testing.T) {
		data.EXPECT().GetGroups(uint(100), uint(100)).Return(internal.GroupsResponse{
			Total: 0,
		}, nil).Times(1)
		response, err := handler.GetGroups(2, 100)
		assert.NoError(t, err)
		assert.Equal(t, internal.GroupsResponse{
			Total:   0,
			Page:    2,
			PerPage: 100,
		}, response)
	})

	t.Run("error on missing page", func(t *testing.T) {
		response, err := handler.GetGroups(0, 100)
		assert.Equal(t, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, fmt.Errorf("page %d  or per page %d is not valid", 0, 100)), err)
		assert.Equal(t, internal.GroupsResponse{}, response)
	})

	t.Run("error on missing perPage", func(t *testing.T) {
		response, err := handler.GetGroups(1, 0)
		assert.Equal(t, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, fmt.Errorf("page %d  or per page %d is not valid", 1, 0)), err)
		assert.Equal(t, internal.GroupsResponse{}, response)
	})
}

func TestGetGroupUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	data := mock.NewMockGroupData(mockCtrl)
	handler := service.NewGroupService(data)

	t.Run("get group users successfully", func(t *testing.T) {
		data.EXPECT().GetUsersByGroupID(uint(1), uint(0), uint(100)).Return(internal.UsersResponse{
			Total: 0,
		}, nil).Times(1)
		response, err := handler.GetUsersByGroupID(1, 1, 100)
		assert.NoError(t, err)
		assert.Equal(t, internal.UsersResponse{
			Total:   0,
			Page:    1,
			PerPage: 100,
		}, response)
	})

	t.Run("error on missing page", func(t *testing.T) {
		response, err := handler.GetUsersByGroupID(1, 0, 100)
		assert.Equal(t, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, fmt.Errorf("page %d  or per page %d is not valid", 0, 100)), err)
		assert.Equal(t, internal.UsersResponse{}, response)
	})

	t.Run("error on missing perPage", func(t *testing.T) {
		response, err := handler.GetUsersByGroupID(1, 1, 0)
		assert.Equal(t, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, fmt.Errorf("page %d  or per page %d is not valid", 1, 0)), err)
		assert.Equal(t, internal.UsersResponse{}, response)
	})
}
