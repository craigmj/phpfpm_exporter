package main

import (
	"flag"
	// "fmt"

	"phpfpm"
)

func main() {
	statusUrl := flag.String("status.url", "http://localhost/status?json", "URL to retrieve fpm-status. Must be in JSON format (end in ?json)")
	listenAddress := flag.String("listen.address", "127.0.0.1:9099", "Address on which to serve metrics")
	scrapeInterval := flag.String("scrape.interval", "5m", "Interval between fetching stats from php-fpm")
	flag.Parse()

	if err := phpfpm.StartScraper(*statusUrl, *scrapeInterval); nil != err {
		panic(err)
	}
	phpfpm.WebServer(*listenAddress)
}
