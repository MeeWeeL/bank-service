package handlers

import (
	"bank-service/src/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	accountService   *services.AccountService
	analyticsService *services.AnalyticsService
	logger          *logrus.Logger
}

func NewAccountHandler(accountService *services.AccountService, analyticsService *services.AnalyticsService, logger *logrus.Logger) *AccountHandler {
	return &AccountHandler{
		accountService:   accountService,
		analyticsService: analyticsService,
		logger:          logger,
	}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	
	account, err := h.accountService.CreateAccount(userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to create account")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, account)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.ParseUint(vars["accountId"], 10, 64)
	
	account, err := h.accountService.GetAccount(uint(accountID))
	if err != nil {
		h.logger.WithError(err).Error("account not found")
		respondWithError(w, http.StatusNotFound, "account not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, account)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    accountID, _ := strconv.ParseUint(vars["accountId"], 10, 64)
    
    var req struct {
        Amount float64 `json:"amount"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "invalid request")
        return
    }

    // Проверка владения счетом
    userID := r.Context().Value("userID").(uint)
    if _, err := h.accountService.GetByIDAndUser(uint(accountID), userID); err != nil {
        respondWithError(w, http.StatusForbidden, "access denied")
        return
    }

    if err := h.accountService.Deposit(uint(accountID), req.Amount); err != nil {
        h.logger.WithError(err).Error("deposit failed")
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"status": "deposit successful"})
}

func (h *AccountHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    accountID, _ := strconv.Atoi(vars["accountId"])
    days, _ := strconv.Atoi(r.URL.Query().Get("days"))
    
    if days < 1 || days > 365 {
        respondWithError(w, http.StatusBadRequest, "Invalid days parameter")
        return
    }

    balance, err := h.analyticsService.PredictBalance(uint(accountID), days)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, map[string]float64{"predicted_balance": balance})
}