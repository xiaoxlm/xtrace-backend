package busi_group_metrics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBusiGroupMetrics_parseExpr(t *testing.T) {

	metrics, err := newBusiGroupMetricsExpr(`100 * (1 - sum by (host_ip)(increase(node_cpu_seconds_total{mode="idle",IBN="$IBN",host_ip="$host_ip"}[15s])) / sum by (host_ip)(increase(node_cpu_seconds_total{IBN="$IBN",host_ip="$host_ip"}[15s])))`, "算网A", []string{"10.10.1.84", "10.10.1.85"})
	if err != nil {
		t.Fatal(err)
	}

	expr := metrics.getParsedExpr()
	fmt.Print(expr)

	assert.Equal(t, `100 * (1 - sum by (host_ip)(increase(node_cpu_seconds_total{mode="idle",IBN="算网A",host_ip=~"10.10.1.84|10.10.1.85"}[15s])) / sum by (host_ip)(increase(node_cpu_seconds_total{IBN="算网A",host_ip=~"10.10.1.84|10.10.1.85"}[15s])))`, expr)
}
