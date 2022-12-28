package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"kafka-reconsign/model"
	"kafka-reconsign/repositories"
	"log"
	"reflect"

	"github.com/Shopify/sarama"
)

type Service interface {
	Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error
}

type service struct {
	orderRepo repositories.ReconcileRepositoryDB
}

func New(orderRepo repositories.ReconcileRepositoryDB) Service {
	return &service{orderRepo: orderRepo}
}

func (s *service) Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error {
	switch topic {
	case reflect.TypeOf(model.PaymentcallBack{}).Name():
		err := paymentCallbackProcess(s, rawMessage)
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	case reflect.TypeOf(model.InsuranceCallBack{}).Name():
		err := insuranceCallbackProcess(s, rawMessage)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	default:
		return nil

	}

}

func paymentCallbackProcess(s *service, rawMessage string) (err error) {
	kafkaMsg := &model.PaymentcallBack{}
	err = json.Unmarshal([]byte(rawMessage), kafkaMsg)
	if err != nil {
		return err
	}

	payment := repositories.Reconcile{
		TransactionRefID:             kafkaMsg.TransactionRefID,
		Status:                       kafkaMsg.Status,
		PaymentInfoAmount:            kafkaMsg.PaymentInfoAmount,
		PaymentInfoWebAdditionalInfo: kafkaMsg.PaymentInfoWebAdditionalInfo,
		PartnerInfoName:              kafkaMsg.PartnerInfoName,
		PartnerInfoDeeplinkUrl:       kafkaMsg.PartnerInfoDeeplinkUrl,
		PaymentPlatform:              kafkaMsg.PaymentPlatform,
	}

	if payment.TransactionRefID == "" {
		return errors.New("missing transaction Ref-ID")
	}

	if foundID, err := s.orderRepo.CheckNullReconcile(payment.TransactionRefID); err != nil {
		return errors.New("invalid Ref-ID error")
	} else {
		if foundID {
			if err = s.orderRepo.UpdateReconcile(payment); err != nil {
				return errors.New("update reconcile error")
			}
		} else {
			if err = s.orderRepo.SaveReconcile(payment); err != nil {
				return errors.New("save reconcile error")
			}
		}
	}

	return nil
}

func insuranceCallbackProcess(s *service, rawMessage string) (err error) {
	kafkaMsg := &model.InsuranceCallBack{}

	err = json.Unmarshal([]byte(rawMessage), kafkaMsg)
	if err != nil {
		return err
	}

	insurance := repositories.Reconcile{
		TransactionRefID: kafkaMsg.RefID,
		IdCard:           kafkaMsg.IdCard,
		PlanCode:         kafkaMsg.PlanCode,
		PlanName:         kafkaMsg.PlanName,
		EffectiveDate:    kafkaMsg.EffectiveDate,
		ExpireDate:       kafkaMsg.ExpireDate,
		IssueDate:        kafkaMsg.IssueDate,
		InsuranceStatus:  kafkaMsg.InsuranceStatus,
		TotalSumInsured:  kafkaMsg.TotalSumInsured,
		ProductOwner:     kafkaMsg.ProductOwner,
		PlanType:         kafkaMsg.PlanType,
	}

	if insurance.TransactionRefID == "" {
		return errors.New("missing transaction Ref-ID")
	}

	if foundID, err := s.orderRepo.CheckNullReconcile(kafkaMsg.RefID); err != nil {
		return errors.New("invalid Ref-ID error")
	} else {
		if foundID {
			if err = s.orderRepo.UpdateReconcile(insurance); err != nil {
				return errors.New("update reconcile error")
			}
		} else {
			if err = s.orderRepo.SaveReconcile(insurance); err != nil {
				return errors.New("save reconcile error")
			}
		}
	}

	return nil
}
