package domain

type KeyValueRepository interface {
	Put(key GoroutineId, value string) error

	Get(key GoroutineId) (string, error)

	Delete(key GoroutineId) error
}
