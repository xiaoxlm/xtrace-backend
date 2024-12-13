package loki

import (
	"fmt"
	"testing"
)

func TestQueryLoki(t *testing.T) {
	lokiURL := "http://10.10.1.84:3100"
	query := `{app="dlrover"} |~ "(?i)NCCL INFO" `
	//query := `{app="dlrover"} |= "NCCL ERROR"`
	var start int64 = 1733976840000000000
	var end int64 = 1733987640000000000
	resp, err := QueryLoki(lokiURL, query, start, end)
	if err != nil {
		t.Fatal(err)
	}

	var ncclErrorCount int
	for _, result := range resp.Data.Result {
		ncclErrorCount += len(result.Values)
	}

	fmt.Println("ncclErrorCount:", ncclErrorCount)
}
