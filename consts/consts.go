package consts

const (
	Buy        = "Buy"
	Sell       = "Sell"
	TimeLayout = "2006-01-02 15:04:05"
)

// GetAllTimeframes
//
// returns "1", "5", "15", "30", "60", "240", "D", "W", "M"
func GetAllTimeframes() []string {
	return []string{"1", "5", "15", "30", "60", "240", "D", "W", "M"}
}

func SingleIndicatorHeader() []string {
	return []string{"Time", "Value"}
}
