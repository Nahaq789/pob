package vo

type Count struct {
	value int
}

func NewCount(v int) Count {
	return Count{value: v}
}

func (c Count) Decrement() Count {
	v := c.value - 1
	return Count{value: v}
}

func (c Count) Consume(n int) Count {
	v := c.value - n
	return Count{value: v}
}

func (c Count) IsEmpty() bool {
	return c.value == 0
}

func (c Count) Increment() Count {
	return Count{value: c.value + 1}
}

func (c Count) Recover(n int) Count {
	return Count{value: c.value + n}
}

func (c Count) Value() int {
	return c.value
}
