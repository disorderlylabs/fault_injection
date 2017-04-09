package main

import (
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	//otlog "github.com/opentracing/opentracing-go/log"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"log"
	"net/http"
	"os"
	"fault_injection/http/core"
)

const (
	hostPort          = "127.0.0.1:10000"
	collectorEndpoint = "http://localhost:10000/collect"
	sameSpan          = true
	traceID128Bit     = true
)

func service1(w http.ResponseWriter, r *http.Request) {
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(r.Header))
	span := core.CheckAndStartSpan(err, "service_1", spCtx)
	defer span.Finish()

	///Requesting service 2
	svc2_req, _ := http.NewRequest("GET", "http://localhost:8080/service2", nil)

	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc2_req.Header))
	core.CheckInjectError(err, r)

	_, err = http.DefaultClient.Do(svc2_req)
	core.CheckRequestError(err, "service2", r)


	//requesting service 3
	svc3_req, _ := http.NewRequest("GET", "http://localhost:8081/service3", nil)

	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc3_req.Header))
	core.CheckInjectError(err, r)

	_, err = http.DefaultClient.Do(svc3_req)
	core.CheckRequestError(err, "service3", r)

	fmt.Println()
} //end service1

func service2(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.CheckAndStartSpan(err, "service2", spCtx)
	defer sp.Finish()

	fmt.Println()
}



func main() {
	collector, err := zipkin.NewHTTPCollector(collectorEndpoint)
	if err != nil {
		fmt.Printf("unable to create a collector: %+v", err)
		os.Exit(-1)
	}

	recorder := zipkin.NewRecorder(collector, false, hostPort, "server1")
	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)

	if err != nil {
		fmt.Printf("unable to create tracer: %+v", err)
		os.Exit(-1)
	}


	var port = flag.Int("port", 8080, "Example app port.")
	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/service1", core.HandlerDecorator(service1))
	mux.HandleFunc("/service2", core.HandlerDecorator(service2))

	fmt.Printf("Listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))

}
