package consumer

import (
	"context"
	"encoding/json"
	"kafka-reconsign/model"
	"kafka-reconsign/repositories"
	"log"
	"reflect"

	"github.com/Shopify/sarama"
)

type Service struct {
	orderRepo repositories.ReconcileRepositoryDB
}

func New(orderRepo repositories.ReconcileRepositoryDB) *Service {
	return &Service{orderRepo: orderRepo}
}

func (s *Service) Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error {
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

func paymentCallbackProcess(s *Service, rawMessage string) (err error) {
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

	if foundID, err := s.orderRepo.CheckNullReconcile(payment.TransactionRefID); err != nil {
		return err
	} else {
		if foundID {
			if err = s.orderRepo.UpdateReconcile(payment); err != nil {
				return err
			}
		} else {
			if err = s.orderRepo.SaveReconcile(payment); err != nil {
				return err
			}
		}
	}

	return nil
}

func insuranceCallbackProcess(s *Service, rawMessage string) (err error) {
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

	if foundID, err := s.orderRepo.CheckNullReconcile(kafkaMsg.RefID); err != nil {
		log.Println("Check reconcile errror :", err)
		return err
	} else {
		if foundID {
			if err = s.orderRepo.UpdateReconcile(insurance); err != nil {
				log.Println("update reconcile error:", err)
				return err
			}
		} else {
			if err = s.orderRepo.SaveReconcile(insurance); err != nil {
				log.Println("save reconcile error:", err)
				return err
			}
		}
	}

	return nil
}
