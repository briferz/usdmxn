package fixer

import (
	"encoding/json"
	"fmt"
	"github.com/briferz/usdmxn/servers/controller/ratesource"
	"github.com/briferz/usdmxn/servers/model"
	"github.com/briferz/usdmxn/shared/cache"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	fixerBaseEndpoint = "http://data.fixer.io/api/latest"
	fixerQuerySymbols = "symbols=USD,MXN"
	fixerQueryFormat  = "format=1"

	fixerSourceCacheId = "fixer"
	fixerProviderName  = "FIXER"
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

func (s *fixerRateSource) obtainNewData() (responseModel, error) {
	resp, err := http.Get(s.requestUrl)
	if err != nil {
		return responseModel{}, fmt.Errorf("doing API request: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseModel{}, fmt.Errorf("reading request body: %s", err)
	}

	var rModel responseModel
	err = json.Unmarshal(body, &rModel)

	return rModel, err
}

func (s *fixerRateSource) requestModelToProvider(m responseModel) (p model.Provider, err error) {
	if !m.Success {
		err = fmt.Errorf("api response rated as unsuccessful")
		return
	}
	mxnRate, ok := m.Rates["MXN"]
	if !ok {
		err = fmt.Errorf("api response did not contain MXN rate")
		return
	}

	usdRate, ok := m.Rates["USD"]
	if !ok {
		err = fmt.Errorf("api response did not contain USD rate")
		return
	}

	if usdRate == 0 {
		err = fmt.Errorf("USD rate was 0, so it's invalid")
		return
	}

	p.Value = mxnRate / usdRate
	p.LastUpdated = time.Now()
	p.ProviderName = fixerProviderName
	return
}

func (s *fixerRateSource) saveInCache(p model.Provider) error {
	newData, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshaling new provider data into json bytes: %s", err)
	}

	err = s.cache.Set(s.CacheId(), newData)
	if err != nil {
		return fmt.Errorf("setting new data into cache provider: %s", err)
	}
	return nil
}

func (s *fixerRateSource) update() error {
	rModel, err := s.obtainNewData()
	if err != nil {
		return err
	}

	p, err := s.requestModelToProvider(rModel)
	if err != nil {
		return err
	}

	return s.saveInCache(p)

}

func (s *fixerRateSource) runUpdates(period time.Duration) error {
	err := s.update()
	if err != nil {
		return fmt.Errorf("on first update: %s", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		for range time.NewTicker(period).C {
			err = s.update()
			if err != nil {
				log.Printf("unable to update rate source: %s", err)
			}
		}
	}()
	wg.Wait()
	return nil
}

func New(cache cache.Interface, updatePeriod time.Duration) (ratesource.Interface, error) {

	if cache == nil || updatePeriod <= 0 {
		log.Panicf("bad parameters: %v / %v", cache, updatePeriod)
	}

	apiKey, ok := fixerAPIKey()
	if !ok {
		return nil, fmt.Errorf("the environment variable API key required for Fixer Service is not set")
	}

	s := &fixerRateSource{
		cache:      cache,
		requestUrl: fixerEndpoint(apiKey),
	}

	err := s.runUpdates(updatePeriod)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *fixerRateSource) GetSource() (p model.Provider, err error) {

	defer func() {
		p.ProviderName = fixerProviderName
		if err != nil {
			p.Error = err.Error()
		}
	}()

	bytes, err := s.cache.Get(s.CacheId())
	if err != nil {
		err = fmt.Errorf("querying cache for id=%s: %s", s.CacheId(), err)
		p.Error = err.Error()
		return
	}

	err = json.Unmarshal(bytes, &p)
	return
}
