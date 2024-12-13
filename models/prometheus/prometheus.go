package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/api"
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

	conf := api.Config{
		Address: p.addr,
	}
	logrus.Debug("address ====================== ", p.addr)
	if query.QueryTime.IsZero() {
		query.QueryTime = time.Now()
	}
	c, err := api.NewClient(conf)
	if err != nil {
		return nil, err
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
