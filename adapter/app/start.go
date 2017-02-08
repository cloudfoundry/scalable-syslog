package app

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/scalable-syslog/adapter/internal/service"
	"github.com/cloudfoundry-incubator/scalable-syslog/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// StartAdapter starts the health endpoint and gRPC service.
func StartAdapter(opts ...AdapterOption) (actualHealth, actualService string) {
	conf := setupConfig(opts)

	actualHealth = startHealthServer(conf.healthAddr)
	actualService = startAdapterService(conf.serviceAddr, conf.serviceCreds)

	return actualHealth, actualService
}

func startHealthServer(hostport string) string {
	l, err := net.Listen("tcp", hostport)
	if err != nil {
		log.Fatalf("Unable to setup Health endpoint (%s): %s", hostport, err)
	}

	server := http.Server{
		Addr:         hostport,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"drainCount": 0}`))
	})

	go func() {
		log.Fatalf("Health server closing: %s", server.Serve(l))
	}()

	return l.Addr().String()
}

func startAdapterService(hostport string, creds credentials.TransportCredentials) string {
	lis, err := net.Listen("tcp", hostport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	adapterService := service.New()
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
	)
	v1.RegisterAdapterServer(grpcServer, adapterService)

	go func() {
		log.Fatalf("failed to serve: %v", grpcServer.Serve(lis))
	}()

	return lis.Addr().String()
}

// AdapterOption is a type that will manipulate a config
type AdapterOption func(c *config)

// WithHealthAddr sets the address for the health endpoint to bind to.
func WithHealthAddr(addr string) func(*config) {
	return func(c *config) {
		c.healthAddr = addr
	}
}

// WithServiceAddr sets the address for the gRPC service to bind to.
func WithServiceAddr(addr string) func(*config) {
	return func(c *config) {
		c.serviceAddr = addr
	}
}

// WithServiceTLSConfig sets the TLS config for the adapter TLS mutual auth.
func WithServiceTLSConfig(cfg *tls.Config) func(*config) {
	return func(c *config) {
		c.serviceCreds = credentials.NewTLS(cfg)
	}
}

type config struct {
	healthAddr   string
	serviceAddr  string
	serviceCreds credentials.TransportCredentials
}

func setupConfig(opts []AdapterOption) *config {
	conf := config{
		healthAddr:  ":8080",
		serviceAddr: ":443",
	}

	for _, o := range opts {
		o(&conf)
	}

	return &conf
}
