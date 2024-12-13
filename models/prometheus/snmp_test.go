package prometheus

import (
	"context"
	"github.com/ccfos/nightingale/v6/pkg/util"
	"testing"
)

func TestSNMP_ListHardWareInfo(t *testing.T) {
	snmp, err := NewSNMP("http://10.10.1.84:9090")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	got, err := snmp.ListHardWareInfo(ctx)
	if err != nil {
		t.Fatal(err)
	}

	util.LogJSON(got)
}
