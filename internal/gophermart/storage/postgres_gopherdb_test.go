package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/models"
)

// Переменные для создания тестовых пользователей
var (
	testUser1 = tUser{login: "testuser1", password: "testpassword"}
	testUser2 = tUser{login: "testuser2", password: "testpassword"}
)

// Переменные для ID тестовых пользователей
var (
	testUser1ID uuid.UUID
	testUser2ID uuid.UUID
)

var testOrder int64 = 1234567890

func TestCreateOrder(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	// Создаем тестовых пользователей
	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	assert.NoError(t, err)
	testUser2ID, err := auth.CreateUser(context.Background(), testUser2.login, testUser2.password)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		testUser      uuid.UUID
		otherTestUser uuid.UUID
		want          int64
		wantErr       error
	}{
		{
			name:     "success",
			testUser: testUser1ID,
			want:     testOrder,
			wantErr:  nil,
		},
		{
			name:     "duplicate current user",
			testUser: testUser1ID,
			want:     testOrder,
			wantErr:  ErrDuplicate,
		},
		{
			name:     "duplicate other user",
			testUser: testUser2ID,
			want:     testOrder,
			wantErr:  ErrDuplicateOtherUser,
		},
	}

	for _, test := range testCases {
		ctx := context.WithValue(context.Background(), models.KeyUserID, test.testUser)
		got, err := gophermart.CreateOrder(ctx, testOrder)
		if test.wantErr != nil {
			assert.ErrorIs(t, err, test.wantErr)
		}
		assert.NotNil(t, got)
	}
}

func TestUpdateOrder(t *testing.T) {
	testUpdateData := &UpdateOrder{Number: testOrder, Status: 2, Accrual: 100}

	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	ctx := context.WithValue(context.Background(), models.KeyUserID, testUser1ID)
	_, err = gophermart.CreateOrder(ctx, testOrder)

	err = gophermart.UpdateOrder(context.Background(), testUpdateData)
	assert.NoError(t, err)
}

func TestUpdateOrderStatus(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	ctx := context.WithValue(context.Background(), models.KeyUserID, testUser1ID)
	_, err = gophermart.CreateOrder(ctx, testOrder)

	err = gophermart.UpdateOrderStatus(context.Background(), testOrder, 3)
	assert.NoError(t, err)
}

func TestGetOrderByNumber(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	ctx := context.WithValue(context.Background(), models.KeyUserID, testUser1ID)
	_, err = gophermart.CreateOrder(ctx, testOrder)

	order, err := gophermart.GetOrderByNumber(context.Background(), testOrder)
	assert.NoError(t, err)
	assert.NotNil(t, order)
}

func TestDeleteOrderByNumber(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	ctx := context.WithValue(context.Background(), models.KeyUserID, testUser1ID)
	_, err = gophermart.CreateOrder(ctx, testOrder)

	err = gophermart.DeleteOrderByNumber(context.Background(), testOrder)
	assert.NoError(t, err)
}
