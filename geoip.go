package main

import (
	"encoding/json"
	"flag"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
)

type Response struct {
	Data  interface{} `json:",omitempty"`
	Error string      `json:",omitempty"`
}

func main() {
	var dbName = flag.String("db", "GeoLite2-City.mmdb", "File name of MaxMind GeoIP2 and GeoLite2 database")
	var lookup = flag.String("lookup", "city", "Specify which value to look up. Can be 'city' or 'country' depending on which database you load.")
	var listen = flag.String("listen", ":5000", "Listen address and port, for instance 127.0.0.1:5000")
	var threads = flag.Int("threads", runtime.NumCPU(), "Number of threads to use. Defaults to number of detected cores")
	var pretty = flag.Bool("pretty", false, "Should output be formatted with newlines and intentation")
	flag.Parse()

	runtime.GOMAXPROCS(*threads)
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
		// We don't need the body
		req.Body.Close()

		// Prepare the response and queue sending the result.
		res := &Response{}
		defer func() {
			var j []byte
			var err error
			if prettyL {
				j, err = json.MarshalIndent(res, "", "  ")
			} else {
				j, err = json.Marshal(res)
			}
			if err != nil {
				log.Fatal(err)
			}
			w.Write(j)
		}()

		ipText := req.URL.Query().Get("ip")
		if ipText == "" {
			ipText = strings.Trim(req.URL.Path, "/")
		}
		ip := net.ParseIP(ipText)
		if ip == nil {
			res.Error = "unable to decode ip"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if lookupCity {
			result, err := db.City(ip)
			if err != nil {
				res.Error = err.Error()
				return
			}
			res.Data = result
		} else {
			result, err := db.Country(ip)
			if err != nil {
				res.Error = err.Error()
				return
			}
			res.Data = result
		}
	})

	log.Println("Listening on " + *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
