package handlers

import (
	"bank-service/src/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CreditHandler struct {
	creditService *services.CreditService
	logger        *logrus.Logger
}

func NewCreditHandler(creditService *services.CreditService, logger *logrus.Logger) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
		logger:        logger,
	}
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	var req struct {
		AccountID uint    `json:"account_id"`
		Amount    float64 `json:"amount"`
		Rate      float64 `json:"rate"`
		Period    int     `json:"period"` // в месяцах
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.AccountID == 0 || req.Amount <= 0 || req.Rate <= 0 || req.Period <= 0 {
		respondWithError(w, http.StatusBadRequest, "missing or invalid fields")
		return
	}

	credit, err := h.creditService.CreateCredit(userID, req.AccountID, req.Amount, req.Rate, req.Period)
	if err != nil {
		h.logger.WithError(err).Error("failed to create credit")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, credit)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	creditID, err := strconv.ParseUint(vars["creditId"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid credit ID")
		return
	}

	schedule, err := h.creditService.GetPaymentSchedule(userID, uint(creditID))
	if err != nil {
		h.logger.WithError(err).Error("failed to get payment schedule")
		respondWithError(w, http.StatusNotFound, "credit or schedule not found")
		return
	}

	respondWithJSON(w, http.StatusOK, schedule)
}

func (h *CreditHandler) GetCreditsByAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	accountID, err := strconv.ParseUint(vars["accountId"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid account ID")
		return
	}

	// Проверим, что аккаунт принадлежит пользователю
	account, err := h.creditService.AccountService().GetByIDAndUser(uint(accountID), userID)
	if err != nil || account == nil {
		respondWithError(w, http.StatusForbidden, "access denied")
		return
	}

	credits, err := h.creditService.GetCreditsByAccount(uint(accountID))
	if err != nil {
		h.logger.WithError(err).Error("failed to get credits by account")
		respondWithError(w, http.StatusInternalServerError, "failed to get credits")
		return
	}

	respondWithJSON(w, http.StatusOK, credits)
}
