package server

import (
	"fmt"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	_interface "github.com/mashmorsik/L0/web"
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

	html := _interface.DisplayOrder(*cachedOrder)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err = w.Write([]byte(html))
	if err != nil {
		log.Errf("failed to write HTML response: %s", err)
	}
}
