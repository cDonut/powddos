package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/cdonut/powddos/pkg/pow"
)

var (
	requests int64
	levelCap int64
	tsExpire int64
	data     string
	port     string
)

func main() {
	var err error
	port = os.Getenv("SERVER_PORT")
	data = os.Getenv("POW_DATA")

	levelCap, err = strconv.ParseInt(os.Getenv("LEVEL_CAP"), 10, 64)
	if err != nil || levelCap <= 0 {
		log.Print("Can't start wisdom server: unsupported LEVEL_CAP variable")
		os.Exit(1)
	}

	tsExpire, err = strconv.ParseInt(os.Getenv("TS_EXPIRE_SEC"), 10, 64)
	if err != nil || levelCap <= 0 {
		log.Print("Can't start wisdom server: unsupported TS_EXPIRE_SEC variable")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", getWisdom)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), powMiddleware(mux))
	if err != nil {
		log.Printf("Can't start wisdom server: %s", err.Error())
		os.Exit(1)
	}
}

func getWisdom(w http.ResponseWriter, r *http.Request) {
	//work hard for 5 sec
	time.Sleep(5 * time.Second)
	fmt.Fprint(w, "what a wisdom!")
}

func powMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt64(&requests, 1)
			defer atomic.AddInt64(&requests, -1)

			level := requests / levelCap
			if level > 20 {
				level = 20
			}

			if level > 0 && !pow.CheckSolution(r.Header.Get("X-Solution"), level, data, tsExpire) {
				w.Header().Set("X-Challenge", pow.NewChallenge(level, data).String())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
