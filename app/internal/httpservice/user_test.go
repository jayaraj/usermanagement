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
	createUserObj     = `{"name":"%s","email":"%s","password":"%s"}`
	responseUserObj   = `{"id":%d,"name":"%s","email":"%s"}`
	updateUserObj     = `{"name":"%s","email":"%s"}`
	responseUsersObj  = `{"users":[{"id":%d,"name":"%s","email":"%s"}],"total":%d,"page":%d,"perPage":%d}`
	changePasswordObj = `{"password":"%s"}`
)

func TestCreateUserHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userService := mock.NewMockUserService(mockCtrl)
	router := gin.Default()
	router.POST("/", httpservice.CreateUserHandler(userService))

	tests := []struct {
		name     string
		request  string
		status   int
		response string
		setup    func()
	}{
		{
			name:     "create user successfully",
			request:  fmt.Sprintf(createUserObj, "test", "test@gmail.com", "12345678"),
			status:   http.StatusCreated,
			response: fmt.Sprintf(responseUserObj, 1, "test", "test@gmail.com"),
			setup: func() {
				request := internal.UserRequest{
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "12345678",
				}
				userService.EXPECT().CreateUser(request).Return(internal.UserResponse{
					ID:    1,
					Name:  "test",
					Email: "test@gmail.com",
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
			request:  fmt.Sprintf(createUserObj, "", "test@gmail.com", "12345678"),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'CreateUser.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
			setup:    func() {},
		},
		{
			name:     "fail on wrong email",
			request:  fmt.Sprintf(createUserObj, "test", "test.com", "12345678"),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'CreateUser.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
			setup:    func() {},
		},
		{
			name:     "fail on small password",
			request:  fmt.Sprintf(createUserObj, "test", "test@gmail.com", "123"),
			status:   http.StatusBadRequest,
			response: `{"message":"Key: 'CreateUser.Password' Error:Field validation for 'Password' failed on the 'min' tag"}`,
			setup:    func() {},
		},
		{
			name:     "fail on service error",
			request:  fmt.Sprintf(createUserObj, "test", "test@gmail.com", "12345678"),
			status:   http.StatusBadRequest,
			response: `{"message":"Invalid User Request : test"}`,
			setup: func() {
				request := internal.UserRequest{
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "12345678",
				}
				userService.EXPECT().CreateUser(request).Return(internal.UserResponse{},
					serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:     "fail on unknown error",
			request:  fmt.Sprintf(createUserObj, "test", "test@gmail.com", "12345678"),
			status:   http.StatusInternalServerError,
			response: `{"message":"test"}`,
			setup: func() {
				request := internal.UserRequest{
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "12345678",
				}
				userService.EXPECT().CreateUser(request).Return(internal.UserResponse{}, errors.New("test")).Times(1)
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

func TestUpdateUserHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userService := mock.NewMockUserService(mockCtrl)
	router := gin.Default()
	router.PUT("/users/:id", httpservice.UpdateUserHandler(userService))

	tests := []struct {
		name    string
		request string
		status  int
		setup   func()
	}{
		{
			name:    "update user successfully",
			request: fmt.Sprintf(updateUserObj, "test", "test@gmail.com"),
			status:  http.StatusOK,
			setup: func() {
				request := internal.UpdateUserRequest{
					ID:    1,
					Name:  "test",
					Email: "test@gmail.com",
				}
				userService.EXPECT().UpdateUser(request).Return(nil).Times(1)
			},
		},
		{
			name:    "fail on body unmarshal",
			request: `{"name":"test",`,
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on wrong email",
			request: fmt.Sprintf(updateUserObj, "test", "test.com"),
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "update only name",
			request: fmt.Sprintf(updateUserObj, "test", ""),
			status:  http.StatusOK,
			setup: func() {
				request := internal.UpdateUserRequest{
					ID:   1,
					Name: "test",
				}
				userService.EXPECT().UpdateUser(request).Return(nil).Times(1)
			},
		},
		{
			name:    "update only email",
			request: fmt.Sprintf(updateUserObj, "", "test@gmail.com"),
			status:  http.StatusOK,
			setup: func() {
				request := internal.UpdateUserRequest{
					ID:    1,
					Email: "test@gmail.com",
				}
				userService.EXPECT().UpdateUser(request).Return(nil).Times(1)
			},
		},
		{
			name:    "fail on service error",
			request: fmt.Sprintf(updateUserObj, "test", ""),
			status:  http.StatusBadRequest,
			setup: func() {
				request := internal.UpdateUserRequest{
					ID:   1,
					Name: "test",
				}
				userService.EXPECT().UpdateUser(request).
					Return(serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:    "fail on unknown error",
			request: fmt.Sprintf(updateUserObj, "test", ""),
			status:  http.StatusInternalServerError,
			setup: func() {
				request := internal.UpdateUserRequest{
					ID:   1,
					Name: "test",
				}
				userService.EXPECT().UpdateUser(request).Return(errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(test.request))
			test.setup()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.status, recorder.Code)
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userService := mock.NewMockUserService(mockCtrl)
	router := gin.Default()
	router.DELETE("/users/:id", httpservice.DeleteUserHandler(userService))

	tests := []struct {
		name   string
		status int
		setup  func()
	}{
		{
			name:   "Delete user successfully",
			status: http.StatusOK,
			setup: func() {
				userService.EXPECT().DeleteUser(uint(1)).Return(nil).Times(1)
			},
		},
		{
			name:   "fail on service error",
			status: http.StatusBadRequest,
			setup: func() {
				userService.EXPECT().DeleteUser(uint(1)).Return(serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:   "fail on unknown error",
			status: http.StatusInternalServerError,
			setup: func() {
				userService.EXPECT().DeleteUser(uint(1)).Return(errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/users/1", nil)
			test.setup()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.status, recorder.Code)
		})
	}
}

func TestGetUsersHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userService := mock.NewMockUserService(mockCtrl)
	router := gin.Default()
	router.GET("/", httpservice.GetUsersHandler(userService))

	tests := []struct {
		name     string
		status   int
		query    string
		response string
		setup    func()
	}{
		{
			name:     "Get users successfully",
			status:   http.StatusOK,
			query:    "/?page=1&perPage=100",
			response: fmt.Sprintf(responseUsersObj, 1, "test", "test@gmail.com", 1, 1, 100),
			setup: func() {
				userService.EXPECT().GetUsers(uint(1), uint(100)).Return(internal.UsersResponse{
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
			response: fmt.Sprintf(responseUsersObj, 1, "test", "test@gmail.com", 1, 1, 10),
			setup: func() {
				userService.EXPECT().GetUsers(uint(1), uint(10)).Return(internal.UsersResponse{
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
			query:    "/?page=1&perPage=100",
			response: `{"message":"Invalid User Request : test"}`,
			setup: func() {
				userService.EXPECT().GetUsers(uint(1), uint(100)).Return(internal.UsersResponse{}, serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:     "fail on unknown error",
			status:   http.StatusInternalServerError,
			query:    "/?page=1&perPage=100",
			response: `{"message":"test"}`,
			setup: func() {
				userService.EXPECT().GetUsers(uint(1), uint(100)).Return(internal.UsersResponse{}, errors.New("test")).Times(1)
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

func TestChangePasswordHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userService := mock.NewMockUserService(mockCtrl)
	router := gin.Default()
	router.PUT("/:id/password", httpservice.ChangePasswordHandler(userService))

	tests := []struct {
		name    string
		request string
		status  int
		setup   func()
	}{
		{
			name:    "change password successfully",
			request: fmt.Sprintf(changePasswordObj, "1234567890"),
			status:  http.StatusOK,
			setup: func() {
				userService.EXPECT().ChangePassword(uint(1), "1234567890").Return(nil).Times(1)
			},
		},
		{
			name:    "fail on body unmarshal",
			request: `{"password":"12321312",`,
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on small password",
			request: fmt.Sprintf(changePasswordObj, "123"),
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on missing password",
			request: fmt.Sprintf(changePasswordObj, ""),
			status:  http.StatusBadRequest,
			setup:   func() {},
		},
		{
			name:    "fail on service error",
			request: fmt.Sprintf(changePasswordObj, "1234567890"),
			status:  http.StatusBadRequest,
			setup: func() {
				userService.EXPECT().ChangePassword(uint(1), "1234567890").
					Return(serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("test"))).Times(1)
			},
		},
		{
			name:    "fail on unknown error",
			request: fmt.Sprintf(changePasswordObj, "1234567890"),
			status:  http.StatusInternalServerError,
			setup: func() {
				userService.EXPECT().ChangePassword(uint(1), "1234567890").
					Return(errors.New("test")).Times(1)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/1/password", strings.NewReader(test.request))
			test.setup()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.status, recorder.Code)
		})
	}
}
