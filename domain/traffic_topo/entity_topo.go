package traffic_topo

import "github.com/ccfos/nightingale/v6/models"

type EntityTopo struct {
	model *models.TrafficTopo

	ibn string
	
	parents []*VoTraffic
}
