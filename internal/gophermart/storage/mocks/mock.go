// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	storage "github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuth) CreateUser(ctx context.Context, login, password string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, login, password)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthMockRecorder) CreateUser(ctx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuth)(nil).CreateUser), ctx, login, password)
}

// GetUser mocks base method.
func (m *MockAuth) GetUser(ctx context.Context, login, passwrod string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, login, passwrod)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAuthMockRecorder) GetUser(ctx, login, passwrod interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuth)(nil).GetUser), ctx, login, passwrod)
}

// GetUserByID mocks base method.
func (m *MockAuth) GetUserByID(ctx context.Context, userID uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockAuthMockRecorder) GetUserByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockAuth)(nil).GetUserByID), ctx, userID)
}

// MockGophermartDB is a mock of GophermartDB interface.
type MockGophermartDB struct {
	ctrl     *gomock.Controller
	recorder *MockGophermartDBMockRecorder
}

// MockGophermartDBMockRecorder is the mock recorder for MockGophermartDB.
type MockGophermartDBMockRecorder struct {
	mock *MockGophermartDB
}

// NewMockGophermartDB creates a new mock instance.
func NewMockGophermartDB(ctrl *gomock.Controller) *MockGophermartDB {
	mock := &MockGophermartDB{ctrl: ctrl}
	mock.recorder = &MockGophermartDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGophermartDB) EXPECT() *MockGophermartDBMockRecorder {
	return m.recorder
}

// CreateBalanceOperation mocks base method.
func (m *MockGophermartDB) CreateBalanceOperation(ctx context.Context, operation int64, order string, userID uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBalanceOperation", ctx, operation, order, userID)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBalanceOperation indicates an expected call of CreateBalanceOperation.
func (mr *MockGophermartDBMockRecorder) CreateBalanceOperation(ctx, operation, order, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBalanceOperation", reflect.TypeOf((*MockGophermartDB)(nil).CreateBalanceOperation), ctx, operation, order, userID)
}

// CreateOrder mocks base method.
func (m *MockGophermartDB) CreateOrder(ctx context.Context, order string, userID uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, order, userID)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockGophermartDBMockRecorder) CreateOrder(ctx, order, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockGophermartDB)(nil).CreateOrder), ctx, order, userID)
}

// DecrementBalance mocks base method.
func (m *MockGophermartDB) DecrementBalance(ctx context.Context, decrementSum int64, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecrementBalance", ctx, decrementSum, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecrementBalance indicates an expected call of DecrementBalance.
func (mr *MockGophermartDBMockRecorder) DecrementBalance(ctx, decrementSum, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecrementBalance", reflect.TypeOf((*MockGophermartDB)(nil).DecrementBalance), ctx, decrementSum, userID)
}

// DeleteBalanceOperationByOrder mocks base method.
func (m *MockGophermartDB) DeleteBalanceOperationByOrder(ctx context.Context, order string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBalanceOperationByOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBalanceOperationByOrder indicates an expected call of DeleteBalanceOperationByOrder.
func (mr *MockGophermartDBMockRecorder) DeleteBalanceOperationByOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBalanceOperationByOrder", reflect.TypeOf((*MockGophermartDB)(nil).DeleteBalanceOperationByOrder), ctx, order)
}

// DeleteOrderByNumber mocks base method.
func (m *MockGophermartDB) DeleteOrderByNumber(ctx context.Context, order string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrderByNumber", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrderByNumber indicates an expected call of DeleteOrderByNumber.
func (mr *MockGophermartDBMockRecorder) DeleteOrderByNumber(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrderByNumber", reflect.TypeOf((*MockGophermartDB)(nil).DeleteOrderByNumber), ctx, order)
}

// GetBalanceByUserID mocks base method.
func (m *MockGophermartDB) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*storage.BalanceItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalanceByUserID", ctx, userID)
	ret0, _ := ret[0].(*storage.BalanceItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalanceByUserID indicates an expected call of GetBalanceByUserID.
func (mr *MockGophermartDBMockRecorder) GetBalanceByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalanceByUserID", reflect.TypeOf((*MockGophermartDB)(nil).GetBalanceByUserID), ctx, userID)
}

// GetBalanceOperationByOrder mocks base method.
func (m *MockGophermartDB) GetBalanceOperationByOrder(ctx context.Context, order string) (*storage.BalanceOperationItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalanceOperationByOrder", ctx, order)
	ret0, _ := ret[0].(*storage.BalanceOperationItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalanceOperationByOrder indicates an expected call of GetBalanceOperationByOrder.
func (mr *MockGophermartDBMockRecorder) GetBalanceOperationByOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalanceOperationByOrder", reflect.TypeOf((*MockGophermartDB)(nil).GetBalanceOperationByOrder), ctx, order)
}

// GetBalanceOperationByUser mocks base method.
func (m *MockGophermartDB) GetBalanceOperationByUser(ctx context.Context, userID uuid.UUID) ([]*storage.BalanceOperationItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalanceOperationByUser", ctx, userID)
	ret0, _ := ret[0].([]*storage.BalanceOperationItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalanceOperationByUser indicates an expected call of GetBalanceOperationByUser.
func (mr *MockGophermartDBMockRecorder) GetBalanceOperationByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalanceOperationByUser", reflect.TypeOf((*MockGophermartDB)(nil).GetBalanceOperationByUser), ctx, userID)
}

// GetOrderByNumber mocks base method.
func (m *MockGophermartDB) GetOrderByNumber(ctx context.Context, order string) (*storage.OrderItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByNumber", ctx, order)
	ret0, _ := ret[0].(*storage.OrderItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByNumber indicates an expected call of GetOrderByNumber.
func (mr *MockGophermartDBMockRecorder) GetOrderByNumber(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByNumber", reflect.TypeOf((*MockGophermartDB)(nil).GetOrderByNumber), ctx, order)
}

// GetOrdersByUser mocks base method.
func (m *MockGophermartDB) GetOrdersByUser(ctx context.Context, userID uuid.UUID) ([]*storage.OrderItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUser", ctx, userID)
	ret0, _ := ret[0].([]*storage.OrderItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUser indicates an expected call of GetOrdersByUser.
func (mr *MockGophermartDBMockRecorder) GetOrdersByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUser", reflect.TypeOf((*MockGophermartDB)(nil).GetOrdersByUser), ctx, userID)
}

// IncrementBalance mocks base method.
func (m *MockGophermartDB) IncrementBalance(ctx context.Context, incrementSum int64, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementBalance", ctx, incrementSum, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementBalance indicates an expected call of IncrementBalance.
func (mr *MockGophermartDBMockRecorder) IncrementBalance(ctx, incrementSum, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementBalance", reflect.TypeOf((*MockGophermartDB)(nil).IncrementBalance), ctx, incrementSum, userID)
}

// UpdateBalanceOperation mocks base method.
func (m *MockGophermartDB) UpdateBalanceOperation(ctx context.Context, order string, operationState storage.BalanceOperationState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalanceOperation", ctx, order, operationState)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBalanceOperation indicates an expected call of UpdateBalanceOperation.
func (mr *MockGophermartDBMockRecorder) UpdateBalanceOperation(ctx, order, operationState interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalanceOperation", reflect.TypeOf((*MockGophermartDB)(nil).UpdateBalanceOperation), ctx, order, operationState)
}

// UpdateOrder mocks base method.
func (m *MockGophermartDB) UpdateOrder(ctx context.Context, order *storage.UpdateOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockGophermartDBMockRecorder) UpdateOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockGophermartDB)(nil).UpdateOrder), ctx, order)
}

// UpdateOrderStatus mocks base method.
func (m *MockGophermartDB) UpdateOrderStatus(ctx context.Context, order string, orderStatus storage.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", ctx, order, orderStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockGophermartDBMockRecorder) UpdateOrderStatus(ctx, order, orderStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockGophermartDB)(nil).UpdateOrderStatus), ctx, order, orderStatus)
}

// MockOrders is a mock of Orders interface.
type MockOrders struct {
	ctrl     *gomock.Controller
	recorder *MockOrdersMockRecorder
}

// MockOrdersMockRecorder is the mock recorder for MockOrders.
type MockOrdersMockRecorder struct {
	mock *MockOrders
}

// NewMockOrders creates a new mock instance.
func NewMockOrders(ctrl *gomock.Controller) *MockOrders {
	mock := &MockOrders{ctrl: ctrl}
	mock.recorder = &MockOrdersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrders) EXPECT() *MockOrdersMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrders) CreateOrder(ctx context.Context, order string, userID uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, order, userID)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrdersMockRecorder) CreateOrder(ctx, order, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrders)(nil).CreateOrder), ctx, order, userID)
}

// DeleteOrderByNumber mocks base method.
func (m *MockOrders) DeleteOrderByNumber(ctx context.Context, order string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrderByNumber", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrderByNumber indicates an expected call of DeleteOrderByNumber.
func (mr *MockOrdersMockRecorder) DeleteOrderByNumber(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrderByNumber", reflect.TypeOf((*MockOrders)(nil).DeleteOrderByNumber), ctx, order)
}

// GetOrderByNumber mocks base method.
func (m *MockOrders) GetOrderByNumber(ctx context.Context, order string) (*storage.OrderItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByNumber", ctx, order)
	ret0, _ := ret[0].(*storage.OrderItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByNumber indicates an expected call of GetOrderByNumber.
func (mr *MockOrdersMockRecorder) GetOrderByNumber(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByNumber", reflect.TypeOf((*MockOrders)(nil).GetOrderByNumber), ctx, order)
}

// GetOrdersByUser mocks base method.
func (m *MockOrders) GetOrdersByUser(ctx context.Context, userID uuid.UUID) ([]*storage.OrderItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUser", ctx, userID)
	ret0, _ := ret[0].([]*storage.OrderItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUser indicates an expected call of GetOrdersByUser.
func (mr *MockOrdersMockRecorder) GetOrdersByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUser", reflect.TypeOf((*MockOrders)(nil).GetOrdersByUser), ctx, userID)
}

// UpdateOrder mocks base method.
func (m *MockOrders) UpdateOrder(ctx context.Context, order *storage.UpdateOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockOrdersMockRecorder) UpdateOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockOrders)(nil).UpdateOrder), ctx, order)
}

// UpdateOrderStatus mocks base method.
func (m *MockOrders) UpdateOrderStatus(ctx context.Context, order string, orderStatus storage.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", ctx, order, orderStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockOrdersMockRecorder) UpdateOrderStatus(ctx, order, orderStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockOrders)(nil).UpdateOrderStatus), ctx, order, orderStatus)
}

// MockBalance is a mock of Balance interface.
type MockBalance struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceMockRecorder
}

// MockBalanceMockRecorder is the mock recorder for MockBalance.
type MockBalanceMockRecorder struct {
	mock *MockBalance
}

// NewMockBalance creates a new mock instance.
func NewMockBalance(ctrl *gomock.Controller) *MockBalance {
	mock := &MockBalance{ctrl: ctrl}
	mock.recorder = &MockBalanceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalance) EXPECT() *MockBalanceMockRecorder {
	return m.recorder
}

// CreateBalanceOperation mocks base method.
func (m *MockBalance) CreateBalanceOperation(ctx context.Context, operation int64, order string, userID uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBalanceOperation", ctx, operation, order, userID)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBalanceOperation indicates an expected call of CreateBalanceOperation.
func (mr *MockBalanceMockRecorder) CreateBalanceOperation(ctx, operation, order, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBalanceOperation", reflect.TypeOf((*MockBalance)(nil).CreateBalanceOperation), ctx, operation, order, userID)
}

// DecrementBalance mocks base method.
func (m *MockBalance) DecrementBalance(ctx context.Context, decrementSum int64, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecrementBalance", ctx, decrementSum, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecrementBalance indicates an expected call of DecrementBalance.
func (mr *MockBalanceMockRecorder) DecrementBalance(ctx, decrementSum, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecrementBalance", reflect.TypeOf((*MockBalance)(nil).DecrementBalance), ctx, decrementSum, userID)
}

// DeleteBalanceOperationByOrder mocks base method.
func (m *MockBalance) DeleteBalanceOperationByOrder(ctx context.Context, order string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBalanceOperationByOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBalanceOperationByOrder indicates an expected call of DeleteBalanceOperationByOrder.
func (mr *MockBalanceMockRecorder) DeleteBalanceOperationByOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBalanceOperationByOrder", reflect.TypeOf((*MockBalance)(nil).DeleteBalanceOperationByOrder), ctx, order)
}

// GetBalanceByUserID mocks base method.
func (m *MockBalance) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*storage.BalanceItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalanceByUserID", ctx, userID)
	ret0, _ := ret[0].(*storage.BalanceItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalanceByUserID indicates an expected call of GetBalanceByUserID.
func (mr *MockBalanceMockRecorder) GetBalanceByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalanceByUserID", reflect.TypeOf((*MockBalance)(nil).GetBalanceByUserID), ctx, userID)
}

// GetBalanceOperationByOrder mocks base method.
func (m *MockBalance) GetBalanceOperationByOrder(ctx context.Context, order string) (*storage.BalanceOperationItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalanceOperationByOrder", ctx, order)
	ret0, _ := ret[0].(*storage.BalanceOperationItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalanceOperationByOrder indicates an expected call of GetBalanceOperationByOrder.
func (mr *MockBalanceMockRecorder) GetBalanceOperationByOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalanceOperationByOrder", reflect.TypeOf((*MockBalance)(nil).GetBalanceOperationByOrder), ctx, order)
}

// GetBalanceOperationByUser mocks base method.
func (m *MockBalance) GetBalanceOperationByUser(ctx context.Context, userID uuid.UUID) ([]*storage.BalanceOperationItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalanceOperationByUser", ctx, userID)
	ret0, _ := ret[0].([]*storage.BalanceOperationItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalanceOperationByUser indicates an expected call of GetBalanceOperationByUser.
func (mr *MockBalanceMockRecorder) GetBalanceOperationByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalanceOperationByUser", reflect.TypeOf((*MockBalance)(nil).GetBalanceOperationByUser), ctx, userID)
}

// IncrementBalance mocks base method.
func (m *MockBalance) IncrementBalance(ctx context.Context, incrementSum int64, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementBalance", ctx, incrementSum, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementBalance indicates an expected call of IncrementBalance.
func (mr *MockBalanceMockRecorder) IncrementBalance(ctx, incrementSum, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementBalance", reflect.TypeOf((*MockBalance)(nil).IncrementBalance), ctx, incrementSum, userID)
}

// UpdateBalanceOperation mocks base method.
func (m *MockBalance) UpdateBalanceOperation(ctx context.Context, order string, operationState storage.BalanceOperationState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalanceOperation", ctx, order, operationState)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBalanceOperation indicates an expected call of UpdateBalanceOperation.
func (mr *MockBalanceMockRecorder) UpdateBalanceOperation(ctx, order, operationState interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalanceOperation", reflect.TypeOf((*MockBalance)(nil).UpdateBalanceOperation), ctx, order, operationState)
}
