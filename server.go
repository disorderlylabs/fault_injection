package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentracing/basictracer-go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
		
	"fault_injection/fi"
	"fault_injection/pb"
	"fault_injection/dapperish"
)

func main() {
	var (
		grpcAddr         = flag.String("grpc.addr", ":8085", "gRPC (HTTP) listen address")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")


	// Business domain.
	var service fi.Service
	{
		service = fi.NewBasicService()
	}
	
	stdopentracing.InitGlobalTracer(basictracer.New(dapperish.NewTrivialRecorder("fi")))
	var tracer = stdopentracing.GlobalTracer()
	serverSpan := tracer.StartSpan("Server")
	defer serverSpan.Finish()

	// Endpoint domain.
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = fi.MakeSumEndpoint(service)
		sumEndpoint = opentracing.TraceServer(tracer, "Sum")(sumEndpoint)
	}
	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = fi.MakeConcatEndpoint(service)
		concatEndpoint = opentracing.TraceServer(tracer, "Concat")(concatEndpoint)
	}
	
	
	endpoints := fi.Endpoints{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
	
	
	

	// Mechanical domain.
	errc := make(chan error)
	ctx := context.Background()

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	
	// gRPC transport.
	go func() {
		logger := log.NewContext(logger).With("transport", "gRPC")

		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := fi.MakeGRPCServer(ctx, endpoints, tracer, logger)
		s := grpc.NewServer()
		pb.RegisterAddServer(s, srv)

		logger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()


	// Run!
	logger.Log("exit", <-errc)
}
