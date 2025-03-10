// HTTP обработчики с улучшенным логгированием и обработкой ошибок
package http

import (
	"context"
	"encoding/json"
	"errors"
	"flexible_transfer_backend/internal/core/domain"
	"flexible_transfer_backend/internal/core/usecase"
	"flexible_transfer_backend/internal/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type ExchangeHandler struct {
	uc     *usecase.ExchangeUseCase
	logger logger.Logger
}

func NewExchangeHandler(uc *usecase.ExchangeUseCase, l logger.Logger) *ExchangeHandler {
	return &ExchangeHandler{uc: uc, logger: l}
}

...