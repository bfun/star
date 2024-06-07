package main

import "testing"

func Test_parseDtaSvcFromExpression(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name        string
		args        args
		wantDtaName string
		wantSvcName string
	}{
		// TODO: Add test cases.
		{"CASE1", args{`$stdprocode = "5700", $stdmsgtype = "0200", $__DDTA_NAME = "CNP2_CPUT", $__DSVC_NAME = "02005700"`}, "CNP2_CPUT", "02005700"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDtaName, gotSvcName := parseDtaSvcFromExpression(tt.args.input)
			if gotDtaName != tt.wantDtaName {
				t.Errorf("parseDtaSvcFromExpression() gotDtaName = %v, want %v", gotDtaName, tt.wantDtaName)
			}
			if gotSvcName != tt.wantSvcName {
				t.Errorf("parseDtaSvcFromExpression() gotSvcName = %v, want %v", gotSvcName, tt.wantSvcName)
			}
		})
	}
}
