package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"

	"github.com/chenpkg/balancer/pkg/proxy"
)

// Config configuration details of balancer
type Config struct {
	SSLCertificateKey   string      `yaml:"ssl_certificate_key"`
	Location            []*Location `yaml:"location"`
	Schema              string      `yaml:"schema"`
	Port                int         `yaml:"port"`
	SSLCertificate      string      `yaml:"ssl_certificate"`
	HealthCheck         bool        `yaml:"tcp_health_check"`
	HealthCheckInterval uint        `yaml:"health_check_interval"`
	MaxAllowed          uint        `yaml:"max_allowed"`
}

// Location routing details of balancer
type Location struct {
	Pattern     string   `yaml:"pattern"`
	ProxyPass   []string `yaml:"proxy_pass"`
	BalanceMode string   `yaml:"balance_mode"`
}

// ReadConfig read configuration from `fileName` file
func ReadConfig(fileName string) (*Config, error) {
	in, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Print print config details
func (c *Config) Print() {
	fmt.Printf("Schema: %s\nPort: %d\nHealth Check: %v\nLocation:\n",
		c.Schema, c.Port, c.HealthCheck)
	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\tProxy Pass: %s\n\tMode: %s\n\n",
			l.Pattern, l.ProxyPass, l.BalanceMode)
	}
}

// Validation verify the configuration details of the balancer
func (c *Config) Validation() error {
	if c.Schema != "http" && c.Schema != "https" {
		return fmt.Errorf("the schema \"%s\" not supported", c.Schema)
	}
	if len(c.Location) == 0 {
		return errors.New("the details of location cannot be null")
	}
	if c.Schema == "https" && (len(c.SSLCertificate) == 0 || len(c.SSLCertificateKey) == 0) {
		return errors.New("the https proxy requires ssl_certificate_key and ssl_certificate")
	}
	if c.HealthCheckInterval < 1 {
		return errors.New("health_check_interval must be greater than 0")
	}
	return nil
}

// maxAllowedMiddleware limit goroutine number
func maxAllowedMiddleware(n uint) mux.MiddlewareFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acquire()
			defer release()
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	err = config.Validation()
	if err != nil {
		log.Fatalf("verify config error: %s", err)
	}

	router := mux.NewRouter()
	for _, l := range config.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, l.BalanceMode)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		// start health check
		if config.HealthCheck {
			httpProxy.HealthCheck(config.HealthCheckInterval)
		}
		router.Handle(l.Pattern, httpProxy)
	}
	if config.MaxAllowed > 0 {
		router.Use(maxAllowedMiddleware(config.MaxAllowed))
	}
	svr := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}

	// print config detail
	config.Print()

	// listen and serve
	if config.Schema == "http" {
		err := svr.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	} else if config.Schema == "https" {
		err := svr.ListenAndServeTLS(config.SSLCertificate, config.SSLCertificateKey)
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	}
}
