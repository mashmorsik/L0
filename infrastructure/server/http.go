package server

import (
	"encoding/json"
	"fmt"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"net/http"
)

type HTTPServer struct {
	o order.CreateOrder
}

func NewServer(o order.CreateOrder) *HTTPServer {
	return &HTTPServer{o: o}
}

func (s *HTTPServer) StartServer() {
	http.HandleFunc("/order", s.getOrderInfoHandler)

	fmt.Println("HTTPServer is listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errf("server can't serve")
		return
	}
}

func (s *HTTPServer) getOrderInfoHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")

	cachedOrder, err := s.o.GetOrderFromCache(orderID)
	if err != nil {
		log.Infof("no order in cache, orderID: %s", orderID)
		return
	}

	orderJSON, err := json.Marshal(cachedOrder)
	if err != nil {
		http.Error(w, "Failed to marshal order information to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(orderJSON)
}
