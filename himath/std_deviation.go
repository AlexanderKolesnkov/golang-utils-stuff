package himath

import "math"

func CalculateStandardDeviation(values []float64) float64 {
	return StandardDeviation(values, Mean(values))
}

// Mean вычисляет среднее арифметическое последовательности чисел.
func Mean(floatSlice []float64) float64 {
	mean := 0.0
	for _, v := range floatSlice {
		mean += v // Суммирование всех элементов слайса
	}

	return mean / float64(len(floatSlice)) // Деление суммы на количество элементов для получения среднего значения
}

// StandardDeviation вычисляет стандартное отклонение последовательности чисел.
func StandardDeviation(floatSlice []float64, mean float64) float64 {
	res := math.Sqrt(varianceStandard(floatSlice, mean)) // Использование корня квадратного из дисперсии

	if math.IsNaN(res) {
		res = 0 // Возвращает 0, если результат является NaN
	}

	return res
}

// varianceStandard вычисляет дисперсию последовательности чисел по стандартной формуле.
func varianceStandard(floatSlice []float64, mean float64) float64 {
	n := float64(len(floatSlice))
	if n < 2 {
		return 0 // Возвращает 0, если элементов меньше двух, поскольку дисперсия не определена
	}

	variance := 0.0
	for _, v := range floatSlice {
		variance += (v - mean) * (v - mean) // Суммирование квадратов разностей от среднего
	}

	return variance / n // Деление суммы на количество элементов
}
