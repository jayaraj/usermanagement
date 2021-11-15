package httpservice

import (
	"net/http"
	"strconv"
	"usermanagement/app/internal"
	"usermanagement/app/internal/serviceerror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateGroup struct {
	Name string `json:"name" validate:"required"`
}

type UpdateGroup struct {
	Name string `json:"name" validate:"required"`
}

type AddUser struct {
	UserID uint `json:"user_id" validate:"required"`
}

func CreateGroupHandler(grpService internal.GroupService) gin.HandlerFunc {
	mapCreateGroupRequest := func(request CreateGroup) internal.GroupRequest {
		return internal.GroupRequest{
			Name: request.Name,
		}
	}
	return func(c *gin.Context) {
		var request CreateGroup
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		v := validator.New()
		if err := v.Struct(request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		response, err := grpService.CreateGroup(mapCreateGroupRequest(request))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.JSON(http.StatusCreated, response)
	}
}

func UpdateGroupHandler(grpService internal.GroupService) gin.HandlerFunc {
	mapUpdateGroupRequest := func(id uint, request UpdateGroup) internal.UpdateGroupRequest {
		return internal.UpdateGroupRequest{
			ID:   id,
			Name: request.Name,
		}
	}
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		var request UpdateGroup
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		v := validator.New()
		if err := v.Struct(request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = grpService.UpdateGroup(mapUpdateGroupRequest(uint(id), request))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func DeleteGroupHandler(grpService internal.GroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = grpService.DeleteGroup(uint(id))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func GetGroupUsersHandler(grpService internal.GroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		perPage, err := strconv.ParseUint(c.DefaultQuery("perPage", "10"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		response, err := grpService.GetUsersByGroupID(uint(id), uint(page), uint(perPage))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetGroupsHandler(grpService internal.GroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		perPage, err := strconv.ParseUint(c.DefaultQuery("perPage", "10"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		response, err := grpService.GetGroups(uint(page), uint(perPage))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func AddUserHandler(grpService internal.GroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		var request AddUser
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		v := validator.New()
		if err := v.Struct(request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = grpService.AddUser(request.UserID, uint(id))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func RemoveUserHandler(grpService internal.GroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		userid, err := strconv.ParseUint(c.Param("userid"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = grpService.RemoveUser(uint(id), uint(userid))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.Status(http.StatusOK)
	}
}
