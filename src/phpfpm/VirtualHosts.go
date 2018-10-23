package phpfpm

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

type VirtualHost struct {
	Name string
	URL  string `yaml:"url,omitempty"`
	FCGI string `yaml:"fcgi,omitempty"`
}

type VirtualHosts struct {
	Hosts []VirtualHost
}

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
