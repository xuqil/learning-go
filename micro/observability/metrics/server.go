package metrics

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"leanring-go/micro/observability"
	"time"
)

type ServerMetricsBuilder struct {
	Namespace string
	Subsystem string
	Port      int
}

func (b *ServerMetricsBuilder) Build() grpc.UnaryServerInterceptor {
	addr := observability.GetOutboundIP()
	if b.Port != 0 {
		addr = fmt.Sprintf("%s:%d", addr, b.Port)
	}
	reqGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      "active_request_cnt",
		Help:      "当前正在处理的请求数量",
		ConstLabels: map[string]string{
			"component": "server",
			"address":   addr,
		},
	}, []string{"service"})
	errCnt := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      "err_request_cnt",
		Help:      "请求错误总数",
		ConstLabels: map[string]string{
			"component": "server",
			"address":   addr,
		},
	}, []string{"service"})
	response := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      "resp_duration",
		Help:      "响应耗时",
		ConstLabels: map[string]string{
			"component": "server",
			"address":   addr,
		},
	}, []string{"service"})
	prometheus.MustRegister(reqGauge)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		reqGauge.WithLabelValues(info.FullMethod).Add(1)
		defer func() {
			reqGauge.WithLabelValues(info.FullMethod).Add(-1)
			if err != nil {
				errCnt.WithLabelValues(info.FullMethod).Add(1)
			}
			response.WithLabelValues(info.FullMethod).
				Observe(float64(time.Now().Sub(startTime).Milliseconds()))
		}()
		resp, err = handler(ctx, req)
		return
	}
}
