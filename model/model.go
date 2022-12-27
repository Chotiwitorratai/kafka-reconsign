package model

import (
	"reflect"
	"time"
)

var Topics = []string{
	reflect.TypeOf(PaymentcallBack{}).Name(),
	reflect.TypeOf(InsuranceCallBack{}).Name(),
}

type PaymentcallBack struct {
	TransactionRefID             string
	Status                       string
	TransactionCreatedTimestamp  time.Time
	PaymentInfoAmount            float32
	PaymentInfoWebAdditionalInfo string
	PartnerInfoName              string
	PartnerInfoDeeplinkUrl       string
	PaymentPlatform              string
}

type InsuranceCallBack struct {
	RefID           string
	IdCard          string
	PlanCode        string
	PlanName        string
	EffectiveDate   time.Time
	ExpireDate      time.Time
	IssueDate       time.Time
	InsuranceStatus string
	TotalSumInsured float32
	ProductOwner    string
	PlanType        string
}
