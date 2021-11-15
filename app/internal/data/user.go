package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
	"usermanagement/app/internal"
	"usermanagement/app/internal/serviceerror"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
)

type userDataService struct {
	db *gorm.DB
}

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name      string
	Password  string
	Salt      string
	Email     string `sql:"index"`
}

func NewUserService(db *gorm.DB) *userDataService {
	db.AutoMigrate(&User{})
	return &userDataService{
		db: db,
	}
}

func (u *userDataService) CreateUser(request internal.UserRequest) (response internal.UserResponse, err error) {
	if request.Email == "" || request.Name == "" || request.Password == "" {
		return response, serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("missing create user fields"))
	}
	var count int64
	err = u.db.Model(&User{}).Where("email = ?", request.Email).Count(&count).Error
	if err != nil {
		return response, errors.Wrap(err, "get user with email count failed")
	}
	if count > 0 {
		return response, serviceerror.NewServiceError(serviceerror.DuplicateUser, fmt.Errorf("user with email %s is present", request.Email))
	}

	salt := u.randomString()
	user := User{
		Name:     request.Name,
		Email:    request.Email,
		Password: u.encodePassword(request.Password, salt),
		Salt:     salt,
	}
	err = u.db.Create(&user).Error
	if err != nil {
		return response, errors.Wrap(err, "create user failed")
	}
	response = internal.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return response, err
}

func (u *userDataService) UpdateUser(request internal.UpdateUserRequest) (err error) {
	if request.ID == 0 || (request.Email == "" && request.Name == "") {
		return serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("missing update user fields"))
	}
	user := User{
		ID: request.ID,
	}
	update := make(map[string]interface{})
	if request.Email != "" {
		var count int64
		err = u.db.Model(&User{}).Where("email = ?", request.Email).Count(&count).Error
		if err != nil {
			return errors.Wrap(err, "get user with email count failed")
		}
		if count > 0 {
			return serviceerror.NewServiceError(serviceerror.DuplicateUser, fmt.Errorf("user with email %s is present", request.Email))
		}
		update["email"] = request.Email
	}
	if request.Name != "" {
		update["name"] = request.Name
	}
	err = u.db.Model(&user).Updates(update).Error
	if err != nil {
		return errors.Wrap(err, "update user failed")
	}
	return err
}

func (u *userDataService) ChangePassword(userID uint, password string) (err error) {
	if userID == 0 {
		return serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("user_id=0 for change password"))
	}
	user := User{
		ID: userID,
	}
	err = u.db.First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return serviceerror.NewServiceError(serviceerror.UserNotFound, fmt.Errorf("user %d not found", userID))
	}
	if err != nil {
		return errors.Wrap(err, "get user failed")
	}
	err = u.db.Model(&user).Updates(User{Password: u.encodePassword(password, user.Salt)}).Error
	if err != nil {
		return errors.Wrap(err, "update user failed")
	}
	return err
}

func (u *userDataService) DeleteUser(id uint) (err error) {
	if id == 0 {
		return serviceerror.NewServiceError(serviceerror.InvalidUserRequest, errors.New("user_id is 0 for delete user"))
	}
	err = u.db.Where("user_id = ?", id).Delete(&UserGroup{}).Error
	if err != nil {
		return errors.Wrap(err, "delete user group failed")
	}
	err = u.db.Delete(&User{}, id).Error
	if err != nil {
		return errors.Wrap(err, "delete user failed")
	}
	return err
}

func (u *userDataService) GetUsers(offset uint, limit uint) (response internal.UsersResponse, err error) {
	if limit == 0 || limit > 1000 {
		return response, serviceerror.NewServiceError(serviceerror.InvalidUserRequest, fmt.Errorf("limit %d is not valid for get users", limit))
	}
	var users []User
	err = u.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return response, errors.Wrap(err, "get users failed")
	}
	var count int64
	err = u.db.Model(&User{}).Count(&count).Error
	if err != nil {
		return response, errors.Wrap(err, "get users count failed")
	}
	userResponse := make([]internal.UserResponse, len(users))
	for i, u := range users {
		userResponse[i] = internal.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}
	response = internal.UsersResponse{
		Total: uint(count),
		Users: userResponse,
	}
	return response, err
}

func (u *userDataService) randomString() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (u *userDataService) encodePassword(password string, salt string) string {
	newPasswd := pbkdf2.Key([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(newPasswd)
}
