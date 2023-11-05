package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testUser = struct {
	login    string
	password string
}{
	login:    "testuser",
	password: "testpassword",
}

func TestStorage_CreateUser(t *testing.T) {
	userID, err := testAuthDB.CreateUser(context.Background(), testUser.login, testUser.password)

	assert.NoError(t, err)
	assert.NotNil(t, userID)
}

func TestStorage_GetUser(t *testing.T) {
	userID, err := testAuthDB.GetUser(context.Background(), testUser.login, testUser.password)

	assert.NoError(t, err)
	assert.NotNil(t, userID)
}
