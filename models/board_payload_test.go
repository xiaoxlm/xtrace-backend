package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshalPanel(t *testing.T) {
	j := `{
  "id" : "b93d912c-b3bf-4c6e-a77b-77ea6123ffe9",
  "name" : "CPU Busy",
  "type" : "gauge",
  "links" : [ ],
  "custom" : {
    "calc" : "lastNotNull",
    "textMode" : "value",
    "valueField" : "Value"
  },
  "layout" : {
    "h" : 4,
    "i" : "b93d912c-b3bf-4c6e-a77b-77ea6123ffe9",
    "w" : 5,
    "x" : 4,
    "y" : 0,
    "isResizable" : true
  },
  "options" : {
    "thresholds" : {
      "steps" : [ {
        "type" : "base",
        "color" : "#3FC453",
        "value" : null
      }, {
        "color" : "#FF9919",
        "value" : 60
      }, {
        "color" : "#FF656B",
        "value" : 80
      } ]
    },
    "standardOptions" : {
      "max" : 100,
      "min" : 0,
      "util" : "percent",
      "decimals" : 2
    }
  },
  "targets" : [ {
    "expr" : "100 * (1 - avg(rate(node_cpu_seconds_total{IBN=\"$IBN\", mode=\"idle\", host_ip=\"$host_ip\"}[1m])))",
    "refId" : "A",
    "maxDataPoints" : 240
  } ],
  "version" : "3.0.0",
  "maxPerRow" : 4,
  "description" : "cpu使用率",
  "datasourceCate" : "prometheus",
  "datasourceValue" : "${DS_PROMETHEUS}",
  "transformations" : [ {
    "id" : "organize",
    "options" : { }
  } ]
}`
	p := &Panel{}
	if err := json.Unmarshal([]byte(j), p); err != nil {
		t.Fatal(err)
	}
	fmt.Print(p)
}
