package service_test

import "github.com/stretchr/testify/mock"

type ProduserMock struct {
	mock.Mock
}

func (m *ProduserMock) Produce() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

type PresenterMock struct {
	mock.Mock
}

func (m *PresenterMock) Present(n []string) error {
	args := m.Called(n)
	return args.Error(0)
}
