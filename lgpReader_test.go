package eventlog

import (
	"github.com/k0kubun/pp"
	"testing"
	"time"
)

func TestLgpReader_Offset(t *testing.T) {

	tests := []struct {
		name      string
		lgpFile   string
		readCount int
		offset    int64
		want      int64
	}{
		{
			"1 events",
			"./tests/20210108100000.lgp",
			1,
			0,
			157,
		},
		{
			"5 events",
			"./tests/20210108100000.lgp",
			5,
			157,
			714,
		},
		{
			"all",
			"./tests/20210108100000.lgp",
			9999999,
			0,
			1820103, // full file size
		},
		{
			"big",
			"./tests/big/20210117000000.lgp",
			0,
			0,
			1820103, // full file size
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewLgpReader(tt.lgpFile)
			if err != nil {
				t.Error(err)
			}

			if tt.offset > 0 {
				_, _ = r.Seek(tt.offset)
			}

			events, err := r.Read(tt.readCount, 20*time.Second)
			pp.Println("events", len(events))

			if got := r.Offset(); got != tt.want {
				t.Errorf("Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}
