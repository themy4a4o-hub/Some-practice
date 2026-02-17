package service_test

import (
	"pract/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Run(t *testing.T) {
	Produser := new(ProduserMock)
	Presenter := new(PresenterMock)

	Produser.On("Produce").Return([]string{"http://test.com"}, nil)
	Presenter.On("Present", []string{"http://********"}).Return(nil)

	srv := service.Service{
		Prod: Produser,
		Pres: Presenter,
	}

	err := srv.Run()

	assert.NoError(t, err)
	Produser.AssertExpectations(t)
	Presenter.AssertExpectations(t)
}

func TestService_Run_ProduserError(t *testing.T) {
	Produser := new(ProduserMock)
	Presenter := new(PresenterMock)
	expectedErr := errors.New("Producer failed")
	Produser.On("Produce").Return([]string(nil), expectedErr)

	srv := service.Service{
		Prod: Produser,
		Pres: Presenter,
	}
	err := srv.Run()
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	Produser.AssertExpectations(t)
	Presenter.AssertNotCalled(t, "Present", mock.Anything)
}

func TestService_Run_PresenterError(t *testing.T) {
	Produser := new(ProduserMock)
	Presenter := new(PresenterMock)
	expectedErr := errors.New("Presenter failed")
	Produser.On("Produce").Return([]string{"http://test.com"}, nil)
	Presenter.On("Present", []string{"http://********"}).Return(expectedErr)
	srv := service.Service{
		Prod: Produser,
		Pres: Presenter,
	}
	err := srv.Run()
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

}
