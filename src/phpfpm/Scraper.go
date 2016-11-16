package phpfpm

import (
	"time"

	"github.com/golang/glog"
)

func scrape(url string) {
	s, err := GetFpmStatus(url)
	if nil != err {
		glog.Errorf("ERROR retrieving %s: %s", url, err.Error())
		return
	}
	if err := s.SetMetrics(); nil != err {
		glog.Errorf("ERROR setting metrics: %s", err.Error())
	}
}

func StartScraper(url string, interval string) error {
	pause, err := time.ParseDuration(interval)
	if nil != err {
		return err
	}
	go func() {
		scrape(url)
		for range time.Tick(pause) {
			scrape(url)
		}
	}()
	return nil
}
