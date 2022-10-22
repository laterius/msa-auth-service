package types

import "sync"

type Kv struct {
	inlay map[string]string
	lock  sync.RWMutex
}

type ErrNotFound struct {
	error
}

func NewKv() *Kv {
	return &Kv{
		inlay: make(map[string]string),
		lock:  sync.RWMutex{},
	}
}

func (m *Kv) Get(key string) (string, error) {
	m.lock.RLock()
	el, found := m.inlay[key]
	m.lock.RUnlock()
	if !found {
		return "", ErrNotFound{}
	}
	return el, nil
}

func (m *Kv) Set(key, value string) {
	m.lock.Lock()
	m.inlay[key] = value
	m.lock.Unlock()
}

func (m *Kv) Delete(key string) {
	m.lock.Lock()
	delete(m.inlay, key)
	m.lock.Unlock()
}

func (m *Kv) All() map[string]string {
	return m.inlay
}
