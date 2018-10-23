package phpfpm

import (
	"sync"
	"time"

	"github.com/golang/glog"
)

func scrape(c *Config) {
	if c.VHosts != nil {
		scrapeMany(c.VHosts)
		return
	}
	scrapeHost(VirtualHost{URL: c.URL})
}

func scrapeMany(hosts *VirtualHosts) {
	var wg sync.WaitGroup
	for _, host := range hosts.Hosts {
		wg.Add(1)
		go func(host VirtualHost) {
			defer wg.Done()
			scrapeHost(host)
		}(host)
	}
	wg.Wait()
}

func scrapeHost(host VirtualHost) {
	callback := GetFpmStatusHTTP
	if host.FCGI != "" {
		callback = GetFpmStatusSocket
	}
	s, err := callback(host)
	if nil != err {
		glog.Errorf("ERROR retrieving %s: %s", host.URL, err.Error())
		return
	}
	if err := s.SetMetrics(host.Name); nil != err {
		glog.Errorf("ERROR setting metrics: %s", err.Error())
	}
}

// StartScraper being the scraping process using the supplied configuration.
// Config must not be nil
func StartScraper(c *Config) error {
	go func() {
		scrape(c)
		for range time.Tick(c.Interval) {
			scrape(c)
		}
	}()
	return nil
}
