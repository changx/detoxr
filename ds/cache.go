package ds

import (
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type CacheItem struct {
	answer *dns.Msg
	expire int64
}

func (item *CacheItem) Expire() int64 {
	return item.expire
}

func (item *CacheItem) Answer() *dns.Msg {
	return item.answer
}

type CacheManager struct {
	locker *sync.RWMutex
	store  *sync.Map
}

func msgAsKey(m *dns.Msg) string {
	return fmt.Sprintf("%s-%d", m.Question[0].Name, m.Question[0].Qtype)
}

func NewCacheManager() *CacheManager {
	m := CacheManager{}
	m.store = &sync.Map{}
	m.locker = &sync.RWMutex{}

	return &m
}

func (cache *CacheManager) Get(m *dns.Msg) *dns.Msg {
	key := msgAsKey(m)
	cache.locker.RLock()
	defer cache.locker.RUnlock()
	item_, ok := cache.store.Load(key)
	lg.Infof("get cache with key: %s, item=%+v", key, item_)
	if item_ == nil || !ok {
		return nil
	}
	item := item_.(*CacheItem)
	now := time.Now().Unix()
	ttl := item.expire - now
	if ttl > 0 || item.expire < 0 {
		if item.answer.Answer != nil {
			for _, answer := range item.answer.Answer {
				answer.Header().Ttl = uint32(ttl)
			}
		}
		return item.answer
	}
	lg.Debugf("cache expired: %s", key)
	return nil
}

func (cache *CacheManager) Set(m *dns.Msg) {
	minTTL := uint32(0)
	for _, answer := range m.Answer {
		if minTTL == 0 || answer.Header().Ttl < minTTL {
			minTTL = answer.Header().Ttl
		}
	}

	now := time.Now().Unix()
	item := &CacheItem{
		answer: m,
		expire: now + int64(minTTL),
	}
	key := msgAsKey(m)
	cache.locker.Lock()
	defer cache.locker.Unlock()
	cache.store.Store(key, item)
}

func (cache *CacheManager) SetWithTTL(m *dns.Msg, ttl int64) {
	now := time.Now().Unix()
	item := &CacheItem{
		answer: m,
		expire: now + ttl,
	}
	key := msgAsKey(m)
	cache.locker.Lock()
	defer cache.locker.Unlock()
	cache.store.Store(key, item)
}

func (cache *CacheManager) SetCacheWithKey(m *dns.Msg, ttl int64) {
	item := &CacheItem{
		answer: m,
		expire: ttl,
	}
	cache.locker.Lock()
	defer cache.locker.Unlock()
	key := msgAsKey(m)
	lg.Debugf("cache with key: %s", key)
	cache.store.Store(key, item)
}

func (cache *CacheManager) Delete(m *dns.Msg) {
	cache.locker.Lock()
	defer cache.locker.Unlock()
	key := msgAsKey(m)
	cache.store.Delete(key)
}

func (cache *CacheManager) walkAndClean() {
	cache.locker.Lock()
	defer cache.locker.Unlock()
	now := time.Now().Unix()
	keysToDelete := make([]string, 0)
	cache.store.Range(func(key, value any) bool {
		item := value.(*CacheItem)
		if item.expire > 0 && item.expire-now < 0 {
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

func (cache *CacheManager) Each(callback func(item *CacheItem) bool) {
	cache.store.Range(func(key, value any) bool {
		return callback(value.(*CacheItem))
	})
}

func (cache *CacheManager) CacheSuperVisor() {
	for {
		cache.walkAndClean()
		time.Sleep(5 * time.Minute)
	}
}
