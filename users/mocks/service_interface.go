// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Bruary/staff-scheduling/users/models"
	mock "github.com/stretchr/testify/mock"
)

// ServiceInterface is an autogenerated mock type for the ServiceInterface type
type ServiceInterface struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: _a0, _a1
func (_m *ServiceInterface) CreateUser(_a0 context.Context, _a1 models.CreateUserRequest) *models.CreateUserResponse {
	ret := _m.Called(_a0, _a1)

	var r0 *models.CreateUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateUserRequest) *models.CreateUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CreateUserResponse)
		}
	}

	return r0
}

// DeleteUser provides a mock function with given fields: _a0, _a1
func (_m *ServiceInterface) DeleteUser(_a0 context.Context, _a1 models.DeleteUserRequest) *models.DeleteUserResponse {
	ret := _m.Called(_a0, _a1)

	var r0 *models.DeleteUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, models.DeleteUserRequest) *models.DeleteUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DeleteUserResponse)
		}
	}

	return r0
}

// GetUserByEmail provides a mock function with given fields: _a0, _a1
func (_m *ServiceInterface) GetUserByEmail(_a0 context.Context, _a1 models.GetUserByEmailRequest) *models.GetUserResponse {
	ret := _m.Called(_a0, _a1)

	var r0 *models.GetUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, models.GetUserByEmailRequest) *models.GetUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetUserResponse)
		}
	}

	return r0
}

// GetUserByUID provides a mock function with given fields: ctx, userUID
func (_m *ServiceInterface) GetUserByUID(ctx context.Context, userUID string) *models.GetUserResponse {
	ret := _m.Called(ctx, userUID)

	var r0 *models.GetUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.GetUserResponse); ok {
		r0 = rf(ctx, userUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetUserResponse)
		}
	}

	return r0
}

// NewServiceInterface creates a new instance of ServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServiceInterface {
	mock := &ServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
