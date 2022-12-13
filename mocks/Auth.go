// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"

	gocloak "github.com/Nerzal/gocloak/v12"

	mock "github.com/stretchr/testify/mock"
)

// Auth is an autogenerated mock type for the Auth type
type Auth struct {
	mock.Mock
}

// GetUsersByEmail provides a mock function with given fields: ctx, email
func (_m *Auth) GetUsersByEmail(ctx context.Context, email string) ([]*gocloak.User, error) {
	ret := _m.Called(ctx, email)

	var r0 []*gocloak.User
	if rf, ok := ret.Get(0).(func(context.Context, string) []*gocloak.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*gocloak.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, rq
func (_m *Auth) Login(ctx context.Context, rq domain.LoginUser) (string, error) {
	ret := _m.Called(ctx, rq)

	return ret.String(0), ret.Error(1)
}

// Logout provides a mock function with given fields: ctx, token
func (_m *Auth) Logout(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Register provides a mock function with given fields: ctx, rq
func (_m *Auth) Register(ctx context.Context, rq domain.RegisterUser) (string, error) {
	ret := _m.Called(ctx, rq)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, domain.RegisterUser) string); ok {
		r0 = rf(ctx, rq)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.RegisterUser) error); ok {
		r1 = rf(ctx, rq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendEmail provides a mock function with given fields: ctx, userID
func (_m *Auth) SendEmail(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendVerifyEmail provides a mock function with given fields: ctx, userID
func (_m *Auth) SendVerifyEmail(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserExists provides a mock function with given fields: ctx, email
func (_m *Auth) UserExists(ctx context.Context, email string) bool {
	ret := _m.Called(ctx, email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewAuth interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuth creates a new instance of Auth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuth(t mockConstructorTestingTNewAuth) *Auth {
	mock := &Auth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}