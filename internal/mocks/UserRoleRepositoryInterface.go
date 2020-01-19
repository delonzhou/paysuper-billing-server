// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billingpb "github.com/paysuper/paysuper-proto/go/billingpb"
import context "context"
import mock "github.com/stretchr/testify/mock"

// UserRoleRepositoryInterface is an autogenerated mock type for the UserRoleRepositoryInterface type
type UserRoleRepositoryInterface struct {
	mock.Mock
}

// AddAdminUser provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) AddAdminUser(_a0 context.Context, _a1 *billingpb.UserRole) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.UserRole) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddMerchantUser provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) AddMerchantUser(_a0 context.Context, _a1 *billingpb.UserRole) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.UserRole) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAdminUser provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) DeleteAdminUser(_a0 context.Context, _a1 *billingpb.UserRole) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.UserRole) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteMerchantUser provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) DeleteMerchantUser(_a0 context.Context, _a1 *billingpb.UserRole) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.UserRole) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAdminUserByEmail provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) GetAdminUserByEmail(_a0 context.Context, _a1 string) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.UserRole); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
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

// GetAdminUserById provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) GetAdminUserById(_a0 context.Context, _a1 string) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.UserRole); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
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

// GetAdminUserByUserId provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) GetAdminUserByUserId(_a0 context.Context, _a1 string) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.UserRole); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
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

// GetMerchantOwner provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) GetMerchantOwner(_a0 context.Context, _a1 string) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.UserRole); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
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

// GetMerchantUserByEmail provides a mock function with given fields: _a0, _a1, _a2
func (_m *UserRoleRepositoryInterface) GetMerchantUserByEmail(_a0 context.Context, _a1 string, _a2 string) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *billingpb.UserRole); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMerchantUserById provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) GetMerchantUserById(_a0 context.Context, _a1 string) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.UserRole); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
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

// GetMerchantUserByUserId provides a mock function with given fields: _a0, _a1, _a2
func (_m *UserRoleRepositoryInterface) GetMerchantUserByUserId(_a0 context.Context, _a1 string, _a2 string) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *billingpb.UserRole); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMerchantsForUser provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) GetMerchantsForUser(_a0 context.Context, _a1 string) ([]*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string) []*billingpb.UserRole); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*billingpb.UserRole)
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

// GetSystemAdmin provides a mock function with given fields: _a0
func (_m *UserRoleRepositoryInterface) GetSystemAdmin(_a0 context.Context) (*billingpb.UserRole, error) {
	ret := _m.Called(_a0)

	var r0 *billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context) *billingpb.UserRole); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.UserRole)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsersForAdmin provides a mock function with given fields: _a0
func (_m *UserRoleRepositoryInterface) GetUsersForAdmin(_a0 context.Context) ([]*billingpb.UserRole, error) {
	ret := _m.Called(_a0)

	var r0 []*billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context) []*billingpb.UserRole); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*billingpb.UserRole)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsersForMerchant provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) GetUsersForMerchant(_a0 context.Context, _a1 string) ([]*billingpb.UserRole, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*billingpb.UserRole
	if rf, ok := ret.Get(0).(func(context.Context, string) []*billingpb.UserRole); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*billingpb.UserRole)
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

// UpdateAdminUser provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) UpdateAdminUser(_a0 context.Context, _a1 *billingpb.UserRole) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.UserRole) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMerchantUser provides a mock function with given fields: _a0, _a1
func (_m *UserRoleRepositoryInterface) UpdateMerchantUser(_a0 context.Context, _a1 *billingpb.UserRole) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.UserRole) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
