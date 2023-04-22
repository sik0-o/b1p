package proxy

import "sync"

type List struct {
	sync.RWMutex
	list map[string]struct{}
}

func NewList(proxies []string) *List {
	pl := &List{
		list: make(map[string]struct{}),
	}

	pl.Lock()
	for _, proxy := range proxies {
		pl.list[proxy] = struct{}{}
	}
	pl.Unlock()

	return pl
}

func (pl *List) Pick() string {
	pl.Lock()
	defer pl.Unlock()

	for p, _ := range pl.list {
		delete(pl.list, p)
		return p
	}

	return ""
}

func (pl *List) Fetch() map[string]struct{} {
	pl.RLock()
	defer pl.RUnlock()

	l := map[string]struct{}{}
	for proxy := range pl.list {
		l[proxy] = struct{}{}
	}

	return l
}

func (pl *List) Add(proxy string) {
	pl.Lock()
	defer pl.Unlock()

	pl.list[proxy] = struct{}{}
}

func (pl *List) Remove(proxy string) {
	pl.Lock()
	defer pl.Unlock()

	delete(pl.list, proxy)
}

func (pl *List) Has(proxy string) bool {
	pl.RLock()
	defer pl.RUnlock()

	_, ok := pl.list[proxy]

	return ok
}

func (pl *List) Len() int {
	pl.RLock()
	defer pl.RUnlock()

	return len(pl.list)
}

func (pl *List) Get() string {
	pl.RLock()
	defer pl.RUnlock()

	for proxy, _ := range pl.list {
		return proxy
	}

	return ""
}
