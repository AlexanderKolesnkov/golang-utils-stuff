package himath

import "github.com/AlexanderKolesnkov/golang-utils-stuff/consts"

// CalcStopLossPcnt calculate pcnt of stopLoss
//
// Formula:  (stopLoss - avgPrice) / (avgPrice / 100) * leverage
func CalcStopLossPcnt(stopLoss, avgPrice, leverage float64, side string) float64 {
	pcnt := (stopLoss - avgPrice) / (avgPrice / 100) * leverage
	if side == consts.Sell {
		return -pcnt
	}
	return pcnt
}

// CalcIM calculate initial margin
//
// Formula: value / leverage
func CalcIM(value, leverage float64) float64 {
	return value / leverage
}

// CalcPnlPcnt calculate pnl in percent
//
// Formula: unrealisedPnl / (Value / 100)
func CalcPnlPcnt(unrealisedPnl, IM float64) float64 {
	return unrealisedPnl / (IM / 100)
}
