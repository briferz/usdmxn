package banxico

import (
	"encoding/json"
	"fmt"
	"github.com/briferz/usdmxn/servers/controller/ratesource"
	"github.com/briferz/usdmxn/servers/model"
	"github.com/briferz/usdmxn/shared/cache"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	banxicoBaseEndpoint = "https://www.banxico.org.mx/SieAPIRest/service/v1/series/SF43718/datos/oportuno"

	banxicoSourceCacheId = "banxico"
	banxicoProviderName  = "BANXICO"
)

func banxicoQueryAPIKey(apiKey string) string {
	return fmt.Sprintf("token=%s", apiKey)
}

func banxicoEndpoint(apiKey string) string {
	return banxicoBaseEndpoint + "?" + banxicoQueryAPIKey(apiKey)
}

type banxicoRateSource struct {
	cache      cache.Interface
	requestUrl string
}

func (s *banxicoRateSource) CacheId() string {
	return banxicoSourceCacheId
}

func (s *banxicoRateSource) obtainNewData() (responseModel, error) {
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

func (s *banxicoRateSource) requestModelToProvider(m responseModel) (p model.Provider, err error) {

	if len(m.BMX.Series) == 0 {
		err = fmt.Errorf("api response did not contain any series")
		return
	}

	datos := m.BMX.Series[0].Datos
	if len(datos) == 0 {
		err = fmt.Errorf("api response did not contain any data items")
		return
	}

	strValue := datos[0].Dato

	value, err := strconv.ParseFloat(strValue, 32)
	if err != nil {
		err = fmt.Errorf("error converting dato field from string to float: %s", err)
		return
	}

	p.Value = float32(value)
	p.LastUpdated = time.Now()
	p.ProviderName = banxicoProviderName
	return
}

func (s *banxicoRateSource) saveInCache(p model.Provider) error {
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

func (s *banxicoRateSource) update() error {
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

func (s *banxicoRateSource) runUpdates(period time.Duration) error {
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

	apiKey, ok := banxicoAPIKey()
	if !ok {
		return nil, fmt.Errorf("the environment variable API key required for BANXICO Service is not set")
	}

	s := &banxicoRateSource{
		cache:      cache,
		requestUrl: banxicoEndpoint(apiKey),
	}

	err := s.runUpdates(updatePeriod)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *banxicoRateSource) GetSource() (p model.Provider, err error) {

	defer func() {
		p.ProviderName = banxicoProviderName
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
