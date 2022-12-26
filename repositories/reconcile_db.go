package repositories

import (
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

func (obj reconcileRepositoryDB) GetReconcileByRefID(id string) (reconciles []Reconcile, err error) {
	err = obj.db.Where("transaction_ref_id=?", id).Find(&reconciles).Error
	return reconciles, err
}

func (obj reconcileRepositoryDB) SaveReconcile(reconcile Reconcile) error {
	return obj.db.Create(&Reconcile{}).Error
}

func (obj reconcileRepositoryDB) SaveAlert(alert Alert) error {
	return obj.db.Create(&Alert{}).Error
}
