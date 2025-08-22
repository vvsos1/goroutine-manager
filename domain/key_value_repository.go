package domain

type KeyValueRepository[KEY any, VALUE any] interface {
	Put(key KEY, value VALUE) error

	Get(key KEY) (VALUE, error)
}
