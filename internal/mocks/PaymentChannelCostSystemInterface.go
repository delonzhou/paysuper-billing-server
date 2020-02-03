// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billingpb "github.com/paysuper/paysuper-proto/go/billingpb"
import context "context"
import mock "github.com/stretchr/testify/mock"

// PaymentChannelCostSystemInterface is an autogenerated mock type for the PaymentChannelCostSystemInterface type
type PaymentChannelCostSystemInterface struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, obj
func (_m *PaymentChannelCostSystemInterface) Delete(ctx context.Context, obj *billingpb.PaymentChannelCostSystem) error {
	ret := _m.Called(ctx, obj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.PaymentChannelCostSystem) error); ok {
		r0 = rf(ctx, obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, name, region, country, mccCode, operatingCompanyId
func (_m *PaymentChannelCostSystemInterface) Get(ctx context.Context, name string, region string, country string, mccCode string, operatingCompanyId string) (*billingpb.PaymentChannelCostSystem, error) {
	ret := _m.Called(ctx, name, region, country, mccCode, operatingCompanyId)

	var r0 *billingpb.PaymentChannelCostSystem
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string) *billingpb.PaymentChannelCostSystem); ok {
		r0 = rf(ctx, name, region, country, mccCode, operatingCompanyId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.PaymentChannelCostSystem)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, string) error); ok {
		r1 = rf(ctx, name, region, country, mccCode, operatingCompanyId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *PaymentChannelCostSystemInterface) GetAll(ctx context.Context) (*billingpb.PaymentChannelCostSystemList, error) {
	ret := _m.Called(ctx)

	var r0 *billingpb.PaymentChannelCostSystemList
	if rf, ok := ret.Get(0).(func(context.Context) *billingpb.PaymentChannelCostSystemList); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.PaymentChannelCostSystemList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: ctx, id
func (_m *PaymentChannelCostSystemInterface) GetById(ctx context.Context, id string) (*billingpb.PaymentChannelCostSystem, error) {
	ret := _m.Called(ctx, id)

	var r0 *billingpb.PaymentChannelCostSystem
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.PaymentChannelCostSystem); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.PaymentChannelCostSystem)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultipleInsert provides a mock function with given fields: ctx, obj
func (_m *PaymentChannelCostSystemInterface) MultipleInsert(ctx context.Context, obj []*billingpb.PaymentChannelCostSystem) error {
	ret := _m.Called(ctx, obj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*billingpb.PaymentChannelCostSystem) error); ok {
		r0 = rf(ctx, obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, obj
func (_m *PaymentChannelCostSystemInterface) Update(ctx context.Context, obj *billingpb.PaymentChannelCostSystem) error {
	ret := _m.Called(ctx, obj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.PaymentChannelCostSystem) error); ok {
		r0 = rf(ctx, obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
