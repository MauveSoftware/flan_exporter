// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license
//
// SPDX-License-Identifier: MIT

package main

type NmapRun struct {
	Hosts []HostResult `xml:"host"`
}

type HostResult struct {
	Address struct {
		Addr string `xml:"addr,attr"`
	} `xml:"address"`
	HostNames struct {
		Names []struct {
			Name string `xml:"name,attr"`
		} `xml:"hostname"`
	} `xml:"hostnames"`
	Ports  PortsResult `xml:"ports"`
	Status struct {
		State string `xml:"state,attr"`
	} `xml:"status"`
}

type PortsResult struct {
	Ports []PortResult `xml:"port"`
}

type PortResult struct {
	Protocol string `xml:"protocol,attr"`
	Number   string `xml:"portid,attr"`
	State    struct {
		State string `xml:"state,attr"`
	} `xml:"state"`
	Service struct {
		Name   string `xml:"name,attr"`
		Method string `xml:"method,attr"`
	} `xml:"service"`
	Script ScriptResult `xml:"script"`
}

type ScriptResult struct {
	ID    string `xml:"id,attr"`
	Table struct {
		Tables []Table `xml:"table"`
	} `xml:"table"`
}

type Table struct {
	Elements []TableElement `xml:"elem"`
}

type TableElement struct {
	Key  string `xml:"key,attr"`
	Text string `xml:",chardata"`
}
