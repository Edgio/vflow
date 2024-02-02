package hammer

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsIPFIX struct {
	sendOK          prometheus.Counter
	sendErr         prometheus.Counter
	sendTemplate    prometheus.Counter
	sendTemplateOpt prometheus.Counter
	sendData        prometheus.Counter
}

func NewMetricsIPFIX() *MetricsIPFIX {
	m := MetricsIPFIX{
		sendOK: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ipfix_send_ok"}),
		sendErr: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ipfix_send_err"}),
		sendTemplate: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ipfix_send_template"}),
		sendTemplateOpt: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ipfix_send_template_opt"}),
		sendData: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ipfix_send_data"}),
	}
	prometheus.MustRegister(m.sendOK, m.sendErr, m.sendTemplate, m.sendTemplateOpt, m.sendData)
	return &m
}

type MetricsSFlow struct {
	sendOK   prometheus.Counter
	sendErr  prometheus.Counter
	sendData prometheus.Counter
}

func NewMetricsSFlow() *MetricsSFlow {
	m := MetricsSFlow{
		sendOK: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "sflow_send_ok"}),
		sendErr: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "sflow_send_err"}),
		sendData: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "sflow_send_data"}),
	}
	prometheus.MustRegister(m.sendOK, m.sendErr, m.sendData)
	return &m
}
