package postgres_test

import (
	"context"
	"flag"

	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage/postgres"
)

var flagDatabaseURI string

func init() {
	flag.StringVar(
		&flagDatabaseURI,
		"database-dsn",
		"postgresql://postgres:postgres@localhost:5433/gophermart?sslmode=disable",
		"connection string to database",
	)
}

type CommonSuite struct {
	suite.Suite
	storage *postgres.Storage
}

func (suite *CommonSuite) SetupSuite() {
	var err error
	suite.storage, err = postgres.Connect(&config.Config{DatabaseURI: flagDatabaseURI})
	suite.Require().NoError(err)
	suite.Require().NoError(suite.storage.Up(context.Background()))
}

func (suite *CommonSuite) TearDownSuite() {
	suite.Require().NoError(suite.storage.Down(context.Background()))
	suite.Require().NoError(suite.storage.Close())
}
