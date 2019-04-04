package main

import (
	"flag"
	"github.com/briferz/usdmxn/servers/controller"
	"github.com/briferz/usdmxn/servers/controller/ratesource/banxico"
	"github.com/briferz/usdmxn/servers/controller/ratesource/fixer"
	"github.com/briferz/usdmxn/shared/cache/rediscache"
	"github.com/briferz/usdmxn/shared/redis"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	bindAddr := flag.String("-bind", ":8080", "specifies the listen address for this server")
	flag.Parse()

	client, err := redis.Client()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("success reaching Redis.")

	cache := rediscache.New(client)

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
	mux.Handle("/", exchangeRateController)

	log.Print("Listening...")
	err = http.ListenAndServe(*bindAddr, mux)
	if err != nil {
		log.Printf("unable to bind server to address %s", *bindAddr)
	}

}
