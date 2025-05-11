package middleware

import (
	"encoding/json"
	"bank-service/src/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func AccountOwnershipMiddleware(repo *repositories.AccountRepository, logger *logrus.Logger) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID := r.Context().Value("userID").(uint)
            accountID, _ := strconv.Atoi(mux.Vars(r)["accountId"])
            
            account, err := repo.GetByIDAndUser(uint(accountID), userID)
            if err != nil || account == nil {
                logger.Warnf("Unauthorized access attempt to account %d by user %d", accountID, userID)
                respondWithError(w, http.StatusForbidden, "Access denied")
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Вспомогательные функции ответов
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}