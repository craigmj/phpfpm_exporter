package phpfpm

import (
	"time"

	"github.com/pkg/errors"
)

// Config contains the configuration of phpfpm_exporter
type Config struct {
	VHosts   *VirtualHosts
	Interval time.Duration
	URL      string
}

// NewConfig instantiates a new configuration object
func NewConfig(configPath, interval, url string) (*Config, error) {
	pause, err := time.ParseDuration(interval)
	if nil != err {
		return nil, errors.WithMessage(err, "parsing interval")
	}
	var hosts *VirtualHosts
	if configPath != "" {
		hosts, err = NewVirtualHosts(configPath)
		if err != nil {
			return nil, errors.WithMessage(err, "loading virtual hosts")
		}
	}
	return &Config{
		Interval: pause,
		URL:      url,
		VHosts:   hosts,
	}, nil
}
