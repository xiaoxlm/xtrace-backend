package prometheus

import (
	"context"
	"fmt"
	"github.com/prometheus/common/model"
	"testing"
)

func TestQueryMetricsRate(t *testing.T) {
	ctx := context.Background()
	//query := `avg(hrProcessorLoad{IBN="A",exported_instance="10.10.1.88"})`
	//query := `avg(hrProcessorLoad{IBN="A",exported_instance='10.10.1.88'})`
	//query := `count by (exported_instance,hrDeviceDescr) (hrDeviceDescr{hrDeviceDescr=~".*CPU.*"})`

	//coresQuery := fmt.Sprintf(`hrDeviceDescr{hrDeviceDescr=~".*CPU.*", exported_instance="%s"}`, "10.10.1.88")
	//query := `ifHCOutOctets_total{exported_instance="10.10.1.89"}`
	query := `sysDescr`

	got, err := NewPrometheus("http://10.10.1.84:9090").QueryMetrics(ctx, &PromQuery{
		Query: query,
	})
	if err != nil {
		t.Fatal(err)
	}

	vector := got.(model.Vector)

	var its []string
	for _, sample := range vector {
		for key, value := range sample.Metric {
			if key != LabelIfName {
				continue
			}

			its = append(its, string(value))
		}
	}

	fmt.Println(its)
	fmt.Println(vector)
}
