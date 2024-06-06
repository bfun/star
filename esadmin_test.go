package main

import (
	"testing"
)

func TestParseESAdminFile(t *testing.T) {
	e := ParseESAdminFile()
	var ports, nodes []ESAdminDtaParm
	for _, v := range e.DtaParms {
		for _, p := range v.IPTabItems {
			if p.Port != "" {
				t.Log(v.DtaName, v.IPTabItems)
			}
		}
		if len(v.IPTabItems) > 1 {
			ports = append(ports, v)
		}
		if len(v.Nodes) > 0 {
			t.Log(v.DtaName, v.Nodes)
		}
		if len(v.Nodes) > 1 {
			nodes = append(nodes, v)
		}
	}
	if len(ports) == 0 {
		t.Error("all Dta len(Ports) <= 1")
	} else {
		t.Log(ports)
	}
	if len(nodes) == 0 {
		t.Error("all Dta len(Nodes) <= 1")
	} else {
		t.Log(nodes)
	}
}