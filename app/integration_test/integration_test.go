package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"usermanagement/app/config"

	"github.com/docker/go-connections/nat"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type IntegrationTestSuite struct {
	suite.Suite
	config   config.Config
	postgres testcontainers.Container
	testDB   *gorm.DB
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	suite.config = config.Config{
		Postgres: config.Postgres{
			DBname:   "postgres",
			Username: "postgres",
			Password: "password",
			Sslmode:  "disable",
		},
		Port: 8080,
	}
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", port.Port(), "postgres")
	}
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForSQL(nat.Port("5432/tcp"), "postgres", dbURL).Timeout(time.Second * 5),
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(suite.T(), err)

	host, err := postgres.Host(ctx)
	assert.NoError(suite.T(), err)

	port, err := postgres.MappedPort(ctx, "5432")
	assert.NoError(suite.T(), err)

	suite.config.Postgres.Host = host
	suite.config.Postgres.Port = port.Int()
	suite.postgres = postgres

	db, err := config.NewDatabase(suite.config.Postgres)
	assert.NoError(suite.T(), err)
	suite.testDB = db
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.postgres.Terminate(context.Background())
}
