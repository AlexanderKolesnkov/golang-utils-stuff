package himath

import (
	"fmt"
	"math"
)

// PctChange вычисляет процентные изменения цен.
//
// Эта функция принимает срез значений цен и возвращает срез процентных изменений между
// последовательными ценами. Процентное изменение рассчитывается как
// (текущая цена - предыдущая цена) / предыдущая цена.
//
// Пример:
//
// prices := []float64{100, 105, 110, 120}
// pctChanges := PctChange(prices)
// // pctChanges теперь содержит: [0.05, 0.047619047619047616, 0.09090909090909091]
//
// Параметры:
//   - prices: Срез float64, содержащий значения цен, для которых нужно рассчитать процентные изменения.
//
// Возвращаемое значение:
//   - Срез float64, содержащий процентные изменения между последовательными ценами.
//
// Примечание:
// Если входной срез содержит менее двух элементов, функция вернет пустой срез, так как нет достаточного количества данных для расчета изменений.
//
// Ошибок не возникает.
func PctChange(prices []float64) []float64 {
	var pctChanges []float64
	for i := 1; i < len(prices); i++ {
		change := (prices[i] - prices[i-1]) / prices[i-1]
		pctChanges = append(pctChanges, change)
	}
	return pctChanges
}

/*
RoundNumber округляет число с плавающей точкой (float64) в зависимости от значения целой части числа.

Правила округления:
- Если целая часть числа не равна нулю, функция округляет число до двух знаков после запятой.
- Если целая часть числа равна нулю, функция округляет число до пяти знаков после запятой.

Параметры:
  - number (float64): Входное число, которое необходимо округлить.

Возвращаемое значение:
  - float64: Округленное число в соответствии с заданными правилами округления.

Примеры использования:

	number1 := 123.456789
	rounded1 := RoundNumber(number1)
	// rounded1 будет 123.46, так как целая часть (123) не равна нулю.

	number2 := 0.123456789
	rounded2 := RoundNumber(number2)
	// rounded2 будет 0.12346, так как целая часть равна нулю.

	number3 := 100.0
	rounded3 := RoundNumber(number3)
	// rounded3 будет 100.00, так как целая часть (100) не равна нулю.

	number4 := 0.000123456
	rounded4 := RoundNumber(number4)
	// rounded4 будет 0.00012, так как целая часть равна нулю.

Ограничения:
  - Функция поддерживает только числа, представленные типом float64.
  - В случае чисел, где округление до меньшего количества знаков после запятой приводит к потерям точности, результат может немного отличаться от ожидаемого значения.
*/
func RoundNumber(number float64) float64 {
	intPart := math.Floor(math.Abs(number))

	if intPart != 0 {
		return math.Round(number*100) / 100
	}
	return math.Round(number*10000000) / 10000000
}

/*
RoundToDecimal округляет число с плавающей точкой (float64) до заданного количества знаков после запятой.

Параметры:
  - number (float64): Входное число, которое необходимо округлить.
  - precision (int): Количество знаков после запятой, до которых нужно округлить число.
    Должно быть неотрицательным числом. Если `precision` отрицательное число, функция возвращает ошибку.

Возвращаемое значение:
  - float64: Округленное число, соответствующее заданной точности.
  - error: Ошибка, если параметр `precision` отрицательный.

Правила округления:
- Если `precision` неотрицателен, число округляется до указанного количества знаков после запятой.
- Если `precision` отрицателен, функция возвращает ошибку с сообщением "precision cannot be negative".

Примеры использования:

	number1 := 123.456789
	precision1 := 2
	rounded1, err1 := RoundToDecimal(number1, precision1)
	// rounded1 будет 123.46, так как число округлено до двух знаков после запятой.

	number2 := 1.987654321
	precision2 := 4
	rounded2, err2 := RoundToDecimal(number2, precision2)
	// rounded2 будет 1.9877, так как число округлено до четырех знаков после запятой.

	number3 := 0.000123456
	precision3 := 6
	rounded3, err3 := RoundToDecimal(number3, precision3)
	// rounded3 будет 0.000123, так как число округлено до шести знаков после запятой.

	number4 := 123.456
	precision4 := -1
	rounded4, err4 := RoundToDecimal(number4, precision4)
	// err4 будет ошибкой с сообщением "precision cannot be negative", так как параметр precision отрицательный.

Ограничения:
  - Функция поддерживает только числа, представленные типом float64.
  - Параметр `precision` должен быть неотрицательным. В противном случае возвращается ошибка.
*/
func RoundToDecimal(number float64, precision int) (float64, error) {
	if precision < 0 {
		return 0, fmt.Errorf("precision cannot be negative")
	}

	multiplier := math.Pow(10, float64(precision))
	roundedNumber := math.Round(number*multiplier) / multiplier

	return roundedNumber, nil
}

func RoundTwoDecimal(number float64) float64 {
	multiplier := math.Pow(10, float64(2))
	roundedNumber := math.Round(number*multiplier) / multiplier

	return roundedNumber
}

func RoundFiveDecimal(number float64) float64 {
	multiplier := math.Pow(10, float64(5))
	roundedNumber := math.Round(number*multiplier) / multiplier

	return roundedNumber
}

func CalcInitialMargin(qty, entryPrice, leverage float64) float64 {
	return qty * entryPrice / leverage
}

func CrossZeroFromAbove(prev, curr float64) bool {
	if prev > 0 && curr < 0 {
		return true
	}

	return false
}

func CrossZeroFromBelow(prev, curr float64) bool {
	if prev < 0 && curr > 0 {
		return true
	}

	return false
}

func CrossTwoLines(prev1, curr1, prev2, curr2 float64) bool {
	if (prev1 < prev2 && curr1 > curr2) || (prev1 > prev2 && curr1 < curr2) {
		return true
	}

	return false
}

func Between(value, min, max float64) bool {
	if value >= min && value <= max {
		return true
	}
	return false
}
