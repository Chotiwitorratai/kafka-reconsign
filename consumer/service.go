package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kafka-reconsign/model"
	"kafka-reconsign/repositories"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
)

const (
	lineNotifyAPI = "https://notify-api.line.me/api/notify"
	lineToken     = "HD5RBYffvJAao8Qxm3mX7Zg1CpSQrEGfdVWlumFipuY"
)

type Service struct {
	orderRepo *repositories.OrderRepository
}

func New(orderRepo *repositories.OrderRepository) *Service {
	return &Service{orderRepo: orderRepo}
}

// {"refId": "refId","state": "state","message": "message"}
func (s *Service) Process(topic string, ctx context.Context, kafkaHeader []*sarama.RecordHeader, rawMessage string) error {
	switch topic {
	case reflect.TypeOf(model.PaymentKNEXT{}).Name():
		kafkaMsg := &model.PaymentKNEXT{}
		err := json.Unmarshal([]byte(rawMessage), kafkaMsg)
		if err != nil {
			return err
		}

		payment := repositories.Order{
			RefId:      kafkaMsg.ID,
			NextStatus: kafkaMsg.Status,
			NextUpdate: strconv.Itoa(kafkaMsg.Update),
		}

		_, err = s.orderRepo.SaveOrder(payment)
		if err != nil {
			log.Println(err)
			return err
		}

	case reflect.TypeOf(model.InsuranceData{}).Name():
		kafkaMsg := &model.InsuranceData{}
		err := json.Unmarshal([]byte(rawMessage), kafkaMsg)
		if err != nil {
			return err
		}

		policy := repositories.Order{
			IpCallBack:      kafkaMsg.ID,
			ReconcileStatus: kafkaMsg.Status,
			IpUpdate:        strconv.Itoa(kafkaMsg.Update),
		}

		order, err := s.orderRepo.GetOrderByRef_id(kafkaMsg.ID)
		if err != nil {
			sendLineNotification("fail")
			return err
		}

		if order.NextStatus == "success" && kafkaMsg.Status != "success" {
			sendLineNotification("status missmach")
		} else {
			_, err = s.orderRepo.SaveOrder(policy) //need update table function
		}

	default:
		return nil

	}
	return nil
}

func sendLineNotification(message string) error {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new POST request
	req, err := http.NewRequest("POST", lineNotifyAPI, nil)
	if err != nil {
		return err
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+lineToken)

	// Create a new URL-encoded form
	form := url.Values{}
	form.Add("message", message)

	// Set the body of the request to the form
	req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))

	// Set the content type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the response status code
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send notification: %s", res.Status)
	}

	return nil
}
