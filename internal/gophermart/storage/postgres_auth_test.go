package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type tUser struct {
	login string
	password string
}

var testUser = tUser{
	login:    "testuser",
	password: "testpassword",
}

func TestAuthStorage(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	ctx := context.Background()

	auth := NewAuthPostgres(db)
	want, err := auth.CreateUser(ctx, testUser.login, testUser.password)

	require.NoError(t, err)
	require.NotNil(t, want)

	got, err := auth.GetUser(ctx, testUser.login, testUser.password)

	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, want, got)
}
