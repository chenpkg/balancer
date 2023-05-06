package balancer

import (
	"hash/crc32"
)

// IPHash will choose a host based on the client's IP address
// IP哈希算法，负载均衡器将请求根据IP地址将其定向分发到后端的目标主机中
// 对主机数量进行取模，即 CRC32(IP) % len(hosts)，则可得到代理的主机
type IPHash struct {
	BaseBalancer
}

func init() {
	factories[IPHashBalancer] = NewIPHash
}

// NewIPHash create new IPHash balancer
func NewIPHash(hosts []string) Balancer {
	return &IPHash{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
	}
}

// Balance selects a suitable host according
func (r *IPHash) Balance(key string) (string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(r.hosts))
	return r.hosts[value], nil
}
