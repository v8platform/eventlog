package eventlog

import (
	"context"
	"github.com/k0kubun/pp"
	"os"
	"reflect"
	"testing"
)

func TestLgpReader_Offset(t *testing.T) {
	type fields struct {
		lgpFile string
	}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewLgpReader(tt.lgpFile)
			if err != nil {
				t.Error(err)
			}

			if tt.offset > 0 {
				_, _ = r.SetOffset(tt.offset)
			}
			var events []*Event
			for i := 0; i < tt.readCount; i++ {
				e := r.Read()
				if e == nil {
					break
				}
				events = append(events, e)
			}

			file, _ := os.OpenFile(tt.lgpFile, os.O_RDONLY, 644)
			//if tt.offset > 0 {
			//	_, _ = file.Seek(tt.offset, io.SeekStart)
			//}
			//
			s, _ := file.Stat()
			pp.Println(s.Size())

			//buf := make([]byte, r.Offset()-tt.offset)
			//
			//_, _ = file.Read(buf)
			//pp.Println(string(buf))

			if got := r.Offset(); got != tt.want {
				t.Errorf("Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLgpReader_StreamRead(t *testing.T) {

	tests := []struct {
		name    string
		lgpFile string
		limit   int
		want    int
	}{
		{
			"1 events",
			"./tests/20210108100000.lgp",
			10,
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := NewLgpReader(tt.lgpFile)

			ctx, _ := context.WithCancel(context.Background())
			stream := r.StreamRead(ctx, tt.limit)
			var events []Event
			for event := range stream {
				events = append(events, event)
			}

			if got := len(events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StreamRead() = %v, want %v", got, tt.want)
			}
		})
	}
}
