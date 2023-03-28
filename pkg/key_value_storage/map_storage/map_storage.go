package map_storage

import (
	"fmt"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/key_value_storage"
	"sync"
)

const (
	INITIAL_CAPACITY = 64
)

type MapStorage[KeyT comparable, ValT any] struct {
	mx      *sync.RWMutex
	storage map[KeyT]ValT
}

func NewMapStorage[KeyT comparable, ValT any]() *MapStorage[KeyT, ValT] {
	return &MapStorage[KeyT, ValT]{
		mx:      &sync.RWMutex{},
		storage: make(map[KeyT]ValT, INITIAL_CAPACITY),
	}
}

func (m *MapStorage[KeyT, ValT]) Set(key KeyT, val ValT) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.storage[key] = val
	fmt.Println(m.storage) // TODO: delete
	return nil
}

func (m *MapStorage[KeyT, ValT]) GetByKey(key KeyT) (*ValT, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	val, ok := m.storage[key]
	if !ok {
		return nil, key_value_storage.ErrNoSuchElem
	}
	return &val, nil
}

func (m *MapStorage[KeyT, ValT]) DeleteByKey(key KeyT) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.storage, key)
	return nil
}
