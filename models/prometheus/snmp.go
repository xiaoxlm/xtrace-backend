package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ccfos/nightingale/v6/pkg/util"
	"github.com/ccfos/nightingale/v6/storage"
	"github.com/toolkits/pkg/logger"
	"sync"
)

const (
	LabelIBN              = "IBN"
	LabelExportedInstance = "exported_instance"
	LabelHRDeviceDescr    = "hrDeviceDescr"
	LabelHRStorageDescr   = "hrStorageDescr"
	LabelIfName           = "ifName"
	LabelSysDescr         = "sysDescr"
)

type InstanceMetaInfo struct {
	InstanceIP string
	IBN        string // 网段
}

type SNMP struct {
	prometheus *Prometheus
	instances  []*InstanceMetaInfo
}

func NewSNMP(prometheusAddr string) (*SNMP, error) {
	snmp := &SNMP{
		prometheus: NewPrometheus(prometheusAddr),
	}

	var err error
	snmp.instances, err = snmp.listInstanceIP(context.Background())

	return snmp, err
}

type HardWareInfo struct {
	InstanceMetaInfo
	CPU       *CPUInfo
	Memory    *StorageInfo
	Disk      *StorageInfo
	Interface []string // 端口列表
	SysDesc   *SysDescInfo
}

func (info *HardWareInfo) SetMeta(ctx context.Context, redis storage.Redis) error {
	meta := map[string]interface{}{
		"cpu": map[string]interface{}{
			"cpu_cores": fmt.Sprintf("%d", info.CPU.Cores),
			"name":      info.CPU.DeviceInfo,
		},
		"memory": map[string]interface{}{
			"total": info.Memory.Cap,
		},
		"filesystem": map[string]interface{}{
			"name":       "/dev/sda",
			"mounted_on": "/",
			"kb_size":    info.Disk.Cap,
		},
		"platform": map[string]interface{}{
			"GOOS": info.SysDesc.SysDesc,
		},
	}
	bytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	return storage.MSet(ctx, redis, map[string]interface{}{
		"n9e_extend_meta_" + info.InstanceIP: string(bytes),
	})
}

type CPUInfo struct {
	CPUBasicInfo
	UsagePercent float64
}

type StorageInfo struct {
	StorageBasicInfo
	UsagePercent float64
}

func (snmp *SNMP) ListHardWareInfo(ctx context.Context) ([]*HardWareInfo, error) {
	var (
		cpuMap     map[InstanceMetaInfo]*CPUBasicInfo
		memoryMap  map[InstanceMetaInfo]*StorageBasicInfo
		diskMap    map[InstanceMetaInfo]*StorageBasicInfo
		sysDescMap map[InstanceMetaInfo]*SysDescInfo
	)
	{
		cpuInfo, err := snmp.ListCPUBasicInfo(ctx)
		if err != nil {
			return nil, err
		}
		cpuMap = MapCPUBasicInfo(cpuInfo)

		memoryInfo, err := snmp.ListMemoryBasicInfo(ctx)
		if err != nil {
			return nil, err
		}
		memoryMap = MapStorageBasicInfo(memoryInfo)

		diskInfo, err := snmp.ListDiskBasicInfo(ctx)
		if err != nil {
			return nil, err
		}
		diskMap = MapStorageBasicInfo(diskInfo)

		sysDesc, err := snmp.ListSysDesc(ctx)
		if err != nil {
			return nil, err
		}
		sysDescMap = MapSysDescInfo(sysDesc)
	}

	var hardWareInfos []*HardWareInfo
	for _, ins := range snmp.instances {
		usage, err := snmp.GetUsage(ctx, ins)
		if err != nil {
			return nil, err
		}

		interfaces, err := snmp.GetInterface(ctx, ins)
		if err != nil {
			return nil, err
		}

		hardWareInfos = append(hardWareInfos, &HardWareInfo{
			InstanceMetaInfo: *ins,
			CPU: &CPUInfo{
				CPUBasicInfo: *cpuMap[*ins],
				UsagePercent: usage.CPUPercent,
			},
			Memory: &StorageInfo{
				StorageBasicInfo: *memoryMap[*ins],
				UsagePercent:     usage.MemoryPercent,
			},
			Disk: &StorageInfo{
				StorageBasicInfo: *diskMap[*ins],
				UsagePercent:     usage.DiskPercent,
			},
			Interface: interfaces,
			SysDesc:   sysDescMap[*ins],
		})
	}

	return hardWareInfos, nil
}

type CPUBasicInfo struct {
	InstanceMetaInfo
	DeviceInfo string
	Cores      int
}

func (snmp *SNMP) ListCPUBasicInfo(ctx context.Context) (list []*CPUBasicInfo, error error) {
	query := `count by (exported_instance,hrDeviceDescr,IBN) (hrDeviceDescr{hrDeviceDescr=~".*CPU.*"})`

	vector, err := snmp.prometheus.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, sample := range vector {
		info := &CPUBasicInfo{
			Cores: int(sample.Value),
		}
		for key, value := range sample.Metric {
			if key == LabelIBN {
				info.IBN = string(value)
				continue
			}

			if key == LabelExportedInstance {
				info.InstanceIP = string(value)
				continue
			}

			if key == LabelHRDeviceDescr {
				info.DeviceInfo = string(value)
				continue
			}

		}
		list = append(list, info)
	}

	return
}

type StorageBasicInfo struct {
	InstanceMetaInfo
	Cap string // 容量
}

func (snmp *SNMP) ListMemoryBasicInfo(ctx context.Context) (list []*StorageBasicInfo, error error) {
	query := `hrStorageSize{hrStorageDescr="Physical memory"} * hrStorageAllocationUnits{hrStorageDescr="Physical memory"}`

	return snmp.listStorage(ctx, query)
}

func (snmp *SNMP) ListDiskBasicInfo(ctx context.Context) (list []*StorageBasicInfo, error error) {
	query := `hrStorageSize{hrStorageDescr="/"} * hrStorageAllocationUnits{hrStorageDescr="/"}`

	return snmp.listStorage(ctx, query)
}

type SysDescInfo struct {
	InstanceMetaInfo
	SysDesc string // 操作系统版本
}

func (snmp *SNMP) ListSysDesc(ctx context.Context) ([]*SysDescInfo, error) {
	query := `sysDescr`
	vector, err := snmp.prometheus.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	var list []*SysDescInfo
	for _, sample := range vector {
		info := &SysDescInfo{}

		for key, value := range sample.Metric {
			if key == LabelIBN {
				info.IBN = string(value)
				continue
			}

			if key == LabelExportedInstance {
				info.InstanceIP = string(value)
				continue
			}

			if key == LabelSysDescr {
				info.SysDesc = string(value)
			}
		}

		list = append(list, info)
	}

	return list, nil
}

type Usage struct {
	CPUPercent    float64
	MemoryPercent float64
	DiskPercent   float64
}

func (snmp *SNMP) GetUsage(ctx context.Context, meta *InstanceMetaInfo) (*Usage, error) {
	var (
		usage = &Usage{}
		errs  []error
	)

	// goroutines
	{
		goroutines := sync.WaitGroup{}
		goroutines.Add(3) // cpu, memory, disk
		// cpu
		go func() {
			defer goroutines.Done()

			cpuPercentQuery := fmt.Sprintf(`avg(hrProcessorLoad{IBN="%s",exported_instance='%s'})`, meta.IBN, meta.InstanceIP)
			percent, err := snmp.listUsage(ctx, cpuPercentQuery)
			if err != nil {
				errs = append(errs, err)
				return
			}

			usage.CPUPercent = percent
		}()

		// memory
		go func() {
			defer goroutines.Done()

			memoryPercentQuery := fmt.Sprintf(`(hrStorageUsed{IBN="%s",hrStorageDescr="Physical memory",exported_instance='%s'} * 100 )/(hrStorageSize{IBN="%s",hrStorageDescr="Physical memory",exported_instance='%s'})`, meta.IBN, meta.InstanceIP, meta.IBN, meta.InstanceIP)
			percent, err := snmp.listUsage(ctx, memoryPercentQuery)
			if err != nil {
				errs = append(errs, err)
				return
			}

			usage.MemoryPercent = percent
		}()

		// disk
		go func() {
			defer goroutines.Done()
			diskPercentQuery := fmt.Sprintf(`(hrStorageUsed{IBN="%s",hrStorageDescr="/",exported_instance='%s'} * 100 )/(hrStorageSize{IBN="%s",hrStorageDescr="/",exported_instance='%s'})`, meta.IBN, meta.InstanceIP, meta.IBN, meta.InstanceIP)

			percent, err := snmp.listUsage(ctx, diskPercentQuery)
			if err != nil {
				errs = append(errs, err)
				return
			}

			usage.DiskPercent = percent
		}()

		goroutines.Wait()
	}

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return usage, nil
}

func (snmp *SNMP) GetInterface(ctx context.Context, meta *InstanceMetaInfo) ([]string, error) {
	query := fmt.Sprintf(`ifHCOutOctets_total{IBN="%s",exported_instance="%s"}`, meta.IBN, meta.InstanceIP)

	vector, err := snmp.prometheus.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	var interfaces []string
	for _, sample := range vector {
		for key, value := range sample.Metric {
			if key != LabelIfName {
				continue
			}

			interfaces = append(interfaces, string(value))
		}
	}

	return interfaces, nil
}

func (snmp *SNMP) listStorage(ctx context.Context, query string) (list []*StorageBasicInfo, error error) {
	vector, err := snmp.prometheus.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, sample := range vector {
		info := &StorageBasicInfo{}
		bytes, err := util.ToBytes(sample.Value.String() + "B")
		if err != nil {
			return nil, err
		}
		info.Cap = fmt.Sprintf("%d", bytes) // util.ByteSizeInMi(bytes)

		for key, value := range sample.Metric {
			if key == LabelIBN {
				info.IBN = string(value)
				continue
			}

			if key == LabelExportedInstance {
				info.InstanceIP = string(value)
				continue
			}
		}
		list = append(list, info)
	}

	return
}

func (snmp *SNMP) listUsage(ctx context.Context, query string) (percent float64, err error) {
	vector, err := snmp.prometheus.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	if len(vector) < 1 {
		logger.Warningf("no data of query:%s", query)
		return 0, nil
	}

	return float64(vector[0].Value), nil
}

func (snmp *SNMP) listInstanceIP(ctx context.Context) (list []*InstanceMetaInfo, err error) {
	query := `count by(exported_instance,IBN) (sysUpTime)`

	vector, err := snmp.prometheus.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, sample := range vector {
		info := &InstanceMetaInfo{}
		for key, value := range sample.Metric {
			if key == LabelIBN {
				info.IBN = string(value)
				continue
			}

			if key == LabelExportedInstance {
				info.InstanceIP = string(value)
				continue
			}

		}
		list = append(list, info)
	}

	return
}

func MapCPUBasicInfo(list []*CPUBasicInfo) map[InstanceMetaInfo]*CPUBasicInfo {
	m := make(map[InstanceMetaInfo]*CPUBasicInfo)

	for _, data := range list {
		m[data.InstanceMetaInfo] = data
	}

	return m
}

func MapStorageBasicInfo(list []*StorageBasicInfo) map[InstanceMetaInfo]*StorageBasicInfo {
	m := make(map[InstanceMetaInfo]*StorageBasicInfo)

	for _, data := range list {
		m[data.InstanceMetaInfo] = data
	}

	return m
}

func MapSysDescInfo(list []*SysDescInfo) map[InstanceMetaInfo]*SysDescInfo {
	m := make(map[InstanceMetaInfo]*SysDescInfo)

	for _, data := range list {
		m[data.InstanceMetaInfo] = data
	}

	return m
}
