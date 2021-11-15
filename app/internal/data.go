package internal

//go:generate mockgen -source=data.go  -destination=mock/data.go -package=mock
type UserData interface {
	CreateUser(request UserRequest) (response UserResponse, err error)
	UpdateUser(request UpdateUserRequest) (err error)
	DeleteUser(id uint) (err error)
	GetUsers(offset uint, limit uint) (response UsersResponse, err error)
	ChangePassword(userID uint, password string) (err error)
}

type GroupData interface {
	CreateGroup(request GroupRequest) (response GroupResponse, err error)
	UpdateGroup(request UpdateGroupRequest) (err error)
	DeleteGroup(id uint) (err error)
	GetUsersByGroupID(groupID uint, offset uint, limit uint) (response UsersResponse, err error)
	GetGroups(offset uint, limit uint) (response GroupsResponse, err error)
	AddUser(userID uint, groupID uint) (err error)
	RemoveUser(groupID uint, userID uint) (err error)
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersResponse struct {
	Users   []UserResponse `json:"users"`
	Total   uint           `json:"total"`
	Page    uint           `json:"page"`
	PerPage uint           `json:"perPage"`
}

type GroupResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GroupsResponse struct {
	Groups  []GroupResponse `json:"groups"`
	Total   uint            `json:"total"`
	Page    uint            `json:"page"`
	PerPage uint            `json:"perPage"`
}

type GroupRequest struct {
	Name string `json:"name"`
}

type UpdateGroupRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
