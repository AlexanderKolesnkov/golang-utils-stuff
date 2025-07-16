package formatters

import (
	"fmt"
	"math"
)

func FloatMarkDown(number, precision float64) string {
	return fmt.Sprintf(
		"%.0f,%0.f",
		number,
		math.Abs((math.Trunc(math.Abs(number))-math.Abs(number))*math.Pow(10, precision)),
	)
}
