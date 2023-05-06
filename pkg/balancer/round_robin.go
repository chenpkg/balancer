package balancer

// RoundRobin will select the server in turn from the server to proxy
// 轮询算法，负载均衡器将请求依次分发到后端每一台主机中
type RoundRobin struct {
	BaseBalancer

	i uint64
}

func init() {
	factories[R2Balancer] = NewRoundRobin
}

// NewRoundRobin create new RoundRobin balancer
func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
		i: 0,
	}
}

// Balance selects a suitable host according
func (r *RoundRobin) Balance(_ string) (string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	host := r.hosts[r.i%uint64(len(r.hosts))]
	r.i++
	return host, nil
}
