package repositories

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type reconcileRepositoryDB struct {
	db *gorm.DB
}

func NewReconcileRepositoryDB(db *gorm.DB) ReconcileRepositoryDB {
	db.Table("tbl_purchase_reconcile").AutoMigrate(&Reconcile{})
	db.Table("tbl_purchase_alert").AutoMigrate(&Alert{})
	return reconcileRepositoryDB{db}
}

func (obj reconcileRepositoryDB) CheckNullReconcile(id string) (boo bool, err error) {
	reconciles := []Reconcile{}
	err = obj.db.Table("tbl_purchase_reconcile").Where("transaction_ref_id=?", id).Find(&reconciles).Error
	if len(reconciles) == 1 {
		return true, err
	} else {
		return false, err
	}
}

func (obj reconcileRepositoryDB) SaveReconcile(reconcile Reconcile) error {
	return obj.db.Table("tbl_purchase_reconcile").Create(&reconcile).Error
}

func (obj reconcileRepositoryDB) UpdateReconcile(reconcile Reconcile) error {
	return obj.db.Table("tbl_purchase_reconcile").Where("transaction_ref_id=?", reconcile.TransactionRefID).Updates(&reconcile).Error
}

func (obj reconcileRepositoryDB) GetReconcileFail() (reconcile []Reconcile, err error) {
	currentTime := time.Now()
	oneMin := currentTime.Add(-time.Minute * 1)
	err = obj.db.Table("tbl_purchase_reconcile").Where("( next_status = '' OR insurance_status = '' ) AND created_at <= ?", oneMin).Find(&reconcile).Error
	return reconcile, err
}

func (obj reconcileRepositoryDB) SaveAlert(alert []Alert) error {
	return obj.db.Table("tbl_purchase_alert").Create(&alert).Error
}

func (obj reconcileRepositoryDB) UpdateAlert(alert Alert) error {
	return obj.db.Table("tbl_purchase_alert").Where("ref_id=?", alert.RefId).Updates(&alert).Error
}

func (obj reconcileRepositoryDB) GetAlertFail() (alert []Alert, err error) {
	config()
	hour := viper.GetInt("repository.DayGetAlertFail") * 24
	currentTime := time.Now()
	hourTime := currentTime.Add(-time.Hour * time.Duration(hour))
	err = obj.db.Table("tbl_purchase_alert").Limit(20).Where("status = 'Fail' AND created_at >= ?", hourTime).Find(&alert).Error
	return alert, err
}

func (obj reconcileRepositoryDB) GetAlertFailByID(id string) (boo bool, err error) {
	alert := []Alert{}
	err = obj.db.Table("tbl_purchase_alert").Where("ref_id=?", id).Find(&alert).Error
	if len(alert) == 1 {
		return true, err
	} else {
		return false, err
	}

}
func (obj reconcileRepositoryDB) GetCountAlertFail() (count int64, err error) {
	config()
	var alert []Alert
	hour := viper.GetInt("repository.DayGetCountAlertFail") * 24
	currentTime := time.Now()
	hourTime := currentTime.Add(-time.Hour * time.Duration(hour))
	result := obj.db.Table("tbl_purchase_alert").Where("status = 'Fail' AND created_at >= ?", hourTime).Find(&alert)
	return result.RowsAffected, result.Error
}
func config() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("fatal error config file: %s", err)
	}
}
