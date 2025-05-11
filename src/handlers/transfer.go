package handlers

import (
    "encoding/json"
    "net/http"
    
    "bank-service/src/services"
    "github.com/sirupsen/logrus"
)

type TransferHandler struct {
    accountService *services.AccountService
    logger         *logrus.Logger
}

func NewTransferHandler(service *services.AccountService, logger *logrus.Logger) *TransferHandler {
    return &TransferHandler{
        accountService: service,
        logger:         logger,
    }
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uint)
    
    var req struct {
        FromAccountID uint    `json:"from_account_id"`
        ToAccountID   uint    `json:"to_account_id"`
        Amount        float64 `json:"amount"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.WithError(err).Error("Failed to decode transfer request")
        respondWithError(w, http.StatusBadRequest, "invalid request format")
        return
    }

    if req.FromAccountID == 0 || req.ToAccountID == 0 {
        respondWithError(w, http.StatusBadRequest, "account IDs must be provided")
        return
    }

    // Проверка владения счетом
    if _, err := h.accountService.GetByIDAndUser(req.FromAccountID, userID); err != nil {
        respondWithError(w, http.StatusForbidden, "access denied")
        return
    }

    if req.Amount <= 0 {
        respondWithError(w, http.StatusBadRequest, "amount must be positive")
        return
    }
	
    if err := h.accountService.Transfer(req.FromAccountID, req.ToAccountID, req.Amount); err != nil {
        h.logger.WithError(err).Error("transfer failed")
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}