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

var testOrder string = "1234567890"
var testOrder2 string = "1234567891"

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
		want          string
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
			wantErr:  models.ErrDuplicate,
		},
		{
			name:     "duplicate other user",
			testUser: testUser2ID,
			want:     testOrder,
			wantErr:  models.ErrDuplicateOtherUser,
		},
	}

	for _, test := range testCases {
		got, err := gophermart.CreateOrder(context.Background(), testOrder, test.testUser)
		if test.wantErr != nil {
			assert.ErrorIs(t, err, test.wantErr)
		}
		assert.NotNil(t, got)
	}
}

func TestUpdateOrder(t *testing.T) {
	testUpdateData := &UpdateOrder{Order: testOrder, Status: 2, Accrual: 100}

	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	_, err = gophermart.CreateOrder(context.Background(), testOrder, testUser1ID)

	err = gophermart.UpdateOrder(context.Background(), testUpdateData)
	assert.NoError(t, err)
}

func TestUpdateOrderStatus(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	_, err = gophermart.CreateOrder(context.Background(), testOrder, testUser1ID)

	err = gophermart.UpdateOrderStatus(context.Background(), testOrder, 3)
	assert.NoError(t, err)
}

func TestGetOrderByNumber(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	_, err = gophermart.CreateOrder(context.Background(), testOrder, testUser1ID)

	order, err := gophermart.GetOrderByNumber(context.Background(), testOrder)
	assert.NoError(t, err)
	assert.NotNil(t, order)
}

func TestGetOrdersByUser(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	_, err = gophermart.CreateOrder(context.Background(), testOrder, testUser1ID)
	_, err = gophermart.CreateOrder(context.Background(), testOrder2, testUser1ID)
	assert.NoError(t, err)

	orders, err := gophermart.GetOrdersByUser(context.Background(), testUser1ID)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestDeleteOrderByNumber(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)

	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	ctx := context.WithValue(context.Background(), models.KeyUserID, testUser1ID)
	_, err = gophermart.CreateOrder(ctx, testOrder, testUser1ID)

	err = gophermart.DeleteOrderByNumber(context.Background(), testOrder)
	assert.NoError(t, err)
}

// TestBalanceOperation - комплексный тест с балансовыми операциями
func TestBalanceOperation(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)
	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	_, err = gophermart.CreateOrder(context.Background(), testOrder, testUser1ID)
	_, err = gophermart.CreateOrder(context.Background(), testOrder2, testUser1ID)
	assert.NoError(t, err)

	//Создание балансовой операции
	balanceOperationID, err := gophermart.CreateBalanceOperation(context.Background(), 100, testOrder, testUser1ID)
	assert.NotEqual(t, uuid.Nil, balanceOperationID)
	assert.NoError(t, err)

	//Обновление состояния балансовой операции
	err = gophermart.UpdateBalanceOperation(context.Background(), testOrder, 1)
	assert.NoError(t, err)

	//Получение балансовой операции по номеру заказа
	balanceOperation, err := gophermart.GetBalanceOperationByOrder(context.Background(), testOrder)
	assert.NotNil(t, balanceOperation)
	assert.NoError(t, err)

	//Создаем вторую балансовую операцию
	_, err = gophermart.CreateBalanceOperation(context.Background(), 200, testOrder2, testUser1ID)
	assert.NoError(t, err)

	//Получаем список балансовых операций по текущему пользователю
	balanceOperations, err := gophermart.GetBalanceOperationByUser(context.Background(), testUser1ID)
	assert.NoError(t, err)
	assert.Len(t, balanceOperations, 2)

	//Удаляем балансовую операцию
	err = gophermart.DeleteBalanceOperationByOrder(context.Background(), testOrder)
	assert.NoError(t, err)

	//Еще раз проверяем список балансовых операций по текущему пользователю
	postBalanceOperations, err := gophermart.GetBalanceOperationByUser(context.Background(), testUser1ID)
	assert.NoError(t, err)
	assert.Len(t, postBalanceOperations, 1)
}

// TestBalance - комплексный тест с балансом
func TestBalance(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	gophermart := NewGophermartPostgres(db)
	auth := NewAuthPostgres(db)
	testUser1ID, err := auth.CreateUser(context.Background(), testUser1.login, testUser1.password)
	ctx := context.WithValue(context.Background(), models.KeyUserID, testUser1ID)
	assert.NoError(t, err)

	//Проверяем баланс текущего пользователя
	balance, err := gophermart.GetBalanceByUserID(context.Background(), testUser1ID)
	assert.NoError(t, err)
	assert.NotNil(t, balance)

	//Инкрементем баланс текущего пользователя
	err = gophermart.IncrementBalance(ctx, 200)
	assert.NoError(t, err)

	postBalance, err := gophermart.GetBalanceByUserID(context.Background(), testUser1ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(200), postBalance.CurrentBalance)

	//Декрементим баланс текущего пользователя
	err = gophermart.DecrementBalance(ctx, 100)
	assert.NoError(t, err)

	post2Balance, err := gophermart.GetBalanceByUserID(context.Background(), testUser1ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), post2Balance.CurrentBalance)

	//Проверяем ошибку декрмента если баланс будет <0
	err = gophermart.DecrementBalance(ctx, 300)
	assert.ErrorIs(t, err, models.ErrBalanceBelowZero)
}
