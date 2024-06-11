package main

import (
	"testing"
)

func Test_nesb_txml(t *testing.T) {
	services := nesb_txml()
	tests := []struct {
		name        string
		dtaName     string
		proCode     string
		wantService NesbTxml
	}{
		// TODO: Add test cases.
		{"COPO_SVR.0100d603", "COPO_SVR", "0100d603", NesbTxml{Tag: "", Service: "ABCSTC0001", Url: ""}},
		{"CMOE_SVR.100CABS004", "CMOE_SVR", "100CABS004", NesbTxml{Tag: "", Service: "CABSTCJ01", Url: "/v1/individual-credit-loan-aum/predict"}},
		{"JSON1_SVR.01EDSSLPOL", "JSON1_SVR", "01EDSSLPOL", NesbTxml{Tag: "JSON_SVR", Service: "JSON1_EDSSJ_T", Url: "/edssol/EDSSLPOL"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dta, ok := services[tt.dtaName]
			if !ok {
				t.Errorf("nesb_txml() did not find dta %s", tt.dtaName)
				return
			}
			service, ok := dta[tt.proCode]
			if !ok {
				t.Errorf("nesb_txml() did not find proCode %s.%s", tt.dtaName, tt.proCode)
			}
			if service.Tag != tt.wantService.Tag || service.Service != tt.wantService.Service || service.Url != tt.wantService.Url {
				t.Errorf("nesb_txml() = %v, want %v", service, tt.wantService)
			}
		})
	}
}
