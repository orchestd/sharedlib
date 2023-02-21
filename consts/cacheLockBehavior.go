package consts

type CacheLockBehavior string

const (
	LockAlways       = "LockAlways"
	LockIfIdNotEmpty = "LockIfIdNotEmpty"
	LockOnlyIfExists = "LockOnlyIfExists"
)

func (c CacheLockBehavior) LockAlways() bool {
	return c == LockAlways
}

func (c CacheLockBehavior) LockIfIdNotEmpty() bool {
	return c == LockIfIdNotEmpty
}

func (c CacheLockBehavior) LockOnlyIfExists() bool {
	return c == LockOnlyIfExists
}
