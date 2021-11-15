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

func TestGetUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	data := mock.NewMockUserData(mockCtrl)
	handler := service.NewUserService(data)

	t.Run("get users successfully", func(t *testing.T) {
		data.EXPECT().GetUsers(uint(100), uint(100)).Return(internal.UsersResponse{
			Total: 0,
		}, nil).Times(1)
		response, err := handler.GetUsers(2, 100)
		assert.NoError(t, err)
		assert.Equal(t, internal.UsersResponse{
			Total:   0,
			Page:    2,
			PerPage: 100,
		}, response)
	})

	t.Run("error on missing page", func(t *testing.T) {
		response, err := handler.GetUsers(0, 100)
		assert.Equal(t, serviceerror.NewServiceError(serviceerror.InvalidUserRequest, fmt.Errorf("page %d  or per page %d is not valid", 0, 100)), err)
		assert.Equal(t, internal.UsersResponse{}, response)
	})

	t.Run("error on missing perPage", func(t *testing.T) {
		response, err := handler.GetUsers(1, 0)
		assert.Equal(t, serviceerror.NewServiceError(serviceerror.InvalidUserRequest, fmt.Errorf("page %d  or per page %d is not valid", 1, 0)), err)
		assert.Equal(t, internal.UsersResponse{}, response)
	})
}
