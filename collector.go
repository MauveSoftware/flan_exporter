package main

import (
	"time"

	"github.com/MauveSoftware/flan_exporter/datasource"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const prefix = "flan_"

var (
	reportAgeDesc    *prometheus.Desc
	hostsUpCountDesc *prometheus.Desc
	servicesDesc     *prometheus.Desc
	vulnsDesc        *prometheus.Desc
)

func init() {
	reportAgeDesc = prometheus.NewDesc(prefix+"report_age_seconds", "", nil, nil)
	hostsUpCountDesc = prometheus.NewDesc(prefix+"host_count", "Number of hosts found and scanned", nil, nil)
	servicesDesc = prometheus.NewDesc(prefix+"service_host_count", "Number of hosts per public available service", []string{"name", "port", "protocol"}, nil)
	vulnsDesc = prometheus.NewDesc(prefix+"vuln_host_count", "Number of hosts affected by the CVE", []string{"cve", "cve_level", "is_exploit"}, nil)
}

type collector struct {
	dataSource datasource.DataSource
}

// Describe implements prometheus.Collector interface
func (m *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- reportAgeDesc
	ch <- hostsUpCountDesc
	ch <- servicesDesc
	ch <- vulnsDesc
}

// Collect implements prometheus.Collector interface
func (m *collector) Collect(ch chan<- prometheus.Metric) {
	r, err := m.dataSource.NewestReport()
	if err != nil {
		logrus.Errorf("could not get newest report: %v", err)
		return
	}

	metrics := newReportMetrics()
	for _, f := range r.Files {
		err := metrics.parseReportXML(f.Content)
		if err != nil {
			logrus.Errorf("could not parse report file %s: %v", f.Name, err)
			return
		}
	}

	ch <- prometheus.MustNewConstMetric(reportAgeDesc, prometheus.GaugeValue, float64(time.Since(r.Date).Seconds()))
	ch <- prometheus.MustNewConstMetric(hostsUpCountDesc, prometheus.GaugeValue, float64(metrics.hosts))

	for svc, count := range metrics.services {
		ch <- prometheus.MustNewConstMetric(servicesDesc, prometheus.GaugeValue, float64(count), svc.name, svc.port, svc.protocol)
	}

	for vuln, count := range metrics.vulns {
		exploit := "0"
		if vuln.isExloit {
			exploit = "1"
		}

		ch <- prometheus.MustNewConstMetric(vulnsDesc, prometheus.GaugeValue, float64(count), vuln.cve, vuln.level, exploit)
	}
}
