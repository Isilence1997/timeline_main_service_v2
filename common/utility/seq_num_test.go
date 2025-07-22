package utility

import (
	"testing"
)

func TestGetSeqNum(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first := GetSeqNum()
			second := GetSeqNum()
			if first == second {
				t.Logf("firt:%s", first)
				t.Logf("second:%s", second)
				t.Error("GetSeqNum, first == second")
			}
		})
	}
}
