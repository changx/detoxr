package ds

import (
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type cacheItem struct {
	answer *dns.Msg
	expire int64
}

type cacheManager struct {
	locker *sync.RWMutex
	store  *sync.Map
}

func msgAsKey(m *dns.Msg) string {
	return fmt.Sprintf("%s-%d", m.Question[0].Name, m.Question[0].Qtype)
}

func NewCacheManager() *cacheManager {
	m := cacheManager{}
	m.store = &sync.Map{}
	m.locker = &sync.RWMutex{}

	return &m
}

func (cache *cacheManager) Get(m *dns.Msg) *dns.Msg {
	key := msgAsKey(m)
	cache.locker.RLock()
	defer cache.locker.RUnlock()
	item_, ok := cache.store.Load(key)
	if item_ == nil || !ok {
		return nil
	}
	item := item_.(*cacheItem)
	now := time.Now().Unix()
	ttl := item.expire - now
	if ttl > 0 {
		for _, answer := range item.answer.Answer {
			answer.Header().Ttl = uint32(ttl)
		}
		return item.answer
	}
	fmt.Printf("!!!!!!!!!!!!  cache expired\n")
	return nil
}

func (cache *cacheManager) Set(m *dns.Msg) {
	minTTL := uint32(0)
	for _, answer := range m.Answer {
		if minTTL == 0 || answer.Header().Ttl < minTTL {
			minTTL = answer.Header().Ttl
		}
	}

	fmt.Printf("item : %#v, minTTL: %d", m, minTTL)
	now := time.Now().Unix()
	item := &cacheItem{
		answer: m,
		expire: now + int64(minTTL),
	}
	key := msgAsKey(m)
	cache.locker.Lock()
	defer cache.locker.Unlock()
	cache.store.Store(key, item)
}

func (cache *cacheManager) walkAndClean() {
	cache.locker.Lock()
	defer cache.locker.Unlock()
	fmt.Printf("cache supervisor working\n")
	now := time.Now().Unix()
	keysToDelete := make([]string, 0)
	cache.store.Range(func(key, value any) bool {
		item := value.(*cacheItem)
		if now-item.expire < 0 {
			keysToDelete = append(keysToDelete, key.(string))
		}
		return true
	})

	if len(keysToDelete) > 0 {
		for _, key := range keysToDelete {
			cache.store.Delete(key)
		}
	}
}

func (cache *cacheManager) CacheSuperVisor() {
	for {
		cache.walkAndClean()
		time.Sleep(5 * time.Minute)
	}
}
