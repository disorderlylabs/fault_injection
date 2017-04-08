package main

import (
	"fault_injection/http/core"
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"log"
	"net/http"
	"os"
)

const (
	hostPort          = "127.0.0.1:0"
	collectorEndpoint = "http://localhost:10000/collect"
	sameSpan          = true
	traceID128Bit     = true
)

var collector zipkin.Collector

func service1(w http.ResponseWriter, r *http.Request) {
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(r.Header))
	span := core.Check_and_start_span(err, "service_1", spCtx)
	defer span.Finish()

	request := span.BaggageItem("injectfault")
	fmt.Println("service 1, fault request: " + request)

	span.SetBaggageItem("svc1_msg", "hello_from_svc1")
	span.SetBaggageItem("svc1_svcs_invoked", "2|3|4")

	///Requesting service 2
	svc2_req, _ := http.NewRequest("GET", "http://localhost:8080/svc2", nil)

	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc2_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc2_req)
	core.Check_request_error(err, "service_2", r)

	///equesting service 3
	svc3_req, _ := http.NewRequest("GET", "http://localhost:8080/svc3", nil)

	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc3_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc3_req)
	core.Check_request_error(err, "service_3", r)

	///Requesting service 4
	svc4_req, _ := http.NewRequest("GET", "http://localhost:8081/svc4", nil)

	// Inject the trace information into the HTTP Headers.
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc4_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc4_req)
	core.Check_request_error(err, "service_4", r)

	fmt.Println()
} //end service1

func service2(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "service-2", spCtx)
	defer sp.Finish()

	sp.LogKV("hello_from", "service_2")
	sp.LogFields(otlog.String("service_2_status", "ok"))
	fmt.Println()
}

func service3(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "serice-3", spCtx)
	defer sp.Finish()

	sp.LogKV("hello_from", "service_3")
	sp.LogFields(otlog.String("service_3_status", "ok"))

	sp.SetBaggageItem("svc3_msg", "hello_from_svc3")
	sp.SetBaggageItem("svc3_svcs_invoked", "4")

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
	mux.HandleFunc("/svc1", core.Handler_decorator(service1))
	mux.HandleFunc("/svc2", core.Handler_decorator(service2))
	mux.HandleFunc("/svc3", core.Handler_decorator(service3))

	fmt.Printf("Listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))

}
