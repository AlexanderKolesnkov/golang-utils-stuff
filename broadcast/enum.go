package broadcast

type Message struct {
	Symbol    string
	Timeframe string
	StartTime int64
	Confirm   bool
}
