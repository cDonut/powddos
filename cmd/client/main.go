package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cdonut/powddos/pkg/pow"
)

var attempts uint64 = 1 << 25

type result struct {
	mutex   sync.Mutex
	minTime time.Duration
	maxTime time.Duration
}

func (r *result) update(duration time.Duration) {
	r.mutex.Lock()

	if r.minTime == 0 {
		r.minTime = duration
	}

	if r.maxTime < duration {
		r.maxTime = duration
	} else if r.minTime > duration {
		r.minTime = duration
	}

	r.mutex.Unlock()
}

func main() {
	address := os.Getenv("SERVER_ADDRESS")

	gc, err := strconv.Atoi(os.Getenv("GC_COUNT"))
	if err != nil {
		log.Printf("Failed to parse GC_COUNT variable: %s", err.Error())
	}

	var wg sync.WaitGroup

	client := &http.Client{}
	result := &result{}

	//requests burst
	for i := 0; i < gc; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			start := time.Now()
			w, err := getWisdom(client, address)
			dur := time.Since(start)

			if err != nil {
				log.Printf("Failed request: %s", err.Error())
			} else {
				log.Printf("Got: %s in %s", w, dur)
			}

			result.update(dur)
		}()
	}

	wg.Wait()

	log.Printf("min request time: %s, max request time: %s", result.minTime, result.maxTime)
}

func getWisdom(client *http.Client, address string) (wisdom string, err error) {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return wisdom, fmt.Errorf("Failed to request some wisdom: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return wisdom, fmt.Errorf("Failed to request some wisdom: %w", err)
	}

	defer res.Body.Close()

	//get a wisdom o fail trying
	for {
		if res.StatusCode == http.StatusOK {
			break
		} else {
			c, err := pow.ParseChallenge(res.Header.Get("X-Challenge"))
			if err != nil {
				return wisdom, fmt.Errorf("Failed to parse challenge: %w", err)
			}

			ans, err := c.Solve(attempts)
			if err != nil {
				return wisdom, fmt.Errorf("Failed to solve challenge: %w", err)
			}

			req.Header.Set("X-Solution", ans)

			res, err = client.Do(req)
			if err != nil {
				return wisdom, fmt.Errorf("Failed to request some wisdom: %w", err)
			}
		}
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return wisdom, fmt.Errorf("Failed to read response body: %w", err)
	}

	wisdom = string(bytes)

	return wisdom, nil
}
