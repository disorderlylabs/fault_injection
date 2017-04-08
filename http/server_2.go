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

func service4(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "service-4", spCtx)
	defer sp.Finish()

	sp.SetBaggageItem("svc4_msg", "hello_from_svc4")
	sp.SetBaggageItem("svc4_svcs_invoked", "5|6|7")

	///Requesting service 5
	svc5_req, _ := http.NewRequest("GET", "http://localhost:8081/svc5", nil)

	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc5_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc5_req)
	core.Check_request_error(err, "service_5", r)

	///Requesting service 6
	svc6_req, _ := http.NewRequest("GET", "http://localhost:8081/svc6", nil)

	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc6_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc6_req)
	core.Check_request_error(err, "service_6", r)

	///Requesting service 7
	svc7_req, _ := http.NewRequest("GET", "http://localhost:8081/svc7", nil)

	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc7_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc7_req)
	core.Check_request_error(err, "service_7", r)

	fmt.Println()

}

func service5(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "service-5", spCtx)
	defer sp.Finish()

	sp.LogFields(otlog.String("service_5_status", "ok"))

	fmt.Println()
}

func service6(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "service-6", spCtx)
	defer sp.Finish()

	sp.LogFields(otlog.String("service_6_status", "ok"))

	///Requesting service 8
	svc8_req, _ := http.NewRequest("GET", "http://localhost:8081/svc8", nil)

	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc8_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc8_req)
	core.Check_request_error(err, "service_8", r)

	///Requesting service 8
	svc9_req, _ := http.NewRequest("GET", "http://localhost:8081/svc9", nil)

	// Inject the trace information into the HTTP Headers.
	err = sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.HTTPHeadersCarrier(svc9_req.Header))
	core.Check_inject_error(err, r)

	_, err = http.DefaultClient.Do(svc9_req)
	core.Check_request_error(err, "service_9", r)

	fmt.Println()
}

func service7(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "service-7", spCtx)
	defer sp.Finish()

	sp.LogFields(otlog.String("service_7_status", "ok"))

	fmt.Println()
}

func service8(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "service-8", spCtx)
	defer sp.Finish()

	sp.LogFields(otlog.String("service_8_status", "ok"))

	fmt.Println()
}

func service9(w http.ResponseWriter, r *http.Request) {
	var sp opentracing.Span
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(r.Header))

	sp = core.Check_and_start_span(err, "service-9", spCtx)
	defer sp.Finish()

	sp.LogFields(otlog.String("service_9_status", "ok"))

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

	var port = flag.Int("port", 8081, "Example app port.")
	opentracing.InitGlobalTracer(tracer)

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/svc4", core.Handler_decorator(service4))
	mux.HandleFunc("/svc5", core.Handler_decorator(service5))
	mux.HandleFunc("/svc6", core.Handler_decorator(service6))
	mux.HandleFunc("/svc7", core.Handler_decorator(service7))
	mux.HandleFunc("/svc8", core.Handler_decorator(service8))
	mux.HandleFunc("/svc9", core.Handler_decorator(service9))

	fmt.Printf("Listening on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(addr, mux))

}
