package formatters

import (
	"math"
)

func Round(price float64, decimals int) float64 {
	precision := math.Pow10(decimals)
	return math.Round(price*precision) / precision
}
