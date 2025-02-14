package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"

	pb "pocs/proto"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/reflection"
)

type DemoServiceServer struct {
	pb.UnimplementedDemoServiceServer
}

func newDemoServer() *DemoServiceServer {
	return &DemoServiceServer{}
}

// SayHello implements a interface defined by protobuf.
func (s *DemoServiceServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	customizedCounterMetric.WithLabelValues(request.Name).Inc()
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello %s", request.Name)}, nil
}

var (
	// Create a metrics registry.
	// reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	customizedCounterMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "customized_counter_metric",
			Help: "A customized counter metric",
		},
		[]string{"name"},
	)

	randomGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "random_gauge_metric",
			Help: "A gauge metric with random values",
		},
		[]string{"label"},
	)

	randomHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "random_histogram_metric",
			Help:    "A histogram metric with random values",
			Buckets: prometheus.LinearBuckets(0, 10, 5),
		},
		[]string{"label"},
	)
)

func init() {
	// Register standard server metrics and customized metrics to registry.
	// reg.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	// reg.MustRegister(grpcMetrics, customizedCounterMetric, randomGauge, randomHistogram)
	prometheus.MustRegister(customizedCounterMetric, randomHistogram)
	customizedCounterMetric.WithLabelValues("Test")
}

func GetMetricRegistry() *prometheus.Registry {
	// Return the default registry as a *prometheus.Registry
	return prometheus.DefaultRegisterer.(*prometheus.Registry)
}

// NOTE: Graceful shutdown is missing. Don't use this demo in your production setup.
func main() {
	// Listen an actual port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9093))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Using default Register to create my registry
	GetMetricRegistry().MustRegister(randomGauge)

	// Create a HTTP server for prometheus.
	// httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9091)}
	httpServer := &http.Server{Handler: promhttp.Handler(), Addr: fmt.Sprintf("0.0.0.0:%d", 9091)}

	// Create a gRPC Server with gRPC interceptor.
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	pb.RegisterDemoServiceServer(grpcServer, newDemoServer())
	grpc_prometheus.Register(grpcServer)
	reflection.Register(grpcServer)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start Prometheus HTTP server: %v", err)
		}
	}()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	go func() {
		for {
			value := rand.Float64() * 100
			log.Printf("Observing value: %f", value)
			randomHistogram.WithLabelValues("example").Observe(value)
			randomGauge.WithLabelValues("example").Set(value)
			time.Sleep(5 * time.Second)
		}
	}()

	select {}

	// // Create a new api server.
	// demoServer := newDemoServer()

	// // Register your service.
	// pb.RegisterDemoServiceServer(grpcServer, demoServer)

	// // Initialize all metrics.
	// grpcMetrics.InitializeMetrics(grpcServer)

	// // Start your http server for prometheus.
	// go func() {
	// 	if err := httpServer.ListenAndServe(); err != nil {
	// 		log.Fatal("Unable to start a http server.")
	// 	}
	// }()

	// // Start your gRPC server.
	// log.Fatal(grpcServer.Serve(lis))
	// // Update metrics with random data.
	// go func() {
	// 	for {
	// 		randomGauge.WithLabelValues("example").Set(rand.Float64() * 100)
	// 		randomHistogram.WithLabelValues("example").Observe(rand.Float64() * 100)
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	// // Block forever.
	// select {}
}
