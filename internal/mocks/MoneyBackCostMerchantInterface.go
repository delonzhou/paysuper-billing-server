// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billingpb "github.com/paysuper/paysuper-proto/go/billingpb"
import mock "github.com/stretchr/testify/mock"

// MoneyBackCostMerchantInterface is an autogenerated mock type for the MoneyBackCostMerchantInterface type
type MoneyBackCostMerchantInterface struct {
	mock.Mock
}

// Delete provides a mock function with given fields: obj
func (_m *MoneyBackCostMerchantInterface) Delete(obj *billingpb.MoneyBackCostMerchant) error {
	ret := _m.Called(obj)

	var r0 error
	if rf, ok := ret.Get(0).(func(*billingpb.MoneyBackCostMerchant) error); ok {
		r0 = rf(obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: merchantId, name, payoutCurrency, undoReason, region, country, mccCode, paymentStage
func (_m *MoneyBackCostMerchantInterface) Get(merchantId string, name string, payoutCurrency string, undoReason string, region string, country string, mccCode string, paymentStage int32) (*billingpb.MoneyBackCostMerchantList, error) {
	ret := _m.Called(merchantId, name, payoutCurrency, undoReason, region, country, mccCode, paymentStage)

	var r0 *billingpb.MoneyBackCostMerchantList
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, string, string, int32) *billingpb.MoneyBackCostMerchantList); ok {
		r0 = rf(merchantId, name, payoutCurrency, undoReason, region, country, mccCode, paymentStage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.MoneyBackCostMerchantList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string, string, string, int32) error); ok {
		r1 = rf(merchantId, name, payoutCurrency, undoReason, region, country, mccCode, paymentStage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllForMerchant provides a mock function with given fields: merchantId
func (_m *MoneyBackCostMerchantInterface) GetAllForMerchant(merchantId string) (*billingpb.MoneyBackCostMerchantList, error) {
	ret := _m.Called(merchantId)

	var r0 *billingpb.MoneyBackCostMerchantList
	if rf, ok := ret.Get(0).(func(string) *billingpb.MoneyBackCostMerchantList); ok {
		r0 = rf(merchantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.MoneyBackCostMerchantList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(merchantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *MoneyBackCostMerchantInterface) GetById(id string) (*billingpb.MoneyBackCostMerchant, error) {
	ret := _m.Called(id)

	var r0 *billingpb.MoneyBackCostMerchant
	if rf, ok := ret.Get(0).(func(string) *billingpb.MoneyBackCostMerchant); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.MoneyBackCostMerchant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultipleInsert provides a mock function with given fields: obj
func (_m *MoneyBackCostMerchantInterface) MultipleInsert(obj []*billingpb.MoneyBackCostMerchant) error {
	ret := _m.Called(obj)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*billingpb.MoneyBackCostMerchant) error); ok {
		r0 = rf(obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: obj
func (_m *MoneyBackCostMerchantInterface) Update(obj *billingpb.MoneyBackCostMerchant) error {
	ret := _m.Called(obj)

	var r0 error
	if rf, ok := ret.Get(0).(func(*billingpb.MoneyBackCostMerchant) error); ok {
		r0 = rf(obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
