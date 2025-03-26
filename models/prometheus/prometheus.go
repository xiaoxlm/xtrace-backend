package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/api"
	prometheus_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"time"
)

type Prometheus struct {
	addr string
}

func NewPrometheus(addr string) *Prometheus {
	return &Prometheus{
		addr: addr,
	}
}

type PromQuery struct {
	Query     string    `json:"query"`               //pql语句
	QueryTime time.Time `json:"queryTime,omitempty"` //查询的时间
}

func (p *Prometheus) QueryMetrics(ctx context.Context, query *PromQuery) (model.Value, error) {
	logrus.Debug("address ====================== ", p.addr)

	conf := api.Config{
		Address: p.addr,
	}
	c, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	if query.QueryTime.IsZero() {
		query.QueryTime = time.Now()
	}

	value, _, err := v1.NewAPI(c).Query(ctx, query.Query, query.QueryTime)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (p *Prometheus) Exec(ctx context.Context, query string) (model.Vector, error) {
	data, err := p.QueryMetrics(ctx, &PromQuery{
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	vector := data.(model.Vector)

	return vector, nil
}

type QueryFormItem struct {
	Start int64  `json:"start"` // 开始时间
	End   int64  `json:"end"`   // 结束时间
	Step  uint   `json:"step"`  // 步长
	Query string `json:"query"` // 查询语句
}

func (p *Prometheus) QueryRange(ctx context.Context, query QueryFormItem) (model.Value, error) {
	conf := api.Config{
		Address: p.addr,
	}
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	r := prometheus_v1.Range{
		Start: time.Unix(query.Start, 0),
		End:   time.Unix(query.End, 0),
		Step:  time.Duration(query.Step) * time.Second,
	}

	resp, _, err := prometheus_v1.NewAPI(client).QueryRange(ctx, query.Query, r)

	return resp, err
}

func (p *Prometheus) BatchQueryRange(ctx context.Context, queries []QueryFormItem) ([]model.Value, error) {
	var list []model.Value

	conf := api.Config{
		Address: p.addr,
	}
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	for _, item := range queries {
		r := prometheus_v1.Range{
			Start: time.Unix(item.Start, 0),
			End:   time.Unix(item.End, 0),
			Step:  time.Duration(item.Step) * time.Second,
		}

		resp, _, err := prometheus_v1.NewAPI(client).QueryRange(ctx, item.Query, r)
		if err != nil {
			return nil, err
		}

		list = append(list, resp)
	}

	return list, nil
}
