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
	createGroupObj = `{"name":"%s"}`
	updateGoupObj  = `{"name":"%s"}`
	addUserObj     = `{"user_id":%d}`
)

func (suite *IntegrationTestSuite) TestCreateGroup() {
	dataService := data.NewGroupService(suite.testDB)
	groupService := service.NewGroupService(dataService)
	router := gin.Default()
	router.POST("/", httpservice.CreateGroupHandler(groupService))

	suite.T().Run("create group successfully", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(createGroupObj, "test")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		if recorder.Code == http.StatusCreated {
			var response internal.UserResponse
			err := json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, "test", response.Name)
		}
	})

	suite.T().Run("fail on duplicate group with same name", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(createGroupObj, "test")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	suite.T().Run("recreate deleted group", func(t *testing.T) {
		grp, err := dataService.CreateGroup(internal.GroupRequest{
			Name: "test3",
		})
		assert.NoError(t, err)
		err = suite.testDB.Delete(&data.Group{}, grp.ID).Error
		assert.NoError(t, err)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(createGroupObj, "test3")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
		if recorder.Code == http.StatusCreated {
			var recreatedResponse internal.UserResponse
			err = json.NewDecoder(recorder.Body).Decode(&recreatedResponse)
			assert.NoError(t, err)
			assert.Equal(t, grp.Name, recreatedResponse.Name)
			assert.NotEqual(t, grp.ID, recreatedResponse.ID)
		}
	})

	suite.cleanGroups()
}

func (suite *IntegrationTestSuite) TestUpdateGroup() {
	dataService := data.NewGroupService(suite.testDB)
	grpService := service.NewGroupService(dataService)
	router := gin.Default()
	router.PUT("/groups/:id", httpservice.UpdateGroupHandler(grpService))

	suite.T().Run("update group successfully", func(t *testing.T) {
		response, err := dataService.CreateGroup(internal.GroupRequest{
			Name: "test",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/groups/%d", response.ID),
			strings.NewReader(fmt.Sprintf(updateGoupObj, "update")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			updatedGrp := data.Group{ID: response.ID}
			err = suite.testDB.Find(&updatedGrp).Error
			assert.NoError(t, err)
			assert.Equal(t, "update", updatedGrp.Name)
		}
	})

	suite.T().Run("fail updating same or existing name", func(t *testing.T) {
		response, err := dataService.CreateGroup(internal.GroupRequest{
			Name: "test2",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/groups/%d", response.ID),
			strings.NewReader(fmt.Sprintf(updateGoupObj, "test2")))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		updatedGrp := data.Group{ID: response.ID}
		err = suite.testDB.Find(&updatedGrp).Error
		assert.NoError(t, err)
		assert.Equal(t, "test2", updatedGrp.Name)
	})

	suite.cleanGroups()
}

func (suite *IntegrationTestSuite) TestDeleteGroup() {
	dataService := data.NewGroupService(suite.testDB)
	grpService := service.NewGroupService(dataService)
	router := gin.Default()
	router.DELETE("/groups/:id", httpservice.DeleteGroupHandler(grpService))

	suite.T().Run("delete group successfully", func(t *testing.T) {
		response, err := dataService.CreateGroup(internal.GroupRequest{
			Name: "test",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/groups/%d", response.ID), nil)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			deletedGrp := data.Group{ID: response.ID}
			err = suite.testDB.Find(&deletedGrp).Error
			assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
		}
	})
	suite.cleanGroups()
}

func (suite *IntegrationTestSuite) TestGetGroups() {
	dataService := data.NewGroupService(suite.testDB)
	grpService := service.NewGroupService(dataService)
	router := gin.Default()
	router.GET("/", httpservice.GetGroupsHandler(grpService))

	suite.T().Run("get groups successfully", func(t *testing.T) {
		_, err := dataService.CreateGroup(internal.GroupRequest{
			Name: "test1",
		})
		assert.NoError(t, err)
		_, err = dataService.CreateGroup(internal.GroupRequest{
			Name: "test2",
		})
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/?page=%d&perPage=%d", 1, 100), nil)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			var response internal.GroupsResponse
			err := json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Len(t, response.Groups, 2)
		}
	})

	suite.T().Run("get groups successfully without query params", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			var response internal.GroupsResponse
			err := json.NewDecoder(recorder.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Len(t, response.Groups, 2)
		}
	})

	suite.cleanGroups()
}

func (suite *IntegrationTestSuite) TestAddUser() {
	dataService := data.NewGroupService(suite.testDB)
	groupService := service.NewGroupService(dataService)
	router := gin.Default()
	router.POST("/groups/:id/users", httpservice.AddUserHandler(groupService))
	grp1, grp2, usr1, usr2 := suite.addUsersAndGroups()

	suite.T().Run("add user to group successfully", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/groups/%d/users", grp1.ID), strings.NewReader(fmt.Sprintf(addUserObj, usr1.ID)))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			response, err := dataService.GetUsersByGroupID(grp1.ID, 0, 100)
			assert.NoError(suite.T(), err)
			assert.Len(t, response.Users, 1)
		}
	})

	suite.T().Run("fail on adding user to another group", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/groups/%d/users", grp2.ID), strings.NewReader(fmt.Sprintf(addUserObj, usr1.ID)))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		response, err := dataService.GetUsersByGroupID(grp2.ID, 0, 100)
		assert.NoError(suite.T(), err)
		assert.Len(t, response.Users, 0)
	})

	suite.T().Run("deleting user clean user group", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/groups/%d/users", grp2.ID), strings.NewReader(fmt.Sprintf(addUserObj, usr2.ID)))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		response, err := dataService.GetUsersByGroupID(grp2.ID, 0, 100)
		assert.NoError(suite.T(), err)
		assert.Len(t, response.Users, 1)

		userDataService := data.NewUserService(suite.testDB)
		userDataService.DeleteUser(usr2.ID)

		response, err = dataService.GetUsersByGroupID(grp2.ID, 0, 100)
		assert.NoError(suite.T(), err)
		assert.Len(t, response.Users, 0)
	})

	suite.T().Run("deleting group clean user group", func(t *testing.T) {
		err := dataService.DeleteGroup(grp1.ID)
		assert.NoError(suite.T(), err)

		var count int64
		err = suite.testDB.Model(&data.UserGroup{}).Where("user_id = ? AND group_id = ?", usr1.ID, grp1.ID).Count(&count).Error
		assert.NoError(suite.T(), err)
		assert.Equal(t, count, int64(0))
	})

	suite.cleanUsers()
	suite.cleanGroups()
	suite.cleanUserGroups()
}

func (suite *IntegrationTestSuite) TestRemoveUser() {
	dataService := data.NewGroupService(suite.testDB)
	groupService := service.NewGroupService(dataService)
	router := gin.Default()
	router.POST("/groups/:id/users", httpservice.AddUserHandler(groupService))
	router.DELETE("/groups/:id/users/:userid", httpservice.RemoveUserHandler(groupService))
	grp1, _, usr1, _ := suite.addUsersAndGroups()

	suite.T().Run("add and remove user to/from group successfully", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/groups/%d/users", grp1.ID), strings.NewReader(fmt.Sprintf(addUserObj, usr1.ID)))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			response, err := dataService.GetUsersByGroupID(grp1.ID, 0, 100)
			assert.NoError(suite.T(), err)
			assert.Len(t, response.Users, 1)

			recorder = httptest.NewRecorder()
			req, _ = http.NewRequest("DELETE", fmt.Sprintf("/groups/%d/users/%d", grp1.ID, usr1.ID), nil)
			router.ServeHTTP(recorder, req)
			response, err = dataService.GetUsersByGroupID(grp1.ID, 0, 100)
			assert.NoError(suite.T(), err)
			assert.Len(t, response.Users, 0)
		}
	})

	suite.cleanUsers()
	suite.cleanGroups()
	suite.cleanUserGroups()
}

func (suite *IntegrationTestSuite) TestGetGroupUsers() {
	dataService := data.NewGroupService(suite.testDB)
	groupService := service.NewGroupService(dataService)
	router := gin.Default()
	router.POST("/groups/:id/users", httpservice.AddUserHandler(groupService))
	router.GET("/groups/:id/users", httpservice.GetGroupUsersHandler(groupService))
	grp1, _, usr1, _ := suite.addUsersAndGroups()

	suite.T().Run("get users of a group successfully", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/groups/%d/users", grp1.ID), strings.NewReader(fmt.Sprintf(addUserObj, usr1.ID)))
		router.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		if recorder.Code == http.StatusOK {
			response, err := dataService.GetUsersByGroupID(grp1.ID, 0, 100)
			assert.NoError(suite.T(), err)
			assert.Len(t, response.Users, 1)

			recorder = httptest.NewRecorder()
			req, _ = http.NewRequest("GET", fmt.Sprintf("/groups/%d/users?page=1&perPage=10", grp1.ID), nil)
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusOK, recorder.Code)
			if recorder.Code == http.StatusOK {
				var users internal.UsersResponse
				err = json.NewDecoder(recorder.Body).Decode(&users)
				assert.NoError(t, err)
				assert.Equal(t, users.Total, uint(1))
				assert.Equal(t, users.Page, uint(1))
				assert.Equal(t, users.PerPage, uint(10))
				assert.Len(t, response.Users, 1)
			}
		}
	})

	suite.cleanUsers()
	suite.cleanGroups()
	suite.cleanUserGroups()
}

func (suite *IntegrationTestSuite) addUsersAndGroups() (internal.GroupResponse, internal.GroupResponse, internal.UserResponse, internal.UserResponse) {
	grpDataService := data.NewGroupService(suite.testDB)
	userDataService := data.NewUserService(suite.testDB)
	grp1, err := grpDataService.CreateGroup(internal.GroupRequest{
		Name: "grp1",
	})
	assert.NoError(suite.T(), err)
	grp2, err := grpDataService.CreateGroup(internal.GroupRequest{
		Name: "grp2",
	})
	assert.NoError(suite.T(), err)
	usr1, err := userDataService.CreateUser(internal.UserRequest{
		Name:     "usr1",
		Email:    "usr1@gmail.com",
		Password: "123455664546",
	})
	assert.NoError(suite.T(), err)
	usr2, err := userDataService.CreateUser(internal.UserRequest{
		Name:     "usr2",
		Email:    "usr2@gmail.com",
		Password: "123455664546",
	})
	assert.NoError(suite.T(), err)
	return grp1, grp2, usr1, usr2
}

func (suite *IntegrationTestSuite) cleanGroups() {
	err := suite.testDB.Where("1 = 1").Delete(&data.Group{}).Error
	assert.NoError(suite.T(), err)
}

func (suite *IntegrationTestSuite) cleanUserGroups() {
	err := suite.testDB.Where("1 = 1").Delete(&data.UserGroup{}).Error
	assert.NoError(suite.T(), err)
}
