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
	PaymentInfoAmount            float32       `gorm:"colum:payment_info_amount"`
	PaymentInfoWebAdditionalInfo string    `gorm:"colum:payment_info_web_additional_info"`
	PartnerInfoName              string    `gorm:"colum:partner_info_name"`
	PartnerInfoDeeplinkUrl       string    `gorm:"colum:partner_info_deeplink_url"`
	PaymentPlatform              string    `gorm:"colum:payment_platform"`
	IdCard                       string    `gorm:"colum:id_card"`
	PlanCode                     string    `gorm:"colum:plan_code"`
	PlanName                     string    `gorm:"colum:plan_name"`
	EffectiveDate                time.Time `gorm:"colum:effective_date"`
	ExpireDate                   time.Time    `gorm:"colum:expire_date"`
	IssueDate                    time.Time    `gorm:"colum:issue_date"`
	InsuranceStatus              string    `gorm:"colum:insurance_status"`
	TotalSumInsured              float32    `gorm:"colum:total_sum_insured"`
	ProductOwner                 string    `gorm:"colum:product_owner"`
	PlanType                     string    `gorm:"colum:plan_type"`
}
type InsulerCallback struct {
	gorm.Model
}
type Alert struct {
	gorm.Model
	Messages  string `gorm:"colum:messages"`
	Count     int `gorm:"colum:count"`
	Status    string `gorm:"colum:status;default:Fail"`
	NextAlert time.Time `gorm:"colum:next_alert"`
	RefId     string `gorm:"colum:ref_id"`
	Missing   string `gorm:"colum:missing"`
}

type ReconcileRepository interface {
	CheckNullReconcile(id string) (bool, error)
	SaveReconcile(reconcile Reconcile) error
	UpdateReconcile(reconcile Reconcile) error

	GGetReconcileFail()(reconcile []Reconcile, err error)
	GetAlertFail()(reconcile []Reconcile, err error)

	SaveAlert(alert []Alert) error
	UpdateAlertStatus(alert Alert) error
}
