package main

import (
	"encoding/xml"
	"fmt"

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

type host struct {
	name string
	addr string
}

type reportMetrics struct {
	hosts    uint32
	services map[service][]host
	vulns    map[vuln][]host
}

func newReportMetrics() *reportMetrics {
	return &reportMetrics{
		services: make(map[service][]host),
		vulns:    make(map[vuln][]host),
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

		fmt.Println(h)
		ho := host{
			addr: h.Address.Addr,
		}

		if len(h.HostNames.Names) > 0 {
			ho.name = h.HostNames.Names[0].Name
		}

		m.processPorts(h, ho)
	}
}

func (m *reportMetrics) processPorts(h HostResult, ho host) {
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

		m.services[svc] = append(m.services[svc], ho)

		m.processVulns(po, ho)
	}
}

func (m *reportMetrics) processVulns(p PortResult, ho host) {
	if p.Script.ID != "vulners" {
		return
	}

	for _, t := range p.Script.Table.Tables {
		v := m.vulnFromTable(t)
		m.vulns[v] = append(m.vulns[v], ho)
	}
}

func (m *reportMetrics) vulnFromTable(t Table) vuln {
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

	return vuln
}
