package server

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/web"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

type HTTPServer struct {
	o order.CreateOrder
}

func NewServer(o order.CreateOrder) *HTTPServer {
	return &HTTPServer{o: o}
}

var port, _ = os.LookupEnv("HTTP_SERVER_PORT")

func (s *HTTPServer) StartServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/order", s.getOrderInfoHandler)
	mux.HandleFunc("/", s.defaultHandler)

	log.Infof("HTTPServer is listening on port: %s\n", port)

	// FIXME: default route handler
	err := http.ListenAndServe(port, mux)
	if err != nil {
		return errors.WithMessagef(err, "server can't ListenAndServe http requests")
	}

	return nil
}

func (s *HTTPServer) defaultHandler(w http.ResponseWriter, r *http.Request) {
	html := web.DefaultDisplay()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err := w.Write([]byte(html))
	if err != nil {
		log.Errf("failed to write HTML response: %s", err)
	}
}

func (s *HTTPServer) getOrderInfoHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")

	cachedOrder, err := s.o.GetOrderFromCache(orderID)
	if err != nil {
		log.Infof("no order in cache, orderID: %s", orderID)
		return
	}

	html := web.DisplayOrder(*cachedOrder)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err = w.Write([]byte(html))
	if err != nil {
		log.Errf("failed to write HTML response: %s", err)
	}
}
