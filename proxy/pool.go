package proxy

import (
	"errors"
	"sync"
	"time"
)

type Pool struct {
	sync.RWMutex

	fresh     *List
	actual    *List
	used      *List
	blacklist *List

	stats map[string]int

	recycleTicker     *time.Ticker
	recycleDoneChan   chan bool
	recycleChan       chan bool
	recycleThreashold int

	blacklistChan chan bool
}

func NewPool(proxies []string) *Pool {
	pl := &Pool{
		fresh:     NewList(nil),
		actual:    NewList(nil),
		used:      NewList(nil),
		blacklist: NewList(nil),
		stats:     make(map[string]int),

		recycleDoneChan:   make(chan bool),
		recycleThreashold: 3,

		blacklistChan: make(chan bool),
	}

	for _, proxy := range proxies {
		pl.fresh.Add(proxy)
	}

	return pl
}

func (pp *Pool) recycleProxies(force bool) error {
	// disable check for force flag
	if !force && pp.actual.Len() > 10 {
		return nil
	}

	if pp.fresh.Len() > 0 {
		for proxy, _ := range pp.fresh.Fetch() {
			pp.actual.Add(proxy)
			pp.fresh.Remove(proxy)
		}
	}

	// disable check for force flag
	if !force && pp.actual.Len() > 20 {
		return nil
	}

	stats := pp.Stats()

	if pp.used.Len() == 0 {
		if pp.actual.Len() == 0 {
			return errors.New("empty used ProxyList")
		} else {
			return nil
		}
	}

	// filter where use <3
	for proxy, _ := range pp.used.Fetch() {
		if v, ok := stats[proxy]; !ok || v <= pp.recycleThreashold {
			pp.actual.Add(proxy)
			pp.used.Remove(proxy)
		}
	}
	if pp.actual.Len() > 0 {
		// update threashold
		pp.recycleThreashold++
	}

	return nil
}

func (pp *Pool) BlacklistProxy(proxy string) {
	pp.blacklist.Add(proxy)
	pp.blacklistChan <- true
}

func (pp *Pool) filterFreshList() {
	if pp.fresh.Len() > 0 {
		for proxy, _ := range pp.fresh.Fetch() {
			if pp.blacklist.Has(proxy) {
				pp.fresh.Remove(proxy)
			}
		}
	}
}

func (pp *Pool) filterUsedList() {
	if pp.fresh.Len() > 0 {
		for proxy, _ := range pp.fresh.Fetch() {
			if pp.blacklist.Has(proxy) {
				pp.fresh.Remove(proxy)
			}
		}
	}
}

func (pp *Pool) StartRecylce() error {
	if pp.recycleTicker != nil {
		return errors.New("recycle already started")
	}

	pp.recycleTicker = time.NewTicker(15 * time.Second)
	go func() {
		for {
			select {
			case <-pp.recycleDoneChan:
				return
			case <-pp.recycleTicker.C:
				pp.recycleProxies(false)
			case <-pp.recycleChan:
				pp.recycleProxies(false)
			case <-pp.blacklistChan:

				pp.filterFreshList()
				pp.filterUsedList()
			}
		}
	}()

	return nil
}

func (pp *Pool) StopRecycle() {
	if pp.recycleTicker != nil {
		pp.recycleTicker.Stop()
		pp.recycleTicker = nil
	}

	pp.recycleDoneChan <- true
}

func (pp *Pool) Stats() map[string]int {
	pp.RLock()
	defer pp.RUnlock()

	stats := make(map[string]int)
	for k, v := range pp.stats {
		stats[k] = v
	}

	return stats
}

func (pp *Pool) pick() string {
	if proxy := pp.actual.Pick(); proxy != "" {
		return proxy
	} else if proxy := pp.fresh.Pick(); proxy != "" {
		return proxy
	}

	return ""
}

func (pp *Pool) Add(proxy string) {
	pp.RLock()
	_, isNotFresh := pp.stats[proxy]
	pp.RUnlock()

	if isNotFresh {
		pp.actual.Add(proxy)
	} else {
		pp.fresh.Add(proxy)
	}
}

func (pp *Pool) Proxy() (string, error) {
	if proxy := pp.pick(); proxy != "" {
		pp.Lock()
		if v, ok := pp.stats[proxy]; !ok {
			pp.stats[proxy] = 1
		} else {
			pp.stats[proxy] = v + 1
		}
		pp.Unlock()

		return proxy, nil
	}

	pp.recycleChan <- true

	return "", errors.New("no available proxy in pool")
}

func (pp *Pool) Return(proxy string) {
	if proxy == "" {
		return
	}

	if pp.blacklist.Has(proxy) {
		return
	}

	pp.used.Add(proxy)
}
