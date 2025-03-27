package traffic_topo

import "github.com/lie-flat-planet/httputil"

type trafficDataModel struct {
	HostIP         string                   `json:"hostIP"`
	PortDevice     string                   `json:"portDevice"`
	ParentIP       string                   `json:"parentIP"`
	InMetricsData  httputil.MetricsFromExpr `json:"inMetricsData"`
	OutMetricsData httputil.MetricsFromExpr `json:"outMetricsData"`

	InExpr  string `json:"-"`
	OutExpr string `json:"-"`
}

type NodeTrafficData struct {
	HostIP   string              `json:"hostIP"`
	Type     string              `json:"type"`
	Traffics []*trafficDataModel `json:"traffics"`
}
