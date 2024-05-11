package test

import (
	"github.com/bfun/star"
	"testing"
)

func TestParseProjectFile(t *testing.T) {
	p := star.ParseProjectFile()
	t.Log(p)
	if len(p.PubDtas) == 0 {
		t.Error("Project PubDtas is empty")
	}
	for _, v := range p.PubDtas {
		if v.Name == "" || v.DtaParm == "" || v.Format == "" || v.Service == "" || v.Route == "" {
			t.Error("Project PubDtas parse error")
		}
	}
	if len(p.Apps) == 0 {
		t.Error("Project Apps is empty")
	}
	for _, app := range p.Apps {
		if app.Name == "" || app.DataElem == "" || app.Format == "" {
			t.Error("Project Apps parse error")
		}
		if len(app.SubApps) == 0 {
			t.Log(app.Name, "SubApps is empty")
			continue
		}
		for _, subApp := range app.SubApps {
			if subApp.Name == "" || subApp.AlaParm == "" || subApp.DataElem == "" || subApp.Format == "" || subApp.SvcLogic == "" {
				t.Error("Project SubApps parse error")
			}
			if len(subApp.Dtas) == 0 {
				t.Log(subApp.Name, "Dtas is empty")
				continue
			}
			for _, dta := range subApp.Dtas {
				if dta.Name == "" || dta.DtaParm == "" || dta.Format == "" || dta.Service == "" || dta.Route == "" {
					t.Error("Project Dtas parse error")
				}
			}
		}
	}
}
