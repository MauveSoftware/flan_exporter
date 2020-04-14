package main

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

type vuln struct {
	cve      string
	level    string
	isExloit bool
}

type service struct {
	port     string
	protocol string
	name     string
}

type reportMetrics struct {
	hosts    uint32
	services map[service]uint32
	vulns    map[vuln]uint32
}

func newReportMetrics() *reportMetrics {
	return &reportMetrics{
		services: make(map[service]uint32),
		vulns:    make(map[vuln]uint32),
	}
}

func (m *reportMetrics) parseReportXML(b []byte) error {
	r := &NmapRun{}
	err := xml.Unmarshal(b, r)
	if err != nil {
		return err
	}

	if r == nil {
		return errors.Errorf("no NMAP run was found")
	}

	m.processHosts(r)

	return nil
}

func (m *reportMetrics) processHosts(run *NmapRun) {
	for _, h := range run.Hosts {
		if h.Status.State != "up" {
			continue
		}

		m.hosts++
		m.processPorts(h)
	}
}

func (m *reportMetrics) processPorts(h HostResult) {
	for _, po := range h.Ports.Ports {
		if po.State.State != "open" {
			continue
		}

		svc := service{
			port:     po.Number,
			protocol: po.Protocol,
		}

		if po.Service.Method == "probed" {
			svc.name = po.Service.Name
		}

		m.services[svc]++

		m.processVulns(po)
	}
}

func (m *reportMetrics) processVulns(p PortResult) {
	if p.Script.ID != "vulners" {
		return
	}

	for _, t := range p.Script.Table.Tables {
		m.processVuln(t)
	}
}

func (m *reportMetrics) processVuln(t Table) {
	vuln := vuln{}

	for _, elem := range t.Elements {
		switch elem.Key {
		case "id":
			vuln.cve = elem.Text
			break
		case "cvss":
			vuln.level = elem.Text
			break
		case "is_exploit":
			if elem.Text == "true" {
				vuln.isExloit = true
			}
			break
		}
	}

	m.vulns[vuln]++
}
