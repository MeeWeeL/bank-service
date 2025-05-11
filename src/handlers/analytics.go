package handlers

import (
	"bank-service/src/services"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logrus.Logger
}

func NewAnalyticsHandler(service *services.AnalyticsService, logger *logrus.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: service,
		logger:           logger,
	}
}

func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	// Для примера берем текущий месяц
	year := 0
	month := 0

	query := r.URL.Query()
	if y := query.Get("year"); y != "" {
		year, _ = strconv.Atoi(y)
	}
	if m := query.Get("month"); m != "" {
		month, _ = strconv.Atoi(m)
	}

	if year == 0 || month == 0 {
		// По умолчанию текущий месяц
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	income, expenses, err := h.analyticsService.GetMonthlyIncomeExpenses(userID, year, time.Month(month))
	if err != nil {
		h.logger.WithError(err).Error("failed to get analytics")
		respondWithError(w, http.StatusInternalServerError, "failed to get analytics")
		return
	}

	creditLoad, err := h.analyticsService.GetCreditLoad(userID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get credit load")
		respondWithError(w, http.StatusInternalServerError, "failed to get credit load")
		return
	}

	response := map[string]interface{}{
		"income":      income,
		"expenses":    expenses,
		"credit_load": creditLoad,
	}

	respondWithJSON(w, http.StatusOK, response)
}
