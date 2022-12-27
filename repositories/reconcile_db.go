package repositories

import (
	"fmt"

	"gorm.io/gorm"
)

type reconcileRepositoryDB struct {
	db *gorm.DB
}

func NewReconcileRepositoryDB(db *gorm.DB) reconcileRepositoryDB {
	db.Table("tbl_purchase_reconcile").AutoMigrate(&Reconcile{})
	db.Table("tbl_purchase_alert").AutoMigrate(&Alert{})
	return reconcileRepositoryDB{db}
}

func (obj reconcileRepositoryDB) CheckNullReconcile(id string) (bool, err error) {
	reconciles := []Reconcile{}
	err = obj.db.Table("tbl_purchase_reconcile").Where("transaction_ref_id=?", id).Find(&reconciles).Error
	fmt.Println(reconciles)
	return nil, err
}

func (obj reconcileRepositoryDB) GetReconcileFail() (reconcile []Reconcile, err error) {
	err = obj.db.Where("status = '' || insurance_status = '' ").Find(&reconcile).Error
	return reconcile, err
}

func (obj reconcileRepositoryDB) GetAlertFail() (reconcile []Reconcile, err error) {
	err = obj.db.Where("status = 'Fail'").Find(&reconcile).Error
	return reconcile, err
}

func (obj reconcileRepositoryDB) SaveReconcile(reconcile Reconcile) error {
	return obj.db.Table("tbl_purchase_reconcile").Create(&reconcile).Error
}

func (obj reconcileRepositoryDB) SaveAlert(alert []Alert) error {
	return obj.db.Table("tbl_purchase_alert").Create(&alert).Error
}

func (obj reconcileRepositoryDB) UpdateReconcile(reconcile Reconcile) error {
	return obj.db.Table("tbl_purchase_reconcile").Save(&reconcile).Error
}

func (obj reconcileRepositoryDB) UpdateAlertStatus(alert Alert) error {
	obj.db.Table("tbl_purchase_alert").First(&alert)
	alert.Status = "succes"
	return obj.db.Table("tbl_purchase_alert").Save(&alert).Error
}
