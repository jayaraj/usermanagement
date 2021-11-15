package serviceerror

const (
	InvalidUserRequest      ErrorCode = "Invalid User Request"
	UserNotFound            ErrorCode = "UserNotFound"
	InvalidGroupRequest     ErrorCode = "Invalid Group Request"
	InvalidUserGroupRequest ErrorCode = "Invalid User Group Request"
	DuplicateUser           ErrorCode = "Duplicate User"
	DuplicateGroup          ErrorCode = "Duplicate Group"
)
