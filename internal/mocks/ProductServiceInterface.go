// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import grpc "github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
import mock "github.com/stretchr/testify/mock"

// ProductServiceInterface is an autogenerated mock type for the ProductServiceInterface type
type ProductServiceInterface struct {
	mock.Mock
}

// CountByProjectSku provides a mock function with given fields: _a0, _a1
func (_m *ProductServiceInterface) CountByProjectSku(_a0 string, _a1 string) (int, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, string) int); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: _a0
func (_m *ProductServiceInterface) GetById(_a0 string) (*grpc.Product, error) {
	ret := _m.Called(_a0)

	var r0 *grpc.Product
	if rf, ok := ret.Get(0).(func(string) *grpc.Product); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *ProductServiceInterface) List(_a0 string, _a1 string, _a2 string, _a3 string, _a4 int, _a5 int) (int32, []*grpc.Product) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 int32
	if rf, ok := ret.Get(0).(func(string, string, string, string, int, int) int32); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r0 = ret.Get(0).(int32)
	}

	var r1 []*grpc.Product
	if rf, ok := ret.Get(1).(func(string, string, string, string, int, int) []*grpc.Product); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*grpc.Product)
		}
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: product
func (_m *ProductServiceInterface) Upsert(product *grpc.Product) error {
	ret := _m.Called(product)

	var r0 error
	if rf, ok := ret.Get(0).(func(*grpc.Product) error); ok {
		r0 = rf(product)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
