package accesslimiter

import (
	"encoding/json"
	"github.com/briferz/usdmxn/middleware/tokenmiddleware"
	"github.com/briferz/usdmxn/shared/limiter"
	"log"
	"net/http"
)

func WithAccessLimiter(l limiter.Interface, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get(tokenmiddleware.TokenQuery)

		allow, err := l.Allow(token)
		if err != nil {
			log.Printf("error checkoing for allowance of %s: %s", token, err)
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unable to validate API allowance",
			})
			return
		}

		if !allow {
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "the api has received too many requests for the given token",
			})
			return
		}
		next(w, r)
	}
}
