package repositories

import (
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
	err = obj.db.Table("tbl_purchase_reconcile").Where("status = '' || insurance_status = '' ").Find(&reconcile).Error
	return reconcile, err
}

func (obj reconcileRepositoryDB) SaveAlert(alert []Alert) error {
	return obj.db.Table("tbl_purchase_alert").Create(&alert).Error
}

func (obj reconcileRepositoryDB) UpdateAlert(alert Alert) error {
	return obj.db.Table("tbl_purchase_alert").Where("ref_id=?", alert.RefId).Updates(&alert).Error
}

func (obj reconcileRepositoryDB) GetAlertFail() (alert []Alert, err error) {
	err = obj.db.Table("tbl_purchase_alert").Where("status = 'Fail'").Find(&alert).Error
	return alert, err
}
