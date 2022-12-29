package notification

import (
	"fmt"
	"io/ioutil"
	"kafka-reconsign/repositories"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	lineNotifyAPI = "https://notify-api.line.me/api/notify"
	lineToken     = "HD5RBYffvJAao8Qxm3mX7Zg1CpSQrEGfdVWlumFipuY"
)

type ReconcileJob interface {
	CheckReconcileStatus() error
	CheckAlertStatus() error
}

type reconcileJob struct {
	reconcileRepo repositories.ReconcileRepositoryDB
}

func NewReconcileJob(reconcileRepo repositories.ReconcileRepositoryDB) ReconcileJob {
	return reconcileJob{reconcileRepo: reconcileRepo}
}

func (n reconcileJob) CheckReconcileStatus() error {
	for {
		data := []repositories.Alert{}
		reconcile, err := n.reconcileRepo.GetReconcileFail()
		if err != nil {
			log.Println("get reconcile error", err)
			return err
		}
		if len(reconcile) > 0 {
			for _, rec := range reconcile {
				foundID, err := n.reconcileRepo.GetAlertFailByID(rec.TransactionRefID)
				if err != nil {
					return err
				}
				if !foundID {
					data = append(data, repositories.Alert{
						Messages:  "test",
						Count:     0,
						Status:    "Fail",
						NextAlert: time.Now().Add(5 * time.Minute),
						RefId:     rec.TransactionRefID,
						Missing:   "Test",
					})
				}
			}
			if len(data) > 0 {
				err = n.reconcileRepo.SaveAlert(data)
				if err != nil {
					log.Println("save alert error", err)
					return err
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (n reconcileJob) CheckAlertStatus() error {
	for {
		alert, err := n.reconcileRepo.GetAlertFail()
		if err != nil {
			log.Println("get alert error")
			return err
		}
		if len(alert) > 0 {
			s := fmt.Sprintf("Alert payment not reconcile (%v tnx)", len(alert))
			for i, rec := range alert {
				_ = rec
				s += fmt.Sprintf("\n\nrefID : %v\ninsurer : %v\npayment time : %v\nmissing : %v", alert[i].RefId, " ", alert[i].UpdatedAt, alert[i].Missing)
			}
			err = sendLineNotification(s)
			if err != nil {
				log.Println("fail to notify :", err)
				return err
			}
		}
		time.Sleep(10 * time.Second)
	}
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
