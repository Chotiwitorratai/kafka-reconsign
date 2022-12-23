package repositories

import (
	"gorm.io/gorm"
)

type reconcileRepositoryDB struct {
	db *gorm.DB
}

func NewReconcileRepositoryDB(db *gorm.DB) reconcileRepositoryDB {
	db.Table("payment_callback").AutoMigrate(&PaymentCallback{})
	db.Table("insuler_callback").AutoMigrate(&InsulerCallback{})
	db.Table("alert").AutoMigrate(&Alert{})
	return reconcileRepositoryDB{db}
}

func (obj reconcileRepositoryDB) SavePaymentCallback(paymentCallback PaymentCallback) error {
	return obj.db.Save(paymentCallback).Error
}
