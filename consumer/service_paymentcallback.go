package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/Shopify/sarama"

	"kafka-reconsign/model"
	"kafka-reconsign/repositories"
)

type ServicePayment interface {
	Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error
	PaymentCallbackProcess(rawMessage string) (err error)
	//InsuranceCallbackProcess(rawMessage string) (err error)
}

type servicePayment struct {
	reconcileRepo repositories.ReconcileRepositoryDB
}

func NewServicePayment(reconcileRepo repositories.ReconcileRepositoryDB) ServicePayment {
	return &servicePayment{reconcileRepo: reconcileRepo}
}

func (s *servicePayment) Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error {
	if topic == reflect.TypeOf(model.PaymentcallBack{}).Name() {
		err := s.PaymentCallbackProcess(rawMessage)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
	// switch topic {
	// case reflect.TypeOf(model.PaymentcallBack{}).Name():
	// 	err := s.PaymentCallbackProcess(rawMessage)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}

	// 	return nil
	// case reflect.TypeOf(model.InsuranceCallBack{}).Name():
	// 	err := s.InsuranceCallbackProcess(rawMessage)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	return nil
	// default:
	// 	return nil

	// }

}

func (s *servicePayment) PaymentCallbackProcess(rawMessage string) (err error) {
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

// func (s *service) InsuranceCallbackProcess(rawMessage string) (err error) {
// 	kafkaMsg := &model.InsuranceCallBack{}

// 	err = json.Unmarshal([]byte(rawMessage), kafkaMsg)
// 	if err != nil {
// 		return err
// 	}

// 	cipherIDcard := encrypt(kafkaMsg.IdCard)

// 	insurance := repositories.Reconcile{
// 		TransactionRefID: kafkaMsg.RefID,
// 		IdCard:           cipherIDcard,
// 		PlanCode:         kafkaMsg.PlanCode,
// 		PlanName:         kafkaMsg.PlanName,
// 		EffectiveDate:    kafkaMsg.EffectiveDate,
// 		ExpireDate:       kafkaMsg.ExpireDate,
// 		IssueDate:        kafkaMsg.IssueDate,
// 		InsuranceStatus:  kafkaMsg.InsuranceStatus,
// 		TotalSumInsured:  kafkaMsg.TotalSumInsured,
// 		ProductOwner:     kafkaMsg.ProductOwner,
// 		PlanType:         kafkaMsg.PlanType,
// 	}

// 	if insurance.TransactionRefID == "" {
// 		return errors.New("missing transaction Ref-ID")
// 	}

// 	if foundID, err := s.reconcileRepo.CheckNullReconcile(kafkaMsg.RefID); err != nil {
// 		return errors.New("invalid Ref-ID error")
// 	} else {
// 		if foundID {
// 			if err = s.reconcileRepo.UpdateReconcile(insurance); err != nil {
// 				return errors.New("update reconcile error")
// 			}
// 		} else {
// 			if err = s.reconcileRepo.SaveReconcile(insurance); err != nil {
// 				return errors.New("save reconcile error")
// 			}
// 		}
// 	}

// 	return nil
// }

// func GenerateNonce(length int) string {
// 	randomBytes := make([]byte, 32)
// 	_, err := rand.Read(randomBytes)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
// }
// func GetNonceFromCypher(cypher string) (nonce string, cyphertext string) {
// 	nonce = cypher[len(cypher)-12:]
// 	cyphertext = cypher[:len(cypher)-12]
// 	return nonce, cyphertext
// }

// func createHash(key string) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(key))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// func encrypt(Identifier string) string {
// 	genNonce := GenerateNonce(viper.GetInt("cdi.nonceLength"))
// 	aesKey := createHash(viper.GetString("secrets.cryptoAesKey"))
// 	encryptor := util.NewAES([]byte(aesKey), []byte(genNonce))
// 	citizenEncrypted := encryptor.Encrypt(Identifier)
// 	cypther := citizenEncrypted + genNonce
// 	return cypther
// }

// func decrypt(cyphertext string) (string, error) {
// 	nonce, cypher := GetNonceFromCypher(cyphertext)
// 	if len(nonce) != 12 {
// 		return "", errors.New("nonce not working")
// 	}
// 	aesKey := createHash(viper.GetString("secrets.cryptoAesKey"))
// 	crypto := util.NewAES([]byte(aesKey), []byte(nonce))
// 	fmt.Println(crypto)
// 	citizenDecrypted, err := crypto.Decrypt(cypher)
// 	if err != nil {
// 		return citizenDecrypted, err
// 	}
// 	return citizenDecrypted, nil
// }

//data for test
// payment call back (NEXT)
//{"TransactionRefID":"555test2022i","Status":"success","PaymentInfoAmount":1900,"PaymentInfoWebAdditionalInfo":"testing","PartnerInfoName":"TIP","PartnerInfoDeeplinkUrl":"www.example.com","PaymentPlatform":"paotang"}

// insurance call back
//{"RefID":"555test2022i","IdCard":"12555995336","PlanCode":"PL01","PlanName":"ประกันโควิด","EffectiveDate":"2022-01-25T12:11:56Z","ExpireDate":"2022-01-25T12:11:56Z","IssueDate":"2022-01-25T12:11:56Z","InsuranceStatus":"success","TotalSumInsured":10000,"ProductOwner":"POtest","PlanType":"Base plan"}
