package prometheus

import (
	"context"
	"fmt"
)

const (
	LabelGPU = "gpu"
)

type GPU struct {
	prometheus *Prometheus
}

func NewGPU(prometheusAddr string) *GPU {
	return &GPU{
		prometheus: NewPrometheus(prometheusAddr),
	}
}

func (gpu *GPU) GetCPUDevice(ctx context.Context, meta *InstanceMetaInfo) ([]string, error) {
	query := fmt.Sprintf(`DCGM_FI_DEV_GPU_TEMP{IBN="%s",host_ip="%s"}`, meta.IBN, meta.InstanceIP)

	vector, err := gpu.prometheus.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	var interfaces []string
	for _, sample := range vector {
		for key, value := range sample.Metric {
			if key != LabelGPU {
				continue
			}

			interfaces = append(interfaces, string(value))
		}
	}

	return interfaces, nil
}
