package internal

//go:generate mockgen -source=service.go  -destination=mock/service.go -package=mock
type UserService interface {
	CreateUser(request UserRequest) (response UserResponse, err error)
	UpdateUser(request UpdateUserRequest) (err error)
	DeleteUser(id uint) (err error)
	GetUsers(page uint, perPage uint) (response UsersResponse, err error)
	ChangePassword(userID uint, password string) (err error)
}

type GroupService interface {
	CreateGroup(request GroupRequest) (response GroupResponse, err error)
	UpdateGroup(request UpdateGroupRequest) (err error)
	DeleteGroup(id uint) (err error)
	GetUsersByGroupID(groupID uint, page uint, perPage uint) (response UsersResponse, err error)
	GetGroups(page uint, perPage uint) (response GroupsResponse, err error)
	AddUser(userID uint, groupID uint) (err error)
	RemoveUser(groupID uint, userID uint) (err error)
}
