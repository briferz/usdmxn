package tokenmiddleware

import (
	"encoding/json"
	"github.com/briferz/usdmxn/middleware/tokenmiddleware/tokencreatorvalidator"
	"log"
	"net/http"
)

const tokenQuery = "token"

func WithValidator(v tokencreatorvalidator.Validator, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		token := request.URL.Query().Get(tokenQuery)
		if token == "" {
			writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(map[string]string{
				"error": "no api token was passed",
			})
			return
		}

		ok, err := v.Validate(token)
		if err != nil {
			log.Printf("error validating token %s: %s", token, err)
			writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(map[string]string{
				"error": "unable to validate the token",
			})
			return
		}

		if !ok {
			writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(map[string]string{
				"error": "the api token passed is invalid",
			})
			return
		}

		next(writer, request)
	}
}

func CreateToken(c tokencreatorvalidator.Creator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		newToken, err := c.Create()
		if err != nil {
			log.Printf("error creating token: %s", err)
			writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(map[string]string{
				"error": "unable to create the token",
			})
			return
		}
		json.NewEncoder(writer).Encode(map[string]string{
			"token": newToken,
		})

	}
}
