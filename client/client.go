package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	minutesToRunInput int64
	reqPerMinInput    int64
	serverURLInput    string
)

func init() {
	flag.Int64Var(&minutesToRunInput, "min", 0, "Total minutes that the program must run")
	flag.Int64Var(&reqPerMinInput, "rpm", 1000, "Requests per minute for each URL")
	flag.StringVar(&serverURLInput, "url", "", "Server URL to send requests to")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: client -min 10 \n")
	flag.PrintDefaults()
	os.Exit(2)
}

type serverConfig struct {
	URL           string // Supports only get
	RequestPerMin int    // Req must not be more than 60000
}

func getConfig() []serverConfig {
	if serverURLInput != "" {
		return []serverConfig{
			{
				URL:           serverURLInput,
				RequestPerMin: int(reqPerMinInput),
			},
		}
	}
	return []serverConfig{
		{
			URL:           "http://localhost:8080?error-code=200",
			RequestPerMin: 1000,
		}, {
			URL:           "http://localhost:8080?error-code=400",
			RequestPerMin: 1000,
		}, {
			URL:           "http://localhost:8080?error-code=500",
			RequestPerMin: 1000,
		},
	}
}

func main() {
	log.Println("************* Starting ***********")
	var minutesToRun int64 = 9223372036854775807
	if minutesToRunInput != 0 {
		minutesToRun = minutesToRunInput
	}
	for i := int64(0); i < minutesToRun; i++ {
		log.Printf("Running for %d th minute\n", i+1)
		go triggerRequest()
		time.Sleep(1 * time.Minute)
	}
	time.Sleep(10 * time.Second) // wait for 10 seconds for requests to complete
	log.Println("************* Stopping ***********")
}

func triggerRequest() {
	for _, conf := range getConfig() {
		go triggerRequestForOneConf(conf)
	}
}

func triggerRequestForOneConf(conf serverConfig) {
	sleepTimeMilliSec := int64(60000 / conf.RequestPerMin)
	var latencyMutex = sync.RWMutex{}
	latencySum := time.Duration(0)
	for i := int(0); i < conf.RequestPerMin; i++ {
		go func() {
			lat := getURL(conf.URL)
			latencyMutex.Lock()
			latencySum += lat
			latencyMutex.Unlock()
		}()
		time.Sleep(time.Duration(sleepTimeMilliSec) * time.Millisecond)
	}
	average := latencySum / time.Duration(conf.RequestPerMin)
	log.Printf("Average latency for - URL=%s; RPM=%d; Latency=%v\n",
		conf.URL, conf.RequestPerMin, average)
}

func getURL(url string) time.Duration {
	start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		return time.Since(start)
	}
	defer response.Body.Close()
	io.ReadAll(response.Body)
	return time.Since(start)
}
