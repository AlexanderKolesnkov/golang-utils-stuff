package himath

import (
	"gonum.org/v1/gonum/stat"
	"math"
)

func Correlation(data1, data2 []float64) []float64 {
	if len(data1) != len(data2) {
		panic("DIFFERENT LEN")
	}
	correlation := make([]float64, 0, len(data1))
	for i := 1; i <= len(data1); i++ {
		corr := stat.Correlation(data1[:i], data2[:i], nil)
		if math.IsNaN(corr) {
			correlation = append(correlation, 0)
		} else {
			correlation = append(correlation, corr)
		}

	}
	return correlation
}

func IsPriceBigger3MA(price, ema, emb, ma float64) bool {
	return price > ema && price > emb && price > ma
}

func ZScore(data []float64) []float64 {
	mean := Mean(data)
	stdDev := StandardDeviation(data, mean)

	return FindSpikes(data, mean, stdDev)
}

// FindSpikes FindSpike рассчитывает z-оценку для элемента в массиве по указанному индексу
func FindSpikes(data []float64, mean, stdDev float64) []float64 {
	output := make([]float64, 0, len(data))
	for _, v := range data {
		zScore := (v - mean) / stdDev
		if math.IsNaN((v - mean) / stdDev) {
			zScore = 0
		}

		output = append(output, zScore)
	}
	return output
}

// CentralDerivative
//
// (f'(x) = SMA(i - h) - SMA(i + h)) / 2h
func CentralDerivative(sma []float64) []float64 {
	derivative := make([]float64, len(sma)-2) // Уменьшаем на 2, так как не можем вычислить края
	for i := 1; i < len(sma)-1; i++ {
		derivative[i-1] = (sma[i+1] - sma[i-1]) / 2
	}
	return derivative
}

// varianceGeneral вычисляет дисперсию последовательности чисел по обобщённой формуле.
func varianceGeneral(floatSlice []float64) float64 {
	n := float64(len(floatSlice))
	if n < 2 {
		return 0
	}
	mean := Mean(floatSlice) // Расчёт среднего значения

	variance := 0.0
	for _, v := range floatSlice {
		variance += (v - mean) * (v - mean) // Суммирование квадратов разностей от среднего
	}

	return variance / (n - 1.0) // Деление суммы на количество элементов минус один
}

// LogReturns calculates the logarithmic returns of a slice of prices.
func LogReturns(prices []float64) []float64 {
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		returns[i-1] = math.Log(prices[i] / prices[i-1])
	}
	return returns
}

// PercentFormula вычисляет процентное изменение между начальным и конечным значением.
//
// Формула: (finalValue - initialValue) / (initialValue / 100)
func PercentFormula(finalValue, initialValue float64) float64 {
	if initialValue == 0 {
		return 0 // Возвращает 0, чтобы избежать деления на ноль
	}

	pcnt := (finalValue - initialValue) / (initialValue / 100)

	if finalValue > initialValue {
		pcnt = math.Abs(pcnt)
	}

	if finalValue < initialValue && pcnt > 0 {
		pcnt = pcnt * -1
	}

	return pcnt
}

// DeviancePercent рассчитывает процентное отклонение текущего значения (cur) от среднего значения (mean),
// корректируя это отклонение заданным значением (deviation).
// Если deviation равно 0 (чтобы предотвратить деление на ноль), функция возвращает 0.
// Это может использоваться для измерения, насколько далеко текущее значение находится от среднего, с учетом типичной изменчивости.
func DeviancePercent(cur, mean, deviation float64) float64 {
	if deviation == 0 {
		return 0
	}
	return (cur - mean) / (deviation / 100)
}

func EmaFormula(price, multiplier, EMAp float64) float64 {
	return (price * multiplier) + EMAp*(1-multiplier)
}

func EmaMultiplier(period float64) float64 {
	return 2.0 / (period + 1.0)
}

// CurFormula вычисляет новое значение на основе начального значения и процента изменения.
//
// Формула: pcnt*initialValue/100 + initialValue
func CurFormula(initialValue, pcnt float64) float64 {
	return pcnt*initialValue/100 + initialValue // Применяется для расчёта текущей стоимости с учётом процента изменения
}
