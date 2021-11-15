package docs

import (
	"usermanagement/app/internal"
	"usermanagement/app/internal/httpservice"
)

// swagger:route POST /users users createUserRequest
// Create new user.
// responses:
//   201: createUserResponse
//   400: serviceError
//   500: serviceError

// swagger:route PUT /users/{id} users updateUserRequest
// Update user.
// responses:
//   200:
//   400: serviceError
//   500: serviceError

// swagger:route DELETE /users/{id} users deleteUserRequest
// Delete a user.
// responses:
//   200:
//   400: serviceError
//   500: serviceError

// swagger:route GET /users users getUsersRequest
// Get users.
// responses:
//   200: getUsersResponse
//   400: serviceError
//   500: serviceError

// swagger:route PUT /users/{id}/password users changePwdRequest
// Change password of user.
// responses:
//   200:
//   400: serviceError
//   500: serviceError

// swagger:response createUserResponse
type createUserResponse struct {
	// in:body
	Body internal.UserResponse
}

// swagger:response getUsersResponse
type getUsersResponse struct {
	// in:body
	Body internal.UsersResponse
}

// swagger:response serviceError
type serviceErrorResponse struct {
	// in:body
	Body struct {
		Message string `json:"message"`
	}
}

// swagger:parameters createUserRequest
type createUserRequest struct {
	// in:body
	Body httpservice.CreateUser
}

// swagger:parameters updateUserRequest
type updateUserRequest struct {
	// in: path
	Id uint `json:"id"`
	// in:body
	Body httpservice.UpdateUser
}

// swagger:parameters deleteUserRequest
type deleteUserRequest struct {
	// in: path
	Id uint `json:"id"`
}

// swagger:parameters getUsersRequest
type getUsersRequest struct {
	// in: query
	Page uint `json:"page"`
	// in: query
	PerPage uint `json:"perPage"`
}

// swagger:parameters changePwdRequest
type changePwdRequest struct {
	// in: path
	Id uint `json:"id"`
	// in:body
	Body httpservice.CreatePassword
}
