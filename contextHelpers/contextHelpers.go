package contextHelpers

import "context"

func NewBackgroundContextWithDefaultValues(c context.Context) context.Context {
	return NewBackgroundContextWithValues(c, "versions", "dateNow", "Uber-Trace-Id")
}

func NewBackgroundContextWithValues(c context.Context, keys ...string) context.Context {
	ce := context.Background()
	for _, key := range keys {
		ce = context.WithValue(ce, key, c.Value(key))
	}
	return ce
}
