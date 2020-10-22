package main

import (
	"log"
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
