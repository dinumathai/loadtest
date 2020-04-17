package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	minutesToRun int64
)

func init() {
	flag.Int64Var(&minutesToRun, "min", 0, "Total minutes that the program must run")
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
	return []serverConfig{
		serverConfig{
			URL:           "http://localhost:8080?error-code=200",
			RequestPerMin: 1000,
		},
		serverConfig{
			URL:           "http://localhost:8080?error-code=400",
			RequestPerMin: 1000,
		},
		serverConfig{
			URL:           "http://localhost:8080?error-code=500",
			RequestPerMin: 1000,
		},
	}
}

func main() {
	fmt.Println("************* Starting ***********")
	if minutesToRun == 0 {
		minutesToRun = 9223372036854775807
	}
	for i := int64(0); i < minutesToRun; i++ {
		fmt.Printf("Running for %d th minute\n", i)
		go triggerRequest()
		time.Sleep(1 * time.Minute)
	}
	fmt.Println("************* Stopping ***********")
}

func triggerRequest() {
	for _, conf := range getConfig() {
		go triggerRequestForOneConf(conf)
	}
}

func triggerRequestForOneConf(conf serverConfig) {
	sleepTimeMilliSec := int64(60000 / conf.RequestPerMin)
	for i := int(0); i < conf.RequestPerMin; i++ {
		go getURL(conf.URL)
		time.Sleep(time.Duration(sleepTimeMilliSec) * time.Millisecond)
	}
}

func getURL(url string) {
	response, err := http.Get(url)
	if err == nil {
		defer response.Body.Close()
		ioutil.ReadAll(response.Body)
	}
}
