package consumer_test

import (
	"encoding/json"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
			name:                         "success - new data",
			transactionRefID:             "test-transaction-ref-id",
			status:                       "success",
			transactionCreatedTimestamp:  time.Now(),
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
			transactionCreatedTimestamp:  time.Now(),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       nil,
			newData:                      false,
		},
		{
			name:                         "failure - new Data",
			transactionRefID:             "test-transaction-ref-id",
			status:                       "fail",
			transactionCreatedTimestamp:  time.Now(),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       errors.New(""),
			newData:                      true,
		},
		{
			name:                         "failure - second Data",
			transactionRefID:             "test-transaction-ref-id",
			status:                       "fail",
			transactionCreatedTimestamp:  time.Now(),
			paymentInfoAmount:            100,
			paymentInfoWebAdditionalInfo: "test-payment-info-web-additional-info",
			partnerInfoName:              "test-partner-info-name",
			partnerInfoDeeplinkUrl:       "test-partner-info-deeplink-url",
			paymentPlatform:              "test-payment-platform",
			expect:                       errors.New(""),
			newData:                      false,
		},
		{
			name:                         "failure - no transaction",
			transactionRefID:             "",
			status:                       "fail",
			transactionCreatedTimestamp:  time.Now(),
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
			if c.newData {
				reconcileRepository.On("SaveReconcile", c.transactionRefID).Return(nil)
			} else {
				reconcileRepository.On("UpdateReconcile", c.transactionRefID).Return(nil)
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

			_ = data
			_ = consumerService

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

	cases := []testCases{
		{
			name:            "success - first Data",
			refID:           "test-ref-id1",
			idCard:          "test-id-card",
			planCode:        "test-plan-code",
			planName:        "test-plan-name",
			effectiveDate:   time.Now(),
			expireDate:      time.Now(),
			issueDate:       time.Now(),
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
			effectiveDate:   time.Now(),
			expireDate:      time.Now(),
			issueDate:       time.Now(),
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
			effectiveDate:   time.Now(),
			expireDate:      time.Now(),
			issueDate:       time.Now(),
			insuranceStatus: "failure",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          errors.New("test error"),
			newData:         true,
		},
		{
			name:            "failure - second Data",
			refID:           "fail-test-ref-id",
			idCard:          "fail-test-id-card",
			planCode:        "fail-test-plan-code",
			planName:        "fail-test-plan-name",
			effectiveDate:   time.Now(),
			expireDate:      time.Now(),
			issueDate:       time.Now(),
			insuranceStatus: "failure",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          errors.New("test error"),
			newData:         false,
		},
		{
			name:            "failure - no ref-id",
			refID:           "",
			idCard:          "test-id-card",
			planCode:        "test-plan-code",
			planName:        "test-plan-name",
			effectiveDate:   time.Now(),
			expireDate:      time.Now(),
			issueDate:       time.Now(),
			insuranceStatus: "failure",
			totalSumInsured: 100,
			productOwner:    "test-product-owner",
			planType:        "test-plan-type",
			expect:          errors.New("ref-id is required"),
			newData:         true,
		},
	}

	_ = cases

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			//Arrange
			reconcileRepository := repositories.NewReconcileRepositoryMock()

			reconcileRepository.On("CheckNullReconcile", c.refID).Return(c.newData, nil)
			if c.newData {
				reconcileRepository.On("SaveReconcile", c.refID).Return(nil)
			} else {
				reconcileRepository.On("UpdateReconcile", c.refID).Return(nil)
			}

			consumerService := consumer.New(reconcileRepository)

			payment := repositories.Reconcile{
				TransactionRefID: c.refID,
				IdCard:           c.idCard,
				PlanCode:         c.planCode,
				PlanName:         c.planName,
				EffectiveDate:    c.effectiveDate,
				ExpireDate:       c.expireDate,
				IssueDate:        c.issueDate,
				InsuranceStatus:  c.insuranceStatus,
				TotalSumInsured:  c.totalSumInsured,
				ProductOwner:     c.productOwner,
				PlanType:         c.planType,
			}

			data, err := json.Marshal(payment)
			if err != nil {
				log.Print(err)
			}

			_ = data
			_ = consumerService

			//Assert
			assert.Equal(t, c.expect, err)
		})
	}

}
