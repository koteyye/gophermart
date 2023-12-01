package service_test

import (
	"context"
	"database/sql"
	"flag"

	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/deployments/gophermart/migrations"
	"github.com/sergeizaitcev/gophermart/pkg/postgres"
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
	db *sql.DB
}

func (suite *CommonSuite) SetupSuite() {
	var err error
	suite.db, err = postgres.Connect(flagDatabaseURI)
	suite.Require().NoError(err)
	suite.Require().NoError(migrations.Up(context.Background(), suite.db))
}

func (suite *CommonSuite) TearDownSuite() {
	suite.Require().NoError(migrations.Down(context.Background(), suite.db))
	suite.Require().NoError(suite.db.Close())
}
