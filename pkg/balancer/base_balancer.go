package balancer

import (
	"sync"
)

type BaseBalancer struct {
	mux   sync.RWMutex
	hosts []string
}

// Add new host to the balancer
func (b *BaseBalancer) Add(host string) {
	b.mux.Lock()
	defer b.mux.Unlock()

	for _, h := range b.hosts {
		if h == host {
			return
		}
	}

	b.hosts = append(b.hosts, host)
}

// Remove new host from the balancer
func (b *BaseBalancer) Remove(host string) {
	b.mux.Lock()
	defer b.mux.Unlock()

	for i, h := range b.hosts {
		if h == host {
			b.hosts = append(b.hosts[:i], b.hosts[i+1:]...)
		}
	}
}

func (b *BaseBalancer) Balance(key string) (string, error) {
	return "", nil
}

func (b *BaseBalancer) Inc(_ string) {}

func (b *BaseBalancer) Done(_ string) {}
