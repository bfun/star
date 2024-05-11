package test

import (
	"github.com/bfun/star"
	"testing"
)

func TestParseESAdminFile(t *testing.T) {
	e := star.ParseESAdminFile()
	var dtas []star.ESAdminDtaParm
	for _, v := range e.DtaParms {
		t.Log(v.DtaName, v.Ports)
		if len(v.Ports) > 1 {
			dtas = append(dtas, v)
		}
	}
	if len(dtas) == 0 {
		t.Error("all Dta len(Ports) == 0")
	}
}
