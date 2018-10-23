package phpfpm

import (
	"time"

	"github.com/golang/glog"
)

func scrape(c *Config) {
	if c.VHosts != nil {
		scrapeMany(c.VHosts)
		return
	}
	scrapeHost(GetFpmStatusHTTP, VirtualHost{URL: c.URL})
}

func scrapeMany(hosts *VirtualHosts) {
	for _, host := range hosts.Hosts {
		callback := GetFpmStatusHTTP
		if host.FCGI != "" {
			callback = GetFpmStatusSocket
		}
		// @todo needs to be in goroutine
		scrapeHost(callback, host)
	}
}

func scrapeHost(callback func(VirtualHost) (*FpmStatus, error), host VirtualHost) {
	s, err := callback(host)
	if nil != err {
		glog.Errorf("ERROR retrieving %s: %s", host.URL, err.Error())
		return
	}
	if err := s.SetMetrics(host.Name); nil != err {
		glog.Errorf("ERROR setting metrics: %s", err.Error())
	}
}

func StartScraper(c *Config) error {
	go func() {
		scrape(c)
		for range time.Tick(c.Interval) {
			scrape(c)
		}
	}()
	return nil
}
