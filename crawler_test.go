package goparkruncrawler

import (
	"context"
	"testing"
)

func TestGetRecentRuns(t *testing.T) {
	type args struct {
		ctx       context.Context
		parkrunID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Maksim",
			args:    args{ctx: context.Background(), parkrunID: "6855386"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRecentRuns(tt.args.ctx, tt.args.parkrunID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecentRuns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < 1 {
				t.Error("Empty results")
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("GetRecentRuns() = %v, want %v", got, tt.want)
			// }
		})
	}
}
