package main

type Tx struct {
	cache     *Cache
	staged    map[string][]byte
	committed bool
}

func (tx *Tx) Commit() error {
	if tx.committed {
		return ErrTransactionAlreadyCommitted
	}

	tx.cache.lock.Lock()
	defer tx.cache.lock.Unlock()

	tx.cache.memMap = tx.staged
	tx.committed = true
	return nil
}

func (tx *Tx) Rollback() {
	tx.committed = true
}

func (tx *Tx) Get(key string) ([]byte, error) {
	v, ok := tx.staged[key]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func (tx *Tx) Set(key string, value []byte) {
	tx.staged[key] = value
}
