package phpfpm

import (
	"encoding/json"
	"strings"
)

// ReadSample ?
func ReadSample() *FpmStatus {
	status := FpmStatus{}
	if err := json.NewDecoder(strings.NewReader(sampleJSON)).Decode(&status); nil != err {
		panic(err)
	}
	return &status
}

var sampleJSON = `{"pool":"www",
"process manager":"dynamic",
"start time":1479299112,
"start since":7472,
"accepted conn":1804516,
"listen queue":1,
"max listen queue":2,
"listen queue len":3,
"idle processes":720,
"active processes":2280,
"total processes":3000,
"max active processes":3000,
"max children reached":21,
"slow requests":0}`
