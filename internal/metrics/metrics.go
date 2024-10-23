package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelHandler = "handler"
	labelError   = "error"
)

var (
	totalAcceptedOrders = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "manager_service_total_accepted_orders",
		Help: "total number of accepted orders",
	}, []string{
		labelHandler,
	},
	)

	totalIssuedOrders = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "manager_service_total_issued_orders",
		Help: "total number of issued orders",
	}, []string{
		labelHandler,
	})

	totalRefundedOrders = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "manager_service_total_refunded_orders",
		Help: "total number of refunded orders",
	}, []string{
		labelHandler,
	})

	totalReturnedOrders = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "manager_service_total_returned_orders",
		Help: "total number of returned orders",
	}, []string{
		labelHandler,
	})

	respTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "manager_service_grpc_response_time",
		Help: "gRPC response time",
	}, []string{
		labelHandler,
	})

	totalErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "manager_service_errors_total",
		Help: "total errors",
	}, []string{
		labelHandler,
		labelError,
	})
)

func AddTotalAcceptedOrders(count int, handler string) {
	totalAcceptedOrders.With(prometheus.Labels{
		labelHandler: handler,
	}).Add(float64(count))
}

func AddTotalIssuedOrders(count int, handler string) {
	totalIssuedOrders.With(prometheus.Labels{
		labelHandler: handler,
	}).Add(float64(count))
}

func AddTotalRefundedOrders(count int, handler string) {
	totalRefundedOrders.With(prometheus.Labels{
		labelHandler: handler,
	}).Add(float64(count))
}

func AddTotalReturnedOrders(count int, handler string) {
	totalReturnedOrders.With(prometheus.Labels{
		labelHandler: handler,
	}).Add(float64(count))
}

func ObserveResponseTime(t time.Duration, handler string) {
	respTime.With(prometheus.Labels{
		labelHandler: handler,
	}).Observe(t.Seconds())
}

func IncTotalErrors(handler string, err error) {
	if err == nil {
		return
	}

	totalErrors.With(prometheus.Labels{
		labelHandler: handler,
		labelError:   err.Error(),
	}).Inc()
}
