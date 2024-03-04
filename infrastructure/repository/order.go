package repository

import (
	"database/sql"
	"fmt"
	"github.com/mashmorsik/L0/infrastructure/data"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"time"
)

type OrderRepo struct {
	data *data.Data
}

func NewOrderRepo(data *data.Data) *OrderRepo {
	return &OrderRepo{data: data}
}

func (r *OrderRepo) AddOrderTx(tx *sql.Tx, order models.Order) error {
	sqlAddOrder := `
	INSERT INTO public.order(
	                  uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, 
	                         shard_key, sm_id, date_created, oof_shard
	                  )
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	res, err := tx.Exec(sqlAddOrder, order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	if err != nil {
		log.Errf("can't add order, err: %s", err)
		return err
	}
	ra, _ := res.RowsAffected()
	fmt.Printf("rows affected: %v", ra)

	return nil
}

func (r *OrderRepo) AddDeliveryTx(tx *sql.Tx, o models.Order) error {
	sqlAddDelivery := `
	INSERT INTO public.delivery(
	                  order_id, name, phone, zip, city, address, region, email
	                  )
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	res, err := tx.Exec(sqlAddDelivery, o.OrderUid, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		log.Errf("can't add delivery, err: %s", err)
		return err
	}
	ra, _ := res.RowsAffected()
	fmt.Printf("rows affected: %v", ra)

	return nil
}

func (r *OrderRepo) AddPaymentTx(tx *sql.Tx, o models.Order) error {
	sqlAddPayment := `
	INSERT INTO public.payment(
	                  transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost,
	                           goods_total, custom_fee
	                  )
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	res, err := tx.Exec(sqlAddPayment, o.OrderUid, o.Payment.RequestId, o.Payment.Currency, o.Payment.Provider,
		o.Payment.Amount, time.Unix(o.Payment.PaymentDt, 0), o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal,
		o.Payment.CustomFee)
	if err != nil {
		log.Errf("can't add payment, err: %s", err)
		return err
	}
	ra, _ := res.RowsAffected()
	fmt.Printf("rows affected: %v", ra)

	return nil
}

func (r *OrderRepo) AddOrderItemsTx(tx *sql.Tx, o models.Order) error {
	sqlAddOrderItems := `
	INSERT INTO public.order_item(
	                  order_id, chrt_id, track_number, price, rid, name, sale, size, count, total_price, nm_id, brand,
	                              status
	                  )
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	for _, i := range o.Items {
		res, err := tx.Exec(sqlAddOrderItems, o.OrderUid, i.ChrtId, i.TrackNumber, i.Price, i.Rid, i.Name, i.Sale,
			i.Size, i.Count, i.TotalPrice, i.NmId, i.Brand, i.Status)
		if err != nil {
			log.Errf("can't add order_items, err: %s", err)
			return err
		}
		ra, _ := res.RowsAffected()
		fmt.Printf("rows affected: %v", ra)
	}

	return nil
}

func (r *OrderRepo) GetOrders() ([]*models.Order, error) {
	var ordersID []string
	var orders []*models.Order

	sqlGetOrdersID := `
	SELECT uid 
	FROM public.order`

	rows, err := r.data.Master().Query(sqlGetOrdersID)
	if err != nil {
		log.Errf("can't get all orders, err: %s", err)
		return nil, err
	}

	for rows.Next() {
		var orderID string
		if err = rows.Scan(&orderID); err != nil {
			return nil, err
		}
		ordersID = append(ordersID, orderID)
	}

	for _, oID := range ordersID {
		var order models.Order
		sqlGetOrdersInfo := `
		SELECT uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard 
		FROM public.order 
		WHERE uid = $1`
		rows, err = r.data.Master().Query(sqlGetOrdersInfo, oID)
		if err != nil {
			log.Errf("can't get order info, err: %s", err)
			return nil, err
		}
		for rows.Next() {
			if err = rows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
				&order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard); err != nil {
				return nil, err
			}
		}
		sqlGetDeliveryInfo := `
		SELECT name, phone, zip, city, address, region, email 
		FROM public.delivery
		WHERE order_id = $1`
		rows, err = r.data.Master().Query(sqlGetDeliveryInfo, oID)
		if err != nil {
			log.Errf("can't get delivery info, err: %s", err)
			return nil, err
		}
		for rows.Next() {
			if err = rows.Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
				&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email); err != nil {
				return nil, err
			}
		}
		var paymentDt time.Time
		sqlGetPaymentInfo := `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, 
		       custom_fee 
		FROM public.payment 
		WHERE transaction = $1`
		rows, err = r.data.Master().Query(sqlGetPaymentInfo, oID)
		if err != nil {
			log.Errf("can't get payment info, err: %s", err)
			return nil, err
		}
		for rows.Next() {
			if err = rows.Scan(&order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency,
				&order.Payment.Provider, &order.Payment.Amount, &paymentDt, &order.Payment.Bank,
				&order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee); err != nil {
				return nil, err
			}
		}
		order.Payment.PaymentDt = paymentDt.Unix()

		sqlGetItemsInfo := `
		SELECT chrt_id, track_number, price, rid, name, sale, size, count, total_price, nm_id, brand, status 
		FROM public.order_item 
		WHERE order_id = $1`
		rows, err = r.data.Master().Query(sqlGetItemsInfo, oID)
		if err != nil {
			log.Errf("can't get items info, err: %s", err)
			return nil, err
		}
		for rows.Next() {
			var item models.Item
			if err = rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
				&item.Size, &item.Count, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status); err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *OrderRepo) CreateOrder(order models.Order) error {
	tx, err := r.data.Master().Begin()
	if err != nil {
		log.Errf("can't begin transaction, err: %s", err)
	}
	defer tx.Rollback()

	if err = r.AddOrderTx(tx, order); err != nil {
		log.Errf("can't add order in transaction, err: %s", err)
	}
	if err = r.AddDeliveryTx(tx, order); err != nil {
		log.Errf("can't add delivery in transaction, err: %s", err)
	}
	if err = r.AddPaymentTx(tx, order); err != nil {
		log.Errf("can't add payment in transaction, err: %s", err)
	}
	if err = r.AddOrderItemsTx(tx, order); err != nil {
		log.Errf("can't add order_items in transaction, err: %s", err)
	}

	if err = tx.Commit(); err != nil {
		log.Errf("can't commit transaction")
	}
	return nil
}
