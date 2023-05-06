package balancer

import (
	"hash/crc32"
	"math/rand"
	"sync"
	"time"
)

const _salt = "%#!"

type host struct {
	name string
	load uint64
}

// P2C refer to the power of 2 random choice
// 若请求IP为空，P2C均衡器将随机选择两个代理主机节点，最后选择其中负载量较小的节点
// 若请求IP不为空，P2C均衡器通过对IP地址以及对IP地址加盐进行CRC32哈希计算
// 则会得到两个32bit的值，将其对主机数量进行取模，即CRC32(IP) % len(hosts) 、CRC32(IP + salt) % len(hosts)，
// 最后选择其中负载量较小的节点
type P2C struct {
	mux     sync.RWMutex
	hosts   []*host
	rnd     *rand.Rand
	loadMap map[string]*host
}

func init() {
	factories[P2CBalancer] = NewP2C
}

// NewP2C create new P2C balancer
func NewP2C(hosts []string) Balancer {
	p := &P2C{
		hosts:   []*host{},
		loadMap: make(map[string]*host),
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, h := range hosts {
		p.Add(h)
	}

	return p
}

// Add new host to the balancer
func (p *P2C) Add(hostName string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if _, ok := p.loadMap[hostName]; ok {
		return
	}

	h := &host{name: hostName, load: 0}
	p.hosts = append(p.hosts, h)
	p.loadMap[hostName] = h
}

// Remove new host from the balancer
func (p *P2C) Remove(hostName string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if _, ok := p.loadMap[hostName]; !ok {
		return
	}

	delete(p.loadMap, hostName)

	for i, h := range p.hosts {
		if h.name == hostName {
			p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			return
		}
	}
}

// Balance selects a suitable host according to the key value
func (p *P2C) Balance(key string) (string, error) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if len(p.hosts) == 0 {
		return "", NoHostError
	}

	n1, n2 := p.hash(key)
	hostName := n2
	if p.loadMap[n1].load <= p.loadMap[n2].load {
		hostName = n1
	}

	return hostName, nil
}

// Inc refers to the number of connections to the server `+1`
func (p *P2C) Inc(key string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if h, ok := p.loadMap[key]; ok {
		h.load++
	}
}

// Done refers to the number of connections to the server `-1`
func (p *P2C) Done(key string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if h, ok := p.loadMap[key]; ok && h.load > 0 {
		h.load--
	}
}

func (p *P2C) hash(key string) (string, string) {
	var n1, n2 string
	if len(key) > 0 {
		saltKey := key + _salt
		n1 = p.hosts[crc32.ChecksumIEEE([]byte(key))%uint32(len(p.hosts))].name
		n2 = p.hosts[crc32.ChecksumIEEE([]byte(saltKey))%uint32(len(p.hosts))].name
		return n1, n2
	}
	n1 = p.hosts[p.rnd.Intn(len(p.hosts))].name
	n2 = p.hosts[p.rnd.Intn(len(p.hosts))].name
	return n1, n2
}
