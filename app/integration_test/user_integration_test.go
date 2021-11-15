package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"usermanagement/app/internal"
	"usermanagement/app/internal/data"
	"usermanagement/app/internal/httpservice"
	"usermanagement/app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	createUserObj      = `{"name":"%s","email":"%s","password":"%s"}`
	updateUserObj      = `{"name":"%s","email":"%s"}`
	updateUserNameObj  = `{"name":"%s"}`
	updateUserEmailObj = `{"email":"%s"}`
	changePasswordObj  = `{"password":"%s"}`
)

func (suite *IntegrationTestSuite) TestCreateUser() {
	dataService := data.NewUserService(suite.testDB)
	userService := service.NewUserService(dataService)
	router := gin.Default()
	router.POST("/", httpservice.CreateUserHandler(userService))

	suite.T().Run("create user successfully", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(createUserObj, "test", "test@gmail.com", "12345678")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		if recorder.Code == http.StatusCreated {
			var response internal.UserResponse
			err := json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, "test", response.Name)
			assert.Equal(t, "test@gmail.com", response.Email)
		}
	})

	suite.T().Run("fail on duplicate user with same email", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(createUserObj, "test2", "test@gmail.com", "33333333333")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	suite.T().Run("recreate deleted user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(createUserObj, "test3", "test3@gmail.com", "33333333333")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		if recorder.Code == http.StatusCreated {
			var response internal.UserResponse
			err := json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, "test3", response.Name)
			assert.Equal(t, "test3@gmail.com", response.Email)

			err = suite.testDB.Delete(&data.User{}, response.ID).Error
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(createUserObj, "test3", "test3@gmail.com", "33333333333")))
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusCreated, recorder.Code)
			var recreatedResponse internal.UserResponse
			err = json.NewDecoder(recorder.Body).Decode(&recreatedResponse)
			assert.NoError(t, err)
			assert.Equal(t, response.Email, recreatedResponse.Email)
			assert.NotEqual(t, response.ID, recreatedResponse.ID)
		}
	})

	suite.cleanUsers()
}

func (suite *IntegrationTestSuite) TestUpdateUser() {
	dataService := data.NewUserService(suite.testDB)
	userService := service.NewUserService(dataService)
	router := gin.Default()
	router.PUT("/users/:id", httpservice.UpdateUserHandler(userService))

	suite.T().Run("update user successfully", func(t *testing.T) {
		response, err := dataService.CreateUser(internal.UserRequest{
			Name:     "test",
			Email:    "test@gmail.com",
			Password: "123455664546",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", response.ID),
			strings.NewReader(fmt.Sprintf(updateUserObj, "update", "test.update@gmail.com")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			updatedUser := data.User{ID: response.ID}
			err = suite.testDB.Find(&updatedUser).Error
			assert.NoError(t, err)
			assert.Equal(t, "update", updatedUser.Name)
			assert.Equal(t, "test.update@gmail.com", updatedUser.Email)
		}
	})

	suite.T().Run("update name alone successfully", func(t *testing.T) {
		response, err := dataService.CreateUser(internal.UserRequest{
			Name:     "test2",
			Email:    "test2@gmail.com",
			Password: "dsfsfsdfdsfdsfsdfs",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", response.ID),
			strings.NewReader(fmt.Sprintf(updateUserNameObj, "update2")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			updatedUser := data.User{ID: response.ID}
			err = suite.testDB.Find(&updatedUser).Error
			assert.NoError(t, err)
			assert.Equal(t, "update2", updatedUser.Name)
			assert.Equal(t, "test2@gmail.com", updatedUser.Email)
		}
	})

	suite.T().Run("update email alone successfully", func(t *testing.T) {
		response, err := dataService.CreateUser(internal.UserRequest{
			Name:     "test3",
			Email:    "test3@gmail.com",
			Password: "dsfsfsdfdsfdsfsdfs",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", response.ID),
			strings.NewReader(fmt.Sprintf(updateUserEmailObj, "test3.update@gmail.com")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			updatedUser := data.User{ID: response.ID}
			err = suite.testDB.Find(&updatedUser).Error
			assert.NoError(t, err)
			assert.Equal(t, "test3", updatedUser.Name)
			assert.Equal(t, "test3.update@gmail.com", updatedUser.Email)
		}
	})

	suite.T().Run("fail updating same email", func(t *testing.T) {
		response, err := dataService.CreateUser(internal.UserRequest{
			Name:     "test3",
			Email:    "test3@gmail.com",
			Password: "dsfsfsdfdsfdsfsdfs",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", response.ID),
			strings.NewReader(fmt.Sprintf(updateUserEmailObj, "test3@gmail.com")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		updatedUser := data.User{ID: response.ID}
		err = suite.testDB.Find(&updatedUser).Error
		assert.NoError(t, err)
		assert.Equal(t, "test3", updatedUser.Name)
		assert.Equal(t, "test3@gmail.com", updatedUser.Email)
	})

	suite.cleanUsers()
}

func (suite *IntegrationTestSuite) TestDeleteUser() {
	dataService := data.NewUserService(suite.testDB)
	userService := service.NewUserService(dataService)
	router := gin.Default()
	router.DELETE("/users/:id", httpservice.DeleteUserHandler(userService))

	suite.T().Run("delete user successfully", func(t *testing.T) {
		response, err := dataService.CreateUser(internal.UserRequest{
			Name:     "test",
			Email:    "test@gmail.com",
			Password: "123455664546",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", response.ID), nil)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			deletedUser := data.User{ID: response.ID}
			err = suite.testDB.Find(&deletedUser).Error
			assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
		}
	})
	suite.cleanUsers()
}

func (suite *IntegrationTestSuite) TestChangePassword() {
	dataService := data.NewUserService(suite.testDB)
	userService := service.NewUserService(dataService)
	router := gin.Default()
	router.PUT("/:id/password", httpservice.ChangePasswordHandler(userService))

	suite.T().Run("change password successfully", func(t *testing.T) {
		response, err := dataService.CreateUser(internal.UserRequest{
			Name:     "test",
			Email:    "test@gmail.com",
			Password: "123455664546",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/%d/password", response.ID),
			strings.NewReader(fmt.Sprintf(changePasswordObj, "453453535353435")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			updatedUser := data.User{ID: response.ID}
			err = suite.testDB.Find(&updatedUser).Error
			assert.NoError(t, err)
			assert.Equal(t, "test", updatedUser.Name)
			assert.Equal(t, "test@gmail.com", updatedUser.Email)
		}
	})

	suite.cleanUsers()
}

func (suite *IntegrationTestSuite) TestGetUsers() {
	dataService := data.NewUserService(suite.testDB)
	userService := service.NewUserService(dataService)
	router := gin.Default()
	router.GET("/", httpservice.GetUsersHandler(userService))

	suite.T().Run("get users successfully", func(t *testing.T) {
		_, err := dataService.CreateUser(internal.UserRequest{
			Name:     "test1",
			Email:    "test1@gmail.com",
			Password: "123455664546",
		})
		assert.NoError(t, err)
		_, err = dataService.CreateUser(internal.UserRequest{
			Name:     "test2",
			Email:    "test2@gmail.com",
			Password: "123455664546",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/?page=%d&perPage=%d", 1, 100), nil)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			var response internal.UsersResponse
			err := json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Len(t, response.Users, 2)
		}
	})

	suite.T().Run("get users successfully without query params", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			var response internal.UsersResponse
			err := json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Len(t, response.Users, 2)
		}
	})

	suite.cleanUsers()
}

func (suite *IntegrationTestSuite) cleanUsers() {
	err := suite.testDB.Where("1 = 1").Delete(&data.User{}).Error
	assert.NoError(suite.T(), err)
}
