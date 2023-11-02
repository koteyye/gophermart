package storage_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_storage "github.com/sergeizaitcev/gophermart/internal/gophermart/storage/mocks"
)

func TestStorage_CreateUser(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	s := mock_storage.NewMockAuth(c)

	userID := uuid.New()
	s.EXPECT().CreateUser(gomock.Any(), gomock.Eq("simpleUser"), gomock.Eq("123456")).Return(userID, nil).AnyTimes()
}