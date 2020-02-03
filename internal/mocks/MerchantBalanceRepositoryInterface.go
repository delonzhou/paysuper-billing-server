// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billingpb "github.com/paysuper/paysuper-proto/go/billingpb"
import context "context"
import mock "github.com/stretchr/testify/mock"

// MerchantBalanceRepositoryInterface is an autogenerated mock type for the MerchantBalanceRepositoryInterface type
type MerchantBalanceRepositoryInterface struct {
	mock.Mock
}

// CountByIdAndCurrency provides a mock function with given fields: _a0, _a1, _a2
func (_m *MerchantBalanceRepositoryInterface) CountByIdAndCurrency(_a0 context.Context, _a1 string, _a2 string) (int64, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int64); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByIdAndCurrency provides a mock function with given fields: _a0, _a1, _a2
func (_m *MerchantBalanceRepositoryInterface) GetByIdAndCurrency(_a0 context.Context, _a1 string, _a2 string) (*billingpb.MerchantBalance, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *billingpb.MerchantBalance
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *billingpb.MerchantBalance); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.MerchantBalance)
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

// Insert provides a mock function with given fields: _a0, _a1
func (_m *MerchantBalanceRepositoryInterface) Insert(_a0 context.Context, _a1 *billingpb.MerchantBalance) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.MerchantBalance) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
