package main

import (
	"log"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/willf/bloom"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var fp = .001                                                                        // False positive threshhold (1 out of 1000)
var n = 572611621.0                                                                  // Hash table size (passwords count)
var m = math.Ceil((n * math.Log(fp)) / math.Log(1.0 / math.Pow(2.0, math.Log(2.0)))) // Number of bits in the filter  m = n * log(0.001) / log(1) / 2^log(2)
var k = uint(10)                                                                     // Number of hash functions
var filterSHA = bloom.New(uint(m), k)

// At this moment only one filter is supported. Since filter should be loaded enterily to memory, it is hard to keep both SHA-1 and NTLM filter at the same time.
// That is why only one filter is supported now. 
// TODO: add CLI options to choose and specify runtime options.

//var filterNTLM = bloom.New(uint(m), k)
var hex = "0123456789ABCDEF"


func hexOnly(hash string) bool {
	for _, c := range hash {
		if !strings.Contains(hex, string(c)) {
			return false
		}
	}
	return true
}


func checkSHA(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Access-Control-Allow-Methods", "HEAD")
	w.Header().Set("Content-Security-Policy", "default-src 'self';")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "same-origin")

	vars := mux.Vars(r)
	hash := strings.ToUpper(vars["hash"])

	if len(hash) != 16 || !hexOnly(hash) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else if filterSHA.Test([]byte(hash)) {
		http.Error(w, "PWNED", http.StatusOK)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}


func checkNTLM(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Access-Control-Allow-Methods", "HEAD")
        w.Header().Set("Content-Security-Policy", "default-src 'self';")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("Referrer-Policy", "same-origin")

        vars := mux.Vars(r)
        hash := strings.ToUpper(vars["hash"])

        if len(hash) != 16 || !hexOnly(hash) {
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        } else if filterSHA.Test([]byte(hash)) {
                http.Error(w, "PWNED", http.StatusOK)
        } else {
                http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        }
}


func favicon(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Access-Control-Allow-Methods", "HEAD")
        w.Header().Set("Content-Security-Policy", "default-src 'self';")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("Referrer-Policy", "same-origin")
	w.Header().Set("Content-Type", "image/x-icon")
        w.Header().Set("Cache-Control", "public, max-age=7776000")
        fmt.Fprintln(w, "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABmJLR0QA/wD/AP+gvaeTAAABQUlEQVRYhe2XQU7CUBCGv6CwFpcmbNyw8AgiSzccgiu4IsYFOxWNG29jIuJCDsEF6gGUhBhCYNGf5hFa89q+kZjwJZPQmb+dH3hv2kIx2sAImCpGQKfgtXJzAyyAZUr0rZu31fwH6AFHQB24Um4BtCwNvBN/015K7Vq1N0sD32pST6kdq/ZlaWD9Xxetb1EpZScA/8ZADXhwjtO2oPvTD3ROMB5/aZoV9yENfOqi5x7alrRRSANZq/tD4asPbiBvfoud74LDkuePg7jwYL0IfW40FxgswgH5t+FdSAM1mYg8GkfEMyDoINrz52RNOGtdgu8kK607cD4/A7dAg/jhE+JB9QScAS9GuoQh2VtraKhLqAKXwMQRTpSrWunce8EcOAGawEy5pnJzQ90Gr3LaVSyVs9YlnLL5jtdRzlq3Z3esAPeurVCaZfkmAAAAAElFTkSuQmCC")
}



func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Security-Policy", "default-src 'self';")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "same-origin")
}


func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bytesRead, err := filterSHA.ReadFrom(f)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bytes read from filter: %d\n", bytesRead)

	router := mux.NewRouter()
        router.HandleFunc("/favicon.ico", favicon).Methods("GET")
        router.HandleFunc("/hashes/sha1/{hash}", checkSHA).Methods("GET")
        router.HandleFunc("/hash/sha1/{hash}", checkSHA).Methods("GET")
        router.HandleFunc("/hash/ntlm/{hash}", checkNTLM).Methods("GET")
        router.HandleFunc("/", index).Methods("GET")

	handler := handlers.CombinedLoggingHandler(os.Stdout, handlers.ProxyHeaders(router))

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         "0.0.0.0:9876",
		Handler:      handler,
	}

	log.Fatal(server.ListenAndServe())
}
