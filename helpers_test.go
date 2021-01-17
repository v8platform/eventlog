package eventlog

import "testing"

func TestFrom16To10(t *testing.T) {

	tests := []struct {
		name string
		str  string
		want int64
	}{
		{
			"simple",
			"243af7fe23410",
			637370997290000,
		},
		{
			"simple",
			"243af7fe2d050",
			637370997330000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := From16To10(tt.str); got != tt.want {
				t.Errorf("From16To10() = %v, want %v", got, tt.want)
			}
		})
	}
}
