package main

import (
	"reflect"
	"testing"
)

func Test_tag_type(t *testing.T) {
	tags := tag_type()
	tests := []struct {
		name     string
		wantTags TagType
	}{
		// TODO: Add test cases.
		{"DCBS", TagType{map[string]string{"stdmsgtype": "mesg_typ", "stdprocode": "tran_cod_cntn_aply_cod", "stdtermtrc": "intor_tran_jnal_num", "std400aqid": "intor_syst_cod", "std400ssys": "frnt_end_syst_cod"}, map[string]string{"std400mgid": "mesg_id", "stdrtninfo": "retn_info_desc"}, map[string]string{"stdrefnum": "/*/serv_call_area/nesb_jnal_num"}, map[string]string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTags := tags[tt.name]; !reflect.DeepEqual(gotTags, tt.wantTags) {
				t.Errorf("tag_type() = %v, want %v", gotTags, tt.wantTags)
			}
		})
	}
}
