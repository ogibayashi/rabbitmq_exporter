package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	RegisterExporter("overview", newExporterOverview)
}

var overviewMetricDescription = map[string]prometheus.Gauge{
	"object_totals.channels":                       newGauge("channelsTotal", "Total number of open channels."),
	"object_totals.connections":                    newGauge("connectionsTotal", "Total number of open connections."),
	"object_totals.consumers":                      newGauge("consumersTotal", "Total number of message consumers."),
	"object_totals.queues":                         newGauge("queuesTotal", "Total number of queues in use."),
	"object_totals.exchanges":                      newGauge("exchangesTotal", "Total number of exchanges in use."),
	"queue_totals.messages":                        newGauge("queue_messages_total", "Total number ready and unacknowledged messages in cluster."),
	"queue_totals.messages_ready":                  newGauge("queue_messages_ready_total", "Total number of messages ready to be delivered to clients."),
	"queue_totals.messages_unacknowledged":         newGauge("queue_messages_unacknowledged_total", "Total number of messages delivered to clients but not yet acknowledged."),
	"message_stats.publish_details.rate":           newGauge("message_publish_rate", "Message rate of messages published."),
	"message_stats.confirm_details.rate":           newGauge("message_confirm_rate", "Message rate of messages confirmed."),
	"message_stats.return_unroutable_details.rate": newGauge("message_return_rate", "Message rate of messages returned to publisher as unroutable."),
	"message_stats.disk_reads_details.rate":        newGauge("message_disk_reads_rate", "Message rate of messages have been read from disk."),
	"message_stats.disk_writes_details.rate":       newGauge("message_disk_writes_rate", "Message rate of messages have been written to disk."),
	"message_stats.get_details.rate":               newGauge("message_get_rate", "Message rate of messages delivered in acknowledgement mode in response to basic.get."),
	"message_stats.get_no_ack_details.rate":        newGauge("message_get_no_ack_rate", "Message rate of messages delivered in no-acknowledgement mode in response to basic.get."),
	"message_stats.deliver_details.rate":           newGauge("message_deliver_rate", "Message rate of messages delivered in acknowledgement mode to consumers."),
	"message_stats.deliver_no_ack_details.rate":    newGauge("message_deliver_no_ack_rate", "Message rate of messages delivered in no-acknowledgement mode to consumers."),
	"message_stats.redeliver_details.rate":         newGauge("message_redeliver_rate", "Message rate of messages in deliver_get which had the redelivered flag set."),
	"message_stats.ack_details.rate":               newGauge("message_ack_rate", "Message rate of ack."),
	"message_stats.deliver_get_details.rate":       newGauge("message_deliver_get_rate", "Message rate of sum of deliver/get."),
}

type exporterOverview struct {
	overviewMetrics map[string]prometheus.Gauge
}

func newExporterOverview() Exporter {
	return exporterOverview{
		overviewMetrics: overviewMetricDescription,
	}
}

func (e exporterOverview) String() string {
	return "Exporter overview"
}

func (e exporterOverview) Collect(ch chan<- prometheus.Metric) error {
	rabbitMqOverviewData, err := getMetricMap(config, "overview")

	if err != nil {
		return err
	}

	log.WithField("overviewData", rabbitMqOverviewData).Debug("Overview data")
	for key, gauge := range e.overviewMetrics {
		if value, ok := rabbitMqOverviewData[key]; ok {
			log.WithFields(log.Fields{"key": key, "value": value}).Debug("Set overview metric for key")
			gauge.Set(value)
		}
	}

	for _, gauge := range e.overviewMetrics {
		gauge.Collect(ch)
	}
	return nil
}

func (e exporterOverview) Describe(ch chan<- *prometheus.Desc) {
	for _, gauge := range e.overviewMetrics {
		gauge.Describe(ch)
	}

}
