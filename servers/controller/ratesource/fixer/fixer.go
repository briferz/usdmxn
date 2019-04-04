package fixer

import (
	"encoding/json"
	"fmt"
	"github.com/briferz/usdmxn/servers/controller/ratesource"
	"github.com/briferz/usdmxn/servers/controller/ratesource/fixer/env"
	"github.com/briferz/usdmxn/servers/model"
	"github.com/briferz/usdmxn/shared/cache"
	"log"
	"strings"
	"time"
)

const (
	fixerBaseEndpoint = "http://data.fixer.io/api/latest"
	fixerQuerySymbols = "symbols=USD,MXN"
	fixerQueryFormat  = "format=1"

	fixerSourceCacheId = "fixer"
)

func fixerQueryAPIKey(apiKey string) string {
	return fmt.Sprintf("access_key=%s", apiKey)
}

func fixerEndpoint(apiKey string) string {
	return fixerBaseEndpoint + "?" + strings.Join([]string{fixerQueryAPIKey(apiKey), fixerQuerySymbols, fixerQueryFormat}, "&")
}

type fixerRateSource struct {
	cache      cache.Interface
	requestUrl string
}

func (s *fixerRateSource) CacheId() string {
	return fixerSourceCacheId
}

func (s *fixerRateSource)update()error{
	// ToDo implementation
	return nil
}

func New(cache cache.Interface, updatePeriod time.Duration) (ratesource.Interface, error) {

	if cache == nil || updatePeriod <= 0 {
		log.Panicf("bad parameters: %v / %v", cache, updatePeriod)
	}

	apiKey, ok := env.FixerAPIKey()
	if !ok {
		return nil, fmt.Errorf("the environment variable API key required for Fixer Service is not set")
	}

	//ToDo: spin for periodic updates

	return &fixerRateSource{
		cache:      cache,
		requestUrl: fixerEndpoint(apiKey),
	}, nil

}

func (s *fixerRateSource) GetSource() (p model.Provider, err error) {

	bytes, err := s.cache.Get(s.CacheId())
	if err != nil {
		err = fmt.Errorf("querying cache for id=%s: %s", s.CacheId(), err)
		return
	}

	err = json.Unmarshal(bytes, &p)
	return
}
