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
	oneMiuteAgo := time.Now().Add(-time.Minute * 1)
	err = obj.db.Table("tbl_purchase_reconcile").Where("( next_status = '' OR insurance_status = '' ) AND created_at <= ?", oneMiuteAgo).Find(&reconcile).Error
	return reconcile, err
}

func (obj reconcileRepositoryDB) SaveAlert(alert []Alert) error {
	return obj.db.Table("tbl_purchase_alert").Create(&alert).Error
}

func (obj reconcileRepositoryDB) UpdateAlert(alert Alert) error {
	return obj.db.Table("tbl_purchase_alert").Where("ref_id=?", alert.RefId).Updates(&alert).Error
}

func (obj reconcileRepositoryDB) GetAlertFail() (alert []Alert, err error) {
	CallConfig()
	hourTimeAgo := time.Now().Add(-time.Hour * time.Duration(viper.GetInt("repository.DayGetAlertFail") * 24))
	err = obj.db.Order("created_at desc").Table("tbl_purchase_alert").Limit(20).Where("status = 'Fail' AND created_at >= ?", hourTimeAgo).Find(&alert).Error
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
	CallConfig()
	var alert []Alert
	hourTimeAgo := time.Now().Add(-time.Hour * time.Duration(viper.GetInt("repository.DayGetCountAlertFail") * 24))
	result := obj.db.Table("tbl_purchase_alert").Where("status = 'Fail' AND created_at >= ?", hourTimeAgo).Find(&alert)
	return result.RowsAffected, result.Error
}
func CallConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("fatal error config file: %s", err)
	}
}
