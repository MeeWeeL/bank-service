package handlers

import (
	"bank-service/src/services"
	"encoding/json"
	"net/http"
	
	"github.com/sirupsen/logrus"
)

type CardHandler struct {
	cardService *services.CardService
	logger      *logrus.Logger
}

func NewCardHandler(service *services.CardService, logger *logrus.Logger) *CardHandler {
	return &CardHandler{
		cardService: service,
		logger:      logger,
	}
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value("userID")
	userID, ok := userIDValue.(uint)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "invalid user ID in context")
		return
	}
	
	var request struct {
		AccountID uint   `json:"account_id"`
		CVV       string `json:"cvv"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if request.AccountID == 0 {
		respondWithError(w, http.StatusBadRequest, "account ID must be provided")
		return
	}

	if len(request.CVV) < 3 || len(request.CVV) > 4 {
		respondWithError(w, http.StatusBadRequest, "CVV must be 3 or 4 digits")
		return
	}

	if !isDigits(request.CVV) {
		respondWithError(w, http.StatusBadRequest, "CVV must contain only digits")
		return
	}
	
	card, err := h.cardService.GenerateCard(userID, request.AccountID, request.CVV)
	if err != nil {
		h.logger.WithError(err).Error("failed to create card")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, card)
}

func isDigits(s string) bool {
    for _, r := range s {
        if r < '0' || r > '9' {
            return false
        }
    }
    return true
}

func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	
	cards, err := h.cardService.GetUserCards(userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get cards")
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, cards)
}