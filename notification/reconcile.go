package notification

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"kafka-reconsign/repositories"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	lineNotifyAPI = "https://notify-api.line.me/api/notify"
	lineToken     = "HD5RBYffvJAao8Qxm3mX7Zg1CpSQrEGfdVWlumFipuY"
)

var reconcileTime time.Duration
var notificationTime time.Duration
var startHour int
var endHour int

type ReconcileJob interface {
	CheckReconcileStatus() error
	CheckAlertStatus() error
}

type reconcileJob struct {
	reconcileRepo repositories.ReconcileRepositoryDB
}

func NewReconcileJob(reconcileRepo repositories.ReconcileRepositoryDB) ReconcileJob {
	initConfig()
	return reconcileJob{reconcileRepo: reconcileRepo}
}

func (n reconcileJob) CheckReconcileStatus() error {
	for {
		newID := false
		data := []repositories.Alert{}
		alertID := []string{}

		reconcile, err := getReconcileFailures(n)
		if err != nil {
			log.Println("get reconcile error", err)
			return err
		}
		alert, err := getAlertFailures(n)
		if err != nil {
			return err
		}

		for _, rec := range reconcile {
			alertID = append(alertID, rec.TransactionRefID)
			existRefID, err := n.reconcileRepo.GetAlertFailByID(rec.TransactionRefID)
			if err != nil {
				return err
			}
			if !existRefID {
				newID = true
				missing := ""
				if rec.NextStatus == "Success" {
					missing = "insuranceCallback"
				} else {
					missing = "paymentNext"
				}

				data = append(data, repositories.Alert{
					Messages:    "test",
					Count:       0,
					Status:      "Fail",
					NextAlert:   time.Now().Add(5 * time.Minute),
					RefId:       rec.TransactionRefID,
					Missing:     missing,
					Insurer:     rec.PartnerInfoName,
					PaymentTime: rec.TransactionCreatedTimestamp,
				})
			}
		}
		if newID {
			err = n.reconcileRepo.SaveAlert(data)
			if err != nil {
				fmt.Println(data)
				return errors.New("save alert error")
			}
		}

		for _, alrt := range alert {
			if !contains(alertID, alrt.RefId) {
				_, err = updateAlertStatus(n, alrt.RefId, "success", 0)
				if err != nil {
					return err
				}
			}

		}
		time.Sleep(reconcileTime * time.Second)
	}
}

func (n reconcileJob) CheckAlertStatus() error {
	for {
		alert, err := getAlertFailures(n)
		if err != nil {
			return err
		}
		failcount, err := n.reconcileRepo.GetCountAlertFail()
		if err != nil {
			return err
		}

		if len(alert) > 0 {
			s := fmt.Sprintf("Alert payment not reconcile (%v tnx)", failcount)
			for i, a := range alert {
				str, err := updateAlertStatus(n, a.RefId, "Fail", a.Count+1)
				if err != nil {
					return err
				}
				if i <= 2 {
					t := a.PaymentTime.Format("2006-01-02 15:04:05")
					str += fmt.Sprintf("\npayment time : %v\nmissing : %v", t, a.Missing)
					s += str
				}
			}
			if workingHour() {
				log.Println("start notification")
				err = sendLineNotification(s)
				if err != nil {
					log.Println("fail to notify :", err)
					return err
				}
			}
		}
		time.Sleep(notificationTime * time.Second)
	}
}

func initConfig() {
	reconcileTime = viper.GetDuration("checkReconcileInterval")
	notificationTime = viper.GetDuration("notificationInterval")
	startHour = viper.GetInt("schedule.startHour")
	endHour = viper.GetInt("schedule.endHour")
}

func workingHour() bool {
	now := time.Now()
	currentHour := now.Hour()
	if currentHour >= startHour && currentHour < endHour {
		return true
	}
	return false
}

func contains(s []string, valueToCheck string) bool {
	for _, v := range s {
		if v == valueToCheck {
			return true
		}
	}
	return false
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

func sendLineNotification(message string) error {
	req, err := http.NewRequest("POST", lineNotifyAPI, bytes.NewBuffer([]byte("message="+message)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+lineToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
