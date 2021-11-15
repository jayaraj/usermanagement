package serviceerror

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ErrorCode string

type ServiceError struct {
	Code ErrorCode
	Err  error
}

func (s *ServiceError) Error() string {
	return string(s.Code) + " : " + s.Err.Error()
}

func NewServiceError(code ErrorCode, err error) error {
	return &ServiceError{
		Code: code,
		Err:  err,
	}
}

func AbortOnError(c *gin.Context, err error) {
	var srvError *ServiceError
	if ok := errors.As(err, &srvError); ok {
		log.WithError(err).Error("service error")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	log.WithError(err).Error("unknown error")
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
}
