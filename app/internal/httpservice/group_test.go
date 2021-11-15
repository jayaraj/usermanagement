package httpservice_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"usermanagement/app/internal"
	"usermanagement/app/internal/httpservice"
	"usermanagement/app/internal/mock"
	"usermanagement/app/internal/serviceerror"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	createGroupObj    = `{"name":"%s"}`
	responseGroupObj  = `{"id":%d,"name":"%s"}`
	updateGroupObj    = `{"name":"%s"}`
	responseGroupsObj = `{"groups":[{"id":%d,"name":"%s"}],"total":%d,"page":%d,"perPage":%d}`
	addUserObj        = `{"user_id":%d}`
)

func TestCreateGroupHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	groupService := mock.NewMockGroupService(mockCtrl)
	router := gin.Default()
	router.POST("/", httpservice.CreateGroupHandler(groupService))

	tests := []struct {
		name     string
		request  string
		status   int
		response string
		setup    func()
	}{
		{
			name:     "create group successfully",
			request:  fmt.Sprintf(createGroupObj, "test"),
			status:   http.StatusCreated,
			response: fmt.Sprintf(responseGroupObj, 1, "test"),
			setup: func() {
				request := internal.GroupRequest{
					Name: "test",
				}
				groupService.EXPECT().CreateGroup(request).Return(internal.GroupResponse{
					ID:   1,
					Name: "test",
				}, nil).Times(1)
			},
		},
		{
			name:     "fail on body unmarshal",
			request:  `{"name":"test",`,
			status:   http.StatusBadRequest,
			response: `{"message":"unexpected EOF"}`,
			setup:    func() {},
		},
		{
			name:     "fail on missing name",
			request:  fmt.Sprintf(createGroupObj, ""),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'CreateGroup.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
			setup:    func() {},
		},
		{
			name:     "fail on service error",
			request:  fmt.Sprintf(createGroupObj, "test"),
			status:   http.StatusBadRequest,
			response: `{"message":"Invalid Group Request : test"}`,
			setup: func() {
				request := internal.GroupRequest{
					Name: "test",
				}
				groupService.EXPECT().CreateGroup(request).Return(internal.GroupResponse{},
					serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:     "fail on unknown error",
			request:  fmt.Sprintf(createGroupObj, "test"),
			status:   http.StatusInternalServerError,
			response: `{"message":"test"}`,
			setup: func() {
				request := internal.GroupRequest{
					Name: "test",
				}
				groupService.EXPECT().CreateGroup(request).Return(internal.GroupResponse{}, errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/", strings.NewReader(test.request))
			test.setup()
			router.ServeHTTP(recorder, req)
			response, err := ioutil.ReadAll(recorder.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.status, recorder.Code)
			assert.Equal(t, test.response, string(response))
		})
	}
}

func TestUpdateGroupHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	groupService := mock.NewMockGroupService(mockCtrl)
	router := gin.Default()
	router.PUT("/groups/:id", httpservice.UpdateGroupHandler(groupService))

	tests := []struct {
		name    string
		request string
		status  int
		setup   func()
	}{
		{
			name:    "update group successfully",
			request: fmt.Sprintf(updateGroupObj, "test"),
			status:  http.StatusOK,
			setup: func() {
				request := internal.UpdateGroupRequest{
					ID:   1,
					Name: "test",
				}
				groupService.EXPECT().UpdateGroup(request).Return(nil).Times(1)
			},
		},
		{
			name:    "fail on body unmarshal",
			request: `{"name":"test",`,
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on missing name",
			request: fmt.Sprintf(updateGroupObj, ""),
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on service error",
			request: fmt.Sprintf(updateGroupObj, "test"),
			status:  http.StatusBadRequest,
			setup: func() {
				request := internal.UpdateGroupRequest{
					ID:   1,
					Name: "test",
				}
				groupService.EXPECT().UpdateGroup(request).
					Return(serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("test"))).Times(1)

			},
		},
		{
			name:    "fail on unknown error",
			request: fmt.Sprintf(updateGroupObj, "test"),
			status:  http.StatusInternalServerError,
			setup: func() {
				request := internal.UpdateGroupRequest{
					ID:   1,
					Name: "test",
				}
				groupService.EXPECT().UpdateGroup(request).
					Return(errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/groups/1", strings.NewReader(test.request))
			test.setup()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.status, recorder.Code)
		})
	}
}

func TestDeleteGroupHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	groupService := mock.NewMockGroupService(mockCtrl)
	router := gin.Default()
	router.DELETE("/groups/:id", httpservice.DeleteGroupHandler(groupService))

	tests := []struct {
		name   string
		status int
		setup  func()
	}{
		{
			name:   "Delete group successfully",
			status: http.StatusOK,
			setup: func() {
				groupService.EXPECT().DeleteGroup(uint(1)).Return(nil).Times(1)
			},
		},
		{
			name:   "fail on service error",
			status: http.StatusBadRequest,
			setup: func() {
				groupService.EXPECT().DeleteGroup(uint(1)).Return(serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:   "fail on unknown error",
			status: http.StatusInternalServerError,
			setup: func() {
				groupService.EXPECT().DeleteGroup(uint(1)).Return(errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/groups/1", nil)
			test.setup()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.status, recorder.Code)
		})
	}
}

func TestGetGroupsHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	groupService := mock.NewMockGroupService(mockCtrl)
	router := gin.Default()
	router.GET("/", httpservice.GetGroupsHandler(groupService))

	tests := []struct {
		name     string
		status   int
		query    string
		response string
		setup    func()
	}{
		{
			name:     "Get groups successfully",
			status:   http.StatusOK,
			query:    "/?page=1&perPage=100",
			response: fmt.Sprintf(responseGroupsObj, 1, "test", 1, 1, 100),
			setup: func() {
				groupService.EXPECT().GetGroups(uint(1), uint(100)).Return(internal.GroupsResponse{
					Groups: []internal.GroupResponse{{
						ID:   1,
						Name: "test",
					}},
					Total:   1,
					Page:    1,
					PerPage: 100,
				}, nil).Times(1)
			},
		},
		{
			name:     "missing page",
			status:   http.StatusBadRequest,
			query:    "/?page=&perPage=100",
			response: `{"message":"strconv.ParseUint: parsing \"\": invalid syntax"}`,
			setup:    func() {},
		},
		{
			name:     "missing perPage",
			status:   http.StatusBadRequest,
			query:    "/?page=1&perPage=",
			response: `{"message":"strconv.ParseUint: parsing \"\": invalid syntax"}`,
			setup:    func() {},
		},
		{
			name:     "default query params",
			status:   http.StatusOK,
			query:    "/",
			response: fmt.Sprintf(responseGroupsObj, 1, "test", 1, 1, 10),
			setup: func() {
				groupService.EXPECT().GetGroups(uint(1), uint(10)).Return(internal.GroupsResponse{
					Groups: []internal.GroupResponse{{
						ID:   1,
						Name: "test",
					}},
					Total:   1,
					Page:    1,
					PerPage: 10,
				}, nil).Times(1)
			},
		},
		{
			name:     "fail on service error",
			status:   http.StatusBadRequest,
			query:    "/?page=1&perPage=100",
			response: `{"message":"Invalid Group Request : test"}`,
			setup: func() {
				groupService.EXPECT().GetGroups(uint(1), uint(100)).Return(internal.GroupsResponse{},
					serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:     "fail on unknown error",
			status:   http.StatusInternalServerError,
			query:    "/?page=1&perPage=100",
			response: `{"message":"test"}`,
			setup: func() {
				groupService.EXPECT().GetGroups(uint(1), uint(100)).Return(internal.GroupsResponse{}, errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", test.query, nil)
			test.setup()
			router.ServeHTTP(recorder, req)
			response, err := ioutil.ReadAll(recorder.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.status, recorder.Code)
			assert.Equal(t, test.response, string(response))
		})
	}
}

func TestGetGroupUsersHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	groupService := mock.NewMockGroupService(mockCtrl)
	router := gin.Default()
	router.GET("/:id/users", httpservice.GetGroupUsersHandler(groupService))

	tests := []struct {
		name     string
		status   int
		query    string
		response string
		setup    func()
	}{
		{
			name:     "Get group users successfully",
			status:   http.StatusOK,
			query:    "/1/users?page=1&perPage=100",
			response: fmt.Sprintf(responseUsersObj, 1, "test", "test@gmail.com", 1, 1, 100),
			setup: func() {
				groupService.EXPECT().GetUsersByGroupID(uint(1), uint(1), uint(100)).Return(internal.UsersResponse{
					Users: []internal.UserResponse{{
						ID:    1,
						Name:  "test",
						Email: "test@gmail.com",
					}},
					Total:   1,
					Page:    1,
					PerPage: 100,
				}, nil).Times(1)
			},
		},
		{
			name:     "missing page",
			status:   http.StatusBadRequest,
			query:    "/1/users?page=&perPage=100",
			response: `{"message":"strconv.ParseUint: parsing \"\": invalid syntax"}`,
			setup:    func() {},
		},
		{
			name:     "missing perPage",
			status:   http.StatusBadRequest,
			query:    "/1/users?page=1&perPage=",
			response: `{"message":"strconv.ParseUint: parsing \"\": invalid syntax"}`,
			setup:    func() {},
		},
		{
			name:     "default query params",
			status:   http.StatusOK,
			query:    "/1/users",
			response: fmt.Sprintf(responseUsersObj, 1, "test", "test@gmail.com", 1, 1, 10),
			setup: func() {
				groupService.EXPECT().GetUsersByGroupID(uint(1), uint(1), uint(10)).Return(internal.UsersResponse{
					Users: []internal.UserResponse{{
						ID:    1,
						Name:  "test",
						Email: "test@gmail.com",
					}},
					Total:   1,
					Page:    1,
					PerPage: 10,
				}, nil).Times(1)
			},
		},
		{
			name:     "fail on service error",
			status:   http.StatusBadRequest,
			query:    "/1/users?page=1&perPage=100",
			response: `{"message":"Invalid Group Request : test"}`,
			setup: func() {
				groupService.EXPECT().GetUsersByGroupID(uint(1), uint(1), uint(100)).
					Return(internal.UsersResponse{}, serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:     "fail on unknown error",
			status:   http.StatusInternalServerError,
			query:    "/1/users?page=1&perPage=100",
			response: `{"message":"test"}`,
			setup: func() {
				groupService.EXPECT().GetUsersByGroupID(uint(1), uint(1), uint(100)).
					Return(internal.UsersResponse{}, errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", test.query, nil)
			test.setup()
			router.ServeHTTP(recorder, req)
			response, err := ioutil.ReadAll(recorder.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.status, recorder.Code)
			assert.Equal(t, test.response, string(response))
		})
	}
}

func TestAddUserHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	groupService := mock.NewMockGroupService(mockCtrl)
	router := gin.Default()
	router.POST("/:id/users", httpservice.AddUserHandler(groupService))

	tests := []struct {
		name    string
		request string
		status  int
		setup   func()
	}{
		{
			name:    "add user successfully",
			request: fmt.Sprintf(addUserObj, 3),
			status:  http.StatusOK,
			setup: func() {
				groupService.EXPECT().AddUser(uint(3), uint(1)).Return(nil).Times(1)
			},
		},
		{
			name:    "fail on body unmarshal",
			request: `{"user_id":"1",`,
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on missing user_id",
			request: `{"user_id":""}`,
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on service error",
			request: fmt.Sprintf(addUserObj, 3),
			status:  http.StatusBadRequest,
			setup: func() {
				groupService.EXPECT().AddUser(uint(3), uint(1)).
					Return(serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:    "fail on unknown error",
			request: fmt.Sprintf(addUserObj, 3),
			status:  http.StatusInternalServerError,
			setup: func() {
				groupService.EXPECT().AddUser(uint(3), uint(1)).
					Return(errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/1/users", strings.NewReader(test.request))
			test.setup()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.status, recorder.Code)
		})
	}
}

func TestRemoveUserHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	groupService := mock.NewMockGroupService(mockCtrl)
	router := gin.Default()
	router.DELETE("/groups/:id/users/:userid", httpservice.RemoveUserHandler(groupService))

	tests := []struct {
		name   string
		status int
		setup  func()
	}{
		{
			name:   "remove user successfully",
			status: http.StatusOK,
			setup: func() {
				groupService.EXPECT().RemoveUser(uint(1), uint(2)).Return(nil).Times(1)
			},
		},
		{
			name:   "fail on service error",
			status: http.StatusBadRequest,
			setup: func() {
				groupService.EXPECT().RemoveUser(uint(1), uint(2)).Return(serviceerror.NewServiceError(serviceerror.InvalidGroupRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:   "fail on unknown error",
			status: http.StatusInternalServerError,
			setup: func() {
				groupService.EXPECT().RemoveUser(uint(1), uint(2)).Return(errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/groups/1/users/2", nil)
			test.setup()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.status, recorder.Code)
		})
	}
}
