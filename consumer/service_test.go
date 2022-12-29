package consumer_test

import (
	"encoding/json"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"kafka-reconsign/consumer"
	"kafka-reconsign/repositories"
)

func TestPaymentCallbackProcess(t *testing.T) {

	type testCases struct {
		name                         string
		transactionRefID             string
		status                       string
		transactionCreatedTimestamp  time.Time
		paymentInfoAmount            float32
		paymentInfoWebAdditionalInfo string
		partnerInfoName              string
		partnerInfoDeeplinkUrl       string
		paymentPlatform              string
		expect                       error
		newData                      bool
	}

	cases := []testCases{
		{
			name:                         "success - first data",
			transactionRefID:             "test-transaction-ref-id",
			status:                       "success",
			transactionCreatedTimestamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       nil,
			newData:                      true,
		},
		{
			name:                         "success - second data",
			transactionRefID:             "test-transaction-ref-id",
			status:                       "success",
			transactionCreatedTimestamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       nil,
			newData:                      false,
		},
		{
			name:                         "failure - first Data",
			transactionRefID:             "test-transaction-ref-id",
			status:                       "fail",
			transactionCreatedTimestamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       nil,
			newData:                      true,
		},
		{
			name:                         "failure - second Data",
			transactionRefID:             "test-transaction-ref-id",
			status:                       "fail",
			transactionCreatedTimestamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       nil,
			newData:                      false,
		},
		{
			name:                         "failure - no transaction",
			transactionRefID:             "",
			status:                       "fail",
			transactionCreatedTimestamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       errors.New("missing transaction Ref-ID"),
			newData:                      true,
		},
	}

	_ = cases

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			//Arrange
			reconcileRepository := repositories.NewReconcileRepositoryMock()
			reconcileRepository.On("CheckNullReconcile", c.transactionRefID).Return(c.newData, nil)
			reconcileData := repositories.Reconcile{
				Model:                        gorm.Model{},
				TransactionRefID:             c.transactionRefID,
				Status:                       c.status,
				TransactionCreatedTimestamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				PaymentInfoAmount:            c.paymentInfoAmount,
				PaymentInfoWebAdditionalInfo: c.paymentInfoWebAdditionalInfo,
				PartnerInfoName:              c.partnerInfoName,
				PartnerInfoDeeplinkUrl:       c.partnerInfoDeeplinkUrl,
				PaymentPlatform:              c.paymentPlatform,
				IdCard:                       "",
				PlanCode:                     "",
				PlanName:                     "",
				EffectiveDate:                time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				ExpireDate:                   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				IssueDate:                    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				InsuranceStatus:              "",
				TotalSumInsured:              0,
				ProductOwner:                 "",
				PlanType:                     "",
			}
			if c.newData {
				reconcileRepository.On("UpdateReconcile", reconcileData).Return(nil)
			} else {
				reconcileRepository.On("SaveReconcile", reconcileData).Return(nil)
			}

			consumerService := consumer.New(reconcileRepository)

			payment := repositories.Reconcile{
				TransactionRefID:             c.transactionRefID,
				Status:                       c.status,
				PaymentInfoAmount:            c.paymentInfoAmount,
				PaymentInfoWebAdditionalInfo: c.paymentInfoWebAdditionalInfo,
				PartnerInfoName:              c.partnerInfoName,
				PartnerInfoDeeplinkUrl:       c.partnerInfoDeeplinkUrl,
				PaymentPlatform:              c.paymentPlatform,
			}

			data, err := json.Marshal(payment)
			if err != nil {
				log.Print(err)
			}

			//Act
			err = consumerService.PaymentCallbackProcess(string(data))

			//Assert
			assert.Equal(t, c.expect, err)
		})
	}

}
func TestInsuranceCallbackProcess(t *testing.T) {

	type testCases struct {
		name            string
		refID           string
		idCard          string
		planCode        string
		planName        string
		effectiveDate   time.Time
		expireDate      time.Time
		issueDate       time.Time
		insuranceStatus string
		totalSumInsured float32
		productOwner    string
		planType        string
		expect          error
		newData         bool
	}

	type insurance struct {
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

	cases := []testCases{
		{
			name:            "success - first Data",
			refID:           "test-ref-id1",
			idCard:          "test-id-card",
			planCode:        "test-plan-code",
			planName:        "test-plan-name",
			effectiveDate:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			expireDate:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			issueDate:       time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			insuranceStatus: "success",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          nil,
			newData:         true,
		},
		{
			name:            "success - second Data",
			refID:           "test-ref-id2",
			idCard:          "test-id-card",
			planCode:        "test-plan-code",
			planName:        "test-plan-name",
			effectiveDate:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			expireDate:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			issueDate:       time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			insuranceStatus: "success",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          nil,
			newData:         false,
		},
		{
			name:            "failure - first Data",
			refID:           "fail-test-ref-id",
			idCard:          "fail-test-id-card",
			planCode:        "fail-test-plan-code",
			planName:        "fail-test-plan-name",
			effectiveDate:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			expireDate:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			issueDate:       time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			insuranceStatus: "failure",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          nil,
			newData:         true,
		},
		{
			name:            "failure - second Data",
			refID:           "fail-test-ref-id",
			idCard:          "fail-test-id-card",
			planCode:        "fail-test-plan-code",
			planName:        "fail-test-plan-name",
			effectiveDate:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			expireDate:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			issueDate:       time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			insuranceStatus: "failure",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          nil,
			newData:         false,
		},
		{
			name:            "failure - no ref-id",
			refID:           "",
			idCard:          "test-id-card",
			planCode:        "test-plan-code",
			planName:        "test-plan-name",
			effectiveDate:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			expireDate:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			issueDate:       time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			insuranceStatus: "failure",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          errors.New("missing transaction Ref-ID"),
			newData:         true,
		},
	}

	_ = cases

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			//Arrange
			reconcileRepository := repositories.NewReconcileRepositoryMock()
			insuranceData := repositories.Reconcile{
				Model:                        gorm.Model{},
				TransactionRefID:             c.refID,
				Status:                       "",
				TransactionCreatedTimestamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				PaymentInfoAmount:            0,
				PaymentInfoWebAdditionalInfo: "",
				PartnerInfoName:              "",
				PartnerInfoDeeplinkUrl:       "",
				PaymentPlatform:              "",
				IdCard:                       c.idCard,
				PlanCode:                     c.planCode,
				PlanName:                     c.planName,
				EffectiveDate:                time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				ExpireDate:                   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				IssueDate:                    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				InsuranceStatus:              c.insuranceStatus,
				TotalSumInsured:              c.totalSumInsured,
				ProductOwner:                 c.productOwner,
				PlanType:                     c.planType,
			}
			reconcileRepository.On("CheckNullReconcile", c.refID).Return(c.newData, nil)
			if c.newData {
				reconcileRepository.On("UpdateReconcile", insuranceData).Return(nil)
			} else {
				reconcileRepository.On("SaveReconcile", insuranceData).Return(nil)
			}

			consumerService := consumer.New(reconcileRepository)

			insurance := insurance{
				RefID:           c.refID,
				IdCard:          c.idCard,
				PlanCode:        c.planCode,
				PlanName:        c.planName,
				EffectiveDate:   c.effectiveDate,
				ExpireDate:      c.expireDate,
				IssueDate:       c.issueDate,
				InsuranceStatus: c.insuranceStatus,
				TotalSumInsured: c.totalSumInsured,
				ProductOwner:    c.productOwner,
				PlanType:        c.planType,
			}

			data, err := json.Marshal(insurance)
			if err != nil {
				log.Print(err)
			}

			err = consumerService.InsuranceCallbackProcess(string(data))

			//Assert
			assert.Equal(t, c.expect, err)
		})
	}

}
