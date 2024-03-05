package repository

import (
	"database/sql"
	"github.com/mashmorsik/L0/pkg/models"
)

type Repository interface {
	AddOrderTx(tx *sql.Tx, order models.Order) error
	AddDeliveryTx(tx *sql.Tx, o models.Order) error
	AddPaymentTx(tx *sql.Tx, o models.Order) error
	AddOrderItemsTx(tx *sql.Tx, o models.Order) error
	GetOrdersIDTx(tx *sql.Tx) ([]string, error)
	GetOrderInfo(tx *sql.Tx, orderID string, model models.Order) (*models.Order, error)
	GetDeliveryInfo(tx *sql.Tx, orderID string, model models.Order) (*models.Order, error)
	GetPaymentInfo(tx *sql.Tx, orderID string, model models.Order) (*models.Order, error)
	GetItemsInfo(tx *sql.Tx, orderID string, model models.Order) (*models.Order, error)
	CreateOrder(order models.Order) error
	GetOrders() ([]*models.Order, error)
}
