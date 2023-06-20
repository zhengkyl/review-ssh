package common

type Cacheable interface {
	Review | Film | Show
}

type CacheInfo[T Cacheable] struct {
	Loading bool
	Data    T
}

type Cache[T Cacheable] map[int]CacheInfo[T]

func (c Cache[T]) Get(id int) (bool, bool, T) {
	res, exists := c[id]

	if !exists {
		var t T
		return false, false, t
	}

	return !res.Loading, res.Loading, res.Data
}

func (c Cache[T]) Set(id int, data T) {
	c[id] = CacheInfo[T]{false, data}
}

func (c Cache[T]) SetLoading(id int) {
	var t T
	c[id] = CacheInfo[T]{true, t}
}

func (c Cache[T]) Delete(id int) {
	delete(c, id)
}
