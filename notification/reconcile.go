package notification

import (
	"errors"
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
		reconcile, err := getReconcileFailures(n)
		if err != nil {
			log.Println("get reconcile error", err)
			return err
		}

		for _, rec := range reconcile {
			foundID, err := n.reconcileRepo.GetAlertFailByID(rec.TransactionRefID)
			if err != nil {
				return err
			}
			if !foundID {
				if rec.Status == "Success" && rec.InsuranceStatus == "Success" {
					_, err = updateAlertStatus(n, rec.TransactionRefID, rec.Status, 0)
					if err != nil {
						return err
					}
				} else {
					var missing string
					if rec.Status == "Success" {
						missing = "paymentCallback"
					} else {
						missing = "insuranceCallback"
					}
					err = SaveAlert(n, rec.TransactionRefID, missing)
					if err != nil {
						return err
					}
				}
			}
		}
		time.Sleep(60 * time.Second)
	}
}

func (n reconcileJob) CheckAlertStatus() error {
	for {
		alert, err := getAlertFailures(n)
		if err != nil {
			return err
		}
		if len(alert) > 0 {
			s := fmt.Sprintf("Alert payment not reconcile (%v tnx)", len(alert))
			for _, rec := range alert {
				str, err := updateAlertStatus(n, rec.RefId, "Fail", rec.Count+1)
				if err != nil {
					return err
				}
				t := rec.CreatedAt.Format("2006-01-02 15:04:05")
				str += fmt.Sprintf("\npayment time : %v\nmissing : %v", t, rec.Missing)
				s += str
			}
			err = sendLineNotification(s)
			if err != nil {
				log.Println("fail to notify :", err)
				return err
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func getReconcileFailures(n reconcileJob) ([]repositories.Reconcile, error) {
	reconcile, err := n.reconcileRepo.GetReconcileFail()
	if err != nil {
		return nil, errors.New("get reconcile failures error")
	}
	return reconcile, nil
}

func getAlertFailures(n reconcileJob) ([]repositories.Alert, error) {
	alert, err := n.reconcileRepo.GetAlertFail()
	if err != nil {
		return nil, errors.New("get alert failures error")
	}
	return alert, nil
}

func getAlertExists(n reconcileJob, id string) (bool, error) {
	foundID, err := n.reconcileRepo.GetAlertFailByID(id)
	if err != nil {
		return false, errors.New("get alert exist error")
	}
	return foundID, nil
}

func updateAlertStatus(n reconcileJob, id string, status string, count int) (string, error) {
	s := repositories.Alert{
		Status: status,
		RefId:  id,
	}
	err := n.reconcileRepo.UpdateAlert(s)

	str := fmt.Sprintf("\n\nrefID : %v\ninsurer : %v", s.RefId, " ")
	if err != nil {
		return "", errors.New("update alert status error")
	}
	return str, nil
}

func SaveAlert(n reconcileJob, id string, missing string) error {
	data := []repositories.Alert{
		{
			Messages:  "test",
			Count:     0,
			Status:    "Fail",
			NextAlert: time.Now().Add(5 * time.Minute),
			RefId:     id,
			Missing:   missing,
		},
	}
	err := n.reconcileRepo.SaveAlert(data)
	if err != nil {
		return errors.New("save alert error")
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
