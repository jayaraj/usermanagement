package config

import (
	"usermanagement/app/internal/httpservice"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *AppConfiguration) initialiseRoutes() {
	a.engine.Use(gin.Recovery())
	a.engine.Use(cors.Default())
	v1 := a.engine.Group("api/v1")
	a.addV1Routes(v1)
}

func (a *AppConfiguration) addV1Routes(router *gin.RouterGroup) {
	users := router.Group("/users")
	a.addUserRouters(users)
	groups := router.Group("/groups")
	a.addGroupRouters(groups)
}

func (a *AppConfiguration) addUserRouters(router *gin.RouterGroup) {
	router.POST("", httpservice.CreateUserHandler(a.userService))
	router.PUT("/:id", httpservice.UpdateUserHandler(a.userService))
	router.DELETE("/:id", httpservice.DeleteUserHandler(a.userService))
	router.GET("", httpservice.GetUsersHandler(a.userService))
	router.PUT("/:id/password", httpservice.ChangePasswordHandler(a.userService))
}

func (a *AppConfiguration) addGroupRouters(router *gin.RouterGroup) {
	router.POST("", httpservice.CreateGroupHandler(a.groupService))
	router.PUT("/:id", httpservice.UpdateGroupHandler(a.groupService))
	router.DELETE("/:id", httpservice.DeleteGroupHandler(a.groupService))
	router.GET("/:id/users", httpservice.GetGroupUsersHandler(a.groupService))
	router.GET("", httpservice.GetGroupsHandler(a.groupService))
	router.POST("/:id/users", httpservice.AddUserHandler(a.groupService))
	router.DELETE("/:id/users/:userid", httpservice.RemoveUserHandler(a.groupService))
}
