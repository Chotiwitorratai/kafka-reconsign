package repositories

import (
	"time"

	"gorm.io/gorm"
)

type Reconcile struct {
	gorm.Model
	TransactionRefID             string    `gorm:"colum:transaction_ref_id"`
	Status                       string    `gorm:"colum:status"`
	TransactionCreatedTimestamp  time.Time `gorm:"colum:transaction_created_timestamp"`
	PaymentInfoAmount            int       `gorm:"colum:payment_info_amount"`
	PaymentInfoWebAdditionalInfo string    `gorm:"colum:payment_info_web_additional_info"`
	PartnerInfoName              string    `gorm:"colum:partner_info_name"`
	PartnerInfoDeeplinkUrl       string    `gorm:"colum:partner_info_deeplink_url"`
	PaymentPlatform              string    `gorm:"colum:payment_platform"`
	IdCard                       string    `gorm:"colum:id_card"`
	PlanCode                     string    `gorm:"colum:plan_code"`
	PlanName                     string    `gorm:"colum:plan_name"`
	EffectiveDate                time.Time `gorm:"colum:effective_date"`
	ExpireDate                   string    `gorm:"colum:expire_date"`
	IssueDate                    string    `gorm:"colum:issue_date"`
	InsuranceStatus              string    `gorm:"colum:insurance_status"`
	TotalSumInsured              string    `gorm:"colum:total_sum_insured"`
	ProductOwner                 string    `gorm:"colum:product_owner"`
	PlanType                     string    `gorm:"colum:plan_type"`
}
type InsulerCallback struct {
	gorm.Model
}
type Alert struct {
	gorm.Model
	Messages  string `gorm:"colum:messages"`
	Count     string `gorm:"colum:count"`
	Status    string `gorm:"colum:status;default:Fail"`
	NextAlert string `gorm:"colum:next_alert"`
	RefId     string `gorm:"colum:ref_id"`
	Missing   string `gorm:"colum:missing"`
}

type ReconcileRepository interface {
	GetReconcileByRefID(id string)()
	GetReconcileStatus()()
	GetAlertFail()()

	SaveReconcile(reconcile Reconcile) (string, error)
	SaveAlert(alert Alert) (string, error)

	UpdateReconcile()()
	UpdateAlertStatus()()
}
