package consumer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/Shopify/sarama"

	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"kafka-reconsign/model"
	"kafka-reconsign/repositories"
)

type Service interface {
	Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error
	PaymentCallbackProcess(rawMessage string) (err error)
	InsuranceCallbackProcess(rawMessage string) (err error)
}

type service struct {
	reconcileRepo repositories.ReconcileRepositoryDB
}

func New(reconcileRepo repositories.ReconcileRepositoryDB) Service {
	return &service{reconcileRepo: reconcileRepo}
}

func (s *service) Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error {
	switch topic {
	case reflect.TypeOf(model.PaymentcallBack{}).Name():
		err := s.PaymentCallbackProcess(rawMessage)
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	case reflect.TypeOf(model.InsuranceCallBack{}).Name():
		err := s.InsuranceCallbackProcess(rawMessage)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	default:
		return nil

	}

}

func (s *service) PaymentCallbackProcess(rawMessage string) (err error) {
	kafkaMsg := &model.PaymentcallBack{}
	err = json.Unmarshal([]byte(rawMessage), kafkaMsg)
	if err != nil {
		return err
	}

	payment := repositories.Reconcile{
		TransactionRefID:             kafkaMsg.TransactionRefID,
		NextStatus:                   kafkaMsg.Status,
		PaymentInfoAmount:            kafkaMsg.PaymentInfoAmount,
		PaymentInfoWebAdditionalInfo: kafkaMsg.PaymentInfoWebAdditionalInfo,
		PartnerInfoName:              kafkaMsg.PartnerInfoName,
		PartnerInfoDeeplinkUrl:       kafkaMsg.PartnerInfoDeeplinkUrl,
		PaymentPlatform:              kafkaMsg.PaymentPlatform,
	}

	if payment.TransactionRefID == "" {
		return errors.New("missing transaction Ref-ID")
	}

	if foundID, err := s.reconcileRepo.CheckNullReconcile(payment.TransactionRefID); err != nil {
		return errors.New("invalid Ref-ID error")
	} else {
		if foundID {
			if err = s.reconcileRepo.UpdateReconcile(payment); err != nil {
				return errors.New("update reconcile error")
			}
		} else {
			if err = s.reconcileRepo.SaveReconcile(payment); err != nil {
				return errors.New("save reconcile error")
			}
		}
	}

	return nil
}

func (s *service) InsuranceCallbackProcess(rawMessage string) (err error) {
	kafkaMsg := &model.InsuranceCallBack{}

	err = json.Unmarshal([]byte(rawMessage), kafkaMsg)
	if err != nil {
		return err
	}

	ciphertext, err := encrypt([]byte(kafkaMsg.IdCard), "password")
	if err != nil {
		return errors.New("error encrypt id card")
	}
	cipherIDcard := fmt.Sprintf("%v", ciphertext)
	insurance := repositories.Reconcile{
		TransactionRefID: kafkaMsg.RefID,
		IdCard:           cipherIDcard,
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

	if foundID, err := s.reconcileRepo.CheckNullReconcile(kafkaMsg.RefID); err != nil {
		return errors.New("invalid Ref-ID error")
	} else {
		if foundID {
			if err = s.reconcileRepo.UpdateReconcile(insurance); err != nil {
				return errors.New("update reconcile error")
			}
		} else {
			if err = s.reconcileRepo.SaveReconcile(insurance); err != nil {
				return errors.New("save reconcile error")
			}
		}
	}

	return nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
