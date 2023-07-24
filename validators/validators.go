package validators

import (
	"time"
)

func IsZeroValue(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v == ""
	case int, float64:
		return v == 0
	case time.Time:
		return v.IsZero()
	default:
		return false
	}
}
