package main

import (
	"encoding/json"
	"flag"
	"github.com/klauspost/geoip-service/geoip2"
	"github.com/pmylund/go-cache"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

//go:generate ffjson --nodecoder $GOFILE

// ffjson: nodecoder
type ResponseCity struct {
	Data  *geoip2.City `json:",omitempty"`
	Error string       `json:",omitempty"`
}

// ffjson: nodecoder
type ResponseCountry struct {
	Data  *geoip2.Country `json:",omitempty"`
	Error string          `json:",omitempty"`
}

func main() {
	var dbName = flag.String("db", "GeoLite2-City.mmdb", "File name of MaxMind GeoIP2 and GeoLite2 database")
	var lookup = flag.String("lookup", "city", "Specify which value to look up. Can be 'city' or 'country' depending on which database you load.")
	var listen = flag.String("listen", ":5000", "Listen address and port, for instance 127.0.0.1:5000")
	var threads = flag.Int("threads", runtime.NumCPU(), "Number of threads to use. Defaults to number of detected cores")
	var pretty = flag.Bool("pretty", false, "Should output be formatted with newlines and intentation")
	var cacheSecs = flag.Int("cache", 0, "How many seconds should requests be cached. Set to 0 to disable")
	var originPolicy = flag.String("origin", "*", `Value sent in the 'Access-Control-Allow-Origin' header. Set to "" to disable.`)
	serverStart := time.Now().Format(http.TimeFormat)

	flag.Parse()

	runtime.GOMAXPROCS(*threads)
	var memCache *cache.Cache
	if *cacheSecs > 0 {
		memCache = cache.New(time.Duration(*cacheSecs)*time.Second, 1*time.Second)
	}
	lookupCity := true
	if *lookup == "country" {
		lookupCity = false
	} else if *lookup != "city" {
		log.Fatalf("lookup parameter should be either 'city', or 'country', it is '%s'", *lookup)
	}

	db, err := geoip2.Open(*dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Loaded database " + *dbName)

	// We dereference this to avoid a pretty big penalty under heavy load.
	prettyL := *pretty

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var ipText string
		// Prepare the response and queue sending the result.
		var cached []byte
		var returnError string
		var result interface{} = nil

		defer func() {
			var j []byte
			var err error
			if cached != nil {
				j = cached
			} else {
				city, ok := result.(*geoip2.City)
				if ok {
					res := ResponseCity{Data: city, Error: returnError}
					if prettyL {
						j, err = json.MarshalIndent(res, "", "  ")
					} else {
						j, err = res.MarshalJSON()
					}
				} else {
					country, _ := result.(*geoip2.Country)
					res := ResponseCountry{Data: country, Error: returnError}
					if prettyL {
						j, err = json.MarshalIndent(res, "", "  ")
					} else {
						j, err = res.MarshalJSON()
					}
				}
				if err != nil {
					log.Fatal(err)
				}
			}
			if memCache != nil && cached == nil {
				memCache.Set(ipText, j, 0)
			}
			w.Write(j)
		}()

		// Set headers
		if *originPolicy != "" {
			w.Header().Set("Access-Control-Allow-Origin", *originPolicy)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", serverStart)

		ipText = req.URL.Query().Get("ip")
		if ipText == "" {
			ipText = strings.Trim(req.URL.Path, "/")
		}
		ip := net.ParseIP(ipText)
		if ip == nil {
			returnError = "unable to decode ip"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if memCache != nil {
			v, found := memCache.Get(ipText)
			if found {
				cached = v.([]byte)
				return
			}
		}
		if lookupCity {
			result, err = db.City(ip)
			if err != nil {
				returnError = err.Error()
				return
			}
		} else {
			result, err = db.Country(ip)
			if err != nil {
				returnError = err.Error()
				return
			}
		}
	})

	log.Println("Listening on " + *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
