package prometheus

import (
	"fmt"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/models/prometheus"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/storage"
	"github.com/toolkits/pkg/logger"
)

const (
	CollectLabelIBN        = "IBN"
	CollectLabelDeviceType = "device_type"
)

type Prometheus struct {
	prometheusAddr string
	redis          storage.Redis
}

func NewPrometheus(ctx *ctx.Context, redis storage.Redis) (*Prometheus, error) {
	sources, err := models.GetDatasourcesGetsBy(ctx, "prometheus", "", "", "enabled")

	if err != nil {
		return nil, err
	}
	if len(sources) < 1 {
		return nil, fmt.Errorf("prometheus data source is empty")
	}

	return &Prometheus{
		prometheusAddr: sources[0].HTTPJson.Url,
		redis:          redis,
	}, nil
}

// TODO 优化 交换机
func (prom *Prometheus) SwitchTarget(ctx *ctx.Context) ([]*models.Target, error) {
	snmp, err := prometheus.NewSNMP(prom.prometheusAddr)
	if err != nil {
		return nil, err
	}

	switchInfos, err := snmp.ListHardWareInfo(ctx.Ctx)
	if err != nil {
		return nil, err
	}

	var list []*models.Target
	for i, info := range switchInfos {
		if info.InstanceIP == "" {
			continue
		}

		list = append(list, &models.Target{
			Id:         int64(100 + i),
			Ident:      info.InstanceIP,
			TagsMap:    map[string]string{"device_type": "Switch", "IBN": info.IBN},
			HostIp:     info.InstanceIP,
			EngineName: "default",
			HostTags:   []string{"device_type=Switch", fmt.Sprintf("IBN=%s", info.IBN)},
			TargetUp:   2,
			MemUtil:    info.Memory.UsagePercent,
			CpuNum:     info.CPU.Cores,
			CpuUtil:    info.CPU.UsagePercent,
			Arch:       "-",
			OS:         "-",
			Offset:     1,
			RemoteAddr: info.InstanceIP,
			Interface:  info.Interface,
		})

		go func() {
			if err = info.SetMeta(ctx.Ctx, prom.redis); err != nil {
				logger.Errorf("SetMeta err: %v", err)
			}
		}()
	}

	return list, nil
}

func (prom *Prometheus) SupplyGPUDevices(ctx *ctx.Context, targetList []*models.Target) error {
	gpuMetrics := prometheus.NewGPU(prom.prometheusAddr)

	for _, target := range targetList {
		if len(target.TagsMap) > 0 {
			ibn, ibnOK := target.TagsMap[CollectLabelIBN]
			_, deviceTypeOK := target.TagsMap[CollectLabelDeviceType]

			if ibnOK && deviceTypeOK {
				devices, err := gpuMetrics.GetCPUDevice(ctx.Ctx, &prometheus.InstanceMetaInfo{
					InstanceIP: target.HostIp,
					IBN:        ibn,
				})
				if err != nil {
					return err
				}

				target.GPUDevices = devices
				continue
			}
		}
	}

	return nil
}