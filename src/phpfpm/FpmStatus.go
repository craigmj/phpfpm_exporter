package phpfpm

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	AcceptedConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "accepted_connections_count",
		Help: "Number of connections accepted",
	}, []string{"pool"})
	ListenQueue = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "listen_queue",
		Help: "Listen queue size",
	}, []string{"pool", "metric"})
	ProcessesCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "processes_count",
		Help: "Number of processes in the pool",
	}, []string{"pool", "state"})
	MaxChildrenReachedCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "max_children_reached",
		Help: "Maximum number of child processes reached",
	}, []string{"pool"})
)

func init() {
	prometheus.MustRegister(AcceptedConnections)
	prometheus.MustRegister(ListenQueue)
	prometheus.MustRegister(ProcessesCount)
	prometheus.MustRegister(MaxChildrenReachedCount)
}

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
}

func (f *FpmStatus) SetMetrics() error {
	AcceptedConnections.WithLabelValues(f.Pool).Set(float64(f.AcceptedConn))
	ListenQueue.WithLabelValues(f.Pool, "current").Set(float64(f.ListenQueue))
	ListenQueue.WithLabelValues(f.Pool, "max").Set(float64(f.MaxListenQueue))
	ListenQueue.WithLabelValues(f.Pool, "len").Set(float64(f.ListenQueueLen))
	ProcessesCount.WithLabelValues(f.Pool, "idle").Set(float64(f.IdleProcesses))
	ProcessesCount.WithLabelValues(f.Pool, "active").Set(float64(f.ActiveProcesses))
	ProcessesCount.WithLabelValues(f.Pool, "total").Set(float64(f.TotalProcesses))
	ProcessesCount.WithLabelValues(f.Pool, "max_active").Set(float64(f.MaxActiveProcesses))
	MaxChildrenReachedCount.WithLabelValues(f.Pool).Set(float64(f.MaxChildrenReached))
	return nil
}

// GetFpmStatus retrieves the FpmStatus from the server
func GetFpmStatus(url string) (*FpmStatus, error) {
	res, err := http.Get(url)
	if nil != err {
		return nil, err
	}
	defer res.Body.Close()
	status := FpmStatus{}
	if err := json.NewDecoder(res.Body).Decode(&status); nil != err {
		return nil, err
	}
	return &status, nil
}
