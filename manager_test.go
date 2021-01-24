package eventlog

import (
	"context"
	"testing"
	"time"
)

func TestManager_Watch(t *testing.T) {

	type args struct {
		ctx    context.Context
		folder string
		ticker time.Duration
	}
	tests := []struct {
		name    string
		options ManagerOptions
		args    args
		wantErr bool
	}{
		{
			"simple",
			ManagerOptions{
				PoolSize: 10,
				BulkSize: 10,
			},
			args{
				ctx:    context.Background(),
				folder: "./tests",
				ticker: time.Second * 30,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.args.ctx, tt.options)
			if err := m.Watch(tt.args.folder); (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
			}

			m.Wait()
		})
	}
}
