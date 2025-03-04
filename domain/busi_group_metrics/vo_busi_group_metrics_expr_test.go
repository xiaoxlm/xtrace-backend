package busi_group_metrics

import (
	"fmt"
	"testing"
)

func TestBusiGroupMetrics_parseExpr(t *testing.T) {
	metrics, err := newBusiGroupMetricsExpr(`100 * (1 - avg(rate(node_cpu_seconds_total{IBN="$IBN", mode="idle", host_ip="$host_ip"}[1m])))`, "算网A", []string{"10.10.1.84", "10.10.1.85"})
	if err != nil {
		t.Fatal(err)
	}

	expr := metrics.getParsedExpr()
	fmt.Println(expr)

}
