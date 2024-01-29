package destination

import (
	"testing"
	"time"
)

func TestUDLReport_String(t *testing.T) {
	tests := []struct {
		name   string
		report UDLReport
		want   string
	}{
		{
			name: "valid",
			report: UDLReport{
				Entries: []UDLEntry{SampleUDLEntry()},
			},
			want: "09002030405.000 854.324972 -806.523053 7049.922417 6.895812284 -2.628367346 -1.133733106\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.report.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func SampleUDLEntry() UDLEntry {
	return UDLEntry{
		Timestamp: time.Date(2009, 01, 02, 03, 04, 05, 06, time.UTC).Format(udlTimeLayout),
		Position: UDLPosition{
			X: "854.324972",
			Y: "-806.523053",
			Z: "7049.922417",
		},
		Velocity: UDLVelocity{
			X: "6.895812284",
			Y: "-2.628367346",
			Z: "-1.133733106",
		},
	}

}
