package main

import (
	"flag"

	// @todo this should be a proper importable path
	"github.com/rthornton128/phpfpm_exporter/src/phpfpm"
)

func main() {
	path := flag.String("config.path", "", "path to configuration for multiple virtual hosts")
	statusURL := flag.String("status.url", "http://localhost/status?json", "URL to retrieve fpm-status. Must be in JSON format (end in ?json)")
	scapeInterval := flag.String("scrape.interval", "5m", "Interval between fetching stats from php-fpm")
	listenAddress := flag.String("listen.address", "127.0.0.1:9099", "Address on which to serve metrics")
	flag.Parse()

	c, err := phpfpm.NewConfig(*path, *scapeInterval, *statusURL)
	if err != nil {
		panic(err)
	}
	if err := phpfpm.StartScraper(c); nil != err {
		panic(err)
	}
	phpfpm.WebServer(*listenAddress)
}
