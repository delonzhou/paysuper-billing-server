// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billingpb "github.com/paysuper/paysuper-proto/go/billingpb"
import context "context"
import mock "github.com/stretchr/testify/mock"

// ProjectRepositoryInterface is an autogenerated mock type for the ProjectRepositoryInterface type
type ProjectRepositoryInterface struct {
	mock.Mock
}

// CountByMerchantId provides a mock function with given fields: _a0, _a1
func (_m *ProjectRepositoryInterface) CountByMerchantId(_a0 context.Context, _a1 string) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6
func (_m *ProjectRepositoryInterface) Find(_a0 context.Context, _a1 string, _a2 string, _a3 []int32, _a4 int64, _a5 int64, _a6 []string) ([]*billingpb.Project, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5, _a6)

	var r0 []*billingpb.Project
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []int32, int64, int64, []string) []*billingpb.Project); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*billingpb.Project)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, []int32, int64, int64, []string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindCount provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *ProjectRepositoryInterface) FindCount(_a0 context.Context, _a1 string, _a2 string, _a3 []int32) (int64, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []int32) int64); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, []int32) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: _a0, _a1
func (_m *ProjectRepositoryInterface) GetById(_a0 context.Context, _a1 string) (*billingpb.Project, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.Project
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.Project); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.Project)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *ProjectRepositoryInterface) Insert(_a0 context.Context, _a1 *billingpb.Project) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.Project) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MultipleInsert provides a mock function with given fields: _a0, _a1
func (_m *ProjectRepositoryInterface) MultipleInsert(_a0 context.Context, _a1 []*billingpb.Project) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*billingpb.Project) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *ProjectRepositoryInterface) Update(_a0 context.Context, _a1 *billingpb.Project) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.Project) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
