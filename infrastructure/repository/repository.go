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
	CreateOrder(order models.Order) error
}
