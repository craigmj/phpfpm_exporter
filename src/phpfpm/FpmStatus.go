package phpfpm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	fcgiclient "github.com/tomasen/fcgi_client"
)

var (
	acceptedConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "phpfpm_acceptedconnections_count",
		Help: "Number of connections accepted",
	}, []string{"app", "pool"})
	listenQueue = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "phpfpm_listenqueue_size",
		Help: "Listen queue size",
	}, []string{"app", "pool", "metric"})
	processesCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "phpfpm_processes_count",
		Help: "Number of processes in the pool",
	}, []string{"app", "pool", "state"})
	maxChildrenReachedCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "phpfpm_maxchildren_count",
		Help: "Maximum number of child processes reached",
	}, []string{"app", "pool"})
	slowRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "phpfpm_slow_requests",
		Help: "Slow requests",
	}, []string{"app", "pool"})
)

// FCGIDialTimeout holds the max timeout for socket connections
const FCGIDialTimeout = time.Duration(time.Second * 5)

func init() {
	prometheus.MustRegister(acceptedConnections)
	prometheus.MustRegister(listenQueue)
	prometheus.MustRegister(processesCount)
	prometheus.MustRegister(maxChildrenReachedCount)
	prometheus.MustRegister(slowRequests)
	prometheus.Unregister(prometheus.NewProcessCollector(os.Getpid(), ""))
	prometheus.Unregister(prometheus.NewGoCollector())
}

// FpmStatus represents the JSON data exported by php-fpm
type FpmStatus struct {
	Pool               string `json:"pool"`
	ProcessManager     string `json:"process manager"`
	StartTime          int64  `json:"start time"`
	StartSince         int64  `json:"start since"`
	AcceptedConn       int64  `json:"accepted conn"`
	ListenQueue        int64  `json:"listen queue"`
	MaxListenQueue     int64  `json:"max listen queue"`
	ListenQueueLen     int64  `json:"listen queue len"`
	IdleProcesses      int64  `json:"idle processes"`
	ActiveProcesses    int64  `json:"active processes"`
	TotalProcesses     int64  `json:"total processes"`
	MaxActiveProcesses int64  `json:"max active processes"`
	MaxChildrenReached int64  `json:"max children reached"`
	SlowRequests       int64  `json:"slow requests"`
}

// SetMetrics assigns a new set of metrics for the given host for export to
// Prometheus.
func (f *FpmStatus) SetMetrics(host string) error {
	acceptedConnections.WithLabelValues(host, f.Pool).Set(float64(f.AcceptedConn))
	listenQueue.WithLabelValues(host, f.Pool, "current").Set(float64(f.ListenQueue))
	listenQueue.WithLabelValues(host, f.Pool, "max").Set(float64(f.MaxListenQueue))
	listenQueue.WithLabelValues(host, f.Pool, "len").Set(float64(f.ListenQueueLen))
	processesCount.WithLabelValues(host, f.Pool, "idle").Set(float64(f.IdleProcesses))
	processesCount.WithLabelValues(host, f.Pool, "active").Set(float64(f.ActiveProcesses))
	processesCount.WithLabelValues(host, f.Pool, "total").Set(float64(f.TotalProcesses))
	processesCount.WithLabelValues(host, f.Pool, "max_active").Set(float64(f.MaxActiveProcesses))
	maxChildrenReachedCount.WithLabelValues(host, f.Pool).Set(float64(f.MaxChildrenReached))
	slowRequests.WithLabelValues(host, f.Pool).Set(float64(f.SlowRequests))
	return nil
}

// GetFpmStatusHTTP retrieves the FpmStatus from the server using the HTTP
// protocal
func GetFpmStatusHTTP(h VirtualHost) (*FpmStatus, error) {
	res, err := http.Get(h.URL)
	if nil != err {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received status code: %d", res.StatusCode)
	}
	status := FpmStatus{}
	if err := json.NewDecoder(res.Body).Decode(&status); nil != err {
		return nil, err
	}
	return &status, nil
}

// GetFpmStatusSocket retrieves the FpmStatus from the server using the FCGI
// socket protocal
func GetFpmStatusSocket(h VirtualHost) (*FpmStatus, error) {
	if h.URL == "" {
		h.URL = "/status"
	}
	u, err := url.Parse(h.FCGI)
	if err != nil {
		return nil, errors.WithMessage(err, "parsing FCGI socket path")
	}
	fcgi, err := fcgiclient.DialTimeout(u.Scheme, u.Path, FCGIDialTimeout)
	if err != nil {
		return nil, errors.WithMessage(err, "dialing FCGI socket")
	}
	defer fcgi.Close()
	env := map[string]string{
		"SCRIPT_FILENAME": h.URL,
		"SCRIPT_NAME":     h.URL,
		"QUERY_STRING":    "json",
		"REMOTE_ADDR":     "127.0.0.1",
		"SERVER_SOFTWARE": "go / phpfpm_exporter",
	}
	res, err := fcgi.Get(env)
	if err != nil {
		return nil, errors.WithMessage(err, "fetching from FCGI socket")
	}
	defer res.Body.Close()
	status := FpmStatus{}
	if err := json.NewDecoder(res.Body).Decode(&status); nil != err {
		return nil, err
	}
	return &status, nil
}
