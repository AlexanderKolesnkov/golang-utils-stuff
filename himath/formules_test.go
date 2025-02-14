package himath

import (
	"github.com/AlexanderKolesnkov/golang-utils-stuff/consts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcStopLossPcnt(t *testing.T) {
	type args struct {
		stopLoss float64
		avgPrice float64
		leverage float64
		side     string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "[LONG] Price upper safeZone",
			args: args{
				stopLoss: 0.8231,
				avgPrice: 0.8205,
				leverage: 10,
				side:     consts.Buy,
			},
			want: 3.16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcStopLossPcnt(tt.args.stopLoss, tt.args.avgPrice, tt.args.leverage, tt.args.side)

			if !assert.InDelta(t, 3.16, got, 0.01) {
				t.Errorf("CalcStopLossPcnt() = %v, want %v", got, tt.want)
			}
		})
	}
}
