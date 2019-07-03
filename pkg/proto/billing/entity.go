package billing

import (
	"github.com/globalsign/mgo/bson"
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-recurring-repository/pkg/constant"
	"time"
)

var (
	orderRefundAllowedStatuses = map[int32]bool{
		constant.OrderStatusPaymentSystemComplete: true,
		constant.OrderStatusProjectInProgress:     true,
		constant.OrderStatusProjectComplete:       true,
		constant.OrderStatusProjectPending:        true,
	}

	orderStatusPublicMapping = map[int32]string{
		constant.OrderStatusNew:                         constant.OrderPublicStatusCreated,
		constant.OrderStatusPaymentSystemCreate:         constant.OrderPublicStatusCreated,
		constant.OrderStatusPaymentSystemCanceled:       constant.OrderPublicStatusCanceled,
		constant.OrderStatusPaymentSystemRejectOnCreate: constant.OrderPublicStatusRejected,
		constant.OrderStatusPaymentSystemReject:         constant.OrderPublicStatusRejected,
		constant.OrderStatusProjectReject:               constant.OrderPublicStatusRejected,
		constant.OrderStatusPaymentSystemDeclined:       constant.OrderPublicStatusRejected,
		constant.OrderStatusPaymentSystemComplete:       constant.OrderPublicStatusProcessed,
		constant.OrderStatusProjectComplete:             constant.OrderPublicStatusProcessed,
		constant.OrderStatusRefund:                      constant.OrderPublicStatusRefunded,
		constant.OrderStatusChargeback:                  constant.OrderPublicStatusChargeback,
	}
)

func (m *Merchant) ChangesAllowed() bool {
	return m.Status == pkg.MerchantStatusDraft
}

func (m *Merchant) GetPayoutCurrency() *Currency {
	if m.Banking == nil {
		return nil
	}

	return m.Banking.Currency
}

func (m *Merchant) NeedMarkESignAgreementAsSigned() bool {
	return m.HasMerchantSignature == true && m.HasPspSignature == true &&
		m.Status != pkg.MerchantStatusAgreementSigned
}

func (m *Merchant) CanGenerateAgreement() bool {
	return (m.Status == pkg.MerchantStatusOnReview || m.Status == pkg.MerchantStatusAgreementSigning ||
		m.Status == pkg.MerchantStatusAgreementSigned) && m.Banking != nil && m.Country != "" &&
		m.Contacts != nil && m.Contacts.Authorized != nil
}

func (m *Merchant) CanChangeStatusToSigning() bool {
	return m.Status == pkg.MerchantStatusOnReview && m.Banking != nil && m.Country != "" &&
		m.Contacts != nil && m.Contacts.Authorized != nil
}

func (m *Merchant) IsDeleted() bool {
	return m.Status == pkg.MerchantStatusDeleted
}

func (m *Order) HasEndedStatus() bool {
	return m.PrivateStatus == constant.OrderStatusPaymentSystemReject || m.PrivateStatus == constant.OrderStatusProjectComplete ||
		m.PrivateStatus == constant.OrderStatusProjectReject || m.PrivateStatus == constant.OrderStatusRefund ||
		m.PrivateStatus == constant.OrderStatusChargeback
}

func (m *Order) RefundAllowed() bool {
	v, ok := orderRefundAllowedStatuses[m.PrivateStatus]

	return ok && v == true
}

func (m *Order) FormInputTimeIsEnded() bool {
	t, err := ptypes.Timestamp(m.ExpireDateToFormInput)

	return err != nil || t.Before(time.Now())
}

func (m *Order) GetProjectId() string {
	return m.Project.Id
}

func (m *Project) IsProduction() bool {
	return m.Status == pkg.ProjectStatusInProduction
}

func (m *Project) IsDeleted() bool {
	return m.Status == pkg.ProjectStatusDeleted
}

func (m *Project) NeedChangeStatusToDraft(req *Project) bool {
	if m.Status != pkg.ProjectStatusTestCompleted &&
		m.Status != pkg.ProjectStatusInProduction {
		return false
	}

	if m.CallbackProtocol == pkg.ProjectCallbackProtocolEmpty &&
		req.CallbackProtocol == pkg.ProjectCallbackProtocolDefault {
		return true
	}

	if req.UrlCheckAccount != "" &&
		req.UrlCheckAccount != m.UrlCheckAccount {
		return true
	}

	if req.UrlProcessPayment != "" &&
		req.UrlProcessPayment != m.UrlProcessPayment {
		return true
	}

	return false
}

func (m *OrderUser) IsIdentified() bool {
	return m.Id != "" && bson.IsObjectIdHex(m.Id) == true
}

func (m *Order) GetPublicStatus() string {
	st, ok := orderStatusPublicMapping[m.PrivateStatus]
	if !ok {
		return constant.OrderPublicStatusPending
	}
	return st
}

func (m *Order) GetReceiptUserEmail() string {
	if m.User != nil {
		return m.User.Email
	}
	return ""
}

func (m *Order) GetReceiptUserPhone() string {
	if m.User != nil {
		return m.User.Phone
	}
	return ""
}

func (m *Order) GetCountry() string {
	if m.BillingAddress != nil && m.BillingAddress.Country != "" {
		return m.BillingAddress.Country
	}
	if m.User != nil && m.User.Address != nil && m.User.Address.Country != "" {
		return m.User.Address.Country
	}
	return ""
}

func (m *Order) GetState() string {
	if m.BillingAddress != nil && m.BillingAddress.State != "" {
		return m.BillingAddress.State
	}
	if m.User != nil && m.User.Address != nil && m.User.Address.State != "" {
		return m.User.Address.State
	}
	return ""
}

func (m *Order) SetNotificationStatus(key string, val bool) {
	if m.IsNotificationsSent == nil {
		m.IsNotificationsSent = make(map[string]bool)
	}
	m.IsNotificationsSent[key] = val
}

func (m *Order) GetNotificationStatus(key string) bool {
	if m.IsNotificationsSent == nil {
		return false
	}
	val, ok := m.IsNotificationsSent[key]
	if !ok {
		return false
	}
	return val
}

func (m *PaymentMethod) IsValid() bool {
	return m.ExternalId != "" &&
		m.Currencies != nil &&
		m.Type != "" &&
		m.Group != "" &&
		m.Name != "" &&
		m.TestSettings != nil &&
		m.ProductionSettings != nil
}

func (m *Merchant) HasAuthorizedEmail() bool {
	return m.Contacts != nil && m.Contacts.Authorized != nil && m.Contacts.Authorized.Email != ""
}

func (m *Merchant) GetAuthorizedEmail() string {
	return m.Contacts.Authorized.Email
}

func (m *RoyaltyReport) ChangesAvailable(newStatus string) bool {
	if m.Status == pkg.RoyaltyReportStatusAccepted {
		return false
	}

	if m.Status == pkg.RoyaltyReportStatusNew && newStatus != pkg.RoyaltyReportStatusPending && newStatus != pkg.RoyaltyReportStatusCanceled {
		return false
	}

	if m.Status == pkg.RoyaltyReportStatusPending && newStatus != pkg.RoyaltyReportStatusAccepted &&
		newStatus != pkg.RoyaltyReportStatusDispute {
		return false
	}

	if m.Status == pkg.RoyaltyReportStatusCanceled && newStatus != pkg.RoyaltyReportStatusNew {
		return false
	}

	if m.Status == pkg.RoyaltyReportStatusDispute && newStatus != pkg.RoyaltyReportStatusPending {
		return false
	}

	return true
}
