package httpservice

import (
	"net/http"
	"strconv"
	"usermanagement/app/internal"
	"usermanagement/app/internal/serviceerror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUser struct {
	Name  string `json:"name" validate:"omitempty"`
	Email string `json:"email" validate:"omitempty,email"`
}

type CreatePassword struct {
	Password string `json:"password" validate:"required,min=6"`
}

func CreateUserHandler(userService internal.UserService) gin.HandlerFunc {
	mapCreateUserRequest := func(request CreateUser) internal.UserRequest {
		return internal.UserRequest{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		}
	}
	return func(c *gin.Context) {
		var request CreateUser
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		v := validator.New()
		if err := v.Struct(request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		response, err := userService.CreateUser(mapCreateUserRequest(request))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.JSON(http.StatusCreated, response)
	}
}

func UpdateUserHandler(userService internal.UserService) gin.HandlerFunc {
	mapUpdateUserRequest := func(id uint, request UpdateUser) internal.UpdateUserRequest {
		return internal.UpdateUserRequest{
			ID:    id,
			Name:  request.Name,
			Email: request.Email,
		}
	}
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		var request UpdateUser
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		v := validator.New()
		if err := v.Struct(request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = userService.UpdateUser(mapUpdateUserRequest(uint(id), request))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func DeleteUserHandler(userService internal.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = userService.DeleteUser(uint(id))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func GetUsersHandler(userService internal.UserService) gin.HandlerFunc {
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
		response, err := userService.GetUsers(uint(page), uint(perPage))
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func ChangePasswordHandler(userService internal.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		var request CreatePassword
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		v := validator.New()
		if err := v.Struct(request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = userService.ChangePassword(uint(id), request.Password)
		if err != nil {
			serviceerror.AbortOnError(c, err)
			return
		}
		c.Status(http.StatusOK)
	}
}
