package balancer

import (
	"sync"

	fibHeap "github.com/starwander/GoFibonacciHeap"
)

type LeastLoad struct {
	mux  sync.RWMutex
	heap *fibHeap.FibHeap
}

func init() {
	factories[LeastLoadBalancer] = NewLeastLoad
}

func (h *host) Tag() interface{} { return h.name }
func (h *host) Key() float64     { return float64(h.load) }

// NewLeastLoad create new LeastLoad balancer
func NewLeastLoad(hosts []string) Balancer {
	ll := &LeastLoad{heap: fibHeap.NewFibHeap()}
	for _, h := range hosts {
		ll.Add(h)
	}
	return ll
}

// Add new host to the balancer
func (l *LeastLoad) Add(hostName string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if ok := l.heap.GetValue(hostName); ok != nil {
		return
	}

	l.heap.InsertValue(&host{hostName, 0})
}

// Remove new host from the balancer
func (l *LeastLoad) Remove(hostName string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if ok := l.heap.GetValue(hostName); ok != nil {
		l.heap.Delete(hostName)
	}
}

// Balance selects a suitable host according
func (l *LeastLoad) Balance(_ string) (string, error) {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if l.heap.Num() == 0 {
		return "", NoHostError
	}

	return l.heap.MinimumValue().Tag().(string), nil
}

// Inc refers to the number of connections to the server `+1`
func (l *LeastLoad) Inc(hostName string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}

	h := l.heap.GetValue(hostName)
	h.(*host).load++
	l.heap.IncreaseKeyValue(h)
}

// Done refers to the number of connections to the server `-1`
func (l *LeastLoad) Done(hostName string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}

	h := l.heap.GetValue(hostName)
	h.(*host).load--
	l.heap.IncreaseKeyValue(h)
}