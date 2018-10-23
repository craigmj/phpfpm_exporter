package phpfpm

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

// VirtualHost contains fields representing a virtual host on a server
type VirtualHost struct {
	// Name associated with the host
	Name string

	// In the case of a unix socket host, this should contain the correct
	// path to `/status`. In the case of an HTTP connect, it could contain
	// the full path including the host and schema:
	// - Example: http://localhost/status?json
	//
	// It is important that, just like from the command line, the URL
	// contains the query parameter `?json` at the end.
	URL string `yaml:"url,omitempty"`

	// The full path the socket. Must include the `unix` scheme in the
	// path name:
	// - Example: unix:///var/run/path/to/php.sock
	FCGI string `yaml:"fcgi,omitempty"`
}

// VirtualHosts represents a collection of virtual hosts which should be
// monitored
type VirtualHosts struct {
	Hosts []VirtualHost
}

// NewVirtualHosts reads the yaml configuration file from the given path
// which contains the virtual hosts being scraped
func NewVirtualHosts(path string) (*VirtualHosts, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "reading vhosts configuration")
	}
	v := &VirtualHosts{}
	err = yaml.Unmarshal(data, v)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshalling yaml")
	}
	return v, nil
}
