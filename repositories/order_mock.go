package repositories

import "github.com/stretchr/testify/mock"

type orderRepositoryDBMock struct {
	mock.Mock
}

func NewOrderRepositoryDBMock() *orderRepositoryDBMock {
	return &orderRepositoryDBMock{}
}

func (m *orderRepositoryDBMock) SaveOrder(Order) (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}
func (m *orderRepositoryDBMock) GetOrderByRef_id(string) (*Order, error) {
	args := m.Called()
	return args.Get(0).(*Order), args.Error(1)
}
func (m *orderRepositoryDBMock) ChangeStatus(string, string) (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}
