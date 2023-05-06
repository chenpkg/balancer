package balancer

import (
	"math/rand"
	"time"
)

// Random will randomly select a http server from the server
// 随机算法，负载均衡器将请求随机分发到后端的目标主机中
type Random struct {
	BaseBalancer

	rnd *rand.Rand
}

func init() {
	factories[RandomBalancer] = NewRandom
}

// NewRandom create new Random balancer
func NewRandom(hosts []string) Balancer {
	return &Random{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Balance selects a suitable host according
func (r *Random) Balance(_ string) (string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}
