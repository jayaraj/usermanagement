package docs

import (
	"usermanagement/app/internal"
	"usermanagement/app/internal/httpservice"
)

// swagger:route POST /groups groups createGroupRequest
// Create new group.
// responses:
//   201: createGroupResponse
//   400: serviceError
//   500: serviceError

// swagger:route PUT /groups/{id} groups updateGroupRequest
// Update  a group.
// responses:
//   200:
//   400: serviceError
//   500: serviceError

// swagger:route DELETE /groups/{id} groups deleteGroupRequest
// Delete a group.
// responses:
//   200:
//   400: serviceError
//   500: serviceError

// swagger:route GET /groups groups getGroupsRequest
// Get groups.
// responses:
//   200: getGroupsResponse
//   400: serviceError
//   500: serviceError

// swagger:route GET /groups/{id}/users groups getGroupUsersRequest
// Get group users.
// responses:
//   200: getUsersResponse
//   400: serviceError
//   500: serviceError

// swagger:route POST /groups/{id}/users groups addUserRequest
// Add user to a group.
// responses:
//   200:
//   400: serviceError
//   500: serviceError

// swagger:route DELETE /groups/{id}/users/{userid} groups removeUserRequest
// Remove user from a group.
// responses:
//   200:
//   400: serviceError
//   500: serviceError

// swagger:response createGroupResponse
type createGroupResponse struct {
	// in:body
	Body internal.GroupResponse
}

// swagger:response getGroupsResponse
type getGroupsResponse struct {
	// in:body
	Body internal.GroupsResponse
}

// swagger:parameters createGroupRequest
type createGroupRequest struct {
	// in:body
	Body httpservice.CreateGroup
}

// swagger:parameters updateGroupRequest
type updateGroupRequest struct {
	// in: path
	Id uint `json:"id"`
	// in:body
	Body httpservice.UpdateGroup
}

// swagger:parameters deleteGroupRequest
type deleteGroupRequest struct {
	// in: path
	Id uint `json:"id"`
}

// swagger:parameters getGroupUsersRequest
type getGroupUsersRequest struct {
	// in: path
	Id uint `json:"id"`
	// in: query
	Page uint `json:"page"`
	// in: query
	PerPage uint `json:"perPage"`
}

// swagger:parameters getGroupsRequest
type getGroupsRequest struct {
	// in: query
	Page uint `json:"page"`
	// in: query
	PerPage uint `json:"perPage"`
}

// swagger:parameters addUserRequest
type addUserRequest struct {
	// in: path
	Id uint `json:"id"`
	// in:body
	Body httpservice.AddUser
}

// swagger:parameters removeUserRequest
type removeUserRequest struct {
	// in: path
	Id uint `json:"id"`
	// in: path
	UserID uint `json:"userid"`
}
