package main

import (
	"flag"
	"github.com/briferz/usdmxn/middleware/accesslimiter"
	"github.com/briferz/usdmxn/middleware/tokenmiddleware"
	"github.com/briferz/usdmxn/middleware/tokenmiddleware/tokencreatorvalidator/redistokencreatorvalidator"
	"github.com/briferz/usdmxn/servers/controller"
	"github.com/briferz/usdmxn/servers/controller/ratesource/banxico"
	"github.com/briferz/usdmxn/servers/controller/ratesource/fixer"
	"github.com/briferz/usdmxn/shared/cache/rediscache"
	"github.com/briferz/usdmxn/shared/limiter/redislimiter"
	"github.com/briferz/usdmxn/shared/redis"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	bindAddr := flag.String("-bind", ":8080", "specifies the listen address for this server")
	maxAllowances := flag.Int("-allowances", 1, "the amount of allowances per period of time")
	allowancesPeriod := flag.Duration("-period", time.Second, "the time lapse in which the allowances counter is reset")
	flag.Parse()

	client, err := redis.Client()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("success reaching Redis.")

	cache := rediscache.New(client)
	validatorCreator := redistokencreatorvalidator.New(client)
	limitEnforcer, limiterErrStream := redislimiter.New(client, *allowancesPeriod, int64(*maxAllowances))
	go func() {
		for e := range limiterErrStream {
			log.Printf("error with limiter enforcer: %s", e)
		}
	}()

	fixerSource, err := fixer.New(cache, 1*time.Hour)
	if err != nil {
		log.Printf("when creating fixer exchange rate source: %s", err)
		os.Exit(2)
	}
	log.Print("FIXER data source created successfully")

	banxicoSource, err := banxico.New(cache, 10*time.Second)
	if err != nil {
		log.Printf("when creating banxico exchange rate source: %s", err)
		os.Exit(3)
	}
	log.Print("BANXICO data source created successfully")

	exchangeRateController := controller.New(fixerSource, banxicoSource)

	mux := http.NewServeMux()
	mux.HandleFunc("/token", tokenmiddleware.CreateToken(validatorCreator))
	mux.HandleFunc("/", tokenmiddleware.WithValidator(validatorCreator, accesslimiter.WithAccessLimiter(limitEnforcer, exchangeRateController.ServeHTTP)))

	log.Print("Listening...")
	err = http.ListenAndServe(*bindAddr, mux)
	if err != nil {
		log.Printf("unable to bind server to address %s", *bindAddr)
	}

}
