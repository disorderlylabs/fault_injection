package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/opentracing/basictracer-go"
	stdopentracing "github.com/opentracing/opentracing-go"
	//golog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	
	"github.com/go-kit/kit/log"

	"fault_injection/fi"
	"fault_injection/dapperish"
	grpcclient "fault_injection/grpc_client"	
)

func main() {
	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an addsvc. This presumption is reflected in
	// the addcli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.

	var (
		grpcAddr         = flag.String("grpc.addr", "", "gRPC (HTTP) address of addsvc")
		method           = flag.String("method", "sum", "sum, concat")
	)
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "usage: addcli [flags] <a> <b>\n")
		os.Exit(1)
	}

	// This is a demonstration client, which supports multiple tracers.
	// Your clients will probably just use one tracer.
	
	stdopentracing.InitGlobalTracer(basictracer.New(dapperish.NewTrivialRecorder("fi")))
	var tracer = stdopentracing.GlobalTracer()
	

	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.

	var (
		service fi.Service
		err     error
	)
	
	//sp := stdopentracing.StartSpan("client span") // Start a new root span.
	//defer sp.Finish()
	//ctx := stdopentracing.ContextWithSpan(context.Background(), sp)
	//sp.SetBaggageItem("User", "USER")
	//sp.LogFields(golog.String("user text", "hello"))
	//sp.LogFields(golog.Object("ctx", ctx))
	
	if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		service = grpcclient.New(conn, tracer, log.NewNopLogger())
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	
	
	
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch *method {
	case "sum":
		a, _ := strconv.ParseInt(flag.Args()[0], 10, 64)
		b, _ := strconv.ParseInt(flag.Args()[1], 10, 64)
		v, err := service.Sum(context.Background(), int(a), int(b))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%d + %d = %d\n", a, b, v)

	case "concat":
		a := flag.Args()[0]
		b := flag.Args()[1]
		v, err := service.Concat(context.Background(), a, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%q + %q = %q\n", a, b, v)

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}
