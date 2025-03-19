package main

import (
	"log"
	"net/http"

	"github.com/RCNRC/flexible_transfer_backend/internal/config"
	"github.com/RCNRC/flexible_transfer_backend/internal/delivery/http/handler"
	"github.com/RCNRC/flexible_transfer_backend/internal/pkg/logger"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	logger := logger.NewZapLogger(cfg.Environment)

	router := mux.NewRouter()
	handler.NewExchangeHandler(nil, logger).RegisterRoutes(router)

	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}
