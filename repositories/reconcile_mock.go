package repositories

import "github.com/stretchr/testify/mock"

type reconcileRepositoryMock struct {
	mock.Mock
}

func NewReconcileRepositoryMock() *reconcileRepositoryMock {
	return &reconcileRepositoryMock{}
}

func (m *reconcileRepositoryMock) CheckNullReconcile(id string) (bool, error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *reconcileRepositoryMock) SaveReconcile(reconcile Reconcile) error {
	args := m.Called(reconcile)
	return args.Error(0)
}

func (m *reconcileRepositoryMock) UpdateReconcile(reconcile Reconcile) error {
	args := m.Called(reconcile)
	return args.Error(0)
}

func (m *reconcileRepositoryMock) GetReconcileFail() (reconcile []Reconcile, err error) {
	args := m.Called(reconcile)
	return args.Get(0).([]Reconcile), args.Error(1)
}

func (m *reconcileRepositoryMock) SaveAlert(alert []Alert) error {
	args := m.Called(alert)
	return args.Error(0)
}

func (m *reconcileRepositoryMock) UpdateAlert(alert Alert) error {
	args := m.Called(alert)
	return args.Error(0)
}

func (m *reconcileRepositoryMock) GetAlertFail() (alert []Alert, err error) {
	args := m.Called(alert)
	return args.Get(0).([]Alert), args.Error(1)
}

func (m *reconcileRepositoryMock) GetAlertFailByID(id string) (boo bool, err error) {
	args := m.Called(id)
	return args.Get(0).(bool), args.Error(1)
}
