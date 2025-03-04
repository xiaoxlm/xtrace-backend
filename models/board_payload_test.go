package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"testing"
)

func TestGetPanelContent(t *testing.T) {

	type args struct {
		ctx     *ctx.Context
		id      int64
		panelID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPanelContent(tt.args.ctx, tt.args.id, tt.args.panelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPanelContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPanelContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
