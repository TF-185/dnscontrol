package diff2

import (
	"strings"
	"testing"

	"github.com/StackExchange/dnscontrol/v3/models"
	"github.com/kylelemons/godebug/diff"
)

func TestNewCompareConfig(t *testing.T) {
	type args struct {
		origin   string
		existing models.Records
		desired  models.Records
		compFn   ComparableFunc
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{
			name: "one",
			args: args{
				origin:   "f.com",
				existing: models.Records{testDataAA1234},
				desired:  models.Records{testDataAA1234clone},
				compFn:   nil,
			},
			want: `
ldata:
  ldata[00]: laba.f.com
             tdata[0]: "A" e(1, 1) d(1, 1)
labelMap: len=1 map[laba.f.com:true]
keyMap:   len=1 map[{laba.f.com A}:true]
existing: ["1.2.3.4"]
desired: ["1.2.3.4"]
origin: f.com
compFn: <nil>
		`,
		},

		{
			name: "cnameAdd",
			args: args{
				origin:   "f.com",
				existing: models.Records{e11mx, d12},
				desired:  models.Records{e11},
				compFn:   nil,
			},
			want: `
ldata:
  ldata[00]: labh.f.com
             tdata[0]: "A" e(1, 1) d(0, 0)
             tdata[1]: "MX" e(1, 1) d(0, 0)
             tdata[2]: "CNAME" e(0, 0) d(1, 1)
labelMap: len=1 map[labh.f.com:true]
keyMap:   len=3 map[{labh.f.com A}:true {labh.f.com CNAME}:true {labh.f.com MX}:true]
existing: ["22 ttt" "1.2.3.4"]
desired: ["labd"]
origin: f.com
compFn: <nil>
		`,
		},

		{
			name: "cnameDel",
			args: args{
				origin:   "f.com",
				existing: models.Records{e11},
				desired:  models.Records{d12, e11mx},
				compFn:   nil,
			},
			want: `
ldata:
  ldata[00]: labh.f.com
             tdata[0]: "CNAME" e(1, 1) d(0, 0)
             tdata[1]: "A" e(0, 0) d(1, 1)
             tdata[2]: "MX" e(0, 0) d(1, 1)
labelMap: len=1 map[labh.f.com:true]
keyMap:   len=3 map[{labh.f.com A}:true {labh.f.com CNAME}:true {labh.f.com MX}:true]
existing: ["labd"]
desired: ["1.2.3.4" "22 ttt"]
origin: f.com
compFn: <nil>
		`,
		},

		{
			name: "two",
			args: args{
				origin:   "f.com",
				existing: models.Records{testDataAA1234, testDataCCa},
				desired:  models.Records{testDataAA1234clone},
				compFn:   nil,
			},
			want: `
ldata:
  ldata[00]: laba.f.com
             tdata[0]: "A" e(1, 1) d(1, 1)
  ldata[01]: labc.f.com
             tdata[0]: "CNAME" e(1, 1) d(0, 0)
labelMap: len=2 map[laba.f.com:true labc.f.com:true]
keyMap:   len=2 map[{laba.f.com A}:true {labc.f.com CNAME}:true]
existing: ["1.2.3.4" "laba"]
desired: ["1.2.3.4"]
origin: f.com
compFn: <nil>
`,
		},

		{
			name: "apex",
			args: args{
				origin:   "f.com",
				existing: models.Records{testDataAA1234, testDataApexMX1aaa},
				desired:  models.Records{testDataAA1234clone, testDataApexMX22bbb},
				compFn:   nil,
			},
			want: `
ldata:
  ldata[00]: f.com
             tdata[0]: "MX" e(1, 1) d(1, 1)
  ldata[01]: laba.f.com
             tdata[0]: "A" e(1, 1) d(1, 1)
labelMap: len=2 map[f.com:true laba.f.com:true]
keyMap:   len=2 map[{f.com MX}:true {laba.f.com A}:true]
existing: ["1.2.3.4" "1 aaa"]
desired: ["1.2.3.4" "22 bbb"]
origin: f.com
compFn: <nil>
`,
		},

		{
			name: "many",
			args: args{
				origin:   "f.com",
				existing: models.Records{testDataAA1234, testDataAMX10a, testDataCCa, testDataEA15},
				desired:  models.Records{testDataAA1234clone, testDataAA12345, testDataAMX20b, d3, d4},
				compFn:   nil,
			},
			want: `
ldata:
  ldata[00]: laba.f.com
             tdata[0]: "A" e(1, 1) d(2, 2)
             tdata[1]: "MX" e(1, 1) d(1, 1)
  ldata[01]: labc.f.com
             tdata[0]: "CNAME" e(1, 1) d(0, 0)
  ldata[02]: labe.f.com
             tdata[0]: "A" e(1, 1) d(2, 2)
labelMap: len=3 map[laba.f.com:true labc.f.com:true labe.f.com:true]
keyMap:   len=4 map[{laba.f.com A}:true {laba.f.com MX}:true {labc.f.com CNAME}:true {labe.f.com A}:true]
existing: ["1.2.3.4" "10 laba" "laba" "10.10.10.15"]
desired: ["1.2.3.4" "1.2.3.5" "20 labb" "10.10.10.95" "10.10.10.96"]
origin: f.com
compFn: <nil>
`,
		},

		{
			name: "all",
			args: args{
				origin:   "f.com",
				existing: models.Records{testDataAA1234, testDataAMX10a, testDataCCa, testDataEA15, e4, e5, e6, e7, e8, e9, e10, e11},
				desired:  models.Records{testDataAA1234clone, testDataAA12345, testDataAMX20b, d3, d4, d5, d6, d7, d8, d9, d10, d11, d12},
				compFn:   nil,
			},
			want: `
ldata:
  ldata[00]: laba.f.com
             tdata[0]: "A" e(1, 1) d(2, 2)
             tdata[1]: "MX" e(1, 1) d(1, 1)
  ldata[01]: labc.f.com
             tdata[0]: "CNAME" e(1, 1) d(0, 0)
  ldata[02]: labe.f.com
             tdata[0]: "A" e(4, 4) d(4, 4)
  ldata[03]: labf.f.com
             tdata[0]: "TXT" e(0, 0) d(1, 1)
  ldata[04]: labg.f.com
             tdata[0]: "NS" e(4, 4) d(4, 4)
  ldata[05]: labh.f.com
             tdata[0]: "CNAME" e(1, 1) d(0, 0)
             tdata[1]: "A" e(0, 0) d(1, 1)
labelMap: len=6 map[laba.f.com:true labc.f.com:true labe.f.com:true labf.f.com:true labg.f.com:true labh.f.com:true]
keyMap:   len=8 map[{laba.f.com A}:true {laba.f.com MX}:true {labc.f.com CNAME}:true {labe.f.com A}:true {labf.f.com TXT}:true {labg.f.com NS}:true {labh.f.com A}:true {labh.f.com CNAME}:true]
existing: ["1.2.3.4" "10 laba" "laba" "10.10.10.15" "10.10.10.16" "10.10.10.17" "10.10.10.18" "10.10.10.15" "10.10.10.16" "10.10.10.17" "10.10.10.18" "labd"]
desired: ["1.2.3.4" "1.2.3.5" "20 labb" "10.10.10.95" "10.10.10.96" "10.10.10.97" "10.10.10.98" "\"foo\"" "10.10.10.10" "10.10.10.15" "10.10.10.16" "10.10.10.97" "1.2.3.4"]
origin: f.com
compFn: <nil>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := NewCompareConfig(tt.args.origin, tt.args.existing, tt.args.desired, tt.args.compFn)
			got := strings.TrimSpace(cc.String())
			tt.want = strings.TrimSpace(tt.want)
			if got != tt.want {
				d := diff.Diff(got, tt.want)
				t.Errorf("NewCompareConfig() = \n%s\n", d)
			}
		})
	}
}

func Test_mkCompareBlobs(t *testing.T) {
	type args struct {
		rc *models.RecordConfig
		f  func(*models.RecordConfig) string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "simple",
			args: args{
				rc: makeRec("label", "A", "1.1.1.1"),
			},
			want:  "1.1.1.1",
			want1: "1.1.1.1 ttl=300",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := mkCompareBlobs(tt.args.rc, tt.args.f)
			if got != tt.want {
				t.Errorf("mkCompareBlobs() got = %q, want %q", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("mkCompareBlobs() got1 = %q, want %q", got1, tt.want1)
			}
		})
	}
}
