package repositories

import "github.com/stretchr/testify/mock"

type orderRepositoryDBMock struct {
	mock.Mock
}

func NewOrderRepositoryDBMock() *orderRepositoryDBMock {
	return &orderRepositoryDBMock{}
}

func (m *orderRepositoryDBMock) GetReconcileByRefID(id string) (reconciles []Reconcile, err error) {
	args := m.Called(id)
	return args.Get(0).([]Reconcile), args.Error(1)
}
func (m *orderRepositoryDBMock) SaveReconcile(reconcile Reconcile) error {
	args := m.Called(reconcile)
	return args.Error(0)
}

// func (m *orderRepositoryDBMock) ChangeStatus(string, string) (string, error) {
// 	args := m.Called()
// 	return args.Get(0).(string), args.Error(1)
// }
