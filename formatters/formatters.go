package formatters

import (
	"fmt"
	"math"
	"strings"
)

func countDecimalLen(price float64) int {
	return len(strings.Split(fmt.Sprint(price), ".")[1])
}

func Round(price float64, decimals int) float64 {
	decimalsAmount := countDecimalLen(price)
	if decimalsAmount < decimals {
		decimals = decimalsAmount
	}
	precision := math.Pow10(decimals)
	return math.Round(price*precision) / precision
}
