package controller

import (
	"encoding/json"
	"github.com/briferz/usdmxn/servers/controller/ratesource"
	"github.com/briferz/usdmxn/servers/model"
	"log"
	"net/http"
)

var errorProvider = model.Provider{}

type ratesEndpoint struct {
	sources []ratesource.Interface
}

func New(sources ...ratesource.Interface) http.Handler {
	return &ratesEndpoint{
		sources: sources,
	}
}

func (e *ratesEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var providers []model.Provider
	for i, source := range e.sources {
		provider, err := source.GetSource()
		if err != nil {
			log.Printf("error getting source at index %d: %s", i, err)
		}
		providers = append(providers, provider)
	}

	w.Header().Add("Content-Type", "application/json; Charset=UTF-8")
	json.NewEncoder(w).Encode(model.Rates{Rates: providers})

}
