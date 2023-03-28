package key_value_storage

import (
	"errors"
)

var (
	ErrNoSuchElem = errors.New("there are no elements by such key in storage")
)

type KeyValueStorage[KeyT comparable, ValT any] interface {
	Set(key KeyT, val ValT) error
	GetByKey(key KeyT) (*ValT, error)
}
