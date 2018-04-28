package chatroot

import (
	"sync"
)

type MyCookie struct {
	cookie map[string]int
	sync.RWMutex
}

var MyCookies = MyCookie{cookie: make(map[string]int)}

func (m *MyCookie) Set(key string, v int) {
	m.Lock()
	defer m.Unlock()
	m.cookie[key] = v
}

func (m *MyCookie) Get(key string) int {
	m.RLock()
	defer m.RUnlock()
	if val, ok := m.cookie[key]; ok{
		return val
	} else {
		return 0
	}
}

func DelCookie(key string)  {
	delete(MyCookies.cookie, key)
}



