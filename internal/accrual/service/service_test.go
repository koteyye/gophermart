package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
	"github.com/sergeizaitcev/gophermart/internal/accrual/service"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
	mock_storage "github.com/sergeizaitcev/gophermart/internal/accrual/storage/mocks"
)

var tOrderNum = "49927398716"

var tOrder = &storage.OrderOut{
	OrderNumber: tOrderNum,
	Status:      "PROCESSED",
	Accrual:     10000,
}

var exOrder = &models.OrderOut{
	Number:  tOrderNum,
	Status:  "PROCESSED",
	Accrual: 10000,
}

var tMatchName = "testItem"

var tMatch = &storage.MatchOut{
	MatchID:   uuid.New(),
	MatchName: tMatchName,
	Reward:    1000,
	Type:      "PERCENT",
}

var tMatchStorage = &storage.Match{
	MatchName: tMatchName,
	Reward:    1000,
	Type:      0,
}

var tMatchInput = &models.Match{
	MatchName:  tMatchName,
	Reward:     1000,
	RewardType: "%",
}

func TestGetOrder(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockStorage := mock_storage.NewMockStorage(ctrl)

	t.Run("getOrder", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				GetOrderByNumber(gomock.Any(), tOrderNum).
				Return(tOrder, (error)(nil))

			order, err := srv.GetOrder(ctx, tOrderNum)
			assert.NoError(t, err)
			assert.Equal(t, exOrder, order)
		})

		t.Run("empty", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				GetOrderByNumber(gomock.Any(), tOrderNum).
				Return(&storage.OrderOut{}, storage.ErrNotFound)

			order, err := srv.GetOrder(ctx, tOrderNum)
			assert.Error(t, err)
			assert.Equal(t, &models.OrderOut{}, order)
		})
	})
}

func TestCheckMatch(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockStorage := mock_storage.NewMockStorage(ctrl)

	t.Run("checkMatch", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				GetMatchByName(gomock.Any(), tMatchName).
				Return(&storage.MatchOut{}, storage.ErrNotFound)

			err := srv.CheckMatch(ctx, tMatchName)
			assert.ErrorIs(t, err, storage.ErrNotFound)
		})

		t.Run("noEmpty", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				GetMatchByName(gomock.Any(), tMatchName).
				Return(tMatch, (error)(nil))

			err := srv.CheckMatch(ctx, tMatchName)
			assert.NoError(t, err)
		})
	})
}

func TestCreateMatch(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockStorage := mock_storage.NewMockStorage(ctrl)

	t.Run("createMatch", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				CreateMatch(gomock.Any(), tMatchStorage).
				Return(uuid.New(), (error)(nil))

			err := srv.CreateMatch(ctx, tMatchInput)
			assert.NoError(t, err)
		})

		t.Run("duplicate", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				CreateMatch(gomock.Any(), tMatchStorage).
				Return(uuid.Nil, storage.ErrDuplicate)

			err := srv.CreateMatch(ctx, tMatchInput)
			assert.ErrorIs(t, err, storage.ErrDuplicate)
		})
	})
}

func TestCheckOrder(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockStorage := mock_storage.NewMockStorage(ctrl)

	t.Run("getOrder", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				GetOrderByNumber(gomock.Any(), tOrderNum).
				Return(&storage.OrderOut{}, storage.ErrNotFound)

			err := srv.CheckOrder(ctx, tOrderNum)
			assert.NoError(t, err)
		})

		t.Run("noempty", func(t *testing.T) {
			srv := service.NewService(mockStorage)

			mockStorage.EXPECT().
				GetOrderByNumber(gomock.Any(), tOrderNum).
				Return(tOrder, (error)(nil))

			err := srv.CheckOrder(ctx, tOrderNum)
			assert.ErrorIs(t, err, service.ErrOrderRegistered)
		})
	})
}
