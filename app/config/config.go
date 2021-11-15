package config

import (
	"fmt"
	"usermanagement/app/internal/data"
	"usermanagement/app/internal/service"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

type Postgres struct {
	Host     string
	Port     int
	DBname   string
	Username string
	Password string
	Sslmode  string
}

type Config struct {
	Postgres Postgres
	Port     int
}

func initializeServices(appConfig *AppConfiguration) {
	db, err := NewDatabase(appConfig.config.Postgres)
	if err != nil {
		log.WithField("err", err).Fatal("intialising DB")
	}
	userData := data.NewUserService(db)
	appConfig.userService = service.NewUserService(userData)

	groupData := data.NewGroupService(db)
	appConfig.groupService = service.NewGroupService(groupData)
}

func NewDatabase(config Postgres) (db *gorm.DB, err error) {
	args := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.DBname, config.Password, config.Sslmode)
	if db, err = gorm.Open("postgres", args); err != nil {
		return db, err
	}
	return db, err
}
