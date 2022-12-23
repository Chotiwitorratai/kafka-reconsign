package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Order struct {
	RefId           string `db:"ref_id"`
	NextStatus      string `db:"next_status"`
	NextUpdate      string `db:"next_update"`
	IpCallBack      string `db:"ip_call_back"`
	IpUpdate        string `db:"ip_update"`
	ReconcileStatus string `db:"reconcile_status"`
}

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepositoryDB(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) SaveOrder(order Order) (string, error) {
	query := "INSERT INTO [repo].[dbo].[order]([ref_id],[next_status],[next_update],[ip_call_back],[ip_update],[reconcile_status]) VALUES (?,?,?,?,?,?)"
	insertResult, err := r.db.ExecContext(context.Background(), query, order.RefId, order.NextStatus, order.NextUpdate, order.IpCallBack, order.IpUpdate, order.ReconcileStatus)
	print(insertResult)
	if err != nil {
		return "fail", err
	}
	return "success", nil
}
func (r *OrderRepository) GetOrderByRef_id(id string) (*Order, error) {
	order := Order{}
	query := "SELECT * FROM order WHERE ref_id = ?;"
	_, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
func (r *OrderRepository) ChangeStatus(id string, status string) (string, error) {
	query := "UPDATE order SET reconcile_status = ? WHERE id = ?;`"
	_, err := r.db.Query(query, id, status)
	if err != nil {
		return "fail", err
	}
	return "success", nil
}
