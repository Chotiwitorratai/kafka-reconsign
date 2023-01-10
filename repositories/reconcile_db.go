package repositories

import (
	"time"

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
	currentTime := time.Now()
	oneDay := currentTime.Add(-time.Hour * 24)
	err = obj.db.Table("tbl_purchase_alert").Limit(20).Where("status = 'Fail' AND created_at >= ?",oneDay).Find(&alert).Error
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
func (obj reconcileRepositoryDB) GetCountAlertFail() (int){
	return 5
}
