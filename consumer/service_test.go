package consumer_test

import (
	"errors"
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

			_ = payment
			_ = consumerService

			//Assert
			assert.Equal(t, c.expect, err)
		})
	}

}
