package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/chenpkg/balancer/pkg/balancer"
)

var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-Ip")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
)

var (
	ReverseProxy = "Balancer-Reverse-Proxy"
)

type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer

	mux   sync.RWMutex
	alive map[string]bool
}

// NewHTTPProxy create  new reverse proxy with url and balancer algorithm
func NewHTTPProxy(targetHosts []string, algorithm string) (*HTTPProxy, error) {
	var (
		hosts   = make([]string, 0)
		hostMap = make(map[string]*httputil.ReverseProxy)
		alive   = make(map[string]bool)
	)

	for _, targetHost := range targetHosts {
		u, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(u)

		originDirector := proxy.Director
		// 修改请求
		proxy.Director = func(req *http.Request) {
			originDirector(req)
			req.Header.Set(XProxy, ReverseProxy)
			req.Header.Set(XRealIP, GetIP(req))
		}

		host := GetHost(u)
		alive[host] = true
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}

	lb, err := balancer.Build(algorithm, hosts)
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{
		hostMap: hostMap,
		lb:      lb,
		alive:   alive,
	}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy causes panic :%s", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.(error).Error()))
		}
	}()

	// 根据 IP 选择合适的主机
	host, err := h.lb.Balance(GetIP(r))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("balance error: %s", err.Error())))
		return
	}

	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.hostMap[host].ServeHTTP(w, r)
}
